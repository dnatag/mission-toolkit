---
description: "Execute current mission with status tracking"
---

## Prerequisites

**CRITICAL:** Run `m mission check --command m.apply` to validate mission state before execution.

1. **Execute Check**: Run `m mission check --command m.apply` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED with m.apply execution" → Continue with execution
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use template `.mission/libraries/displays/error-no-mission.md`

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

### Step 3: Generate Commit Message
1. **Read Mission Context**: Extract MISSION_ID, TYPE, TRACK, INTENT, SCOPE from mission.md
2. **Generate Conventional Commit**:
   - **Type**: WET → `feat`, DRY → `refactor` (override if clearly fix/docs/test/chore)
   - **Scope**: Extract from primary file/module in SCOPE (e.g., `auth`, `api`, `cli`)
   - **Title**: Imperative mood, max 72 chars, capitalize first letter, no period
   - **Description**: Explain what and why (not how), wrap at 72 chars
   - **Footer**: Add `Mission-ID: {{MISSION_ID}}`, `Track: {{TRACK}}`, `Type: {{TYPE}}`
3. **Update mission.md**: Replace `## COMMIT_MESSAGE` section with generated message

**CRITICAL - When to Regenerate**:
- After initial implementation (always)
- After ANY code changes (polish pass, bug fixes, user-requested changes)
- After verification fixes that modify code
- Description must reflect ALL changes made, not just latest

**Format Example**:
```
feat(auth): Add JWT authentication middleware

Implement token-based authentication for API endpoints using JWT.
Includes middleware for token validation and user context injection.

Mission-ID: M-20250115-001
Track: 2
Type: WET
```

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS] | m.apply 3: Generate Commit Message | [commit type and scope]"

### Step 4: Status Handling
- **On Success**: Keep `status: active`, use success template
- **On Failure**: Change `status: active` to `status: failed` and run `git checkout .`

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.apply 4: Status Handling | [final outcome]"

**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Success**: Use template `.mission/libraries/displays/apply-success.md` with variables:
- {{CHANGE_TITLE}} = The commit message title (e.g., "feat(auth): Add JWT authentication")
- {{CHANGE_DESCRIPTION}} = The commit message description
- {{CHANGE_DETAILS}} = 4 bullet points with implementation → reasoning format:
  - {{IMPLEMENTATION_DETAIL}} → {{REASONING}}
  - {{KEY_FILES_CHANGED}} → {{FILE_NECESSITY}}
  - {{TECHNICAL_APPROACH}} → {{APPROACH_RATIONALE}}
  - {{ADDITIONAL_CHANGES}} → {{CHANGE_NECESSITY}}

**Failure**: Use template `.mission/libraries/displays/apply-failure.md`