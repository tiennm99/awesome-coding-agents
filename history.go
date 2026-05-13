package main

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

type Snapshot struct {
	Date  string         `json:"date"`
	Stars map[string]int `json:"stars"`
}

func appendHistory(path string, stats []Stat) (map[string]int, error) {
	today := time.Now().UTC().Format("2006-01-02")
	current := Snapshot{Date: today, Stars: map[string]int{}}
	for _, s := range stats {
		current.Stars[s.NameWithOwner] = s.Stars
	}

	snapshots, err := readSnapshots(path)
	if err != nil {
		return nil, err
	}

	deltas := computeDeltas(snapshots, current)

	// drop any pre-existing snapshot for today, then append current
	kept := snapshots[:0]
	for _, s := range snapshots {
		if s.Date != today {
			kept = append(kept, s)
		}
	}
	kept = append(kept, current)

	return deltas, writeSnapshots(path, kept)
}

func readSnapshots(path string) ([]Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []Snapshot
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var s Snapshot
		if err := json.Unmarshal(line, &s); err != nil {
			// skip malformed lines rather than fail the whole run
			continue
		}
		out = append(out, s)
	}
	return out, scanner.Err()
}

func writeSnapshots(path string, snapshots []Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, s := range snapshots {
		if err := enc.Encode(s); err != nil {
			return err
		}
	}
	return nil
}

// computeDeltas returns stars-now minus stars-at-or-before-cutoff for each repo.
// Cutoff = 7 days ago UTC. If no snapshot is old enough, delta is omitted.
func computeDeltas(history []Snapshot, current Snapshot) map[string]int {
	deltas := map[string]int{}
	cutoff := time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02")

	var base *Snapshot
	for i := range history {
		if history[i].Date <= cutoff {
			base = &history[i]
		}
	}
	if base == nil {
		return deltas
	}
	for repo, cur := range current.Stars {
		if prev, ok := base.Stars[repo]; ok {
			deltas[repo] = cur - prev
		}
	}
	return deltas
}
