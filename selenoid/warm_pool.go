package selenoid

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	WarmPoolURL            = "http://127.0.0.1:9090"
	warmPoolHealthURL      = WarmPoolURL + "/health"
	warmPoolComposeProject = "selenoid-pool"
	warmPoolDirName        = "warm-pool"
	warmPoolManagedMarker  = ".cm-managed"
	warmPoolComposeFile    = "docker-compose.yml"
	warmPoolConfigFile     = "config.yaml"
)

//go:embed data/warm-pool/docker-compose.yml data/warm-pool/config.yaml
var warmPoolFS embed.FS

type PoolAware struct {
	WarmPool bool
	HotPool  bool
}

func (p PoolAware) poolEnabled() bool {
	return p.WarmPool || p.HotPool
}

func warmPoolDir(configDir string) string {
	return filepath.Join(configDir, warmPoolDirName)
}

func warmPoolMarkerPath(configDir string) string {
	return filepath.Join(warmPoolDir(configDir), warmPoolManagedMarker)
}

func warmPoolManaged(configDir string) bool {
	_, err := os.Stat(warmPoolMarkerPath(configDir))
	return err == nil
}

func warmPoolComposeUpArgs(hot bool) []string {
	args := []string{"-p", warmPoolComposeProject, "-f", warmPoolComposeFile}
	if hot {
		args = append(args, "--profile", "hot")
	}
	return append(args, "up", "-d")
}

func warmPoolComposeDownArgs() []string {
	return []string{"-p", warmPoolComposeProject, "-f", warmPoolComposeFile, "--profile", "hot", "down", "--remove-orphans"}
}

func appendWarmPoolHubArgs(cmd []string) []string {
	if contains(cmd, "-warm-pool-url") || contains(cmd, "-pool-url") {
		return cmd
	}
	return append(cmd, "-warm-pool-url", WarmPoolURL)
}

func writeWarmPoolFiles(configDir string) error {
	dir := warmPoolDir(configDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("warm-pool dir: %v", err)
	}
	for _, name := range []string{warmPoolComposeFile, warmPoolConfigFile} {
		data, err := fs.ReadFile(warmPoolFS, filepath.ToSlash(filepath.Join("data/warm-pool", name)))
		if err != nil {
			return fmt.Errorf("embed %s: %v", name, err)
		}
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("write %s: %v", path, err)
		}
	}
	if err := os.WriteFile(warmPoolMarkerPath(configDir), []byte("cm\n"), 0644); err != nil {
		return fmt.Errorf("warm-pool marker: %v", err)
	}
	return nil
}

var composeExec = defaultComposeExec
var waitHealthy = defaultWaitHealthy
var healthWaitTimeout = 2 * time.Minute
var healthPollInterval = 2 * time.Second

func defaultComposeExec(dir string, args []string, quiet bool) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Dir = dir
	if quiet {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose %v: %v", args, err)
	}
	return nil
}

func defaultWaitHealthy(rawURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last error
	for time.Now().Before(deadline) {
		resp, err := http.Get(rawURL)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			last = fmt.Errorf("HTTP %d", resp.StatusCode)
		} else {
			last = err
		}
		time.Sleep(healthPollInterval)
	}
	if last == nil {
		last = fmt.Errorf("timeout")
	}
	return fmt.Errorf("warm-pool health %s: %v", rawURL, last)
}

func startWarmPoolSidecar(log Logger, configDir string, hot, quiet bool) error {
	log.Titlef("Starting warm-pool sidecar...")
	if err := writeWarmPoolFiles(configDir); err != nil {
		return err
	}
	dir := warmPoolDir(configDir)
	if err := composeExec(dir, warmPoolComposeUpArgs(hot), quiet); err != nil {
		return err
	}
	if err := waitHealthy(warmPoolHealthURL, healthWaitTimeout); err != nil {
		return err
	}
	log.Pointf("Warm-pool orchestrator %s", WarmPoolURL)
	return nil
}

func stopWarmPoolSidecar(log Logger, configDir string) error {
	if !warmPoolManaged(configDir) {
		return nil
	}
	dir := warmPoolDir(configDir)
	if !fileExists(filepath.Join(dir, warmPoolComposeFile)) {
		return nil
	}
	log.Titlef("Stopping warm-pool sidecar...")
	if err := composeExec(dir, warmPoolComposeDownArgs(), log.Quiet); err != nil {
		return err
	}
	_ = os.Remove(warmPoolMarkerPath(configDir))
	return nil
}
