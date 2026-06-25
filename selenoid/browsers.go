package selenoid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "embed"

	"github.com/fatih/color"
)

//go:embed data/browsers-qaguru.json
var embeddedBrowsersJSON []byte

func (c *DockerConfigurator) configureEmbeddedBrowsers() (*SelenoidConfig, error) {
	c.Titlef("Using qa-guru browsers configuration...")
	if len(embeddedBrowsersJSON) == 0 {
		return nil, errors.New("embedded browsers.json is empty")
	}
	var cfg SelenoidConfig
	if err := json.Unmarshal(embeddedBrowsersJSON, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse embedded browsers.json: %v", err)
	}
	if c.DownloadNeeded {
		for _, versions := range cfg {
			for _, version := range versions.Versions {
				if ref, ok := version.Image.(string); ok {
					if !c.pullImage(context.Background(), ref) {
						return nil, fmt.Errorf("failed to pull image %s", ref)
					}
				}
			}
		}
		c.pullVideoRecorderImage()
	}
	return &cfg, os.WriteFile(getSelenoidConfigPath(c.ConfigDir), embeddedBrowsersJSON, 0644)
}

func (c *DockerConfigurator) syncBrowsersFromFile(path string) (*SelenoidConfig, error) {
	c.Titlef(`Syncing browsers configuration from "%v"...`, color.GreenString(path))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read browsers.json from %s: %v", path, err)
	}
	var cfg SelenoidConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse browsers.json from %s: %v", path, err)
	}
	if c.DownloadNeeded {
		for _, versions := range cfg {
			for _, version := range versions.Versions {
				if ref, ok := version.Image.(string); ok {
					if !c.pullImage(context.Background(), ref) {
						return nil, fmt.Errorf("failed to pull image %s from %s", ref, path)
					}
				}
			}
		}
		c.pullVideoRecorderImage()
	}
	return &cfg, os.WriteFile(getSelenoidConfigPath(c.ConfigDir), data, 0644)
}
