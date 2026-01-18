---
name: "m.apply"
description: "Execute current mission with two-pass implementation and polish"
---

## Prerequisites

**CRITICAL:** Run `m mission check --context apply` to validate mission state before execution.

1. **Execute Check**: Run `m mission check --context apply` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED with m.apply execution" → Continue with execution
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use file read tool to load template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Executor**. Implement the current mission using a two-pass approach:
1. **First Pass (Implementation)**: Implement functionality following PLAN steps
2. **Second Pass (Polish)**: Refine code quality with automatic rollback on failure

**Required Output Format:** Use the exact success and failure format below. Do not create custom summaries.

## Execution Steps

### Step 0: Load Governance (Required)

**CRITICAL:** Use file read tool to read `.mission/governance.md` to prevent data corruption.

### Step 1: Update Status & Create Checkpoint
1. **Update Status**: Execute `m mission update --status active`
2. **Create Initial Checkpoint**: Execute `m checkpoint create` to save clean state
   - Returns checkpoint name (e.g., `MISS-20260103-143022-0`)
   - **On Checkpoint Creation Failure**:
     - **If error contains "unstaged changes" or "working directory not clean"**:
       - Run `m log --step "Update Status" "Checkpoint creation failed due to unstaged changes. Using template."`
       - Execute `m mission update --status failed`
       - Use file read tool to load template `.mission/libraries/displays/checkpoint-failure-unstaged.md` with variables:
         - {{UNSTAGED_FILES}} = Extract file names from git error message (or "Multiple files" if extraction fails)
         - {{MISSION_ID}} = Current mission ID
       - Halt execution
     - **For other checkpoint errors**:
       - Run `m log --step "Update Status" "Checkpoint creation failed: <error>. Aborting mission."`
       - Execute `m mission update --status failed`
       - Display error and halt
3. **Log**: Run `m log --step "Update Status" "Status active, checkpoint created"`

### Step 2: First Pass (Implementation)
1. **Verify SCOPE Files**: Check all files listed in SCOPE section exist before modifying
2. **Follow PLAN with Step Tracking**: Execute each step in the PLAN section:
   - **After each PLAN step**: Run `m mission mark-complete --step [N] --status SUCCESS --message "Completed: [step description]"`
   - **On step failure**: Run `m log --step "Plan Step [N]" "Failed: [step description] - [error details]"`
3. **Scope Enforcement**: Only modify files listed in SCOPE
4. **Run Verification**: Execute the VERIFICATION command
5. **On Verification Failure**: Attempt to fix issues and re-run verification (iterate until passing or unable to fix)
6. **If Unable to Fix**: Proceed to Step 4 (Status Handling) with failure
7. **Log**: Run `m log --step "First Pass" "[files modified, verification result, fix attempts if any]"`

### Step 3: Second Pass (Polish) - Runs after Step 2 succeeds

**Important:** This step runs after Step 2 succeeds. Polish improves code quality with automatic rollback protection.

1. **Create Polish Checkpoint**: Execute `m checkpoint create` to save first pass state
   - Returns checkpoint name (e.g., `MISS-20260103-143022-1`)
   - **On Checkpoint Creation Failure**:
     - Run `m log --step "Polish Pass" "Checkpoint creation failed: <error>. Skipping polish."`
     - Skip polish pass entirely
     - Continue to Step 4 with first pass code
     - Add footer to commit message: `Polish-Skipped: checkpoint-creation-failed`

2. **Enhanced Polish Review**: Analyze all modified code from Step 2 and apply comprehensive quality improvements:
   - **Code Structure**: Idiomatic patterns, language conventions
   - **Readability**: Clear variable names, logical flow, reduced complexity
   - **Robustness**: Error handling, edge cases, input validation, defensive programming
   - **Performance**: Algorithmic efficiency, resource usage, caching opportunities
   - **Maintainability**: Documentation, comments for complex logic, consistent formatting
   - **Testing**: Robust coverage for new logic, edge cases, error conditions, integration points
   - **Security**: Input sanitization, secure defaults, vulnerability prevention
   - **Standards Compliance**: Project conventions, linting rules, best practices

