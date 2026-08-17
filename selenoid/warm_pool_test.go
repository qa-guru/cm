package selenoid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestPoolEnabled(t *testing.T) {
	assert.False(t, PoolAware{}.poolEnabled())
	assert.True(t, PoolAware{WarmPool: true}.poolEnabled())
	assert.True(t, PoolAware{HotPool: true}.poolEnabled())
	assert.True(t, PoolAware{WarmPool: true, HotPool: true}.poolEnabled())
}

func TestWarmPoolComposeArgs(t *testing.T) {
	warm := warmPoolComposeUpArgs(false)
	assert.Equal(t, []string{"-p", "selenoid-warm", "-f", "docker-compose.yml", "up", "-d"}, warm)
	assert.NotContains(t, warm, "--profile")

	hot := warmPoolComposeUpArgs(true)
	assert.Equal(t, []string{"-p", "selenoid-warm", "-f", "docker-compose.yml", "--profile", "hot", "up", "-d"}, hot)

	down := warmPoolComposeDownArgs()
	assert.Contains(t, down, "down")
	assert.Contains(t, down, "--profile")
	assert.Contains(t, down, "hot")
}

func TestAppendWarmPoolHubArgs(t *testing.T) {
	cmd := appendWarmPoolHubArgs([]string{"-conf", "/etc/selenoid/browsers.json"})
	assert.Equal(t, []string{"-conf", "/etc/selenoid/browsers.json", "-warm-pool-url", WarmPoolURL}, cmd)
	assert.Equal(t, "http://127.0.0.1:9090", WarmPoolURL)

	already := []string{"-warm-pool-url", "http://example.invalid:1"}
	assert.Equal(t, already, appendWarmPoolHubArgs(already))

	aliased := []string{"-pool-url", "http://example.invalid:1"}
	assert.Equal(t, aliased, appendWarmPoolHubArgs(aliased))
}

func TestWarmPoolEmbedHasNoBuildAndHotProfile(t *testing.T) {
	compose, err := warmPoolFS.ReadFile("data/warm-pool/docker-compose.yml")
	assert.NoError(t, err)
	text := string(compose)
	assert.NotRegexp(t, `(?m)^\s+build:`, text)
	assert.Contains(t, text, "qaguru/selenoid-pool:min")
	assert.Contains(t, text, `profiles: ["hot"]`)
	assert.Contains(t, text, "127.0.0.1:9090:9090")
	assert.Contains(t, text, "warm-chrome-1")
	assert.Contains(t, text, "hot-chrome-min-1")
	assert.Contains(t, text, "hot-chrome-1")
	assert.Contains(t, text, "hot-pw-1")

	cfg, err := warmPoolFS.ReadFile("data/warm-pool/config.yaml")
	assert.NoError(t, err)
	assert.Contains(t, string(cfg), "pool-chrome-1")
	assert.Contains(t, string(cfg), "pool-hot-chrome-min-1")
	assert.Contains(t, string(cfg), "pool-hot-chrome-1")
	assert.Contains(t, string(cfg), "pool-hot-pw-1")
	assert.Contains(t, string(cfg), "pool: hot")
}

func TestWriteWarmPoolFiles(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, writeWarmPoolFiles(dir))
	assert.True(t, warmPoolManaged(dir))
	assert.True(t, fileExists(filepath.Join(warmPoolDir(dir), warmPoolComposeFile)))
	assert.True(t, fileExists(filepath.Join(warmPoolDir(dir), warmPoolConfigFile)))
	body, err := os.ReadFile(filepath.Join(warmPoolDir(dir), warmPoolComposeFile))
	assert.NoError(t, err)
	assert.NotRegexp(t, `(?m)^\s+build:`, string(body))
}

func TestStartWarmPoolSidecarRunsCompose(t *testing.T) {
	origExec := composeExec
	origWait := waitHealthy
	t.Cleanup(func() {
		composeExec = origExec
		waitHealthy = origWait
	})

	var gotDir string
	var gotArgs []string
	composeExec = func(dir string, args []string, quiet bool) error {
		gotDir = dir
		gotArgs = append([]string{}, args...)
		return nil
	}
	waitHealthy = func(string, time.Duration) error { return nil }

	dir := t.TempDir()
	log := Logger{Quiet: true}
	assert.NoError(t, startWarmPoolSidecar(log, dir, false, true))
	assert.Equal(t, warmPoolDir(dir), gotDir)
	assert.Equal(t, warmPoolComposeUpArgs(false), gotArgs)
	assert.True(t, warmPoolManaged(dir))

	assert.NoError(t, startWarmPoolSidecar(log, dir, true, true))
	assert.Equal(t, warmPoolComposeUpArgs(true), gotArgs)
}

func TestStopWarmPoolSidecarSkipsWhenNotManaged(t *testing.T) {
	origExec := composeExec
	t.Cleanup(func() { composeExec = origExec })
	called := false
	composeExec = func(string, []string, bool) error {
		called = true
		return nil
	}
	assert.NoError(t, stopWarmPoolSidecar(Logger{Quiet: true}, t.TempDir()))
	assert.False(t, called)
}

func TestStopWarmPoolSidecarRunsDown(t *testing.T) {
	origExec := composeExec
	t.Cleanup(func() { composeExec = origExec })
	var gotArgs []string
	composeExec = func(dir string, args []string, quiet bool) error {
		gotArgs = append([]string{}, args...)
		return nil
	}
	dir := t.TempDir()
	assert.NoError(t, writeWarmPoolFiles(dir))
	assert.NoError(t, stopWarmPoolSidecar(Logger{Quiet: true}, dir))
	assert.Equal(t, warmPoolComposeDownArgs(), gotArgs)
	assert.False(t, warmPoolManaged(dir))
}

func TestWaitHealthyOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	t.Cleanup(srv.Close)
	origPoll := healthPollInterval
	healthPollInterval = time.Millisecond
	t.Cleanup(func() { healthPollInterval = origPoll })
	assert.NoError(t, defaultWaitHealthy(srv.URL+"/health", time.Second))
}

func TestWaitHealthyRejectsNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)
	origPoll := healthPollInterval
	healthPollInterval = time.Millisecond
	t.Cleanup(func() { healthPollInterval = origPoll })
	err := defaultWaitHealthy(srv.URL+"/health", 20*time.Millisecond)
	assert.Error(t, err)
}

func TestHubUsesHostNetworkOnlyWithPool(t *testing.T) {
	cold := &DockerConfigurator{}
	assert.False(t, cold.hubUsesHostNetwork())
	warm := &DockerConfigurator{PoolAware: PoolAware{WarmPool: true}}
	assert.True(t, warm.hubUsesHostNetwork())
	hot := &DockerConfigurator{PoolAware: PoolAware{HotPool: true}}
	assert.True(t, hot.hubUsesHostNetwork())
}

func TestDockerConfiguratorCopiesPoolFlags(t *testing.T) {
	c := &DockerConfigurator{PoolAware: PoolAware{WarmPool: true, HotPool: true}}
	assert.True(t, c.poolEnabled())
	assert.True(t, c.hubUsesHostNetwork())
	cmd := appendWarmPoolHubArgs([]string{"-conf", "/etc/selenoid/browsers.json", "-container-network", networkName})
	assert.Contains(t, cmd, "-warm-pool-url")
	assert.Contains(t, cmd, WarmPoolURL)
}
