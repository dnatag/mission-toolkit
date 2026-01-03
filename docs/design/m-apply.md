# Design: `m.apply` - Mission Execution with Commit Message Generation

**Status**: 🚧 Future Enhancement (design complete, not implemented)  
**Last Updated**: 2026-01-03  
**Implementation**: Prompt template exists, CLI not implemented

## 1. Problem Statement
1.  **Lack of Polish**: AI-generated code often functions correctly ("First Draft") but lacks professional polish, idiomatic elegance, or optimization. Users currently rely on external tools to refine this code.
2.  **Architectural Inconsistency**: The current `m.apply` prompt contains excessive workflow logic ("Thick Agent"), making it brittle, hard to debug, and inconsistent with the newer `m.plan` implementation which uses a "Thick Client, Thin Agent" approach.
3.  **Non-deterministic State Management**: Mission status transitions are handled by shell scripts instead of robust CLI commands.
4.  **Missing Commit Message Generation**: No automatic generation of conventional commit messages during execution.

## 2. Implemented Solution
1.  **Two-Pass Execution**: First pass implements functionality, second pass (Polish) always runs to refine code quality.
2.  **Structured Execution Workflow**: Five-step process with validation, execution, polish, commit generation, and status handling.
3.  **"Thick Client, Thin Agent" Pattern**: AI orchestrates workflow using shell scripts for deterministic operations.
4.  **Automatic Commit Message Generation**: Conventional commit messages generated after polish pass and stored in mission.md.
5.  **Execution Logging**: Structured logging of all execution steps for debugging and observability.

## 3. Architecture Overview

### 3.1 Responsibility Division

#### CLI Responsibilities (Thick Client)
**Deterministic operations, state management, data validation**
- Mission state validation and transitions
- File existence and safety checks  
- Logging and audit trails
- Template loading and basic file I/O

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
| **1. Update Status** | • `m mission update --status active` - Transition to active<br>• Execution log initialization | • Update mission status<br>• Log status change |
| **2. First Pass (Implementation)** | • Execution logging | • Verify SCOPE files exist<br>• Read `.mission/mission.md` and implement PLAN<br>• Enforce SCOPE constraints during file modification<br>• Execute verification command and interpret results<br>• Stop if verification fails |
| **3. Second Pass (Polish)** | • Execution logging | • Review implemented code for quality improvements<br>• Apply idiomatic patterns and optimizations<br>• Improve readability and maintainability<br>• Re-run verification command<br>• Stop if verification fails |
| **4. Generate Commit Message** | • Execution logging | • Extract mission context (MISSION_ID, TYPE, TRACK, INTENT, SCOPE)<br>• Generate conventional commit message reflecting ALL changes<br>• Update mission.md `## COMMIT_MESSAGE` section<br>• Regenerate if code changes after initial generation |
| **5. Status Handling** | • `m mission update --status failed` on failure<br>• `git checkout .` on failure<br>• Execution logging | • Load and populate display templates<br>• Format final output for user<br>• Decide success or failure based on verification |

### 3.3 Step-by-Step Workflow Details

#### Step 1: Update Status
1. **CLI**: Execute `m mission update --status active`
2. **CLI**: Log status change via `m log`

#### Step 2: First Pass (Implementation)
1. **AI**: Verify all SCOPE files exist
2. **AI**: Read `.mission/mission.md` and extract PLAN steps
3. **AI**: Implement changes, enforcing SCOPE file constraints
4. **CLI**: Log implementation progress via `m log`
5. **AI**: Execute verification command from mission.md
6. **AI**: If verification fails, update status to `failed` and stop

#### Step 3: Second Pass (Polish) - ALWAYS RUNS
1. **AI**: Review all modified code from Step 2
2. **AI**: Apply quality improvements:
   - Idiomatic patterns and language conventions
   - Code readability and clarity
   - Performance optimizations
   - Error handling improvements
   - Documentation and comments where needed
3. **CLI**: Log polish progress via `m log`
4. **AI**: Re-execute verification command from mission.md
5. **AI**: If verification fails, update status to `failed` and stop

