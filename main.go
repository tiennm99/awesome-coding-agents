package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("update failed: %v", err)
	}
}

func run() error {
	agents, err := loadAgents("data/agents.yml")
	if err != nil {
		return err
	}
	if len(agents) == 0 {
		return fmt.Errorf("no agents in data/agents.yml")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN env var required")
	}

	stats, err := fetchStats(token, agents)
	if err != nil {
		return err
	}

	deltas, err := appendHistory("data/history.jsonl", stats)
	if err != nil {
		return err
	}

	if err := renderReadme("templates/readme.tmpl", "README.md", stats, deltas); err != nil {
		return err
	}

	fmt.Printf("updated %d agents\n", len(stats))
	return nil
}
