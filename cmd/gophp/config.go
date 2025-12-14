package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the .gophp.yaml configuration
type Config struct {
	Service string `yaml:"service"`
	Source  string `yaml:"source"`
	Output  struct {
		Dir    string `yaml:"dir"`
		LibDir string `yaml:"lib_dir"`
	} `yaml:"output"`
}

// loadConfig loads the .gophp.yaml configuration file
func loadConfig() (*Config, error) {
	data, err := os.ReadFile(".gophp.yaml")
	if err != nil {
		return nil, fmt.Errorf("读取 .gophp.yaml 失败：%w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 .gophp.yaml 失败：%w", err)
	}

	// Set defaults
	if config.Output.Dir == "" {
		config.Output.Dir = "dist"
	}
	if config.Output.LibDir == "" {
		config.Output.LibDir = "dist/lib"
	}

	return &config, nil
}
