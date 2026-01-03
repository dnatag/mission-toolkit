---
description: "Complete current mission and update project tracking"
---

## Prerequisites

**CRITICAL:** Run `m mission check --command m.complete` to validate mission state before completion.

1. **Execute Check**: Run `m mission check --command m.complete` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED with m.complete execution" → Continue with completion
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Completor**. Finalize the current mission and update project tracking for continuous improvement.

## Execution Steps

Before generating output, read `.mission/governance.md`.

**MUST LOG:** Use file read tool to check if `.mission/execution.log` exists. If file doesn't exist, use file read tool to load template `libraries/scripts/init-execution-log.md`, then use file write tool to create the log file.

### Step 1: Mission Status Check
1. **Parse CLI Output**: Extract `mission_status` from `m mission check --command m.complete` JSON output (already validated in Prerequisites)
2. **Route by Status**: 
   - `mission_status: "active"` or `"completed"` → Success completion workflow
   - `mission_status: "failed"` → Failure completion workflow
   - Other status → Error (use error template)

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 1: Mission Status Check | [status found, workflow selected]"

### Step 2A: Success Completion Workflow
**For `status: active` missions:**

1. **Extract Mission Data**: Read mission.md to get id, track, type, SCOPE files, PLAN steps
2. **Calculate Metrics**: 
   - DURATION = Estimate based on mission complexity (Track 2: ~30min, Track 3: ~60min)
   - FILE_COUNT = Count files actually modified during execution
   - COMPLETED_STEPS = Count completed PLAN items
   - TOTAL_STEPS = Count total PLAN items
   - QUALITY_SCORE = (COMPLETED_STEPS / TOTAL_STEPS) * 100
   - VERIFICATION_STATUS = Check if VERIFICATION command was run successfully
3. **Update Mission**: Change `status: active` to `status: completed` and add `completed_at: YYYY-MM-DDTHH:MM:SS.sssZ`
4. **Create Metrics**: Use template `.mission/libraries/metrics/completion.md` with calculated variables
5. **Archive Mission**: Use script `.mission/libraries/scripts/archive-completed.md`
6. **Clean Up**: Remove `.mission/mission.md` after archiving

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 2A: Success Completion | [validation result, archived to location]"

### Step 2B: Failure Completion Workflow
**For `status: failed` missions:**

1. **Extract Mission Data**: Read mission.md to get id, track, type, failure_reason
2. **Calculate Metrics**:
   - DURATION = Estimate time spent before failure
   - FILE_COUNT = Count files in SCOPE (attempted files)
   - COMPLETED_STEPS = Count completed PLAN items before failure
   - TOTAL_STEPS = Count total PLAN items
   - FAILURE_REASON = Extract from mission.md or execution.log
3. **Update Mission**: Change `status: failed` to `status: completed` and add `completed_at: YYYY-MM-DDTHH:MM:SS.sssZ` and `failure_reason: [reason]`
4. **Create Metrics**: Use template `.mission/libraries/metrics/completion.md` with calculated variables
5. **Archive Mission**: Use script `.mission/libraries/scripts/archive-completed.md`
6. **Clean Up**: Remove `.mission/mission.md` after archiving

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 2B: Failure Completion | [failure reason, archived to location]"

### Step 3: Project Tracking Updates
**For both success and failure:**

1. **Update Backlog**: Check `.mission/backlog.md` for matching intent and mark as completed if found

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 3A: Update Backlog | [backlog updated or no matching intent found]"

2. **Refresh Metrics**: Use script `.mission/libraries/scripts/refresh-metrics.md` to update `.mission/metrics.md`

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.complete 3B: Refresh Metrics | [metrics updated with aggregate data]"

### Step 4: Final Cleanup
**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS] | m.complete 4: Final Cleanup | All steps completed, archiving execution log"

1. **Archive Execution Log**: Copy `.mission/execution.log` to `.mission/completed/{{MISSION_ID}}-execution.log`
2. **Clean Up**: Remove `.mission/execution.log` after archiving

**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Success Completion**: Use template `.mission/libraries/displays/complete-success.md` with variables:
- {{MISSION_ID}} = From mission.md id field
- {{DURATION}} = Estimated time (e.g., "45 minutes")
- {{FILE_COUNT}} = Number of files modified
- {{TRACK}} = Mission track
- {{MISSION_TYPE}} = WET/DRY
- {{VERIFICATION_STATUS}} = PASSED/FAILED/SKIPPED
- {{COMPLETED_STEPS}} = Number of completed steps
- {{TOTAL_STEPS}} = Total number of steps
- {{QUALITY_SCORE}} = Calculated quality percentage

**Failure Completion**: Use template `.mission/libraries/displays/complete-failure.md` with variables:
- {{MISSION_ID}} = From mission.md id field
- {{DURATION}} = Estimated time (e.g., "45 minutes")
- {{FAILURE_REASON}} = Reason for failure
- {{TRACK}} = Mission track
- {{MISSION_TYPE}} = WET/DRY
- {{COMPLETED_STEPS}} = Number of completed steps
- {{TOTAL_STEPS}} = Total number of steps
- {{FILE_COUNT}} = Number of files in scope

**Metrics Template**: Use `.mission/libraries/metrics/completion.md` with variables:
- {{MISSION_ID}} = From mission.md id field
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