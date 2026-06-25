# Awesome Coding Agents

> Curated ranking of AI agent coding tools, sorted by GitHub stars.
> Updated daily by GitHub Actions.

**Last updated:** 2026-06-25 03:51 UTC · **Tracked:** 19 repos

| # | Repo | Stars | Δ7d | Language | Last push | Description |
|---|------|------:|----:|----------|-----------|-------------|
| 1 | [anomalyco/opencode](https://github.com/anomalyco/opencode) | 178.2k | +2474 | TypeScript | 2026-06-25 | The open source coding agent. |
| 2 | [anthropics/claude-code](https://github.com/anthropics/claude-code) | 134.2k | +1146 | Python | 2026-06-24 | Claude Code is an agentic coding tool that lives in your terminal, understands your codebase, and helps you code faster by executing routine tasks, explaining complex code, and handling git workflows - all through natural language commands. |
| 3 | [google-gemini/gemini-cli](https://github.com/google-gemini/gemini-cli) | 105.5k | +166 | TypeScript | 2026-06-25 | An open-source AI agent that brings the power of Gemini directly into your terminal. |
| 4 | [openai/codex](https://github.com/openai/codex) | 93.5k | +1643 | Rust | 2026-06-25 | Lightweight coding agent that runs in your terminal |
| 5 | [zed-industries/zed](https://github.com/zed-industries/zed) | 85.9k | +444 | Rust | 2026-06-25 | Code at the speed of thought – Zed is a high-performance, multiplayer code editor from the creators of Atom and Tree-sitter. |
| 6 | [cline/cline](https://github.com/cline/cline) | 63.9k | +399 | TypeScript | 2026-06-25 | Autonomous coding agent as an SDK, IDE extension, or CLI assistant. |
| 7 | [AntonOsika/gpt-engineer](https://github.com/AntonOsika/gpt-engineer) | 55.2k | -1 | Python | 2025-05-14 | CLI platform to experiment with codegen. Precursor to: https://lovable.dev |
| 8 | [aaif-goose/goose](https://github.com/aaif-goose/goose) | 50.1k | +435 | Rust | 2026-06-25 | an open source, extensible AI agent that goes beyond code suggestions - install, execute, edit, and test with any LLM |
| 9 | [Aider-AI/aider](https://github.com/Aider-AI/aider) | 46.7k | +274 | Python | 2026-05-22 | aider is AI pair programming in your terminal |
| 10 | [continuedev/continue](https://github.com/continuedev/continue) | 34.4k | +481 | TypeScript | 2026-06-24 | open-source coding agent |
| 11 | [TabbyML/tabby](https://github.com/TabbyML/tabby) | 33.7k | +42 | Rust | 2026-03-02 | Self-hosted AI coding assistant |
| 12 | [voideditor/void](https://github.com/voideditor/void) | 28.8k | +7 | TypeScript | 2026-06-02 |  |
| 13 | [charmbracelet/crush](https://github.com/charmbracelet/crush) | 25.7k | +242 | Go | 2026-06-25 | Glamourous agentic coding for all 💘 |
| 14 | [RooCodeInc/Roo-Code](https://github.com/RooCodeInc/Roo-Code) | 24.3k | +37 | TypeScript | 2026-05-15 | Roo Code gives you a whole dev team of AI agents in your code editor. |
| 15 | [kortix-ai/suna](https://github.com/kortix-ai/suna) | 19.9k | +26 | TypeScript | 2026-06-25 | The Company AI Command Center |
| 16 | [SWE-agent/SWE-agent](https://github.com/SWE-agent/SWE-agent) | 19.6k | +68 | Python | 2026-06-22 | SWE-agent takes a GitHub issue and tries to automatically fix it, using your LM of choice. It can also be employed for offensive cybersecurity or competitive coding challenges. [NeurIPS 2024] |
| 17 | [yetone/avante.nvim](https://github.com/yetone/avante.nvim) | 18.0k | +24 | Lua | 2026-06-24 | Use your Neovim like using Cursor AI IDE! |
| 18 | [stackblitz/bolt.new](https://github.com/stackblitz/bolt.new) | 16.4k | +12 | TypeScript | 2024-12-17 | Prompt, run, edit, and deploy full-stack web applications. -- bolt.new -- Help Center: https://support.bolt.new/ -- Community Support: https://discord.com/invite/stackblitz |
| 19 | [plandex-ai/plandex](https://github.com/plandex-ai/plandex) | 15.5k | +16 | Go | 2025-10-03 | Open source AI coding agent. Designed for large projects and real world tasks. |

---

## How it works

1. `data/agents.yml` is the curated source list.
2. A daily GitHub Actions workflow (`.github/workflows/update.yml`) runs the Go updater.
3. The updater fetches live repo metadata via the GitHub GraphQL API in one batched query.
4. Star counts are appended to `data/history.jsonl` for 7-day delta computation.
5. This `README.md` is regenerated from `templates/readme.tmpl` and committed back to the repo.

**Δ7d:** Change in stars over the past 7 days; `—` means fewer than 7 days of history.

## Contributing

Add an agent to [`data/agents.yml`](data/agents.yml):

```yaml
agents:
  - owner: github-username-or-org
    repo: repository-name
    category: cli   # cli | ide | extension | library | research | web
```

Open a PR. The next daily run picks it up automatically.

## License

Apache-2.0
