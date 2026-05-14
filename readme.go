package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Row struct {
	Rank          int
	NameWithOwner string
	URL           string
	Stars         int
	Delta7d       int
	HasDelta      bool
	Language      string
	PushedAt      string
	Description   string
	Category      string
}

func renderReadme(tmplPath, outPath string, stats []Stat, deltas map[string]int) error {
	rows := make([]Row, len(stats))
	for i, s := range stats {
		// Deltas are keyed by CanonicalKey (owner/repo from agents.yml).
		delta, has := deltas[s.CanonicalKey]
		rows[i] = Row{
			Rank:          i + 1,
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Delta7d:       delta,
			HasDelta:      has,
			Language:      s.Language,
			PushedAt:      s.PushedAt.Format("2006-01-02"),
			Description:   sanitizeCell(s.Description),
			Category:      s.Category,
		}
	}

	funcs := template.FuncMap{
		"formatDelta": func(d int, has bool) string {
			if !has {
				return "—"
			}
			if d > 0 {
				return fmt.Sprintf("+%d", d)
			}
			if d < 0 {
				return fmt.Sprintf("%d", d)
			}
			return "0"
		},
		"formatStars": func(n int) string {
			switch {
			case n < 1000:
				return fmt.Sprintf("%d", n)
			case n < 1_000_000:
				return fmt.Sprintf("%.1fk", float64(n)/1000)
			default:
				return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
			}
		},
	}

	tmpl, err := template.New("").Funcs(funcs).ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.ExecuteTemplate(f, filepath.Base(tmplPath), map[string]any{
		"Rows":      rows,
		"UpdatedAt": time.Now().UTC().Format("2006-01-02 15:04 UTC"),
		"Total":     len(rows),
	})
}

// sanitizeCell escapes pipe and newline characters so descriptions stay in one table cell.
func sanitizeCell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}
