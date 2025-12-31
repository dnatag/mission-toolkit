# High-Level Design: Refactoring `m.plan` to CLI-Driven Toolbox Architecture

## 1. Problem Statement
The current `m.plan` workflow relies entirely on a "Simulated State Machine" within the LLM prompt. This leads to:
- **Non-deterministic behavior**: AI agents occasionally skip steps or ignore rules.
- **Fragile logic**: Complex branching (Track 1 vs 4) is handled probabilistically.
- **Inconsistent output**: The generated `mission.md` file sometimes deviates from the required format.
- **Black Box Execution**: Hard to debug which part of the planning process failed.
- **Missing Tests**: AI often forgets to include test files in the scope.
- **Non-Idiomatic Code**: AI drifts from project-specific coding standards.
- **Fragmented Logging**: AI and CLI logs might diverge in format.

## 2. Proposed Solution
Transition to a **"Thick Client, Thin Agent"** architecture (Toolbox Approach).
- **AI Role**: Orchestrator. It analyzes intent, identifies scope, and calls specific CLI tools to validate its assumptions before committing.
- **CLI Role**: Logic Engine. It provides discrete subcommands for analysis, validation, generation, and logging.

## 3. Responsibility Assignment

| Task | Owner | Rationale | Command |
| :--- | :--- | :--- | :--- |
| **Pre-execution Check** | **CLI** | Checking mission state and cleaning stale artifacts. | `m plan check` |
| **Clarification Check** | **AI** | Requires semantic understanding of ambiguity. | N/A (Prompt Logic) |
| **Intent Analysis** | **AI** | Requires summarizing natural language. | N/A (Prompt Logic) |
| **Complexity Analysis** | **CLI** | Deterministic rules based on file counts + domain flags. | `m plan analyze` |
| **Test Coverage Check** | **CLI** | Heuristic check for missing test files (Agnostic). | `m plan analyze` |
| **Duplication Analysis** | **AI** | Semantic similarity check. | N/A (Prompt Logic) |
| **Security/Scope Validation** | **CLI** | Strict file system and safety checks. | `m plan validate` |
| **Mission Generation** | **CLI** | Strict file formatting and writing + Guideline Injection. | `m plan generate` |
| **Logging** | **CLI** | Unified logging format for both AI and System. | `m log` |

## 4. Architecture Overview

### 4.1 New Workflow
1.  **Pre-Check**: AI runs `m plan check`.
    - CLI checks mission state, generates `MISSION_ID`, cleans stale `plan.json`.
    - CLI logs: "Started planning session".
2.  **User Intent**: User runs `/m.plan "Refactor login"`
3.  **AI Logging**: AI runs `m log --step "Intent" "Analyzing user request..."`.
4.  **Draft Spec**: AI creates a draft `plan.json` with Intent, Scope, and Domain.
5.  **Complexity & Test Check**: AI runs `m plan analyze --file plan.json`.
    - CLI logs: "Analyzed scope: X files".
    - CLI returns Track, Recommendation, and Warnings (missing tests).
6.  **Validation**: AI runs `m plan validate --file plan.json`.
    - CLI logs: "Validation passed/failed".
7.  **Finalize Spec**: AI updates `plan.json` with Plan Steps and Verification.
8.  **Generation**: AI runs `m plan generate --file plan.json`.
    - CLI logs: "Mission generated".
9.  **Output**: CLI creates `mission.md`.

### 4.2 Component Changes

#### A. New CLI Command: `m plan` (Parent)
Parent command for planning tools.

#### B. Subcommand: `m plan check`
- **Logic**:
    - Check if `.mission/mission.md` exists.
    - **ID Generation**: Generate `MISSION_ID` and store in `.mission/id` (or `plan.json`).
    - **Cleanup**: Remove `.mission/plan.json` if it exists.
    - Return JSON status.

#### C. Subcommand: `m plan analyze`
- **Flags**: `--file` (path to `plan.json`).
- **Logic**:
    - Read `plan.json`.
    - Count files (Base Track).
    - Apply Domain Multipliers.
    - **Test Gap Detection**: Check for missing test files based on patterns.
    - **Output**: JSON with `track`, `reason`, `recommendation` ("proceed"/"decompose"), `warnings`.

#### D. Subcommand: `m plan validate`
- **Flags**: `--file` (path to `plan.json`).
- **Logic**:
    - Read `plan.json`.
    - **Scope Check**: Fail on traversal/absolute paths. Warn on missing files.
    - **Verification Check**: Fail on banned patterns.
    - Return JSON output.

#### E. Subcommand: `m plan generate`
- **Flags**: `--file` (path to `plan.json`).
- **Logic**:
    - Read `plan.json`.
    - **Guideline Injection**: Append `.mission/guidelines.md` content.
    - Render `mission.md`.

#### F. New CLI Command: `m log`
- **Flags**: `--level` (INFO/SUCCESS/ERROR), `--step` (string), `message` (arg).
- **Logic**:
    - Read `MISSION_ID` from `.mission/id` (or find active mission).
    - Append formatted line to `.mission/execution.log`.

#### G. Refactored Prompt: `m.plan.md`
Prompt becomes a workflow script using `m plan *` and `m log` commands.

## 5. Implementation Plan

### Phase 1: Core Infrastructure
1.  Create `internal/logger` package.
2.  Create `cmd/log.go` (`m log` command).
3.  Create `cmd/plan.go` (Parent command).

### Phase 2: Check & Analysis Tools
1.  Define `PlanSpec` struct.
2.  Implement `check` subcommand (ID generation).
3.  Implement `analyze` subcommand (Complexity + Test Gap).

### Phase 3: Validation & Generation Tools
1.  Implement `validate` subcommand.
2.  Implement `generate` subcommand (Guideline Injection).

### Phase 4: Prompt Migration
1.  Rewrite `m.plan.md`.
2.  Test end-to-end.

## 6. Success Metrics
- **Unified Logs**: `execution.log` is consistent and readable.
- **Test Coverage**: AI consistently includes test files.
- **Idiomatic Code**: Guidelines are injected.
- **Safety**: Pre-check prevents overwriting.
