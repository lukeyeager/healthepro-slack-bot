package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds static configuration loaded from config.yaml.
type Config struct {
	OrgID       int    `yaml:"org_id"`
	MenuID      int    `yaml:"menu_id"`
	EveningCron string `yaml:"evening_cron"`
	MorningCron string `yaml:"morning_cron"`
	DBPath      string `yaml:"db_path"`
	Timezone    string `yaml:"timezone"`
}

// Load reads and parses a YAML config file, applying defaults for optional fields.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := &Config{
		EveningCron: "0 19 * * 0-4",
		MorningCron: "0 6 * * 1-5",
		DBPath:      "/data/menus.db",
		Timezone:    "America/Chicago",
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.OrgID == 0 {
		return nil, fmt.Errorf("org_id is required")
	}
	if cfg.MenuID == 0 {
		return nil, fmt.Errorf("menu_id is required")
	}

	return cfg, nil
}
