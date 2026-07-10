package selenoid

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestEmbeddedBrowsersJSON(t *testing.T) {
	assert.NotEmpty(t, embeddedBrowsersJSON, "go:embed data/browsers.json must be non-empty")

	var cfg SelenoidConfig
	assert.NoError(t, json.Unmarshal(embeddedBrowsersJSON, &cfg))

	for _, name := range []string{"chrome", "firefox", "msedge", "playwright-chromium"} {
		browser, ok := cfg[name]
		assert.Truef(t, ok, "embedded catalog missing %q", name)
		assert.NotEmptyf(t, browser.Default, "%s.default", name)
		assert.NotEmptyf(t, browser.Versions, "%s.versions", name)
	}
	assert.Equal(t, "149.0", cfg["chrome"].Default)
	assert.Equal(t, "1.61.1", cfg["playwright-chromium"].Default)
}

func TestConfigureEmbeddedBrowsersWritesConfig(t *testing.T) {
	dir := t.TempDir()
	c := &DockerConfigurator{
		Logger:         Logger{Quiet: true},
		ConfigDirAware: ConfigDirAware{ConfigDir: dir},
		DownloadAware:  DownloadAware{DownloadNeeded: false},
	}
	cfg, err := c.configureEmbeddedBrowsers()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "149.0", (*cfg)["chrome"].Default)

	written, err := os.ReadFile(filepath.Join(dir, "browsers.json"))
	assert.NoError(t, err)
	assert.Equal(t, embeddedBrowsersJSON, written)
}
