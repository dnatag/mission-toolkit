# Debug Workflow Design

> **Separation of Concerns**: m.debug investigates, m.plan plans, m.apply executes

## Overview

The `m.debug` command addresses the gap between the linear feature-development workflow and the exploratory nature of debugging. Instead of adding complexity to m.plan, we introduce a dedicated investigation phase that produces a structured diagnosis.

## Flow

```
m.debug "symptom" → diagnosis.md → m.plan (consumes diagnosis) → m.apply → m.complete
```

For obvious bugs, users can skip m.debug and go directly to m.plan.

## m.debug Responsibilities

1. **Symptom Gathering** — Collect error messages, logs, reproduction steps
2. **Hypothesis Generation** — Rank likely causes with confidence levels
3. **Investigation Plan** — Propose what to check (files, logs, tests)
4. **Root Cause Confirmation** — Verify hypothesis before planning fix
5. **Output** — `.mission/diagnosis.md` with structured findings

## diagnosis.md Structure

```markdown
---
id: DIAG-20260121-230342
status: confirmed | investigating | inconclusive
confidence: high | medium | low
created: 2026-01-21T23:03:42
---

## SYMPTOM
[User-reported issue]

## INVESTIGATION
- [x] Checked auth.go:45 - nil pointer on empty session
- [ ] Reviewed recent commits to auth package

## HYPOTHESES
1. **[HIGH]** Session not initialized before access
2. **[LOW]** Race condition in concurrent requests

## ROOT CAUSE
Session middleware skipped for /api/login endpoint

## AFFECTED FILES
- internal/auth/session.go
- internal/middleware/auth.go

## RECOMMENDED FIX
Add session initialization to login handler

## REPRODUCTION
`curl -X POST localhost:8080/api/login -d '{}'` → 500
```

## m.plan Integration

When `.mission/diagnosis.md` exists with `status: confirmed`:

1. **Skip intent clarification** — diagnosis already gathered context
2. **Pre-populate SCOPE** — from `AFFECTED FILES` section
3. **Seed PLAN** — use `RECOMMENDED FIX` as starting point
4. **Link diagnosis** — add `diagnosis: DIAG-XXXXXX` to mission.md frontmatter
5. **Archive together** — diagnosis.md moves to `completed/` with mission.md

## Status Lifecycle

```
investigating → confirmed → consumed (by m.plan)
            ↘ inconclusive (needs more info or escalation)
```

## Design Decisions

### Why a Separate Command?

| Approach | Rejected Because |
|----------|------------------|
| Add `--mode debug` to m.plan | m.plan already ~150 lines with multiple branches |
| Keyword detection in m.plan | Fragile — "fix typo" ≠ "fix crash" |
| Inline investigation in m.apply | Mixes concerns, harder to audit |

### Why diagnosis.md?

- **Audit trail** — Documents investigation for future reference
- **Handoff point** — Clean interface between debug and plan phases
- **Reusable** — Useful even if user decides not to fix immediately
- **Archivable** — Preserved with completed missions

## Implementation Notes

1. **Auto-detection** — m.plan checks for diagnosis.md automatically, no flag needed
2. **Staleness warning** — Warn if diagnosis.md >24h old when consumed
3. **Iteration support** — m.debug can update diagnosis.md as investigation progresses
4. **CLI commands needed**:
   - `m diagnosis create --symptom "..."`
   - `m diagnosis update --status confirmed`
   - `m diagnosis check` (for m.plan to query)

## Trade-offs

| PROs | CONs |
|------|------|
| Clean separation of concerns | Two commands for bug fixes |
| m.plan stays focused | Users must learn when to use which |
| Creates investigation audit trail | Extra file to manage |
| Optional — can skip for obvious bugs | Diagnosis can go stale |
| Structured handoff to planning | |
