# Template Versioning & Dynamic Refresh

> Version-aware templates that detect staleness and guide users to refresh after `m` upgrades.

## Problem Statement

Templates and governance are deployed once during `m init` and never updated. When `m` is upgraded, deployed templates become stale — they may reference old commands or miss new capabilities (e.g., Beads integration commands). There's no mechanism to detect or resolve this drift.

## Requirements

- Persist init config (AI type, global mode, version) in `.mission/config.json`
- On every `m` command, check deployed template version against binary version
- On mismatch, print warning with exact re-init command (e.g., `Templates outdated (v2.3.1 → v2.4.0). Run 'm init --ai q' to update.`)
- `m init` always overwrites all templates — no customization preservation
- Version stamp embedded in deployed templates for traceability

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Refresh trigger | Warning on mismatch, user re-runs `m init` | User controls when/how to refresh; knows their setup |
| Config persistence | `.mission/config.json` with ai_type, global, version | Enables exact re-init command in warning |
| Customization | Always overwrite on `m init` | Simple, predictable; templates are managed by `m` |

## Config File Format

Written to `.mission/config.json` during `m init`:

```json
{
  "ai_type": "q",
  "global": false,
  "version": "v2.3.1"
}
```

## Warning Behavior

On every `m` command (except `m init` and when no project is initialized):

```
⚠️  Templates outdated (v2.3.1 → v2.4.0). Run 'm init --ai q' to update.
```

For global mode:

```
⚠️  Templates outdated (v2.3.1 → v2.4.0). Run 'm init --ai q --global' to update.
```

No warning when:
- `.mission/` doesn't exist (no project initialized)
- `.mission/config.json` doesn't exist (legacy project, pre-versioning)
- Versions match

## Task Breakdown

### Task 1: Create config persistence

- **Objective**: Create `pkg/config/config.go` with `Config` struct (`AIType`, `Global`, `Version`) and `Read`/`Write` functions for `.mission/config.json`.
- **Guidance**: `Read` returns zero-value config (no error) if file doesn't exist — backward compatible with pre-versioning projects. `Write` creates `.mission/` directory if needed.
- **Test**: Round-trip read/write. Missing file returns zero-value config. Malformed JSON returns error.
- **Demo**: Config file created and read back correctly.

### Task 2: Write config during `m init`

- **Objective**: Update `cmd/init.go` to write `.mission/config.json` with current AI type, global flag, and `version.Version` after deploying templates.
- **Guidance**: Write config as the last step of `m init`, after all templates are deployed. This ensures config only exists if init succeeded.
- **Test**: After `m init --ai q`, config.json contains `{"ai_type":"q","global":false,"version":"v2.3.1"}`. After `m init --ai claude --global`, config.json contains `{"ai_type":"claude","global":true,"version":"v2.3.1"}`.
- **Demo**: `m init --ai q` creates config file alongside templates.

### Task 3: Add version mismatch check to root command

- **Objective**: In `cobra.OnInitialize`, read `.mission/config.json`, compare version against `version.Version`. If mismatch, print warning to stderr with exact command.
- **Guidance**: Only warn if `.mission/config.json` exists (skip for `m init` itself and when no project is initialized). Build command string from persisted config: `m init --ai <type>` or `m init --ai <type> --global`. Print to stderr so it doesn't interfere with JSON output parsing.
- **Test**: Version match → no warning. Mismatch → warning printed to stderr. No config file → no warning. No `.mission/` directory → no warning.
- **Demo**: After upgrading `m`, any command shows warning with exact re-init command.

## Files Affected

### New Files
- `pkg/config/config.go` — Config struct and Read/Write functions

### Modified Files
- `cmd/init.go` — Write config after template deployment
- `cmd/root.go` — Version mismatch check in `cobra.OnInitialize`

### Runtime Files
- `.mission/config.json` — Persisted init configuration