3. **Re-run Verification**: Execute the VERIFICATION command again

4. **Handle Polish Verification**:
   - **On Success**: 
     - Execute `m checkpoint create` to save polished state
     - Returns checkpoint name (e.g., `MISS-20260103-143022-2`)
     - Run `m log --step "Polish Pass" "Polish applied successfully, verification passed"`
     - Continue to Step 4
   - **On Failure**: 
     - Execute `m checkpoint restore <checkpoint-name>` to rollback polish changes
     - **On Restore Success**:
       - Run `m log --step "Polish Pass" "Polish verification failed, rolled back to first pass"`
       - Continue to Step 4 with first pass code
     - **On Restore Failure**:
       - Run `m log --step "Polish Pass" "Restore failed but first pass code exists"`
       - Display warning about manual cleanup
       - Continue to Step 4 with current state (first pass code should still be present)

### Step 4: Status Handling

**On Any Failure (Step 1 or Step 2)**:
1. Execute `m checkpoint restore <initial-checkpoint> --all` to revert all changes (if checkpoints exist)
2. Execute `m mission update --status failed`
3. Run `m log --step "Status Handling" "Mission failed, all changes reverted"`
4. **Analyze Failure**: Determine failure type and provide guidance:
   - **Step 1 Checkpoint Creation Failure**: Environment issue, retry unlikely to help
   - **Step 2 Verification Failure**: Implementation issue, review error and retry
5. Use file read tool to load template `.mission/libraries/displays/apply-failure.md` with variables:
   - {{FAILURE_REASON}} = Brief summary of what failed (e.g., "Verification failed: test errors", "Checkpoint creation failed")
   - {{RETRY_ADVICE}} = "Retry /m.apply" or "Fix environment first" or "Review mission plan"

**On Success (Step 2 passed, Step 3 completed or skipped)**:
1. Execute `m mission update --status executed`
2. Run `m log --step "Status Handling" "Mission execution complete"``
3. Use file read tool to load template `.mission/libraries/displays/apply-success.md` with variables:
   - {{CHANGE_DETAILS}} = 4 bullet points with implementation → reasoning format:
     - {{IMPLEMENTATION_DETAIL}} → {{REASONING}}
     - {{KEY_FILES_CHANGED}} → {{FILE_NECESSITY}}
     - {{TECHNICAL_APPROACH}} → {{APPROACH_RATIONALE}}
     - {{ADDITIONAL_CHANGES}} → {{CHANGE_NECESSITY}}
   - {{CHECKPOINT_0}} = Initial checkpoint name from Step 1 (e.g., "MISS-20260103-143022-0")
   - {{CHECKPOINT_1}} = Polish checkpoint name from Step 3.1 (e.g., "MISS-20260103-143022-1", or "N/A" if polish skipped)
   - {{CHECKPOINT_2}} = Final checkpoint name from Step 3.4 (e.g., "MISS-20260103-143022-2", or "N/A" if polish failed/skipped)

## Error Handling

### Polish Checkpoint Restore Failure
- Log warning: `m log --step "Polish Pass" "Restore failed but first pass code exists"`
- Display warning to user:
  ```
  Warning: Polish rollback failed, but first pass code should still be present.
  If you see unexpected changes, manually restore:
  - m checkpoint restore {{CHECKPOINT_1}}
  Review .mission/execution.log for details.
  ```
- Continue with mission completion (first pass code passed verification)

### Critical Checkpoint Restore Failure (Step 1 or Step 2)
- Mark mission as failed: `m mission update --status failed`
- Display manual recovery steps:
  ```
  Checkpoint restore failed. Manual recovery required:
  1. Run: git reset --soft HEAD
  2. Review staged changes and unstage/discard as needed
  3. Review .mission/execution.log for details
  4. Re-run: /m.apply
  ```

### Verification Command Crash
- Treat as verification failure
- Run `m log --step "Verification" "Command crashed: exit_code=<code> stderr=<output>"`
- Follow normal failure path (Step 2 crash → mission failed, Step 3 crash → rollback polish)

**CRITICAL**: Use file read tool to load templates from `.mission/libraries/` for consistent output.