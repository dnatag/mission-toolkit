---
description: "Finalize the mission by generating a commit message, metrics, and creating a consolidated commit"
---

## Prerequisites

**CRITICAL:** Run `m mission check --context complete` to validate mission state before execution.

1. **Execute Check**: Run `m mission check --context complete` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED with m.complete execution" → Continue with execution
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Completer**. Your job is to finalize the mission by generating a conventional commit message, creating a single consolidated commit, generating mission metrics, and cleaning up the mission artifacts.

## Execution Steps

### Step 1: Generate Commit Message
1. **Load Template**: Use file read tool to load template `libraries/scripts/generate-commit-message.md`
2. **Read Mission Context**: Extract MISSION_ID, TYPE, TRACK, INTENT, SCOPE from `mission.md`.
3. **Analyze Code Changes**: Review the final state of the code in the SCOPE files to understand the changes.
4. **Generate Message**: Follow the template rules to populate all variables, creating a comprehensive and accurate commit message.
5. **Log**: Run `m log --step "Generate Commit" "[commit type and scope]"`

### Step 2: Create Final Commit
1. **Execute Commit**: Run `m checkpoint commit` with the generated commit message from the previous step.
   - The command will be: `m checkpoint commit -m "your generated commit message"`
   - **On Commit Failure**:
     - If the error is "no changes to commit", this is a critical failure. The mission should have produced changes. Mark the mission as failed and display an error.
     - For other errors, display the error and halt.
2. **Log**: Run `m log --step "Final Commit" "Consolidated commit created"`

### Step 3: Generate Mission Metrics
1. **Load Template**: Use file read tool to load template `internal/templates/mission/metrics.md`.
2. **Gather Data**:
   - **Final Commit Hash**: Use the hash from the previous step.
   - **Checkpoints**: Count the number of checkpoints created during the mission (e.g., by listing tags).
   - **Duration**: This will be calculated by the CLI, but you can note the start/end times if available.
3. **AI Reflection**: Analyze the mission's execution. Consider the initial plan, the number of checkpoints (as a proxy for rework), and the final solution.
4. **Generate Content**: Populate the `metrics.md` template with the gathered data and your reflection in the "What I Learned" section.
5. **Write to File**: Write the generated content to `.mission/metrics.md`.
6. **Log**: Run `m log --step "Generate Metrics" "Metrics file created"`

### Step 4: Finalize Mission
1. **Update Status**: Execute `m mission update --status completed`
2. **Archive Mission**: Execute `m mission archive`
3. **Log**: Run `m log --step "Finalize" "Mission completed and archived"`
4. **Display Success**: Use template `.mission/libraries/displays/complete-success.md` with the final commit hash.

## Error Handling

### Commit Message Generation Failure
- If you are unable to generate a commit message for any reason, halt and ask the user for guidance.

### Final Commit Failure
- If `m checkpoint commit` fails, display the error message and halt. Do not proceed to the finalization step. The user may need to manually intervene.
- **CRITICAL**: Do not run `m mission archive` if the commit fails.
