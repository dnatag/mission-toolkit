---
description: "Create a formal mission.md file from user's intent"
---

## User Input

```text
$ARGUMENTS
```

## Interactive Prompt

**Required:** Run `m check "$ARGUMENTS"`. Parse JSON output and check `next_step` field:
- If `next_step` says "ASK_USER" → Output the message and STOP
- If `next_step` says "PROCEED" → Continue with execution

## Prerequisites

**Required:** Run `m mission check` to validate mission state before planning.

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
    - AI reads `.mission/` files; `m` CLI writes them
    - Use `m` commands for all mission state modifications
    - Parse JSON output and follow conditional logic
3.  **TEMPLATE INTERPRETATION**
    - `m analyze *` commands output JSON with `template_path` in `.mission/templates/`
    - Parse JSON, use file read tool to load the template file, follow its instructions
    - Never display template content to user - it's for your analysis only
    - `m mission *` and `m backlog *` commands output JSON for programmatic parsing

## Execution Steps

### Step 0: Setup & Load Governance (Required)

1. **Setup Templates**: Execute `mkdir -p .mission/templates` to create template directory
2. **Load Governance**: Use file read tool to read `.mission/governance.md` before proceeding.

If governance.md is not loaded, stop and report error.

### Step 1: Intent & Clarification

1.  **Analyze Intent**: `m analyze intent "$ARGUMENTS"` → Parse JSON, read `template_path`, follow template
    *   If analysis concludes "AMBIGUOUS" → **STOP**, ask user to clarify
2.  **Create Mission**: `m mission create --intent "[REFINED_INTENT]"`
3.  **Check Clarity**: `m analyze clarify` → Parse JSON, read `template_path`, follow template
    *   "✅ INTENT CLEAR" → Proceed to Step 2
    *   "⚠️ PROCEEDING WITH ASSUMPTIONS" → Display assumptions, proceed to Step 2
    *   "🛑 CLARIFICATION NEEDED" → Display questions, **STOP**. When user responds, `m mission update --section intent --content "[REFINED_INTENT]"`, then proceed
4.  **Log**: `m log --step "Intent" "Intent analyzed and refined"`

### Step 2: Context Analysis

1.  **Analyze Scope**: `m analyze scope` → Parse JSON, read `template_path`, follow template
    *   `m mission update --section scope --item "[file1]" --item "[file2]" ...`
2.  **Analyze Test Requirements**: `m analyze test` → Parse JSON, read `template_path`, follow template
    *   If needed: `m mission update --section scope --append --item "[test_file]" ...`
3.  **Duplication Analysis & WET→DRY Decision (Rule of Three)**:
    *   `m backlog list --include refactor` → Review existing patterns to avoid creating duplicate pattern IDs
    *   `m analyze duplication` → Parse JSON, read `template_path`, follow template
    *   For each pattern detected:
        - Check if semantically similar pattern already exists in backlog (reuse existing pattern-id)
        - `m backlog add "[description]" --type refactor --pattern-id "[pattern-id]"`
        - CLI auto-increments count if pattern exists, or creates with count=1
    *   **Determine Mission Type based on pattern counts**:
        - No duplication detected → `type=WET`
        - All pattern counts < 3 → `type=WET`
        - Any pattern count >= 3 → `type=DRY` (refactor that pattern)
        - User explicitly requests refactor → `type=DRY`
    *   `m mission update --frontmatter type=[WET|DRY]`
4.  **Log**: `m log --step "Context" "Context analyzed and mission updated"`

### Step 3: Complexity Analysis

1.  **Run Analysis**: `m analyze complexity` → Parse JSON, read `template_path`, follow template
2.  **Update Mission**: `m mission update --frontmatter track=[N] domains="[list]"`
3.  **React Based on Track**:
    *   **Track 1 (Atomic)**: 
        - Execute `m mission archive --force` to clean up generated mission.md
        - Load `.mission/libraries/displays/plan-atomic.md`, fill with `{{REFINED_INTENT}}` and `{{SUGGESTED_EDIT}}`
        - Display and **STOP**
    *   **Track 4 (Epic)**: 
        - Run `m analyze decompose` → Parse JSON, read `template_path`, follow template for decomposition guidance
        - Decompose intent into sub-intents based on template analysis
        - `m backlog list --exclude refactor --exclude completed` (parse JSON)
        - `m backlog add "[sub-intent]" ... --type decomposed`
        - Execute `m mission archive --force` to clean up generated mission.md
        - Load `.mission/libraries/displays/plan-epic.md`, fill with `{{SUB_INTENTS}}`
        - Display and **STOP**
    *   **Track 2 or 3**: Continue to Step 4
4.  **Log**: `m log --step "Analyze" "Complexity analysis complete. Track: [TRACK]"`

### Step 4: Plan and Validation

1.  **Create Plan**: Generate a numbered, step-by-step implementation plan with clear actions.
    - **Format**: "1. [Action] in [File]", "2. [Action] in [File]", etc.
    - **If Type is WET**: Add note: "Note: Allow duplication for initial implementation (WET principle)".
    - **If Type is DRY**: Add note: "Note: Refactor identified duplication into shared abstraction".
2.  **Define Verification**: Create a safe, executable verification command (e.g., `go test ./...`, `npm test`).
3.  **Update Mission**: 
    *   Execute `m mission update --section plan --item "[Step 1]" --item "[Step 2]" ...` to save plan.
    *   Execute `m mission update --section verification --content "[command]"` to save verification.
4.  **Log**: Run `m log --step "Validate" "Plan created and saved"`

### Step 5: Finalize & Generate

1.  **Finalize**: Execute `m mission finalize` to validate mission.md.
2.  **React to Output**:
    *   If `action: PROCEED` → Mission is valid, continue.
    *   If `action: INVALID` → Display errors and **STOP**.
3.  **Log**: Run `m log --step "Generate" "Mission generated successfully"`
4.  **Final Output**: 
    1.  Use file read tool to load template `.mission/libraries/displays/plan-success.md`.
    2.  Output the filled template with variables:
        - `{{TRACK}}`: From mission frontmatter.
        - `{{MISSION_TYPE}}`: From mission frontmatter.
        - `{{FILE_COUNT}}`: Count of files in `scope`.
        - `{{MISSION_CONTENT}}`: The content of `.mission/mission.md`.
