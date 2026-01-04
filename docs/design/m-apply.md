# Design: `m.apply` - Mission Execution

**Status**: 🚧 Future Enhancement (design complete, not implemented)
**Last Updated**: 2026-01-03
**Implementation**: Prompt template exists, CLI not implemented

## 1. Problem Statement
1.  **Lack of Polish**: AI-generated code often functions correctly ("First Draft") but lacks professional polish, idiomatic elegance, or optimization. Users currently rely on external tools to refine this code.
2.  **Architectural Inconsistency**: The current `m.apply` prompt contains excessive workflow logic ("Thick Agent"), making it brittle, hard to debug, and inconsistent with the newer `m.plan` implementation which uses a "Thick Client, Thin Agent" approach.
3.  **Non-deterministic State Management**: Mission status transitions are handled by shell scripts instead of robust CLI commands.

## 2. Implemented Solution
1.  **Two-Pass Execution**: First pass implements functionality, second pass (Polish) always runs to refine code quality.
2.  **Structured Execution Workflow**: A streamlined process with validation, execution, and polish.
3.  **"Thick Client, Thin Agent" Pattern**: AI orchestrates workflow using CLI tools for deterministic operations.
4.  **Execution Logging**: Structured logging of all execution steps for debugging and observability.

## 3. Architecture Overview

### 3.1 Responsibility Division

#### CLI Responsibilities (Thick Client)
**Deterministic operations, state management, data validation**
- Mission state validation and transitions
- File existence and safety checks
- Logging and audit trails
- Template loading and basic file I/O
- Checkpoint creation and restoration

#### AI Responsibilities (Thin Agent)
**Creative tasks, code analysis, decision-making**
- Code implementation and modification
- Code quality analysis and optimization
- Verification command execution and result interpretation
- Test selection and execution
- Polish decision-making (keep vs revert)

### 3.2 The Workflow (AI-Orchestrated)
Following the "Thick Client, Thin Agent" pattern, AI follows structured steps using CLI tools for deterministic operations:

**Prerequisites**: Run `m mission check --context apply` to validate mission state before execution.

| Step | CLI Responsibilities | AI Responsibilities |
|------|---------------------|--------------------|
| **1. Update Status** | • `m mission update --status active`<br>• `m checkpoint create` | • Update mission status<br>• Create initial checkpoint |
| **2. First Pass (Implementation)** | • Execution logging | • Verify SCOPE files exist<br>• Read `.mission/mission.md` and implement PLAN<br>• Enforce SCOPE constraints during file modification<br>• Execute verification command and interpret results<br>• Stop if verification fails |
| **3. Second Pass (Polish)** | • `m checkpoint create`<br>• `m checkpoint restore` on failure | • Review implemented code for quality improvements<br>• Apply idiomatic patterns and optimizations<br>• Re-run verification command<br>• Rollback on failure |
| **4. Status Handling** | • `m mission update --status failed` on failure<br>• `m checkpoint restore --all` on failure | • Load and populate display templates<br>• Decide success or failure based on verification |

### 3.3 Step-by-Step Workflow Details

#### Step 1: Update Status & Create Checkpoint
1. **CLI**: Execute `m mission update --status active`
2. **CLI**: Execute `m checkpoint create` to save the initial state.
3. **CLI**: Log status change and checkpoint creation.

#### Step 2: First Pass (Implementation)
1. **AI**: Verify all SCOPE files exist.
2. **AI**: Read `.mission/mission.md` and extract PLAN steps.
3. **AI**: Implement changes, enforcing SCOPE file constraints.
4. **CLI**: Log implementation progress.
5. **AI**: Execute verification command from `mission.md`.
6. **AI**: If verification fails, update status to `failed` and stop.

#### Step 3: Second Pass (Polish) - ALWAYS RUNS
1. **CLI**: Execute `m checkpoint create` to save the state after the first pass.
2. **AI**: Review all modified code from Step 2 and apply quality improvements.
3. **CLI**: Log polish progress.
4. **AI**: Re-execute verification command.
5. **CLI**: If verification fails, execute `m checkpoint restore <previous-checkpoint>` to roll back the polish changes.

#### Step 4: Status Handling
1. **AI**: Determine success or failure based on verification results.
2. **CLI**: On failure, execute `m checkpoint restore --all` to revert all changes and update status to `failed`.
3. **AI**: On success, load `apply-success.md` display template. On failure, load `apply-failure.md`.

## 4. Component Changes

### A. Prompt Template: `m.apply.md`
- The prompt is now simplified to focus on the two-pass execution and polish workflow.
- Commit message generation is removed.

### B. CLI Commands
- No new commands are required for `m.apply`. The existing `m mission` and `m checkpoint` commands are sufficient.

## 5. Success Metrics
- ✅ **Architectural Consistency**: `m.apply` is now a focused execution step.
- ✅ **Responsibility Clarity**: `m.apply` handles code changes; `m.complete` handles finalization.
- ✅ **Safety**: Rollback mechanism (`m checkpoint restore`) is preserved.
- ✅ **Simplicity**: The overall workflow is easier to understand and maintain.
