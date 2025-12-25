---
description: "Complete current mission and update project tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist. If `.mission/mission.md` is not found, use template `.mission/libraries/displays/error-no-mission.md`.

## Role & Objective

You are the **Completor**. Finalize the current mission and update project tracking for continuous improvement.

## Execution Steps

Before generating output, read `.mission/governance.md`.

### Step 1: Mission Status Check
1. **Status Check**: Check mission status (`active`, `failed`, or other)
2. **Route by Status**: 
   - `status: active` → Success completion workflow
   - `status: failed` → Failure completion workflow
   - Other status → Error (use error template)

**Log Step 1**: Append to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 1: Mission Status Check | [status found, workflow selected]"

### Step 2A: Success Completion Workflow
**For `status: active` missions:**

1. **Validation**: Verify all PLAN items completed and VERIFICATION passed
2. **Update Mission**: Change `status: active` to `status: completed` and add `completed_at: YYYY-MM-DDTHH:MM:SS.sssZ`
3. **Archive Mission**: Use script `.mission/libraries/scripts/archive-completed.md` (includes execution log)
4. **Clean Up**: Remove `.mission/mission.md` after archiving

**Log Step 2A**: Append to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 2A: Success Completion | [validation result, archived to location]"

### Step 2B: Failure Completion Workflow
**For `status: failed` missions:**

1. **Extract Failure Info**: Determine failure reason and completed steps
2. **Update Mission**: Change `status: failed` to `status: completed` and add `completed_at: YYYY-MM-DDTHH:MM:SS.sssZ` and `failure_reason: [reason]`
3. **Archive Mission**: Use script `.mission/libraries/scripts/archive-completed.md`
4. **Clean Up**: Remove `.mission/mission.md` after archiving

**Log Step 2B**: Append to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 2B: Failure Completion | [failure reason, archived to location]"

### Step 3: Project Tracking Updates
**For both success and failure:**

1. **Update Summary**: Update `.mission/metrics.md` aggregate statistics and RECENT COMPLETIONS
2. **Update Backlog**: Check `.mission/backlog.md` for matching intent and mark as completed if found

**Log Step 3**: Append to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "SUCCESS | m.complete 3: Project Tracking | Updated metrics.md, marked backlog item completed" (or appropriate actual values)

### Step 4: Archive Execution Log
1. **Copy Log**: Copy `.mission/execution.log` to `.mission/completed/{{MISSION_ID}}-execution.log`
2. **Clean Up**: Remove `.mission/execution.log` after archiving

**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Success Completion**: Use template `.mission/libraries/displays/complete-success.md` with variables:
- {{MISSION_ID}} = Track-Type-Timestamp
- {{DURATION}} = Estimated time (e.g., "45 minutes")
- {{FILE_COUNT}} = Number of files modified
- {{TRACK}} = Mission track
- {{MISSION_TYPE}} = WET/DRY
- {{VERIFICATION_STATUS}} = PASSED/FAILED/SKIPPED
- {{COMPLETED_STEPS}} = Number of completed steps
- {{TOTAL_STEPS}} = Total number of steps
- {{QUALITY_SCORE}} = Calculated quality percentage
- {{TIMESTAMP}} = Archive timestamp

**Failure Completion**: Use template `.mission/libraries/displays/complete-failure.md` with variables:
- {{MISSION_ID}} = Track-Type-Timestamp
- {{DURATION}} = Estimated time (e.g., "45 minutes")
- {{FAILURE_REASON}} = Reason for failure
- {{TRACK}} = Mission track
- {{MISSION_TYPE}} = WET/DRY
- {{COMPLETED_STEPS}} = Number of completed steps
- {{TOTAL_STEPS}} = Total number of steps
- {{FILE_COUNT}} = Number of files in scope
- {{TIMESTAMP}} = Archive timestamp

**Metrics Template**: Use `.mission/libraries/metrics/completion.md` with variables:
- {{MISSION_ID}} = Track-Type-Timestamp
- {{COMPLETION_DATE}} = YYYY-MM-DD HH:MM:SS
- {{DURATION_MINUTES}} = Numeric duration
- {{FILES_MODIFIED}} = Actual file count
- {{LINES_ADDED}} = Lines of code added
- {{LINES_REMOVED}} = Lines of code removed
- {{DUPLICATION_FOUND}} = Yes/No
- {{SECURITY_ISSUES}} = None/List of issues
- {{PERFORMANCE_IMPACT}} = Minimal/Moderate/High
- {{PATTERNS_FOUND}} = List of identified patterns
- {{REFACTORING_OPPORTUNITIES}} = List of opportunities
- {{NEXT_MISSIONS}} = Suggested follow-ups
- {{FAILURE_REASON}} = Reason for failure (for failed missions)
- {{SUCCESS_STATUS}} = "SUCCESS" or "FAILED"