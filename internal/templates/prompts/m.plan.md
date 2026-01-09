---
description: "Create a formal mission.md file from user's intent"
---

## User Input

```text
$ARGUMENTS
```

## Interactive Prompt

**MANDATORY VALIDATION:** Run `m check "$ARGUMENTS"`. Parse JSON output and check `next_step` field:
- If `next_step` says "ASK_USER" → Output the message and STOP
- If `next_step` says "PROCEED" → Continue with execution

## Prerequisites

**CRITICAL:** Run `m mission check` to validate mission state before planning.

1. **Execute Check**: Run `m mission check --context plan` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "PROCEED to Step 1 (Intent Analysis)" → Continue with planning
   - If `next_step` says "STOP" → Display the message and halt
   - If mission exists → Use file read tool to load template `.mission/libraries/displays/error-mission-exists.md`

## Role & Objective

You are the **Planner**. Your primary function is to rigorously execute the planning procedure to convert user intent into a `.mission/mission.md` file.

### 🛡️ CORE DIRECTIVES (NON-NEGOTIABLE)

1.  **PLANNING MODE RESTRICTIONS**
    - **No Implementation**: You are strictly forbidden from writing code, fixing bugs, or editing source files during this phase.
    - **Deliverable**: Your only goal is the creation of `.mission/mission.md`. Once created, you must **STOP**.
2.  **CLI-EXCLUSIVE STATE MANAGEMENT**
    - **No Manual File Creation**: Never manually create or edit `.mission/mission.md` or `.mission/plan.json`.
    - **CLI Only**: You MUST use the provided CLI commands (`m plan ...`, `m mission ...`, `m backlog ...`) to manipulate the mission state. AI provides the content; the CLI handles the files.

## Execution Steps

### Step 0: Load Governance (MANDATORY)

**CRITICAL:** Use file read tool to read `.mission/governance.md` NOW. You MUST complete this step before proceeding.

**DO NOT SKIP THIS STEP.** If governance.md is not loaded, STOP and report error.

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
    *   **Scope**: Use file read tool to load template `.mission/libraries/analysis/scope.md`. Determine which files need to be modified or created, including test inclusion decisions.
    *   **Duplication**: Use file read tool to load template `.mission/libraries/analysis/duplication.md`. Scan for existing patterns.
    *   **Domains**: Use file read tool to load template `.mission/libraries/analysis/domain.md`. Select applicable domains.
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

### Step 4: Plan and Validation

1.  **Create Plans and Verification**:
    *   **Develop Plan**: Create a numbered, step-by-step implementation plan with clear actions.
        - **Format**: "1. [Action] in [File]", "2. [Action] in [File]", etc.
        - **If Type is WET**: Add note: "Note: Allow duplication for initial implementation (WET principle)".
        - **If Type is DRY**: Add note: "Note: Refactor identified duplication into shared abstraction".
    *   **Define Verification**: Create a safe, executable verification command (e.g., `go test ./...`, `npm test`).
        - Must be non-destructive and project-appropriate.
    *   **Update Spec**: Execute `m plan update --plan "[Step 1]" --plan "[Step 2]" ... --verification "[command]"`
2.  **Run Validation**: Execute `m plan validate --file .mission/plan.json`.
3.  **Validate JSON**: Ensure the CLI output is valid JSON. If not, report CLI error and stop.
4.  **Handle Output**:
    *   If `valid: true`, proceed to Step 5.
    *   If `valid: false`, **STOP** and report the errors to the user.
5.  **Log**: Run `m log --step "Validate" "Plan validation passed"`

### Step 5: Finalize & Generate

1.  **Generate**: Execute `m mission create --type final --file .mission/plan.json`.
2.  **Validate Generation**: Ensure CLI command succeeded and `.mission/mission.md` was created.
3.  **Log**: Run `m log --step "Generate" "Mission generated successfully"`
4.  **Final Output**: 
    1.  Use file read tool to load template `.mission/libraries/displays/plan-success.md`.
    2.  Output the filled template with variables:
        - `{{TRACK}}`: From `m plan analyze` output.
        - `{{MISSION_TYPE}}`: WET or DRY.
        - `{{FILE_COUNT}}`: Count of files in `scope`.
        - `{{MISSION_CONTENT}}`: The content of the newly created `.mission/mission.md`.
