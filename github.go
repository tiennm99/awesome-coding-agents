package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Stat struct {
	Owner         string
	Repo          string
	Category      string
	Notes         string
	Description   string
	Stars         int
	Language      string
	PushedAt      time.Time
	URL           string
	NameWithOwner string
}

type repoNode struct {
	StargazerCount  int    `json:"stargazerCount"`
	Description     string `json:"description"`
	PrimaryLanguage *struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	PushedAt      time.Time `json:"pushedAt"`
	URL           string    `json:"url"`
	NameWithOwner string    `json:"nameWithOwner"`
}

type graphQLResponse struct {
	Data   map[string]*repoNode `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Path    []any  `json:"path"`
	} `json:"errors"`
}

const repoFields = `
		stargazerCount
		description
		primaryLanguage { name }
		pushedAt
		url
		nameWithOwner
	`

func fetchStats(token string, agents []Agent) ([]Stat, error) {
	var b strings.Builder
	b.WriteString("query {\n")
	for i, a := range agents {
		fmt.Fprintf(&b, "  r%d: repository(owner: %q, name: %q) {%s}\n", i, a.Owner, a.Repo, repoFields)
	}
	b.WriteString("}\n")

	body, err := json.Marshal(map[string]string{"query": b.String()})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "awesome-coding-agents-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graphql HTTP %d: %s", resp.StatusCode, raw)
	}

	var out graphQLResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w (body=%s)", err, raw)
	}
	for _, e := range out.Errors {
		// repos that 404 or are renamed land here; continue with partial data
		fmt.Fprintf(os.Stderr, "graphql warn: %s (path=%v)\n", e.Message, e.Path)
	}

	stats := make([]Stat, 0, len(agents))
	for i, a := range agents {
		node := out.Data[fmt.Sprintf("r%d", i)]
		if node == nil {
			fmt.Fprintf(os.Stderr, "skip %s/%s — no data\n", a.Owner, a.Repo)
			continue
		}
		lang := ""
		if node.PrimaryLanguage != nil {
			lang = node.PrimaryLanguage.Name
		}
		stats = append(stats, Stat{
			Owner:         a.Owner,
			Repo:          a.Repo,
			Category:      a.Category,
			Notes:         a.Notes,
			Description:   node.Description,
			Stars:         node.StargazerCount,
			Language:      lang,
			PushedAt:      node.PushedAt,
			URL:           node.URL,
			NameWithOwner: node.NameWithOwner,
		})
	}
	return stats, nil
}

func sortByStars(stats []Stat) {
	sort.SliceStable(stats, func(i, j int) bool {
		return stats[i].Stars > stats[j].Stars
	})
}
