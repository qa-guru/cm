package selenoid

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/aerokube/selenoid/config"
	ctr "github.com/moby/moby/api/types/container"
	img "github.com/moby/moby/api/types/image"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/docker/go-connections/nat"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/mattn/go-colorable"

	"github.com/aerokube/cm/render/rewriter"
	"github.com/fatih/color"
	. "github.com/fvbommel/sortorder"
)

const (
	Latest                = "latest"
	selenoidImage           = selenoidWrapperImage
	selenoidUIImage         = selenoidUIWrapperImage
	videoRecorderImage      = "selenoid/video-recorder:latest-release"
	selenoidContainerName   = "selenoid"
	ggrUIContainerName      = "ggr-ui"
	selenoidUIContainerName = "selenoid-ui"
	overrideHome            = "OVERRIDE_HOME"
	dockerApiVersion        = "DOCKER_API_VERSION"
	selenoidDockerAPI       = "1.55" // matches github.com/moby/moby/client v0.5.0
)

type SelenoidConfig map[string]config.Versions

type DockerConfigurator struct {
	Logger
	ConfigDirAware
	VersionAware
	DownloadAware
	ArgsAware
	EnvAware
	PortAware
	UserNSAware
	LogsAware
	GracefulAware
	Forceable
	GithubBaseUrl        string
	OS                   string
	Arch                 string
	SelenoidBinaryPath   string
	SelenoidUIBinaryPath string
	RegistryUrl          string
	BrowsersJson         string
	docker               *client.Client
	reg                  *registry.Registry
	authConfig           *dockerAuthConfig
	registryHost         string
}

func NewDockerConfigurator(config *LifecycleConfig) (*DockerConfigurator, error) {
	c := &DockerConfigurator{
		Logger:           Logger{Quiet: config.Quiet},
		ConfigDirAware:   ConfigDirAware{ConfigDir: config.ConfigDir},
		VersionAware:     VersionAware{Version: config.Version},
		DownloadAware:    DownloadAware{DownloadNeeded: config.Download},
		ArgsAware:        ArgsAware{Args: config.Args},
		EnvAware:         EnvAware{Env: config.Env},
		PortAware:        PortAware{Port: config.Port},
		UserNSAware:      UserNSAware{UserNS: config.UserNS},
		LogsAware:        LogsAware{DisableLogs: config.DisableLogs},
		GracefulAware:    GracefulAware{Graceful: config.Graceful, GracefulTimeout: config.GracefulTimeout},
		Forceable:        Forceable{Force: config.Force},
		GithubBaseUrl:    config.GithubBaseUrl,
		OS:               config.OS,
		Arch:             config.Arch,
		SelenoidBinaryPath:   config.SelenoidBinary,
		SelenoidUIBinaryPath: config.SelenoidUIBinary,
		RegistryUrl:      config.RegistryUrl,
		BrowsersJson:     config.BrowsersJson,
	}
	if c.Quiet {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}
	err := c.initDockerClient()
	if err != nil {
		return nil, fmt.Errorf("new configurator: %v", err)
	}
	authConfig, err := c.initAuthConfig()
	if err != nil {
		c.Errorf("Failed to load authentication configuration, using default values: %v", err)
	} else {
		c.authConfig = authConfig
	}
	return c, nil
}

func createCompatibleDockerClient(onVersionSpecified, onVersionDetermined, onUsingDefaultVersion func(string)) (*client.Client, error) {
	dockerApiVersionEnv := os.Getenv(dockerApiVersion)
	if dockerApiVersionEnv != "" {
		onVersionSpecified(dockerApiVersionEnv)
	} else {
		maxMajorVersion, maxMinorVersion := parseVersion(client.MaxAPIVersion)
		minMajorVersion, minMinorVersion := parseVersion("1.24")
		for majorVersion := maxMajorVersion; majorVersion >= minMajorVersion; majorVersion-- {
			for minorVersion := maxMinorVersion; minorVersion >= minMinorVersion; minorVersion-- {
				apiVersion := fmt.Sprintf("%d.%d", majorVersion, minorVersion)
				_ = os.Setenv(dockerApiVersion, apiVersion)
				docker, err := client.New(client.FromEnv)
				if err != nil {
					return nil, err
				}
				if isDockerAPIVersionCorrect(docker) {
					onVersionDetermined(apiVersion)
					return docker, nil
				}
				_ = docker.Close()
			}
		}
		onUsingDefaultVersion(client.MaxAPIVersion)
	}
	return client.New(client.FromEnv)
}