#### Step 4: Generate Commit Message
1. **AI**: Extract mission context (MISSION_ID, TYPE, TRACK, INTENT, SCOPE)
2. **AI**: Generate conventional commit message:
   - **Type**: WET → `feat`, DRY → `refactor` (override if clearly fix/docs/test/chore)
   - **Scope**: Extract from primary file/module in SCOPE
   - **Title**: Imperative mood, max 72 chars, capitalize first letter, no period
   - **Description**: Explain what and why (not how), wrap at 72 chars
   - **Footer**: Add `Mission-ID`, `Track`, `Type`
3. **AI**: Update mission.md `## COMMIT_MESSAGE` section
4. **CLI**: Log commit message generation via `m log`

**Regeneration Rule**: AI MUST regenerate commit message after ANY code changes (polish pass, user-requested adjustments)

#### Step 5: Status Handling
1. **AI**: Determine success or failure based on verification results
2. **Shell**: On failure, execute `git checkout .` to revert changes and update status to `failed`
3. **AI**: Load appropriate display template (`apply-success.md` or `apply-failure.md`)
4. **AI**: Populate template with mission details and commit message
5. **Shell**: Log final outcome to execution log

### 3.4 Commit Message Regeneration Rules

The AI MUST regenerate the commit message in these scenarios:
1. **After initial implementation** - Always generate after Step 2
2. **After ANY code changes** - Polish pass, bug fixes, user-requested changes
3. **After verification fixes** - If verification fails and code is modified to fix
4. **Description must be comprehensive** - Reflect ALL changes made, not just latest

**Commit Message Format**:
- **Type**: WET → `feat`, DRY → `refactor` (override if clearly fix/docs/test/chore)
- **Scope**: Extract from primary file/module in SCOPE (e.g., `auth`, `api`, `cli`)
- **Title**: Imperative mood, max 72 chars, capitalize first letter, no period
- **Description**: Explain what and why (not how), wrap at 72 chars
- **Footer**: Add `Mission-ID`, `Track`, `Type`

### 3.5 Component Changes

#### A. Prompt Template: `m.apply.md`
Structured five-step execution pattern:
- **Step 1: Pre-execution Validation** - Use shell scripts to validate and transition status
- **Step 2: First Pass (Implementation)** - AI implements plan with scope enforcement and verification
- **Step 3: Second Pass (Polish)** - AI always refines code quality after first pass completes
- **Step 4: Generate Commit Message** - AI creates conventional commit message reflecting all changes
- **Step 5: Status Handling** - Use shell scripts and display templates

#### B. Execution Logging
Structured logging using `execution.log` with template `libraries/logs/execution.md`:
- Each step logs outcome with format: `[SUCCESS/FAILED] | m.apply <step>: <name> | <details>`
- Logs archived with completed missions for debugging and analysis

#### C. CLI Commands
Deterministic operations for state management:

**`m mission check --context apply`** (Prerequisites)
- Validates mission.md exists and status is `planned` or `active`
- Returns JSON with `next_step` field

**`m mission update --status <status>`**
- Updates mission status atomically
- Supports: `active`, `failed`, `completed`

## 4. Implementation Status

### 🚧 Future Work
1. Five-step execution workflow with structured logging
2. Mandatory polish pass after first implementation pass
3. Automatic commit message generation in Step 4
4. Shell scripts for state validation and transitions
5. Display templates for success and failure outcomes
6. Execution log archival with completed missions
7. Commit message regeneration rules

### ✅ Current State
- Prompt template exists (`internal/templates/prompts/m.apply.md`)
- CLI commands implemented: `m mission check`, `m mission update`
- Prerequisites section validates mission state before execution
- AI-driven workflow with CLI tools for deterministic operations

## 5. Success Metrics
- ✅ **Architectural Consistency**: `m.apply.md` follows structured step-by-step pattern
- ✅ **Responsibility Clarity**: Clear division between shell scripts (deterministic) and AI (creative) tasks
- ✅ **Two-Pass Execution**: Polish pass always runs after first implementation pass completes
- ✅ **Code Quality**: Mandatory polish pass ensures professional, idiomatic code
- ✅ **Commit Message Generation**: Automatic conventional commit messages with regeneration rules
- ✅ **Execution Logging**: All steps logged with consistent format for debugging
- ✅ **Safety**: Rollback mechanism (`git checkout .`) on verification failure
- ✅ **Display Consistency**: Template-based output for success and failure scenarios
- ✅ **CLI Migration**: Shell scripts replaced with Go CLI commands (`m mission check`, `m mission update`)
