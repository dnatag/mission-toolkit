# Auto-Continue (m.auto)

> Skip review handshakes for simple Track 2 missions

## Problem

The 3-stop handshake workflow (`m.plan → 🤝 review → m.apply → 🤝 review → m.complete`) is overkill for mechanical refactoring and straightforward implementations. Users want to opt out of review stops when they trust the task.

## Design Decisions

### Approach: Flag Setter Prompt

A new `@m.auto` / `/m.auto` prompt sets an `auto` field in mission frontmatter, then follows the normal `m.plan` flow. Existing prompts (`m.plan.md`, `m.apply.md`) check this field at their success exit points and conditionally continue into the next phase.

**Why this approach:**
- Auto intent persisted in mission frontmatter — survives context loss, is inspectable, archived with mission
- Each prompt still loads as the primary prompt for its phase — full AI fidelity
- Minimal duplication — existing prompts are the single source of truth

**Alternatives considered:**
- *Single combined prompt* — reliability concerns with massive prompt, duplication of all three phases
- *Thin orchestrator reading prompts on-demand* — AI follows mid-session file reads less reliably than primary prompts
- *CLI-orchestrated pipeline* — CLI can't drive the AI; the AI does the implementation work
- *Full composable flags (`--auto-apply` × `--auto-complete`)* — over-engineered; `--auto-complete` alone (skip code review but review plan) is not a useful mode

### Three Modes

| Mode | Invocation | Workflow | Use Case |
|---|---|---|---|
| Normal | `@m.plan` | `plan → 🤝 → apply → 🤝 → complete` | Default, full review |
| Auto-apply | `@m.auto --auto=apply "intent"` | `plan → 🤝 → apply → complete` | Review plan, trust execution |
| Full auto | `@m.auto --auto=full "intent"` | `plan → apply → complete` | Mechanical refactoring, trusted tasks |

The user must explicitly pass `--auto=apply` or `--auto=full`. The AI never decides the auto mode — it only reads the flag and sets the frontmatter value accordingly.

### Governance Compliance

The 🤝 Handshake principle states: *"It forces a Handshake before every coding task."* Auto-continue is a **user-initiated opt-in** that extends the existing elasticity precedent (Track 1 already bypasses the entire workflow).

**Guardrails:**
- **User-initiated only** — The AI never self-decides to skip handshakes
- **CLI gates still enforced** — All `m mission check`, verification, and STOP signals respected. Auto-continue only skips the "display success and wait" pause
- **Track 2 only** — CLI clears `auto` field in Status JSON for Track 3+ (cross-cutting concerns need human review). Enforcement is in the CLI, not the prompt, consistent with governance principle that the AI should not self-classify complexity
- **Failure stops the pipeline** — Verification failure halts auto-continue; no auto-retry
- **Pause/restore preserves auto** — The `auto` field persists in frontmatter through pause/restore since it represents user intent

### Flow Diagram

```
@m.auto --auto=apply "refactor X"   @m.auto --auto=full "refactor X"
  │                                   │
  ├─ Set auto=apply                   ├─ Set auto=full
  ├─ Follow m.plan.md                 ├─ Follow m.plan.md
  │    └─ Finalize:                   │    └─ Finalize:
  │         auto=apply, Track 2            auto=full, Track 2
  │         → display plan-success         → continue to m.apply
  │         → STOP (🤝 review plan)   │
  │                                   ├─ Follow m.apply.md
  ├─ User invokes @m.apply            │    └─ On Success:
  │    └─ On Success:                 │         auto=full
  │         auto=apply                         → continue to m.complete
  │         → continue to m.complete  │
  │                                   └─ Follow m.complete.md
  └─ Follow m.complete.md
```

## Implementation Plan (MVP)

### Task 1: Add `Auto` field to Mission struct and frontmatter

**Files:** `pkg/mission/mission.go`, `pkg/mission/writer.go`

- Add `Auto string \`yaml:"auto,omitempty"\`` to `Mission` struct
- Handle `auto` key in `UpdateFrontmatter()` switch statement with validation (only `apply` or `full`)
- Include `auto` in `format()` frontmatter map when non-empty

**Verify:** `m mission update --frontmatter auto=apply` persists field, round-trip read confirms it.

### Task 2: Expose `auto` in CheckService Status JSON with Track 2 enforcement

**Files:** `pkg/mission/check.go`

