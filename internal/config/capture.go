package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Capture uses only the real machine config. PARLEY_HOME and deck/env layers
// cannot enable this private diagnostic facility.
func MachineZcodeCapture() (enabled bool, directory string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, "", err
	}
	directory = filepath.Join(home, ".parley", "captures", "zcode")
	raw, err := os.ReadFile(filepath.Join(home, ".parley", "agents.toml"))
	if os.IsNotExist(err) {
		return false, directory, nil
	}
	if err != nil {
		return false, directory, err
	}
	var cfg struct {
		Defaults struct {
			Capture bool `toml:"capture_zcode_failures"`
		} `toml:"defaults"`
	}
	if err = toml.Unmarshal(raw, &cfg); err != nil {
		return false, directory, err
	}
	return cfg.Defaults.Capture, directory, nil
}
