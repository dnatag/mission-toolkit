# Mission Toolkit

> **"Slow down the process to speed up the understanding"**

AIDD (Atomic Intent-Driven Development) — a workflow that bridges "Vibe Coding" (chaos) and "Spec-Driven Development" (bureaucracy). Forces a 🤝 Handshake before every coding task, keeping changes within human comprehension limits.

❌ Without AIDD: AI modifies 20 files in one go — you lose track and feel alienated from your own codebase

✅ With AIDD: You control atomic missions scoped to 1-5 files — changes you actually understand and own

Works with **Amazon Q CLI**, **Claude Code**, **Kiro CLI**, and **OpenCode**.

## The Commands

| Command | Purpose |
|---------|---------|
| `/m.plan` | Convert intent → structured mission |
| `/m.apply` | Execute the authorized plan |
| `/m.complete` | Archive mission, generate rich commit |
| `/m.debug` | Investigate bugs → diagnosis.md |

*Amazon Q & Kiro: Use `@` prefix (e.g., `@m.plan`). Inline arguments are ignored — just invoke the command and the AI will prompt you for your intent.*

## Installation

```bash
# Homebrew
brew tap dnatag/mission-toolkit && brew install mission-toolkit

# From source (Go 1.21+)
git clone https://github.com/dnatag/mission-toolkit.git && cd mission-toolkit
go build -o m main.go && sudo mv m /usr/local/bin/
```

## Quick Start

```bash
# 1. Initialize (project-specific)
m init --ai q    # or: claude, kiro, opencode

# Or initialize globally (applies to all projects)
m init --ai q --global

# 2. Plan a mission
/m.plan "Add user authentication to the API"

# 3. Review mission.md, then execute
/m.apply

# 4. Review changes, then complete
/m.complete
```

## How It Works

```
m.plan → 🤝 Review → m.apply → 🤝 Review → m.complete
```

1. You define intent, AI proposes scope and plan
2. You authorize the architecture
3. AI implements, you verify
4. System archives mission

## Documentation

- [Core Concepts](docs/concepts.md) — Philosophy, WET→DRY, complexity matrix
- [Workflows](docs/workflows.md) — Mission lifecycle, bugfix workflow, project structure
- [CLI Reference](docs/cli-reference.md) — All commands and options

## License

See [LICENSE](LICENSE) file.

## Versioning

```bash
m version                          # Check current version
./scripts/sync-version.sh v1.0.0   # Update version (maintainers)
```

## Release Process

1. Tag the release:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions automatically builds cross-platform binaries and publishes the release.

### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)
