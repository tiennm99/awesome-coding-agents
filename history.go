package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"
)

type Snapshot struct {
	Date  string         `json:"date"`
	Stars map[string]int `json:"stars"`
}

// canonicalKeyMigrations remaps old history keys to the current canonical
// owner/repo from agents.yml. Add entries here when a tracked repo is renamed
// and history.jsonl contains the old name.
//
// Known renames:
//   - block/goose was formerly tracked as aaif-goose/goose
var canonicalKeyMigrations = map[string]string{
	"aaif-goose/goose": "block/goose",
}

func appendHistory(path string, stats []Stat) (map[string]int, error) {
	today := time.Now().UTC().Format("2006-01-02")

	// Key snapshots by canonical owner/repo from agents.yml, not by the
	// API-returned nameWithOwner, so renames don't orphan historical data.
	current := Snapshot{Date: today, Stars: map[string]int{}}
	for _, s := range stats {
		current.Stars[s.CanonicalKey] = s.Stars
	}

	snapshots, err := readSnapshots(path)
	if err != nil {
		return nil, err
	}

	deltas := computeDeltas(snapshots, current)

	// Drop any pre-existing snapshot for today, then append current.
	kept := slices.DeleteFunc(snapshots, func(s Snapshot) bool { return s.Date == today })
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
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var s Snapshot
		if err := json.Unmarshal(line, &s); err != nil {
			// Log corruption so it's observable; don't silently drop lines.
			fmt.Fprintf(os.Stderr, "history.jsonl line %d: skipping malformed entry: %v\n", lineNum, err)
			continue
		}
		// Apply canonical-key migrations so history keyed under old names is
		// transparently remapped to the current agents.yml canonical key.
		s.Stars = applyMigrations(s.Stars)
		out = append(out, s)
	}
	return out, scanner.Err()
}

// applyMigrations rewrites any deprecated history keys to their current
// canonical form. Old keys are removed; new keys accumulate stars additively
// (in practice the old key had no concurrent new entry, so max is the same).
func applyMigrations(stars map[string]int) map[string]int {
	for old, canonical := range canonicalKeyMigrations {
		if v, ok := stars[old]; ok {
			if stars[canonical] < v {
				stars[canonical] = v
			}
			delete(stars, old)
		}
	}
	return stars
}

// writeSnapshots writes to a temp file then renames atomically so a crash
// mid-write never leaves history.jsonl truncated or partially written.
func writeSnapshots(path string, snapshots []Snapshot) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	writeErr := error(nil)
	for _, s := range snapshots {
		if err := enc.Encode(s); err != nil {
			writeErr = err
			break
		}
	}

	if syncErr := f.Sync(); syncErr != nil && writeErr == nil {
		writeErr = syncErr
	}
	f.Close()

	if writeErr != nil {
		os.Remove(tmp)
		return writeErr
	}

	return os.Rename(tmp, path)
}

// computeDeltas returns stars-now minus stars-at-or-before-cutoff for each
// repo. Cutoff = 7 days ago UTC. The chosen prior snapshot must be within a
// 3-day window of the cutoff; if cron was skipped for more than 10 days the
// delta would be misleadingly labeled "Δ7d", so we return no delta instead.
func computeDeltas(history []Snapshot, current Snapshot) map[string]int {
	deltas := map[string]int{}
	now := time.Now().UTC()
	cutoff := now.AddDate(0, 0, -7).Format("2006-01-02")
	// Accept snapshots in (cutoff - 3 days, cutoff]. A snapshot older than
	// cutoff-3d is too stale to label as a 7-day delta.
	lowerBound := now.AddDate(0, 0, -10).Format("2006-01-02")

	var base *Snapshot
	for i := range history {
		d := history[i].Date
		if d > lowerBound && d <= cutoff {
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
