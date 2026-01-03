# Design: `m.clarify` - Clarification Workflow

**Status**: ✅ Implemented (prompt-based, no CLI)  
**Last Updated**: 2026-01-03  
**Implementation**: `internal/templates/prompts/m.clarify.md`

## 1. Problem Statement
The current `m.clarify` workflow is isolated, manual, and relies heavily on the AI ("Thick Agent") to manage state, validate logic, and generate the final `mission.md` using raw templates. This leads to inconsistency with `m plan`, potential for invalid output, and a lack of rigorous analysis (complexity, risk) after new information is gathered.

## 2. Proposed Solution
Refactor `m.clarify` to adopt the **"Thick Client, Thin Agent"** pattern.
- **Thick Client (CLI)**: Handles deterministic logic, state management, complexity analysis, validation, and file generation.
- **Thin Agent (AI)**: Focuses on understanding context, interacting with the user, and orchestrating the CLI tools.

The `m.clarify` workflow will be redesigned to feed back into the standard `m plan` pipeline, specifically re-entering at **Step 3 (Contextualization)** to ensure rigorous analysis of the clarified intent.

## 3. Architecture Overview

### 3.1 Prerequisites

**All clarification workflows start with**: `m mission check` to validate mission state.

- Checks if mission.md exists with `status: clarifying`
- Returns JSON with `next_step` instructions
- Ensures mission is ready for clarification

### 3.2 Principles
- **Single Source of Truth**: The CLI (`m plan` subcommands) defines the rules for a valid mission.
- **Reuse**: `m.clarify` should reuse the analysis and generation logic of `m plan`.
- **State Awareness**: The CLI knows when a mission is `clarifying` and guides the user/agent accordingly.

### 3.3 Roles & Responsibilities Matrix

| Task | Responsibility | Component | Description |
| :--- | :--- | :--- | :--- |
| **Context Gathering** | **AI (Thin Agent)** | `m.clarify` Prompt | Reads mission context and asks relevant questions. |
| **User Interaction** | **AI (Thin Agent)** | Chat Interface | Interprets user answers, handling ambiguity and nuance. |
| **Spec Management** | **AI (Thin Agent)** | `plan.json` | Updates the JSON spec directly (consistent with `m.plan` Step 4). |
| **Complexity Analysis** | **CLI (Thick Client)** | `m plan analyze` | Deterministic calculation of track and risk. |
| **Safety Validation** | **CLI (Thick Client)** | `m plan validate` | Enforces security and project constraints. |
| **Artifact Generation** | **CLI (Thick Client)** | `m plan generate` | Generates `mission.md` from the validated spec. |
| **Audit Trail** | **CLI (Thick Client)** | `m log` | Structured logging of the execution flow. |

### 3.4 The New Workflow

**Prerequisites**: Run `m mission check` to validate mission state before clarification.

1.  **Trigger**: User runs `m.clarify` (prompt).
2.  **Load**: AI reads `mission.md` (status: `clarifying`) and extracts questions.
3.  **Interact**: AI asks user for answers.
4.  **Contextualize (Re-run Step 3)**:
    - AI re-evaluates **Domains** and **Duplication** based on the clarified intent.
    - *Why*: New details might reveal security implications or existing patterns.
5.  **Update Spec**:
    - AI updates `.mission/plan.json` with the refined intent, scope, and domains.
    - *Pattern*: Reuses the "Draft Spec Creation" pattern from `m.plan` Step 4.
6.  **Pipeline Re-entry (Matches `m.plan` Steps 5-7)**:
    - **Analyze**: AI calls `m plan analyze`. CLI recalculates track/complexity.
    - **Validate**: AI calls `m plan validate`. CLI checks file paths and constraints.
    - **Generate**: AI calls `m plan generate`. CLI overwrites `mission.md` with the final, valid mission.

## 4. Detailed Design

### 4.1 CLI Updates (`cmd/plan.go`)
- **`m plan check`**: Update to detect `clarifying` state and suggest `m.clarify`.
- **`m plan analyze`**: Ensure it can handle `plan.json` updates during clarification.
- **`m plan generate`**: Ensure it can overwrite an existing `mission.md` safely.

### 4.2 Prompt Refactor (`m.clarify.md`)
The prompt will be restructured to:
1.  **Setup**: Check state (`m plan check`), Log start (`m log`).
2.  **Clarification Loop**:
    - Display questions.
    - Get answers.
    - **Contextualize**: Check for new domains/duplication.
    - **Update**: Overwrite `.mission/plan.json` with new info.
3.  **Analysis & Validation**:
    - Run `m plan analyze`.
    - Run `m plan validate`.
4.  **Finalization**:
    - Run `m plan generate`.
    - Log completion (`m log`).

### 4.3 Data Flow
`mission.md (clarifying)` -> **AI (Questions/Answers)** -> `plan.json (updated)` -> **CLI (Analyze/Validate)** -> `mission.md (planned)`

## 5. Workflow Analysis & Improvements

### 5.1 Current Limitations (Thick Agent)
- **Logic in Prompts**: The current `m.clarify` prompt contains logic for complexity scoring and file generation.
- **State Disconnect**: The clarification process happens "outside" the rigorous checks of `m plan`.

### 5.2 Strategic Improvements
- **Unified Pipeline**: Clarification is now a detour that returns to the main `m plan` highway at Step 3.
- **State Persistence**: The `plan.json` artifact becomes the persistent state during clarification, ensuring no data loss if the process is interrupted.
- **Simplicity**: By reusing the `m.plan` pattern for JSON creation, we avoid introducing new CLI commands (`m plan update`) and keep the tool surface area small.

## 6. Success Metrics
- **Consistency**: `mission.md` generated by `m.clarify` has the exact same structure and validity guarantees as one from `m plan`.
- **Robustness**: Invalid paths or complexity spikes are caught by the CLI, not the user.
- **Maintainability**: No new CLI commands required; leverages existing `m plan` infrastructure.
