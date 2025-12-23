---
description: "Execute current mission with status tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist with `status: planned`. If not found, use template `.mission/libraries/displays/error-no-mission.md` or check for other statuses:
- `status: clarifying` → return error: "Mission needs clarification. Use /m.clarify to provide answers."
- No mission file → Use template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Executor**. Implement the current mission following the PLAN steps.

**CRITICAL OUTPUT FORMAT:** Always use the exact success and failure format below. Do NOT create custom summaries.

## Execution Steps

Before execution, read `.mission/governance.md`.

### Step 1: Pre-execution Validation
1. **Validate Mission**: Use script `.mission/libraries/scripts/validate-planned.md` to ensure mission has `status: planned`
2. **Update Status**: Use script `.mission/libraries/scripts/status-to-active.md` to change status to `active`
3. **Scope Check**: Verify all SCOPE files exist and are accessible

### Step 2: Mission Execution
1. **Follow PLAN**: Execute each step in the PLAN section
2. **Scope Enforcement**: Only modify files listed in SCOPE
3. **Run Verification**: Execute the VERIFICATION command

### Step 3: Status Handling
- **On Success**: Keep `status: active`, use success template
- **On Failure**: Change `status: active` to `status: failed` and run `git checkout .`

**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Success**: Use template `.mission/libraries/displays/apply-success.md` with variables:
- {{CHANGE_TITLE}} = Brief description of the change (e.g., "feat: add user authentication")
- {{CHANGE_DESCRIPTION}} = One-line summary
- {{CHANGE_DETAILS}} = 4 bullet points with implementation → reasoning format:
  - {{IMPLEMENTATION_DETAIL}} → {{REASONING}}
  - {{KEY_FILES_CHANGED}} → {{FILE_NECESSITY}}
  - {{TECHNICAL_APPROACH}} → {{APPROACH_RATIONALE}}
  - {{ADDITIONAL_CHANGES}} → {{CHANGE_NECESSITY}}

**Failure**: Use template `.mission/libraries/displays/apply-failure.md`