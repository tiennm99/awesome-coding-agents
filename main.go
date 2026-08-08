package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	check := flag.Bool("check", false, "validate data/agents.yml offline (no network, no token) and exit")
	flag.Parse()

	if *check {
		if err := runCheck("data/agents.yml"); err != nil {
			log.Printf("check failed: %v", err)
			os.Exit(1)
		}
		return
	}

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

	snapshots, deltas7, deltas30, err := appendHistory("data/history.jsonl", stats)
	if err != nil {
		return err
	}

	if err := renderReadme("templates/readme.tmpl", "README.md", stats, deltas7); err != nil {
		return err
	}

	if err := writeSiteData("site/data.json", stats, deltas7, deltas30, snapshots); err != nil {
		return err
	}

	fmt.Printf("updated %d agents\n", len(stats))
	return nil
}