func parseVersion(ver string) (int, int) {
	const point = "."
	pieces := strings.Split(ver, point)
	major, err := strconv.Atoi(pieces[0])
	if err != nil {
		return 0, 0
	}
	minor, err := strconv.Atoi(pieces[1])
	if err != nil {
		return 0, 0
	}
	return major, minor
}

func isDockerAPIVersionCorrect(docker *client.Client) bool {
	ctx := context.Background()
	apiInfo, err := docker.ServerVersion(ctx, client.ServerVersionOptions{})
	if err != nil {
		return false
	}
	return apiInfo.APIVersion == docker.ClientVersion()
}

func (c *DockerConfigurator) initDockerClient() error {
	docker, err := createCompatibleDockerClient(
		func(specifiedApiVersion string) {
			c.Pointf("Using Docker API version: %s", specifiedApiVersion)
		},
		func(determinedApiVersion string) {
			c.Pointf("Your Docker API version is %s", determinedApiVersion)
		},
		func(defaultApiVersion string) {
			c.Pointf("Did not manage to determine your Docker API version - using default version: %s", defaultApiVersion)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to init Docker client: %v", err)
	}
	c.docker = docker
	return nil
}

func (c *DockerConfigurator) initAuthConfig() (*dockerAuthConfig, error) {
	authConfigs, err := loadDockerAuthConfigs()
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(c.RegistryUrl)
	if err != nil {
		return nil, err
	}

	registryHost := u.Host
	if c.RegistryUrl != DefaultRegistryUrl {
		c.registryHost = registryHost
	}
	if cfg, ok := authConfigs[registryHost]; ok {
		c.Titlef(`Loaded authentication data for "%s"`, registryHost)
		return &cfg, nil
	}

	return nil, nil
}

func (c *DockerConfigurator) getRegistryClient() *registry.Registry {
	if c.reg != nil {
		return c.reg
	}

	u := strings.TrimSuffix(c.RegistryUrl, "/")
	username, password := "", ""
	if c.authConfig != nil {
		username, password = c.authConfig.Username, c.authConfig.Password
	}
	reg := &registry.Registry{
		URL: u,
		Client: &http.Client{
			Transport: registry.WrapTransport(http.DefaultTransport, u, username, password),
		},
		Logf: func(format string, args ...interface{}) {
			c.Tracef(format, args...)
		},
	}

	if err := reg.Ping(); err != nil {
		c.Errorf("Docker Registry is not available: %v", err)
		return nil
	}

	c.reg = reg
	return reg
}

func (c *DockerConfigurator) Close() error {
	if c.docker != nil {
		return c.docker.Close()
	}
	return nil
}

func (c *DockerConfigurator) Status() {
	selenoidImage := c.getSelenoidImage()
	if selenoidImage != nil {
		c.Pointf("Using Selenoid image: %s (%s)", selenoidImage.RepoTags[0], selenoidImage.ID)
	} else {
		c.Pointf("Selenoid image is not present")
	}
	configPath := getSelenoidConfigPath(c.ConfigDir)
	c.Pointf("Selenoid configuration directory is %s", c.ConfigDir)
	if fileExists(configPath) {
		c.Pointf("Selenoid configuration file is %s", configPath)
	} else {
		c.Pointf("Selenoid is not configured")
	}
	selenoidContainer := c.getSelenoidContainer()
	if selenoidContainer != nil {
		c.Pointf("Selenoid container is running: %s (%s)", selenoidContainerName, selenoidContainer.ID)
	} else {
		c.Pointf("Selenoid container is not running")
	}
}

func (c *DockerConfigurator) UIStatus() {
	selenoidUIImage := c.getSelenoidUIImage()
	if selenoidUIImage != nil {
		c.Pointf("Using Selenoid UI image: %s (%s)", selenoidUIImage.RepoTags[0], selenoidUIImage.ID)
	} else {
		c.Pointf("Selenoid UI image is not present")
	}
	selenoidUIContainer := c.getSelenoidUIContainer()
	if selenoidUIContainer != nil {
		c.Pointf("Selenoid UI container is running: %s (%s)", selenoidUIContainerName, selenoidUIContainer.ID)
	} else {
		c.Pointf("Selenoid UI container is not running")
	}
}

func (c *DockerConfigurator) IsDownloaded() bool {
	return c.getSelenoidImage() != nil
}

func (c *DockerConfigurator) getSelenoidImage() *img.Summary {
	return c.getImage(selenoidWrapperImage, wrapperImageTag)
}

func (c *DockerConfigurator) IsUIDownloaded() bool {
	return c.getSelenoidUIImage() != nil
}

func (c *DockerConfigurator) getSelenoidUIImage() *img.Summary {
	return c.getImage(selenoidUIWrapperImage, wrapperImageTag)
}

func (c *DockerConfigurator) getImage(name string, version string) *img.Summary {
	result, err := c.docker.ImageList(context.Background(), client.ImageListOptions{})
	if err != nil {
		c.Errorf("Failed to list images: %v", err)
		return nil
	}
	return findMatchingImage(result.Items, name, version)
}

func findMatchingImage(images []img.Summary, name string, version string) *img.Summary {
	sort.Slice(images, func(i, j int) bool {
		return images[i].Created > images[j].Created
	})
	for _, img := range images {
		const colon = ":"
		for _, tag := range img.RepoTags {
			nameAndVersion := strings.Split(tag, colon)
			if len(nameAndVersion) >= 2 {
				imageVersion := nameAndVersion[len(nameAndVersion)-1]
				imageName := strings.TrimSuffix(tag, colon+imageVersion)
				if strings.HasSuffix(imageName, name) && (version == "" || version == Latest || version == imageVersion) {
					return &img
				}
			}
		}
	}
	return nil
}

func (c *DockerConfigurator) Download() (string, error) {
	ref, err := c.downloadWrapperImage(selenoidWrapperImage, "failed to pull Selenoid wrapper image")
	if err != nil {
		return "", err
	}
	if err := c.downloadSelenoidBinary(); err != nil {
		return "", err
	}
	return ref, nil
}

func (c *DockerConfigurator) DownloadUI() (string, error) {
	ref, err := c.downloadWrapperImage(selenoidUIWrapperImage, "failed to pull Selenoid UI wrapper image")
	if err != nil {
		return "", err
	}
	if err := c.downloadSelenoidUIBinary(); err != nil {
		return "", err
	}
	return ref, nil
}

func (c *DockerConfigurator) downloadWrapperImage(imageName, errorMessage string) (string, error) {
	return c.downloadImpl(imageName, wrapperImageTag, errorMessage)
}

func (c *DockerConfigurator) downloadImpl(imageName string, version string, errorMessage string) (string, error) {
	if version == Latest {
		latestVersion := c.getLatestImageVersion(imageName)
		if latestVersion != nil {
			version = *latestVersion
		}
	}
	ref := c.getFullyQualifiedImageRef(imageName)
	if version != Latest {
		ref = imageWithTag(ref, version)
	}
	if !c.pullImage(context.Background(), ref) {
		return "", errors.New(errorMessage)
	}
	return ref, nil
}

func (c *DockerConfigurator) getLatestImageVersion(imageName string) *string {
	tags := c.fetchImageTags(imageName)
	if len(tags) > 0 {
		return &tags[0]
	}
	return nil
}

func (c *DockerConfigurator) IsConfigured() bool {
	return fileExists(getSelenoidConfigPath(c.ConfigDir))
}

func (c *DockerConfigurator) Configure() (*SelenoidConfig, error) {
	err := c.createConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}
	if c.BrowsersJson != "" {
		return c.syncBrowsersFromFile(c.BrowsersJson)
	}
	return c.configureEmbeddedBrowsers()
}

