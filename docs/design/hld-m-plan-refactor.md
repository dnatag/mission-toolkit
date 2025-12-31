# High-Level Design: Refactoring `m.plan` to CLI-Driven Architecture

## 1. Problem Statement
The current `m.plan` workflow relies entirely on a "Simulated State Machine" within the LLM prompt. This leads to:
- **Non-deterministic behavior**: AI agents occasionally skip steps or ignore rules.
- **Fragile logic**: Complex branching (Track 1 vs 4) is handled probabilistically.
- **Inconsistent output**: The generated `mission.md` file sometimes deviates from the required format.

## 2. Proposed Solution
Transition from a **Prompt-Driven** architecture to a **CLI-Driven** architecture (Option 3).
- **AI Role**: Analyst & Estimator. It understands intent and estimates scope.
- **CLI Role**: Enforcer & Generator. It calculates complexity (Tracks), validates inputs, and writes the physical files.

## 3. Architecture Overview

### 3.1 New Workflow
1.  **User Intent**: User runs `/m.plan "Refactor login"`
2.  **AI Analysis**: AI reads the prompt, analyzes the request, and estimates affected files.
3.  **CLI Execution**: AI invokes the new `m plan create` command with structured flags.
    ```bash
    m plan create --intent "Refactor login" --scope "file1.go,file2.go" --type "WET"
    ```
4.  **Deterministic Logic (Go)**:
    - Validates file existence.
    - Calculates `Track` based on file count/rules.
    - Generates `mission.md` using strict Go templates.
5.  **Output**: CLI confirms creation; AI presents the result to the user.

### 3.2 Component Changes

#### A. New CLI Command: `m plan create`
A new subcommand in `cmd/plan.go` that accepts:
- `--intent` (string, required)
- `--scope` (comma-separated strings, required)
- `--type` (enum: WET, DRY, CLARIFICATION)
- `--est-lines` (int, optional, for complexity hints)
- `--parent` (string, optional, for DRY missions)

**Responsibilities:**
- **Validation**: Ensure scope files exist (warn if not).
- **Logic**:
    - `Track 1`: < 1 file, trivial change.
    - `Track 2`: 1-5 files.
    - `Track 3`: 6-9 files.
    - `Track 4`: 10+ files (Error/Decompose).
- **Generation**: Render `internal/templates/mission/wet.md` (etc.) with provided data.

#### B. Refactored Prompt: `m.plan.md`
Drastically simplified prompt that focuses on:
1.  **Analysis**: "Review the user's request."
2.  **Estimation**: "Identify which files need to change."
3.  **Execution**: "Construct and run the `m plan create` command."

*Removed*: All "If Track 1 do X" logic. The CLI handles that now.

#### C. Template Updates
- Ensure `internal/templates/` are compatible with Go's `text/template` rendering (already true, but verify variable names).

## 4. Implementation Plan

### Phase 1: Core CLI Logic
1.  Create `cmd/plan.go`.
2.  Implement `plan create` subcommand with flags.
3.  Port "Complexity Logic" from Markdown to Go struct.
4.  Implement file writing using existing `internal/templates` package.

### Phase 2: Prompt Migration
1.  Rewrite `internal/templates/prompts/m.plan.md`.
2.  Remove "Simulated State Machine" instructions.
3.  Add "Tool Use" instructions for `m plan create`.

### Phase 3: Validation & Testing
1.  Unit test `cmd/plan.go` logic (e.g., ensure 10 files triggers Track 4 error).
2.  End-to-end test with an AI agent to verify it correctly constructs the CLI command.

## 5. Risk Assessment
- **CLI Syntax Errors**: AI might malform the command string (quoting issues).
    - *Mitigation*: robust flag parsing and clear error messages back to the AI.
- **Loss of Nuance**: Rigid track rules might misclassify complex single-file changes.
    - *Mitigation*: Add an `--override-track` flag for manual AI override if justified.

## 6. Success Metrics
- **Zero Format Errors**: `mission.md` is always valid markdown.
- **100% Process Adherence**: Track 4 is *always* rejected/decomposed.
- **Reduced Prompt Size**: `m.plan.md` token count reduced by >40%.
