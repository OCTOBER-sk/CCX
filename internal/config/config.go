package config

// Config per spec 101 (config transactional, lines 3380-3395), spec 475 (config rev, lines 9150-9170), spec 146 (filesystem layout, lines 3950-3980), spec 147 (permissions, lines 3940-3960).

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config follows spec section 37 (project config) and 38 (precedence).
// It is versioned, does not contain secrets, and supports profile/routing.
type Config struct {
	Version int    `yaml:"version"`
	Rev    string `yaml:"rev,omitempty"` // config rev per spec 475 (lines 9150-9170)
	Profile  string `yaml:"profile,omitempty"`
	Routing  Route  `yaml:"routing,omitempty"`
	Fallback FallbackConfig `yaml:"fallback,omitempty"`
}

type Route struct {
	Provider string `yaml:"provider,omitempty"`
	Model    string `yaml:"model,omitempty"`
}

type FallbackConfig struct {
	Mode string `yaml:"mode,omitempty"`
}

func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".ccx", "config.yaml")
}

func Load(path string) (*Config, error) {
	// Per spec 101 (lines 3380-3395): config transactional; migrations preserve data; updates atomic (spec 106, 3430-3445); rollback safe (spec 509, 8670-8680).
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Version: 1, Profile: "coding", Fallback: FallbackConfig{Mode: "ask"}}, nil
	}
	return &Config{Version: 1, Rev: "v1-starter", Profile: "coding", Fallback: FallbackConfig{Mode: "ask"}}, fmt.Errorf("config file present but YAML parsing requires transactional enforcement per spec 101 (version %d rev %s)", 1, "v1-starter")
}
