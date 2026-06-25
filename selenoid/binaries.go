package selenoid

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

func (c *DockerConfigurator) resolvedSelenoidBinaryPath() string {
	if c.SelenoidBinaryPath != "" {
		return c.SelenoidBinaryPath
	}
	return filepath.Join(c.ConfigDir, binDirName, "selenoid")
}

func (c *DockerConfigurator) resolvedSelenoidUIBinaryPath() string {
	if c.SelenoidUIBinaryPath != "" {
		return c.SelenoidUIBinaryPath
	}
	return filepath.Join(c.ConfigDir, binDirName, "selenoid-ui")
}

func (c *DockerConfigurator) downloadSelenoidBinary() error {
	path := c.resolvedSelenoidBinaryPath()
	if fileExists(path) && !c.Force {
		c.Pointf("Selenoid binary already present at %s", path)
		return nil
	}
	if c.SelenoidBinaryPath != "" {
		if !fileExists(path) {
			return fmt.Errorf("selenoid binary not found at %s", path)
		}
		return nil
	}
	return c.downloadBinary(selenoidRepo, path)
}

func (c *DockerConfigurator) downloadSelenoidUIBinary() error {
	path := c.resolvedSelenoidUIBinaryPath()
	if fileExists(path) && !c.Force {
		c.Pointf("Selenoid UI binary already present at %s", path)
		return nil
	}
	if c.SelenoidUIBinaryPath != "" {
		if !fileExists(path) {
			return fmt.Errorf("selenoid ui binary not found at %s", path)
		}
		return nil
	}
	return c.downloadBinary(selenoidUIRepo, path)
}

func (c *DockerConfigurator) downloadBinary(repo, outputPath string) error {
	osName := c.OS
	arch := c.Arch
	if osName == "" {
		osName = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	u, err := getGithubReleaseAssetURL(githubOwner, repo, c.Version, osName, arch, c.GithubBaseUrl)
	if err != nil {
		return fmt.Errorf("failed to get %s/%s download URL: %v", githubOwner, repo, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %v", err)
	}
	c.Titlef("Downloading %s from %s", repo, color.BlueString(u))
	f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := downloadFileWithProgressBar(u, f); err != nil {
		return fmt.Errorf("failed to download %s: %v", repo, err)
	}
	c.Titlef("Successfully downloaded %s to %s", repo, color.GreenString(outputPath))
	return nil
}
