package main

import (
	"encoding/json"
	"os"
	"time"
)

// siteRow is one ranked repo in site/data.json, consumed by site/index.html.
type siteRow struct {
	Key           string `json:"key"` // canonical owner/repo, matches history keys
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
	Stars         int    `json:"stars"`
	Delta7d       int    `json:"delta7d"`
	HasDelta      bool   `json:"hasDelta"`
	Language      string `json:"language"`
	PushedAt      string `json:"pushedAt"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Notes         string `json:"notes,omitempty"`
}

type siteData struct {
	UpdatedAt string     `json:"updatedAt"`
	Rows      []siteRow  `json:"rows"`
	History   []Snapshot `json:"history"`
}

// writeSiteData emits the JSON payload for the GitHub Pages dashboard.
// The file is generated fresh on every updater run and is not committed;
// the Pages deploy step in the workflow picks it up from the working tree.
func writeSiteData(path string, stats []Stat, deltas map[string]int, history []Snapshot) error {
	rows := make([]siteRow, len(stats))
	for i, s := range stats {
		delta, has := deltas[s.CanonicalKey]
		rows[i] = siteRow{
			Key:           s.CanonicalKey,
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Delta7d:       delta,
			HasDelta:      has,
			Language:      s.Language,
			PushedAt:      s.PushedAt.Format("2006-01-02"),
			Description:   s.Description,
			Category:      s.Category,
			Notes:         s.Notes,
		}
	}

	out, err := json.Marshal(siteData{
		UpdatedAt: time.Now().UTC().Format("2006-01-02 15:04 UTC"),
		Rows:      rows,
		History:   history,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
