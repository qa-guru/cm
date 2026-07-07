package selenoid

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type dockerAuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`
}

type dockerConfigFile struct {
	Auths map[string]dockerAuthConfig `json:"auths"`
}

func loadDockerAuthConfigs() (map[string]dockerAuthConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory: %w", err)
	}
	path := filepath.Join(home, ".docker", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]dockerAuthConfig{}, nil
		}
		return nil, fmt.Errorf("read docker config: %w", err)
	}
	var cfg dockerConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse docker config: %w", err)
	}
	if cfg.Auths == nil {
		return map[string]dockerAuthConfig{}, nil
	}
	return cfg.Auths, nil
}