- Add `Auto string \`json:"auto,omitempty"\`` to `Status` struct
- Populate `status.Auto` from `mission.Auto` in `handleExistingMission()`
- Track 2 enforcement: if `mission.Track > 2`, set `status.Auto = ""` regardless of frontmatter value
- If track increases beyond 2 during mission execution, auto is effectively disabled

**Verify:** `m mission check --context apply` returns `"auto": "apply"` for Track 2, omits it for Track 3+.

### Task 3: Create `m.auto.md` prompt template

**Files:** `pkg/templates/prompts/m.auto.md`

- Parse `$ARGUMENTS` for required `--auto=apply` or `--auto=full` flag
- If neither flag is provided → display usage and STOP
- Strip the flag from arguments to extract intent
- After mission creation in m.plan Step 1: run `m mission update --frontmatter auto=[apply|full]`
- Instruct AI to read and follow `m.plan.md` for all remaining steps
- The AI must never infer or default the auto mode — it must come from the user's explicit flag

### Task 4: Add auto-continue logic to prompt exit points

**Files:** `pkg/templates/prompts/m.plan.md`, `pkg/templates/prompts/m.apply.md`

**m.plan.md** — Append after Step 5, item 4 (Final Output) which ends with outputting the filled `plan-success.md` template:

```markdown
### Step 6: Auto-Continue Check

1. **Check Auto Mode**: Run `m mission check --context apply` and parse JSON output
2. **React Based on `auto` field**:
   - If `auto` is `full` and `next_step` says "PROCEED" → Use file read tool to load the m.apply prompt file from the AI-specific prompt directory. Follow all its instructions starting from Step 0, skipping the Prerequisites section (state was just validated).
   - If `auto` is `apply` or empty → **STOP**. User reviews plan and invokes /m.apply manually.
```

**m.apply.md** — Append after Step 4 "On Success" block which ends with loading `apply-success.md` template:

```markdown
### Step 5: Auto-Continue Check

1. **Check Auto Mode**: Run `m mission check --context complete` and parse JSON output
2. **React Based on `auto` field**:
   - If `auto` is `apply` or `full`, and `next_step` says "PROCEED" → Use file read tool to load the m.complete prompt file from the AI-specific prompt directory. Follow all its instructions starting from Step 0, skipping the Prerequisites section (state was just validated).
   - If `auto` is empty → **STOP**. User invokes /m.complete manually.

**On any failure in Step 2-4**: Always **STOP** regardless of `auto` field.
```

**Error handling:** If reading the next prompt file fails, fall back to normal STOP behavior. The plan-success or apply-success template was already displayed, so the user can continue manually.

### Task 5: Deploy `m.auto.md` via template system

**Files:** `pkg/templates/templates.go`

- Add `"m.auto.md"` to `promptFiles` slice in `WriteTemplates()`
- Existing `/m.` → `@m.` prefix replacement handles client-specific prefixes

## Data Model

### Mission Frontmatter

```yaml
---
id: MISS-20260207-143500
type: WET
track: 2
iteration: 1
status: planned
auto: apply             # NEW — "apply" or "full", omitted when empty
---
```

### Auto Field Values

| Value | Meaning | Skips |
|---|---|---|
| (empty/omitted) | Normal handshake workflow | Nothing |
| `apply` | Auto-continue after apply | 🤝 #2 (code review) |
| `full` | Auto-continue entire pipeline | 🤝 #1 (plan review) + 🤝 #2 (code review) |

### Status JSON

```json
{
  "has_active_mission": true,
  "mission_status": "planned",
  "mission_id": "MISS-20260207-143500",
  "auto": "apply",
  "ready": false,
  "message": "Mission is ready for execution or re-execution",
  "next_step": "PROCEED with m.apply execution."
}
```

### Task 6: Update governance, docs, and display templates

**Files:** `pkg/templates/mission/governance.md`, `docs/workflows.md`, `docs/concepts.md`, `README.md`, `pkg/templates/libraries/displays/plan-success.md`, `pkg/templates/libraries/displays/apply-success.md`

- Governance Section II (Elasticity): add auto-continue principle for Track 2
- Governance WORKFLOW section: add auto-continue variant
- `docs/workflows.md`: show auto-continue flow alongside standard workflow
- `README.md`: add `m.auto` to commands table
- `plan-success.md`: show auto-continue hint when `auto` field is set
- `apply-success.md`: show auto-continue hint when `auto` field is set
