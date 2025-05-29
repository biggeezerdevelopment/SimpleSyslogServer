package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port          int    `yaml:"port"`
		LogFile       string `yaml:"log_file"`
		ConsoleOutput bool   `yaml:"console_output"`
	} `yaml:"server"`
	Filter struct {
		Enabled         bool     `yaml:"enabled"`
		AllowedIPs      []string `yaml:"allowed_ips"`
		MinSeverity     int      `yaml:"min_severity"`
		ExcludePatterns []string `yaml:"exclude_patterns"`
	} `yaml:"filter"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Set default values if not specified
	if config.Server.Port == 0 {
		config.Server.Port = 514
	}
	if !config.Server.ConsoleOutput && config.Server.LogFile == "" {
		// Ensure at least one output is enabled
		config.Server.ConsoleOutput = true
	}
	if !config.Filter.Enabled {
		// If filtering is disabled, set most permissive settings
		config.Filter.MinSeverity = 7
		config.Filter.AllowedIPs = []string{"0.0.0.0/0"}
		config.Filter.ExcludePatterns = nil
	} else if config.Filter.MinSeverity < 0 || config.Filter.MinSeverity > 7 {
		config.Filter.MinSeverity = 7
	}

	return &config, nil
}
