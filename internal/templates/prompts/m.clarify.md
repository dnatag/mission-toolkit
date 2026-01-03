---
description: "Handle clarification workflow and update mission"
---

## Prerequisites

**CRITICAL:** Run `m mission check` to validate mission state before clarification.

1. **Execute Check**: Run `m mission check` and parse JSON output
2. **Validate Status**: Check `next_step` field:
   - If `next_step` says "Run the m.clarify prompt to resolve questions" → Continue with clarification
   - If `next_step` says "STOP" → Display the message and halt
   - If no mission exists → Use template `.mission/libraries/displays/error-no-mission.md`

## Role & Objective

You are the **Clarification Handler**. Your goal is to guide the user through answering clarification questions and then use the `m plan` CLI tools to finalize the mission.

## Execution Steps

Before processing, read `.mission/governance.md` and current `.mission/mission.md`.

### Step 1: Display Questions and Collect Answers

1.  **Extract Questions**: Parse `NEED_CLARIFICATION` section from `.mission/mission.md`.
2.  **Display to User**: Show numbered list of questions.
3.  **Collect Responses**: Get user answers.
4.  **Refine Intent**: Combine original intent with user answers to form a `[REFINED_INTENT]`.
5.  **Log**: Run `m log --step "Clarify" "Questions displayed, awaiting user responses"`

**Display Format:**
```
🤔 CLARIFICATION NEEDED

Please provide answers to these questions:

1. [Question 1]
2. [Question 2]
3. [Question 3]

Provide your answers - you can reference questions by number or respond in any clear format.
```

### Step 2: Contextualization (The "How")

*Prerequisite: You must have a clear `[REFINED_INTENT]` from Step 2.*

1.  **Duplication Check**: Use file read tool to load template `.mission/libraries/analysis/duplication.md`. Use it to scan for existing patterns or redundant code.
    *   **Action**: If duplication or refactoring opportunities are detected:
        1.  Use file read tool to load `.mission/backlog.md`.
        2.  Append the pattern and affected files to the `## REFACTORING OPPORTUNITIES` section.
        3.  Note the findings for use in the Plan section later.
2.  **Domain Identification**: Use file read tool to load template `.mission/libraries/analysis/domain.md`. Use it to select applicable domains.
3.  **Log**: Run `m log --step "Context" "Duplication and domain analysis complete"`

### Step 3: Update Plan Specification

1.  **Identify Scope**: Determine which files need to be modified or created based on the `[REFINED_INTENT]`.
2.  **Determine Type**:
    *   If Duplication Check (Step 3) found a `refactor_opportunity`, set `type` to "DRY".
    *   Otherwise, set `type` to "WET".
3.  **Create Draft Spec**: Create or update `.mission/plan.json` with the following structure:
    ```json
    {
      "intent": "[REFINED_INTENT]",
      "type": "[WET or DRY]",
      "scope": ["path/to/file1.go", "path/to/file2.go"],
      "domain": ["security"]
    }
    ```
4.  **Log**: Run `m log --step "Draft" "Updated draft plan.json with refined intent and scope"`

### Step 4: Complexity Analysis

1.  **Run Analysis**: Execute `m plan analyze --file .mission/plan.json`.
2.  **Validate JSON**: Ensure the CLI output is valid JSON. If not, report CLI error and stop.
3.  **Follow Instructions**: Read the `next_step` field from the JSON output and **follow it literally**.
    *   **If `next_step` says UPDATE**: Update `.mission/plan.json` as instructed and retry analysis.
    *   **If `next_step` says STOP (Track 1)**:
        1.  Use file read tool to load template `.mission/libraries/displays/plan-atomic.md`.
        2.  Output the filled `plan-atomic.md` template with `{{REFINED_INTENT}}` and a `{{SUGGESTED_EDIT}}`.
    *   **If `next_step` says STOP (Track 4)**:
        1.  Use file read tool to load `.mission/backlog.md`.
        2.  Decompose the `[REFINED_INTENT]` into 3-5 atomic sub-intents.
        3.  Append sub-intents to the `## DECOMPOSED INTENTS` section of `.mission/backlog.md`.
        4.  Use file read tool to load template `.mission/libraries/displays/plan-epic.md`.
        5.  Output the filled `plan-epic.md` template.
    *   **If `next_step` says PROCEED**: Continue to Step 6.
    *   **If CLI command fails**: Report error and ask user to check CLI installation.
4.  **Log**: Run `m log --step "Analyze" "Complexity analysis complete. Track: [TRACK]"`

### Step 5: Validation

1.  **Run Validation**: Execute `m plan validate --file .mission/plan.json`.
2.  **Validate JSON**: Ensure the CLI output is valid JSON. If not, report CLI error and stop.
3.  **Handle Output**:
    *   If `valid: true`, proceed to Step 7.
    *   If `valid: false`, **STOP** and report the errors to the user. Fix the `plan.json` if possible (e.g., remove invalid files) and retry validation.
    *   **If CLI command fails**: Report error and ask user to check CLI installation.
4.  **Log**: Run `m log --step "Validate" "Plan validation passed"`

### Step 6: Finalize & Generate

1.  **Develop Plan**: Create a step-by-step implementation plan.
    *   **If Type is WET**: Add a note to the plan: "Note: Allow duplication for initial implementation (WET principle)."
    *   **If Type is DRY**: Add a note to the plan: "Note: Refactor identified duplication into shared abstraction."
2.  **Define Verification**: Create a safe verification command (e.g., `go test ./...`).
3.  **Update Spec**: Update `.mission/plan.json` to include `plan` and `verification` fields:
    ```json
    {
      "intent": "...",
      "type": "...",
      "scope": [...],
      "domain": [...],
      "plan": [
        "Step 1: ...",
        "Step 2: ...",
        "Note: ..."
      ],
      "verification": "go test ./..."
    }
    ```
4.  **Generate**: Execute `m mission create --file .mission/plan.json`.
5.  **Validate Generation**: Ensure the CLI command succeeded and `.mission/mission.md` was created.
6.  **Log**: Run `m log --step "Generate" "Mission generated successfully"`
7.  **Final Output**: 
    1.  Use file read tool to load template `.mission/libraries/displays/clarify-success.md`.
    2.  Output the filled `clarify-success.md` template with variables:
        - `{{TRACK}}`: From `m plan analyze` output.
        - `{{MISSION_TYPE}}`: WET or DRY.
        - `{{MISSION_CONTENT}}`: The content of the newly created `.mission/mission.md`.
