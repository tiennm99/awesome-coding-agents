package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Agent struct {
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
	Category string `yaml:"category,omitempty"`
	Notes    string `yaml:"notes,omitempty"`
}

func loadAgents(path string) ([]Agent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg struct {
		Agents []Agent `yaml:"agents"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	for i, a := range cfg.Agents {
		if a.Owner == "" || a.Repo == "" {
			return nil, fmt.Errorf("entry %d missing owner or repo", i)
		}
	}
	return cfg.Agents, nil
}
