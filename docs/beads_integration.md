# Beads Integration for Backlog Management

> Integrate [Beads (`bd`)](https://github.com/steveyegge/beads) as the source of truth for backlog management, with `m backlog` as a thin wrapper preserving the existing prompt workflow.

## Problem Statement

The current `m backlog` uses a flat markdown file (`backlog.md`) as its store. Beads (`bd`) offers a richer, dependency-aware, git-backed graph model that AI agents interact with more natively. We want Beads as source of truth when available, `m backlog` as a thin wrapper preserving the existing prompt workflow, with fallback to the current markdown system. Since AI agents may bypass `m backlog` and directly create/edit `backlog.md`, we need automatic reconciliation on every `m backlog` call.

Additionally, the decompose template already produces dependency information between sub-intents, but this is currently ignored. With Beads' dependency graph, we can wire blocking relationships so `bd ready` tells agents exactly which task to pick up next.

## Requirements

- Beads is source of truth when detected (auto-detect: `bd` on PATH + `.beads/` exists)
- `m backlog` CLI interface stays identical — prompt templates don't change
- Hierarchical model: one Beads epic per backlog type (features, bugfixes, decomposed, refactor, future), items as children — created upfront
- `m backlog list` generates a read-only `backlog.md` snapshot for prompt context
- Every `m backlog` call reconciles: detects rogue `backlog.md` edits, imports new items into Beads, then overwrites with canonical snapshot
- Best-effort rogue parsing: any `- [ ]` item found, infer type from section headers if present, default to `feature`
- Direct `bd` usage by agents is transparent — changes appear naturally in next snapshot
- Fallback to current markdown manager when Beads is unavailable
- Pattern tracking (Rule-of-Three) maps to Beads task notes
- Epic decomposition wires Beads dependency graph via batch `m backlog decompose` command
- `plan-epic.md` display shows visual dependency tree

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Source of truth | Beads | AI agents interact with `bd` more natively than `m backlog` |
| `m backlog` role | Thin wrapper | Preserves existing prompt workflow without template changes |
| Activation | Auto-detect (`bd` on PATH + `.beads/` exists) | No explicit opt-in needed; user runs `bd init` separately |
| Type mapping | Hierarchical (epic per type, items as children) | Clean structure, leverages Beads' native hierarchy |
| Epic creation | Upfront (all 5 on first detection) | Predictable structure, low noise |
| `backlog.md` when Beads active | Read-only snapshot, generated on demand | Keeps prompt templates working; agents that read it see current state |
| Rogue edit handling | Auto-reconcile on every `m backlog` call | Most robust — agents can bypass anything, reconciliation is the answer |
| Rogue parsing | Best-effort (any `- [ ]` item, infer type from sections, default to `feature`) | Handles unpredictable agent-generated markdown |
| Direct `bd` usage | Transparent | Changes appear naturally in next snapshot since Beads is source of truth |
| Dependency wiring | Batch command (`m backlog decompose --json`) | Takes entire decompose JSON, creates items + deps in one shot |
| Dependency display | Visual text-based dependency tree in `plan-epic.md` | Shows blocking relationships clearly to the user |

## Architecture

```
m backlog (cmd/backlog.go)
    │
    ▼
NewProvider(missionDir)
    │
    ├── bd detected ──► BeadsProvider
    │                       │
    │                       ├── Reconcile (on every call)
    │                       │     ├── Rogue backlog.md? → Parse & import into Beads
    │                       │     └── Write canonical backlog.md snapshot
    │                       │
    │                       ├── Decompose (batch create + dependency wiring)
    │                       │     ├── Create child tasks under Decomposed Intents epic
    │                       │     └── bd dep add <child> <blocker> for each dependency
    │                       │
    │                       └── bd CLI (os/exec)
    │
    └── fallback ──► MarkdownProvider (existing manager.go)
                         │
                         └── Decompose (batch create, dependencies ignored)
```

## Type Mapping

| Mission Toolkit Type | Beads Epic Title | Beads Issue Type |
|---------------------|------------------|-----------------|
| `feature` | Features | epic (parent), task (children) |
| `bugfix` | Bugfixes | epic (parent), task (children) |
| `decomposed` | Decomposed Intents | epic (parent), task (children) |
| `refactor` | Refactoring Opportunities | epic (parent), task (children) |
| `future` | Future Enhancements | epic (parent), task (children) |

Pattern tracking for refactor items uses Beads task notes: `PATTERN:<id> COUNT:<n>`.

## Dependency Graph for Epic Decomposition

### Current Flow (Track 4 in m.plan)

1. `m analyze decompose` → produces JSON with `sub_intents[].dependencies`
2. `m backlog add` for each sub-intent → flat, unordered items
3. Agent picks whichever sub-intent they want — dependencies ignored

### Proposed Flow (when Beads active)

1. `m analyze decompose` → same JSON with `sub_intents[].dependencies`
2. `m backlog decompose --json '<decompose output>'` → batch command that:
   a. Creates all sub-intents as Beads child tasks under "Decomposed Intents" epic
   b. Wires `bd dep add <child> <blocker>` blocking relationships from the dependency array
   c. Returns created task IDs and dependency graph
3. Agent runs `bd ready` → gets only tasks with no open blockers
4. `plan-epic.md` displays visual dependency tree

### Decompose JSON Input

The `m backlog decompose` command accepts the same JSON that `m analyze decompose` produces:

```json
{
  "action": "DECOMPOSE",
  "sub_intents": [
    {
      "intent": "Create payment data models and database schema",
      "rationale": "Foundation layer needed by all other components",
      "estimated_files": 3,
      "dependencies": []
    },
    {
      "intent": "Implement payment validation service",
      "rationale": "Core business logic for payment validation",
      "estimated_files": 4,
      "dependencies": ["Create payment data models and database schema"]
    },
    {
      "intent": "Create payment API endpoints",
      "rationale": "API layer depends on service layer",
      "estimated_files": 4,
      "dependencies": ["Implement payment validation service"]
    }
  ],
  "decomposition_rationale": "Split by architectural layers"
}
```

### Dependency Tree Display (plan-epic.md)

```
🏔️ EPIC DETECTED — Decomposed into 4 sub-intents with dependency tracking

📋 DEPENDENCY GRAPH:
  Create payment data models          (bd-a3f8.1) ← no blockers, READY
  └── Implement payment validation    (bd-a3f8.2) ← blocked by .1
      └── Create payment API endpoints(bd-a3f8.3) ← blocked by .2
  Add Stripe integration              (bd-a3f8.4) ← blocked by .2

🚀 READY TO START (no blockers):
  • bd-a3f8.1: Create payment data models

💡 Use 'bd ready' to see available tasks at any time
```

### Fallback Behavior (no Beads)

When Beads is not available, `m backlog decompose` falls back to the current behavior: creates items via `m backlog add --type decomposed` in dependency order. The dependency information is preserved in the item description as a hint (e.g., `"(depends on: Create data models)"`) but not enforced.

## Reconciliation Flow

Runs at the start of every `BeadsProvider` method call:

1. Check if `.mission/backlog.md` exists
2. Compute SHA256 hash, compare against stored hash in `.mission/beads-backlog-hash`
3. If hashes match → no rogue edits, skip to step 6
4. If hashes differ → rogue edit detected:
   a. Parse rogue file using `pkg/md` list extraction
   b. Try recognized section headers (`## FEATURES`, `## BUGFIXES`, etc.) to infer type
   c. Any `- [ ]` items in unrecognized sections default to `feature`
   d. Diff parsed items against current Beads state (by description text match)
   e. Import only genuinely new items via `bd create`
5. Log reconciliation results
6. After command execution, regenerate canonical snapshot from Beads state
7. Update stored hash

## Task Breakdown

### Task 1: Extract BacklogProvider interface

- **Objective**: Define a `BacklogProvider` interface in `pkg/backlog/provider.go` capturing the public API: `List`, `Add`, `AddWithPattern`, `AddMultiple`, `Complete`, `Cleanup`, `GetPatternCount`, `Decompose`.
- **Guidance**: Create the interface file. Existing `BacklogManager` already satisfies the original methods — no changes to `manager.go`. The `Decompose` method is new: on `BacklogManager` it falls back to `AddMultiple` with dependency hints in descriptions. Add compile-time check.
- **Test**: `var _ BacklogProvider = (*BacklogManager)(nil)` compiles. All existing tests pass.
- **Demo**: `go test ./pkg/backlog/...` passes with interface in place.

### Task 2: Wire factory into CLI commands

- **Objective**: Add `NewProvider(missionDir string) BacklogProvider` factory function. Replace all `backlog.NewManager(missionDir)` calls in `cmd/backlog.go` with `backlog.NewProvider(missionDir)`.
- **Guidance**: Factory returns `BacklogManager` for now (behavior identical). Update the 4 existing command handlers. Add new `m backlog decompose --json` command that calls `provider.Decompose()`.
- **Test**: `go test ./...` passes. Existing CLI commands behave identically.
- **Demo**: `m backlog list`, `m backlog add "test" --type feature`, `m backlog decompose --json '<json>'` all work.

### Task 3: Implement Beads detection

- **Objective**: Create `pkg/backlog/beads_detect.go` — check if `bd` is on PATH and `.beads/` directory exists.
- **Guidance**: `exec.LookPath("bd")` + `os.Stat(filepath.Join(projectRoot, ".beads"))`. Accept a `CommandRunner` interface (for testability) and a filesystem path. Return bool.
- **Test**: Table-driven tests covering: both present → true, bd missing → false, .beads missing → false.
- **Demo**: Detection function works in isolation.

### Task 4: Implement BeadsProvider — exec layer and epic management

- **Objective**: Create `pkg/backlog/beads.go` with `BeadsProvider` struct. Implement the `bd` CLI exec layer and upfront epic creation.
- **Guidance**: Define a `CommandRunner` interface (`Run(args ...string) (string, error)`) following `pkg/git/cmd.go` pattern. `BeadsProvider` stores project root and a `CommandRunner`. Implement `ensureEpics()`: reads/creates a `.mission/beads-epics.json` cache mapping type→epicID. On first call, creates 5 epics via `bd create "<Type>" -t epic --json`, stores IDs. On subsequent calls, reads cache.
- **Test**: Mock `CommandRunner` tests verifying correct `bd create` calls and cache file read/write.
- **Demo**: Epic creation and caching works against mock.

### Task 5: Implement BeadsProvider — List and Add

- **Objective**: Implement `List`, `Add`, `AddWithPattern`, `AddMultiple` on `BeadsProvider`.
- **Guidance**:
  - `List`: For each type epic, `bd list --parent <epicID> --json`. Filter by include/exclude params. Format as `- [ ] description` / `- [x] description` to match existing output.
  - `Add`: `bd create "<desc>" -t task --parent <epicID> --json`.
  - `AddWithPattern`: Same as Add, plus `bd update <id> --notes "PATTERN:<id> COUNT:<n>"` for refactor pattern tracking.
  - `AddMultiple`: Loop, create each child under the type epic.
- **Test**: Mock `CommandRunner` tests verifying correct `bd` arguments for each operation.
- **Demo**: List and Add work against mock executor.

### Task 6: Implement BeadsProvider — Complete, Cleanup, GetPatternCount

- **Objective**: Implement remaining core interface methods.
- **Guidance**:
  - `Complete`: List children of relevant epics, find matching item by text, `bd close <id> --reason "Completed"`.
  - `Cleanup`: List closed items across epics. For Beads, "cleanup" can be a no-op (closed items are already closed in Beads) or call `bd admin compact`. Keep it simple — just return count of closed items.
  - `GetPatternCount`: List refactor epic children, search notes for `PATTERN:<id>` marker, extract count.
- **Test**: Mock-based tests for each method.
- **Demo**: Full core `BacklogProvider` interface satisfied by `BeadsProvider`.

### Task 7: Implement BeadsProvider — Decompose with dependency wiring

- **Objective**: Implement `Decompose` on `BeadsProvider` — batch create sub-intents with Beads dependency graph.
- **Guidance**:
  - Parse the decompose JSON input (same schema as `m analyze decompose` output).
  - Create each sub-intent as a child task under the "Decomposed Intents" epic via `bd create`.
  - Build a name→ID map from created tasks.
  - Wire dependencies via `bd dep add <child> <blocker>` using the `dependencies` array (match by intent text against the name→ID map).
  - Return structured result with created IDs and dependency graph for display.
  - On `BacklogManager` (fallback): call `AddMultiple` with dependency hints appended to descriptions (e.g., `"(depends on: Create data models)"`).
- **Test**: Mock `CommandRunner` tests verifying: correct `bd create` calls for each sub-intent, correct `bd dep add` calls matching dependency array, fallback adds items with dependency hints.
- **Demo**: Decompose JSON with 4 sub-intents and 3 dependencies → 4 `bd create` + 3 `bd dep add` calls.

### Task 8: Implement reconciliation

- **Objective**: Create `pkg/backlog/reconcile.go` — detect rogue `backlog.md` edits and import into Beads.
- **Guidance**:
  - On every `BeadsProvider` method call, run reconciliation first.
  - Store a SHA256 hash of the last canonical snapshot in `.mission/beads-backlog-hash`.
  - If `.mission/backlog.md` exists and its hash differs from stored hash → rogue edit detected.
  - Parse rogue file using `pkg/md` list extraction: try recognized sections (`## FEATURES`, etc.) to infer type, then scan for any remaining `- [ ]` items and default to `feature`.
  - Diff parsed items against current Beads state (by description text match). Import only genuinely new items via `bd create`.
  - After reconciliation (and after every command), regenerate canonical snapshot and update hash.
- **Test**: Table-driven tests: no rogue file → no-op; rogue file with new items → items imported; rogue file with existing items → no duplicates; unrecognized sections → items default to feature.
- **Demo**: Create a rogue `backlog.md` with new items, run `m backlog list`, verify items appear in Beads output.

### Task 9: Implement canonical snapshot generation

- **Objective**: Add snapshot generation to `BeadsProvider` — writes `.mission/backlog.md` from Beads state after every operation.
- **Guidance**: Query all epics and children. Format into standard `backlog.md` template (sections per type, `- [ ]` / `- [x]` items). Include header comment: `<!-- MANAGED BY BEADS — Source of truth is 'bd'. Use 'm backlog' commands. -->`. For decomposed items with dependencies, include dependency info in the snapshot (e.g., `- [ ] Create API endpoints (blocked by: Create data models)`). Write file + update hash in `.mission/beads-backlog-hash`.
- **Test**: Verify generated markdown matches expected format. Verify hash file is updated. Verify dependency info appears for decomposed items.
- **Demo**: After `m backlog decompose`, `backlog.md` reflects items with dependency annotations.

### Task 10: Update prompt templates for dependency-aware decomposition

- **Objective**: Update `m.plan` prompt template (Track 4 section) and `plan-epic.md` display template to use `m backlog decompose --json` and show dependency graph.
- **Guidance**:
  - In `m.plan.md` Track 4 section: replace the `m backlog add` loop with a single `m backlog decompose --json '<decompose output>'` call.
  - In `plan-epic.md`: add `{{DEPENDENCY_GRAPH}}` variable showing the visual tree (foundations at top, dependents indented below with arrows). Include `bd ready` hint.
  - Keep backward compatibility: when Beads is not active, the decompose command falls back to flat adds, and the display omits the dependency graph.
- **Test**: Template renders correctly with and without dependency graph data.
- **Demo**: `m.plan` on an epic intent shows dependency tree in output when Beads is active.

### Task 11: Wire BeadsProvider into factory with auto-detection

- **Objective**: Update `NewProvider` to auto-detect Beads and return `BeadsProvider` when available, `BacklogManager` otherwise.
- **Guidance**: Factory calls detection from Task 3. If detected → `BeadsProvider` (which runs `ensureEpics` + reconciliation on first use). Otherwise → `BacklogManager`. Print a one-line log indicating which provider is active.
- **Test**: Verify factory returns correct provider type. End-to-end: with mock `bd`, full cycle of add → list → decompose → complete works through the provider interface.
- **Demo**: `m backlog list` works with both providers. When Beads is active, `backlog.md` is a managed snapshot with dependency info. When not, it's the existing read/write file.

## Files Affected

### New Files
- `pkg/backlog/provider.go` — BacklogProvider interface + factory
- `pkg/backlog/beads_detect.go` — Beads detection logic
- `pkg/backlog/beads.go` — BeadsProvider implementation (exec layer, epic management, CRUD, decompose)
- `pkg/backlog/reconcile.go` — Reconciliation logic
- `pkg/backlog/snapshot.go` — Canonical snapshot generation

### Modified Files
- `cmd/backlog.go` — Use `NewProvider()` instead of `NewManager()`, add `decompose` subcommand
- `pkg/backlog/manager.go` — Add `Decompose` fallback method (flat adds with dependency hints)
- `pkg/templates/prompts/m.plan.md` — Track 4 section uses `m backlog decompose --json`
- `pkg/templates/libraries/displays/plan-epic.md` — Add `{{DEPENDENCY_GRAPH}}` display

### Runtime Files (managed by BeadsProvider)
- `.mission/beads-epics.json` — Cache of epic IDs per type
- `.mission/beads-backlog-hash` — SHA256 hash of last canonical snapshot
- `.mission/backlog.md` — Read-only snapshot (when Beads active)