func (c *DockerConfigurator) fetchImageTags(image string) []string {
	c.Pointf(`Fetching tags for image %v`, color.BlueString(image))
	reg := c.getRegistryClient()
	if reg == nil {
		c.Errorf(`Docker registry client not initialized`)
		return nil
	}
	tags, err := reg.Tags(image)
	if err != nil {
		c.Errorf(`Failed to fetch tags for image "%s": %v`, image, err)
		return nil
	}
	tagsWithoutLatest := filterOutLatest(tags)
	strSlice := Natural(tagsWithoutLatest)
	sort.Sort(sort.Reverse(strSlice))
	return tagsWithoutLatest
}

func filterOutLatest(tags []string) []string {
	var ret []string
	for _, tag := range tags {
		if !strings.HasPrefix(tag, Latest) {
			ret = append(ret, tag)
		}
	}
	return ret
}

func imageWithTag(image string, tag string) string {
	return fmt.Sprintf("%s:%s", image, tag)
}

func (c *DockerConfigurator) pullVideoRecorderImage() {
	c.Titlef("Pulling video recorder image...")
	c.pullImage(context.Background(), c.getFullyQualifiedImageRef(videoRecorderImage))
}

func (c *DockerConfigurator) getFullyQualifiedImageRef(ref string) string {
	if c.registryHost != "" {
		return fmt.Sprintf("%s/%s", c.registryHost, ref)
	}
	return ref
}

