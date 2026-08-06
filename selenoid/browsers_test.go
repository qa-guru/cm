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

	android, ok := cfg["android"]
	assert.True(t, ok, "embedded catalog missing android")
	assert.Equal(t, "16.0", android.Default)
	assert.Contains(t, android.Versions, "5.1")
	assert.NotContains(t, android.Versions, "4.4")
	assert.Equal(t, "qaguru/android:16", android.Versions["16.0"].Image)
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
