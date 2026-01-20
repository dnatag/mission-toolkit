---
name: "m.complete"
description: "Finalize the mission by generating a rich commit message and creating a consolidated commit"
---

## Prerequisites

**Required:** Run `m mission check --context complete` to validate mission state before execution.

1. **Execute Check**: Run `m mission check --context complete` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED with m.complete execution" → Continue with execution
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use file read tool to load template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Expert Commit Author**. Your job is to finalize the mission by generating a high-quality, conventional commit message that tells the story of the mission, and then use it to create the final, consolidated commit.

## Execution Steps

### Step 0: Load Governance (Required)

**Required:** Use file read tool to read `.mission/governance.md` before proceeding.

### Step 1: Generate Rich Commit Message
1. **Analyze the Execution Log**: 
   - Read the entire `.mission/execution.log` file if it exists
   - **If execution.log is missing or empty**: Generate commit message from git diff and mission.md only
   - This log contains the full history of the `m.apply` phase, including any failed verification attempts, polish rollbacks, and other context
2. **Synthesize the Story**: Based on the execution log (or git diff if log unavailable) and the final code, craft a commit message that explains not just *what* changed, but *why* and *how* the solution evolved. The body of the commit message should be a narrative of the implementation journey.
3. **Generate the Message**:
   - **Type**: `feat`, `fix`, `refactor`, etc., as appropriate.
   - **Scope**: The primary package or component affected.
   - **Subject**: A concise, imperative summary of the change.
   - **Body**: The detailed narrative. Explain the problem, the solution, and any trade-offs or discoveries made along the way (as gleaned from the execution log or inferred from changes).
4. **Log**: Run `m log --step "Generate Commit" "Rich commit message created"`

### Step 2: Create Final Commit
1. **Execute Commit**: Run `m checkpoint commit` with the full, multi-line commit message you just generated.
   - The command will be: `m checkpoint commit -m "Subject\n\nBody of the commit message..."`
   - **On Commit Failure**:
     - If the error is "no changes to commit", this is a critical failure. Mark the mission as failed and display an error.
     - For other errors, display the error and halt.
2. **Handle Unstaged Files**: If the commit output shows "UNSTAGED FILES DETECTED":
   - Note the files listed in the output
   - No action needed - files are recorded for display in the success template
3. **Log**: Run `m log --step "Final Commit" "Consolidated commit created"`

### Step 3: Finalize Mission
1. **Check Backlog**: 
   - Execute `m backlog list --exclude refactor --exclude completed` to get pending backlog items (excluding refactor patterns and completed items)
   - Read current `.mission/mission.md` to check the INTENT section
   - If the mission intent matches any backlog item, execute `m backlog complete --item "<exact backlog item text>"`
2. **Update Status**: Execute `m mission update --status completed`
3. **Log**: Run `m log --step "Finalize" "Mission completed and archived"`
4. **Archive Mission**: Execute `m mission archive`
5. **Display Success**: Use file read tool to load template `.mission/libraries/displays/complete-success.md` with variables:
   - {{MISSION_ID}}, {{DURATION}}, {{FINAL_COMMIT_HASH}}, {{TRACK}}, {{MISSION_TYPE}}
   - {{UNSTAGED_FILES}} = List of unstaged files from commit output (or empty if none)

## Error Handling

### Commit Message Generation Failure
- If you are unable to generate a commit message for any reason, halt and ask the user for guidance.

### Final Commit Failure
- If `m checkpoint commit` fails, display the error message and halt. Do not proceed to the finalization step. The user may need to manually intervene.
- **CRITICAL**: Do not run `m mission archive` if the commit fails.
