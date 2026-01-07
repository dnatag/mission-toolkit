---
description: "Create a formal mission.md file from user's intent"
---

## User Input

```text
$ARGUMENTS
```

## Interactive Prompt

**CRITICAL:** Always check if `$ARGUMENTS` is empty or contains only whitespace first.

**MANDATORY VALIDATION:**
```
if ($ARGUMENTS is empty OR contains only whitespace OR equals "$ARGUMENTS"):
    ASK USER IMMEDIATELY
    OUTPUT ONLY: "What is your intent or goal for this task?"
    DO NOT PROCEED WITH ANY OTHER STEPS
    WAIT FOR USER RESPONSE
else:
    Continue with execution steps
```

## Prerequisites

**CRITICAL:** Run `m mission check` to validate mission state before planning.

1. **Execute Check**: Run `m mission check --context plan` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED to Step 2 (Intent Analysis)" → Continue with planning
   - If `next_step` says "STOP" → Display the message and halt
   - If mission exists → Use file read tool to load template `.mission/libraries/displays/error-mission-exists.md`

## Role & Objective

You are the **Planner**. Your goal is to convert the user's intent into a formal `.mission/mission.md` file using the Mission Toolkit CLI. If the intent is ambiguous, you will guide the user through clarification.

**CRITICAL**: 
- Do NOT create `.mission/mission.md` or `.mission/plan.json` manually
- Do NOT edit JSON files directly
- You MUST use CLI commands: `m plan init`, `m plan update`, `m plan analyze`, `m plan validate`, `m mission create`, `m backlog`
- AI provides values, CLI handles all file operations

## Execution Steps

Before generating output, read `.mission/governance.md`.

### Step 1: Intent & Clarification

1.  **Analyze Intent**: Use file read tool to load template `.mission/libraries/analysis/intent.md`. Use it to refine the user's request.
    *   **Decision**: If the output is "AMBIGUOUS", **STOP IMMEDIATELY**. Ask the user to clarify the specific reason for the ambiguity.
2.  **Check Clarity**: Use file read tool to load template `.mission/libraries/analysis/clarification.md`. Run it to check for missing details.
    *   **If output is "✅ INTENT CLEAR"**: Set `[REFINED_INTENT]` and proceed to Step 2.
    *   **If output is "⚠️ PROCEEDING WITH ASSUMPTIONS"**: Display assumptions to user, set `[REFINED_INTENT]`, and proceed to Step 2.
    *   **If output is "🛑 CLARIFICATION NEEDED"**:
        1.  Display questions to user and **STOP**.
        2.  When user responds, combine original intent with answers to form `[REFINED_INTENT]` and restart from Step 2.
3.  **Log**: Run `m log --step "Intent" "Intent analyzed and refined"`

### Step 2: Context Analysis

1.  **Analyze Context**:
    *   **Duplication**: Use file read tool to load template `.mission/libraries/analysis/duplication.md`. Scan for existing patterns.
    *   **Domains**: Use file read tool to load template `.mission/libraries/analysis/domain.md`. Select applicable domains.
    *   **Scope**: Determine which files need to be modified or created based on `[REFINED_INTENT]`.
        *   **Include tests** if: New logic, bug fixes, or critical paths.
        *   **Exclude tests** if: Trivial changes.
2.  **Determine Strategy (WET vs DRY)**:
    *   **If Duplication Detected**:
        1.  Execute `m backlog list` to check if refactoring is already tracked (match by pattern description).
        2.  **If NOT in backlog**:
            - Execute `m backlog add "Refactor [Pattern] in [Files]" --type refactor`.
            - Set `[MISSION_TYPE]` to "WET" (Defer refactor).
        3.  **If ALREADY in backlog**:
            - Set `[MISSION_TYPE]` to "DRY" (Enforce refactor).
    *   **If NO Duplication**: Set `[MISSION_TYPE]` to "WET".
3.  Execute `m plan init --intent "[REFINED_INTENT]" --type [MISSION_TYPE] --scope [file1] --scope [file2] --domain [domain1] ...`
4.  **Log**: Run `m log --step "Context" "Context analyzed and draft plan initialized"`

### Step 3: Complexity Analysis

1.  **Run Analysis**: Execute `m plan analyze --file .mission/plan.json --update`.
2.  **Validate JSON**: Ensure the CLI output is valid JSON. If not, report CLI error and stop.
3.  **Follow Instructions**: Read the `next_step` field from the JSON output and **follow it literally**.
    *   **If `next_step` says STOP (Track 1)**:
        1.  Use file read tool to load template `.mission/libraries/displays/plan-atomic.md`.
        2.  Analyze the intent to generate a safe, specific edit suggestion.
        3.  Output the filled template with `{{REFINED_INTENT}}` and `{{SUGGESTED_EDIT}}`.
    *   **If `next_step` says STOP (Track 4)**:
        1.  Decompose `[REFINED_INTENT]` into atomic sub-intents.
        2.  Execute `m backlog list` to verify no duplicates.
        3.  Execute `m backlog add "[sub-intent 1]" "[sub-intent 2]" ... --type decomposed` (excluding duplicates).
        4.  Use file read tool to load template `.mission/libraries/displays/plan-epic.md`.
        5.  Output the filled template with `{{SUB_INTENTS}}` populated with decomposed intents.
    *   **If `next_step` says PROCEED**: Continue to Step 4.
4.  **Log**: Run `m log --step "Analyze" "Complexity analysis complete. Track: [TRACK]"`

### Step 4: Validation

1.  **Run Validation**: Execute `m plan validate --file .mission/plan.json`.
2.  **Validate JSON**: Ensure the CLI output is valid JSON. If not, report CLI error and stop.
3.  **Handle Output**:
    *   If `valid: true`, proceed to Step 5.
    *   If `valid: false`, **STOP** and report the errors to the user.
4.  **Log**: Run `m log --step "Validate" "Plan validation passed"`

### Step 5: Finalize & Generate

1.  **Develop Plan**: Create a numbered, step-by-step implementation plan with clear actions.
    *   **Format**: "1. [Action] in [File]", "2. [Action] in [File]", etc.
    *   **If Type is WET**: Add note: "Note: Allow duplication for initial implementation (WET principle)".
    *   **If Type is DRY**: Add note: "Note: Refactor identified duplication into shared abstraction".
2.  **Define Verification**: Create a safe, executable verification command (e.g., `go test ./...`, `npm test`).
    *   Must be non-destructive and project-appropriate.
3.  **Update Spec**: Execute `m plan update --plan "[Step 1]" --plan "[Step 2]" ... --verification "[command]"`
4.  **Generate**: Execute `m mission create --type final --file .mission/plan.json`.
5.  **Validate Generation**: Ensure CLI command succeeded and `.mission/mission.md` was created.
6.  **Log**: Run `m log --step "Generate" "Mission generated successfully"`
7.  **Final Output**: 
    1.  Use file read tool to load template `.mission/libraries/displays/plan-success.md`.
    2.  Output the filled template with variables:
        - `{{TRACK}}`: From `m plan analyze` output.
        - `{{MISSION_TYPE}}`: WET or DRY.
        - `{{FILE_COUNT}}`: Count of files in `scope`.
        - `{{MISSION_CONTENT}}`: The content of the newly created `.mission/mission.md`.