// JSONMessage defines a message struct from docker.
type JSONMessage struct {
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ID              string        `json:"id,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` //deprecated
}

// JSONProgress describes a Progress. terminalFd is the fd of the current terminal,
// Start is the initial value for the operation. Current is the current status and
// value of the progress made towards Total. Total is the end value describing when
// we made 100% progress for an operation.
type JSONProgress struct {
	terminalFd uintptr
	Current    int64 `json:"current,omitempty"`
	Total      int64 `json:"total,omitempty"`
	Start      int64 `json:"start,omitempty"`
	// If true, don't show xB/yB
	HideCounts bool   `json:"hidecounts,omitempty"`
	Units      string `json:"units,omitempty"`
}

func (c *DockerConfigurator) pullImage(ctx context.Context, ref string) bool {
	c.Pointf("Pulling image %v", color.BlueString(ref))
	pullOptions := client.ImagePullOptions{}
	if c.authConfig != nil {
		buf, err := json.Marshal(c.authConfig)
		if err != nil {
			c.Errorf("Failed to prepare registry authentication config: %v", err)
		} else {
			pullOptions.RegistryAuth = base64.URLEncoding.EncodeToString(buf)
		}
	}
	resp, err := c.docker.ImagePull(ctx, ref, pullOptions)
	if err != nil {
		c.Errorf(`Failed to pull image "%s": %v`, ref, err)
		return false
	}
	defer resp.Close()

	var row JSONMessage

	scanner := bufio.NewScanner(resp)
	writer := rewriter.New(colorable.NewColorableStdout())

	for _ = ""; scanner.Scan(); {
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			return false
		}

		select {
		case <-ctx.Done():
			{
				c.Errorf(`Pulling "%s" interrupted: %v`, ref, ctx.Err())
				return false
			}
		default:
			{
				if row.Progress != nil {
					if row.Progress.Current != row.Progress.Total {
						_, _ = fmt.Fprintf(writer, "\t[%s]: %s %s\n", row.ID, row.Status, row.ProgressMessage)
					} else {
						_, _ = fmt.Fprint(writer, "\r")
					}
				}

				_ = writer.Flush()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		c.Errorf(`Failed to pull image "%s": %v`, ref, color.RedString("%v", err))
	}
	return true
}

func (c *DockerConfigurator) IsRunning() bool {
	return c.getSelenoidContainer() != nil
}

func (c *DockerConfigurator) getSelenoidContainer() *ctr.Summary {
	return c.getContainer(selenoidContainerName)
}

func (c *DockerConfigurator) IsUIRunning() bool {
	return c.getSelenoidUIContainer() != nil
}

func (c *DockerConfigurator) getSelenoidUIContainer() *ctr.Summary {
	return c.getContainer(selenoidUIContainerName)
}

func (c *DockerConfigurator) getContainer(name string) *ctr.Summary {
	f := client.Filters{}.Add("name", fmt.Sprintf("^/%s$", name))
	result, err := c.docker.ContainerList(context.Background(), client.ContainerListOptions{Filters: f})
	if err != nil {
		return nil
	}
	if len(result.Items) > 0 {
		return &result.Items[0]
	}
	return nil
}

func (c *DockerConfigurator) PrintArgs() error {
	img := c.getSelenoidImage()
	if img == nil {
		return errors.New("Selenoid image is not downloaded: this is probably a bug")
	}
	cfg := &containerConfig{
		Image:     img,
		Cmd:       []string{"--help"},
		PrintLogs: true,
	}
	return c.startContainer(cfg)
}

const (
	videoDirName = "video"
	logsDirName  = "logs"
	networkName  = "selenoid"
)

func (c *DockerConfigurator) Start() error {
	img := c.getSelenoidImage()
	if img == nil {
		return errors.New("selenoid image is not downloaded: this is probably a bug")
	}

	volumeConfigDir := getVolumeConfigDir(c.ConfigDir, selenoidConfigDirElem)
	videoConfigDir := getVolumeConfigDir(filepath.Join(c.ConfigDir, videoDirName), append(selenoidConfigDirElem, videoDirName))
	logsConfigDir := getVolumeConfigDir(filepath.Join(c.ConfigDir, logsDirName), append(selenoidConfigDirElem, logsDirName))
	volumes := []string{
		fmt.Sprintf("%s:/etc/selenoid:ro,Z", volumeConfigDir),
		fmt.Sprintf("%s:/opt/selenoid/video:Z", videoConfigDir),
		fmt.Sprintf("%s:/opt/selenoid/logs:Z", logsConfigDir),
	}
	binaryPath := c.resolvedSelenoidBinaryPath()
	if !fileExists(binaryPath) {
		return fmt.Errorf("selenoid binary not found at %s: run \"cm selenoid download\" first", binaryPath)
	}
	volumes = append(volumes, fmt.Sprintf("%s:%s:ro", binaryPath, selenoidBinaryMountPath))
	const dockerSocket = "/var/run/docker.sock"
	if isWindows() {
		//With two slashes. See https://stackoverflow.com/questions/36765138/bind-to-docker-socket-on-windows
		volumes = append(volumes, fmt.Sprintf("/%s:%s", dockerSocket, dockerSocket))
	} else if fileExists(dockerSocket) {
		volumes = append(volumes, fmt.Sprintf("%s:%s:Z", dockerSocket, dockerSocket))
	}

	cmd := []string{}
	overrideCmd := strings.Fields(c.Args)
	if len(overrideCmd) > 0 {
		cmd = overrideCmd
	}
	if !contains(cmd, "-conf") {
		cmd = append(cmd, "-conf", "/etc/selenoid/browsers.json")
	}
	if !contains(cmd, "-video-output-dir") && isVideoRecordingSupported(c.Logger, c.Version) {
		cmd = append(cmd, "-video-output-dir", "/opt/selenoid/video/")
	}
	if !contains(cmd, "-video-recorder-image") && isVideoRecordingSupported(c.Logger, c.Version) {
		cmd = append(cmd, "-video-recorder-image", c.getFullyQualifiedImageRef(videoRecorderImage))
	}
	if !c.DisableLogs && !contains(cmd, "-log-output-dir") && isLogSavingSupported(c.Logger, c.Version) {
		cmd = append(cmd, "-log-output-dir", "/opt/selenoid/logs/")
	}
	if !contains(cmd, "-container-network") {
		cmd = append(cmd, "-container-network", networkName)
	}
	// qaguru/selenoid image sets ENTRYPOINT to /usr/bin/selenoid; prepending the
	// binary again makes Go flag.Parse stop at the duplicate argv and ignore flags.

	overrideEnv := strings.Fields(c.Env)
	if !strings.Contains(c.Env, "OVERRIDE_VIDEO_OUTPUT_DIR") {
		overrideEnv = append(overrideEnv, fmt.Sprintf("OVERRIDE_VIDEO_OUTPUT_DIR=%s", videoConfigDir))
	}
	cfg := &containerConfig{
		Name:        selenoidContainerName,
		Image:       img,
		HostPort:    c.Port,
		ServicePort: DefaultPort,
		Volumes:     volumes,
		Network:     networkName,
		Cmd:         cmd,
		OverrideEnv: overrideEnv,
		UserNS:      c.UserNS,
	}
	return c.startContainer(cfg)
}

func isVideoRecordingSupported(logger Logger, version string) bool {
	return isVersion(version, ">= 1.4.0", func(version string) {
		logger.Pointf(`Not enabling video feature because specified version "%s" is not semantic`, version)
	})
}

func isVersion(version string, condition string, notSemanticVersionCallback func(string)) bool {
	if version == Latest {
		return true
	}
	constraint, _ := semver.NewConstraint(condition)
	v, err := semver.NewVersion(version)
	if err != nil {
		notSemanticVersionCallback(version)
		return false
	}
	return constraint.Check(v)
}

func isLogSavingSupported(logger Logger, version string) bool {
	return isVersion(version, ">= 1.7.0", func(version string) {
		logger.Pointf(`Not enabling log saving feature because specified version "%s" is not semantic`, version)
	})
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func getVolumeConfigDir(defaultConfigDir string, elem []string) string {
	configDir := chooseVolumeConfigDir(defaultConfigDir, elem)
	if isWindows() { //A bit ugly, but conditional compilation is even worse
		return postProcessPath(configDir)
	}
	return configDir
}

// According to https://stackoverflow.com/questions/34161352/docker-sharing-a-volume-on-windows-with-docker-toolbox
func postProcessPath(path string) string {
	if len(path) >= 2 {
		replacedSlashes := strings.Replace(path, string("\\"), "/", -1)
		re := regexp.MustCompile("([A-Z]):(.+)")
		lowerCaseDriveLetter := strings.ToLower(re.ReplaceAllString(replacedSlashes, "$1"))
		pathTail := re.ReplaceAllString(replacedSlashes, "$2")
		return "/" + lowerCaseDriveLetter + pathTail
	}
	return path
}

func chooseVolumeConfigDir(defaultConfigDir string, elem []string) string {
	overrideHome := os.Getenv(overrideHome)
	if overrideHome != "" {
		return joinPaths(overrideHome, elem)
	}
	return defaultConfigDir
}

func (c *DockerConfigurator) PrintUIArgs() error {
	img := c.getSelenoidUIImage()
	if img == nil {
		return errors.New("selenoid UI image is not downloaded: this is probably a bug")
	}
	cfg := &containerConfig{
		Image:     img,
		Cmd:       []string{"--help"},
		PrintLogs: true,
	}
	return c.startContainer(cfg)
}

func (c *DockerConfigurator) StartUI() error {
	img := c.getSelenoidUIImage()
	if img == nil {
		return errors.New("selenoid ui image is not downloaded: this is probably a bug")
	}

	var cmd, candidates []string
	var selenoidUri string
containers:
	for _, containerName := range []string{
		selenoidContainerName, ggrUIContainerName,
	} {
		if ctr := c.getContainer(containerName); ctr != nil {
			port := uint16(DefaultPort)
			for _, p := range ctr.Ports {
				if p.PrivatePort != 0 {
					port = p.PrivatePort
					break
				}
				if p.PublicPort != 0 {
					port = p.PublicPort
					break
				}
			}
			selenoidUri = fmt.Sprintf("--selenoid-uri=http://%s:%d", containerName, port)
			candidates = []string{containerName}
			break containers
		}
	}
	if selenoidUri == "" && c.getContainer(selenoidContainerName) != nil {
		selenoidUri = fmt.Sprintf("--selenoid-uri=http://%s:%d", selenoidContainerName, DefaultPort)
	}
	overrideCmd := strings.Fields(c.Args)
	if len(overrideCmd) > 0 {
		cmd = overrideCmd
	}
	if !contains(cmd, "--selenoid-uri") {
		cmd = append(cmd, selenoidUri)
	}

	if len(candidates) == 0 {
		c.Errorf("Neither Selenoid nor Ggr UI is started. Selenoid UI may not work.")
	}

	binaryPath := c.resolvedSelenoidUIBinaryPath()
	if !fileExists(binaryPath) {
		return fmt.Errorf("selenoid ui binary not found at %s: run \"cm selenoid-ui download\" first", binaryPath)
	}
	volumeConfigDir := getVolumeConfigDir(c.ConfigDir, selenoidConfigDirElem)
	volumes := []string{
		fmt.Sprintf("%s:/etc/selenoid:ro", volumeConfigDir),
		fmt.Sprintf("%s:%s:ro", binaryPath, selenoidUIBinaryMountPath),
	}
	if !contains(cmd, "-browsers-conf") && !contains(cmd, "--browsers-conf") {
		cmd = append(cmd, "-browsers-conf", "/etc/selenoid/browsers.json")
	}
	if !contains(cmd, "-listen") && !contains(cmd, "--listen") {
		cmd = append(cmd, fmt.Sprintf("-listen=:%d", c.Port))
	}
	// qaguru/selenoid-ui image sets ENTRYPOINT to /selenoid-ui; prepending the
	// binary again makes Go flag.Parse stop at the duplicate argv and ignore flags.

	overrideEnv := strings.Fields(c.Env)
	cfg := &containerConfig{
		Name:        selenoidUIContainerName,
		Image:       img,
		HostPort:    c.Port,
		ServicePort: UIDefaultPort,
		Volumes:     volumes,
		Network:     networkName,
		Cmd:         cmd,
		OverrideEnv: overrideEnv,
		UserNS:      c.UserNS,
	}
	return c.startContainer(cfg)
}

func validateEnviron(envs []string) []string {
	validEnv := []string{}
	for _, e := range envs {
		k := strings.Split(e, "=")
		if len(k[0]) != 0 {
			validEnv = append(validEnv, e)
		}
	}
	return validEnv
}

type containerConfig struct {
	Name        string
	Image       *img.Summary
	HostPort    int
	ServicePort int
	Volumes     []string
	Network     string
	Cmd         []string
	OverrideEnv []string
	UserNS      string
	PrintLogs   bool
}

func (c *DockerConfigurator) startContainer(cfg *containerConfig) error {
	ctx := context.Background()
	env := validateEnviron(os.Environ())
	env = append(env, fmt.Sprintf("TZ=%s", time.Local))
	if len(cfg.OverrideEnv) > 0 {
		env = cfg.OverrideEnv
	}
	if !contains(env, dockerApiVersion) {
		env = append(env, fmt.Sprintf("%s=%s", dockerApiVersion, selenoidDockerAPI))
	}
	servicePortString := strconv.Itoa(cfg.ServicePort)
	port, err := nat.NewPort("tcp", servicePortString)
	if err != nil {
		return fmt.Errorf("failed to init port: %v", err)
	}

	err = c.createNetworkIfNeeded(cfg.Network)
	if err != nil {
		return fmt.Errorf("failed to configure container network: %v", err)
	}
	containerConfig := ctr.Config{
		Hostname: "localhost",
		Image:    cfg.Image.RepoTags[0],
		Env:      env,
	}
	if cfg.ServicePort > 0 {
		exposedPorts, portErr := networkPortSet(map[nat.Port]struct{}{port: {}})
		if portErr != nil {
			return fmt.Errorf("failed to convert exposed ports: %v", portErr)
		}
		containerConfig.ExposedPorts = exposedPorts
	}
	if len(cfg.Cmd) > 0 {
		containerConfig.Cmd = cfg.Cmd
	}
	hostConfig := ctr.HostConfig{
		Binds:       cfg.Volumes,
		NetworkMode: ctr.NetworkMode(networkName),
	}
	if cfg.UserNS != "" {
		mode := ctr.UsernsMode(cfg.UserNS)
		if !mode.Valid() {
			return fmt.Errorf("invalid userns value: %s", cfg.UserNS)
		}
		hostConfig.UsernsMode = mode
	}
	if cfg.PrintLogs {
		containerConfig.Tty = true
	} else {
		hostConfig.RestartPolicy = ctr.RestartPolicy{
			Name: "always",
		}
	}
	if cfg.HostPort > 0 && cfg.ServicePort > 0 {
		hostPortString := strconv.Itoa(cfg.HostPort)
		portBindings, portErr := networkPortMap(nat.PortMap{
			port: {{HostIP: "0.0.0.0", HostPort: hostPortString}},
		})
		if portErr != nil {
			return fmt.Errorf("failed to convert port bindings: %v", portErr)
		}
		hostConfig.PortBindings = portBindings
	}
	created, err := c.docker.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:             cfg.Name,
		Config:           &containerConfig,
		HostConfig:       &hostConfig,
		NetworkingConfig: &network.NetworkingConfig{},
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	_, err = c.docker.ContainerStart(ctx, created.ID, client.ContainerStartOptions{})
	if err != nil {
		_ = c.removeContainer(created.ID)
		return fmt.Errorf("failed to start container: %v", err)
	}
	if cfg.PrintLogs {
		defer c.removeContainer(created.ID)
		r, err := c.docker.ContainerLogs(ctx, created.ID, client.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			return fmt.Errorf("failed to read container logs: %v", err)
		}
		defer r.Close()
		_, _ = io.Copy(os.Stderr, r)
	}
	return nil
}

func (c *DockerConfigurator) createNetworkIfNeeded(networkName string) error {
	ctx := context.Background()
	_, err := c.docker.NetworkInspect(ctx, networkName, client.NetworkInspectOptions{})
	if err != nil {
		_, err = c.docker.NetworkCreate(ctx, networkName, client.NetworkCreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create custom network %s: %v", networkName, err)
		}
	}
	return nil
}

func (c *DockerConfigurator) removeContainer(id string) error {
	ctx := context.Background()
	if c.Graceful {
		timeout := int(c.GracefulTimeout.Seconds())
		_, err := c.docker.ContainerStop(ctx, id, client.ContainerStopOptions{Timeout: &timeout})
		if err == nil {
			_, err = c.docker.ContainerRemove(ctx, id, client.ContainerRemoveOptions{RemoveVolumes: true})
			return err
		}
		return err
	}
	_, err := c.docker.ContainerRemove(ctx, id, client.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	return err
}

func (c *DockerConfigurator) Stop() error {
	sc := c.getSelenoidContainer()
	if sc != nil {
		err := c.removeContainer(sc.ID)
		if err != nil {
			return fmt.Errorf("failed to stop Selenoid container: %v", err)
		}
	}
	return nil
}

func (c *DockerConfigurator) StopUI() error {
	uc := c.getSelenoidUIContainer()
	if uc != nil {
		err := c.removeContainer(uc.ID)
		if err != nil {
			return fmt.Errorf("failed to stop Selenoid UI container: %v", err)
		}
	}
	return nil
}
