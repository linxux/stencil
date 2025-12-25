package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FormatOptions controls which variable formats are enabled
type FormatOptions struct {
	// EnableBraces enables {{var}} format
	EnableBraces bool `json:"enableBraces"`
	// EnableAngleBrackets enables <<var>> format
	EnableAngleBrackets bool `json:"enableAngleBrackets"`
	// EnableUnderscores enables __var__ format
	EnableUnderscores bool `json:"enableUnderscores"`
	// EnablePercent enables %var% format
	EnablePercent bool `json:"enablePercent"`
}

// Config represents the generator configuration
type Config struct {
	// TemplateDir is the source template directory
	TemplateDir string `json:"templateDir"`

	// OutputDir is the target output directory
	OutputDir string `json:"outputDir"`

	// Variables contains key-value pairs for replacement
	Variables map[string]string `json:"variables"`

	// Interactive mode enables prompt for values
	Interactive bool `json:"interactive"`

	// DryRun shows what would be generated without creating files
	DryRun bool `json:"dryRun"`

	// SkipConfirm skips confirmation prompt in interactive mode
	SkipConfirm bool `json:"skipConfirm"`

	// Formats controls which variable formats are enabled
	Formats FormatOptions `json:"formats"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(configPath string, cfg *Config) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		TemplateDir: "./template",
		OutputDir:   "./output",
		Variables:   make(map[string]string),
		Interactive: false,
		DryRun:      false,
		SkipConfirm: false,
		Formats: FormatOptions{
			EnableBraces:        true,
			EnableAngleBrackets: true,
			EnableUnderscores:   true,
			EnablePercent:       true,
		},
	}
}
