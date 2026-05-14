# Contributing

## Adding an Agent

Edit [`data/agents.yml`](../data/agents.yml) and add an entry:

```yaml
agents:
  - owner: github-username-or-org
    repo: repository-name
    category: cli
```

Required fields: `owner`, `repo`, `category`

Optional fields: `notes` (for clarifications or caveats)

## Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `owner` | string | Yes | GitHub user or organization that owns the repo |
| `repo` | string | Yes | Repository name on GitHub |
| `category` | string | Yes | One of the values listed below |
| `notes` | string | No | Additional context or disclaimers |

## Valid Categories

- **cli** — Command-line tools and CLI wrappers
- **ide** — Standalone editors and IDEs
- **extension** — Editor extensions (VS Code, Neovim, etc.)
- **library** — Libraries and SDKs
- **research** — Research papers and proof-of-concept projects
- **web** — Web-based tools and online IDEs

## Handling Duplicates and Changes

**Duplicate repos:** If a repo appears twice (same owner/repo), the second entry is ignored in the next daily run.

**Renamed repos:** If a GitHub repo is renamed after being tracked, update the `repo` field in `data/agents.yml`. The updater will fetch fresh metadata under the new name. Historical data in `data/history.jsonl` remains under the old key and is not migrated.

**Deprecation:** To remove an agent, delete its entry from `data/agents.yml`. The next run will drop it from the README; historical data is preserved.

## PR Review

- Keep PRs to changes in `data/agents.yml` only (do not edit `README.md` or `data/history.jsonl`)
- The daily GitHub Actions workflow (runs at 00:00 UTC) picks up merged PRs automatically
- No manual review required; the updater regenerates the README after your PR merges

For local testing before opening a PR, see [LOCAL_DEV.md](./LOCAL_DEV.md).
