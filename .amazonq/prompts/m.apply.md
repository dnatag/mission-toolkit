---
description: "Execute current mission with status tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist with `status: planned`. If not found, use template `.mission/libraries/displays/error-no-mission.md` or check for other statuses:
- `status: clarifying` → return error: "Mission needs clarification. Use @m.clarify to provide answers."
- No mission file → Use template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Executor**. Implement the current mission following the PLAN steps.

**CRITICAL OUTPUT FORMAT:** Always use the exact success and failure format below. Do NOT create custom summaries.

## Execution Steps

Before execution, read `.mission/governance.md`.

**MUST LOG:** Use file read tool to check if `.mission/execution.log` exists. If file doesn't exist, use file read tool to load template `libraries/scripts/init-execution-log.md`, then use file write tool to create the log file.

### Step 1: Pre-execution Validation
1. **Validate Mission**: Execute script `.mission/libraries/scripts/validate-planned.sh`
2. **Update Status**: Execute script `.mission/libraries/scripts/status-to-active.sh`
3. **Scope Check**: Verify all SCOPE files exist and are accessible

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.apply 1: Pre-execution Validation | [brief outcome]"

### Step 2: Mission Execution
1. **Follow PLAN**: Execute each step in the PLAN section
2. **Scope Enforcement**: Only modify files listed in SCOPE
3. **Run Verification**: Execute the VERIFICATION command

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.apply 2: Mission Execution | [files modified, verification result]"

### Step 3: Status Handling
- **On Success**: Keep `status: active`, use success template
- **On Failure**: Change `status: active` to `status: failed` and run `git checkout .`

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.apply 3: Status Handling | [final outcome]"

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