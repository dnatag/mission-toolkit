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
1.  **Structured Execution Workflow**: Four-step process with validation, execution, commit generation, and status handling.
2.  **"Thick Client, Thin Agent" Pattern**: AI orchestrates workflow using shell scripts for deterministic operations.
3.  **Automatic Commit Message Generation**: Conventional commit messages generated during execution and stored in mission.md.
4.  **Execution Logging**: Structured logging of all execution steps for debugging and observability.

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
Following the `m.plan` pattern, AI follows structured steps using CLI tools for deterministic operations:

| Step | Shell Script Responsibilities | AI Responsibilities |
|------|------------------------------|--------------------|
| **1. Pre-execution Validation** | • `validate-planned.sh` - Validate mission state<br>• `status-to-active.sh` - Transition to active<br>• Execution log initialization | • Check prerequisites (mission.md exists, status is planned)<br>• Verify SCOPE files exist<br>• Log validation outcome |
| **2. Mission Execution** | • Execution logging | • Read `.mission/mission.md` and implement PLAN<br>• Enforce SCOPE constraints during file modification<br>• Execute verification command and interpret results<br>• Stop if verification fails |
| **3. Generate Commit Message** | • Execution logging | • Extract mission context (MISSION_ID, TYPE, TRACK, INTENT, SCOPE)<br>• Generate conventional commit message<br>• Update mission.md `## COMMIT_MESSAGE` section<br>• Regenerate if code changes after initial generation |
| **4. Status Handling** | • Status update (active/failed)<br>• `git checkout .` on failure<br>• Execution logging | • Load and populate display templates<br>• Format final output for user<br>• Decide success or failure based on verification |

### 3.3 Step-by-Step Workflow Details

#### Step 1: Pre-execution Validation
1. **CLI**: Execute `m mission check --context apply` and validate JSON output
2. **AI**: Follow `next_step` field instructions (proceed/stop)
   - If `next_step` says STOP: Use error template and wait for user
   - If `next_step` says PROCEED: Continue to status update
3. **CLI**: Execute `m mission update --status active` if proceeding
4. **CLI**: Log validation completion via `m log`

#### Step 2: Implementation
1. **AI**: Read `.mission/mission.md` and extract PLAN steps
2. **AI**: Implement changes, enforcing SCOPE file constraints
3. **CLI**: Log implementation progress via `m log`
4. **AI**: Execute verification command from mission.md
5. **AI**: If verification fails, update status to `failed` and stop

#### Step 3: Generate Commit Message
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

#### Step 4: Status Handling
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
Structured four-step execution pattern:
- **Step 1: Pre-execution Validation** - Use shell scripts to validate and transition status
- **Step 2: Mission Execution** - AI implements plan with scope enforcement and verification
- **Step 3: Generate Commit Message** - AI creates conventional commit message and stores in mission.md
- **Step 4: Status Handling** - Use shell scripts and display templates

#### B. Execution Logging
Structured logging using `execution.log` with template `libraries/logs/execution.md`:
- Each step logs outcome with format: `[SUCCESS/FAILED] | m.apply <step>: <name> | <details>`
- Logs archived with completed missions for debugging and analysis

#### C. Shell Scripts (`.mission/libraries/scripts/`)
Deterministic operations for state management:

**`validate-planned.sh`**
- Validates mission.md exists and status is `planned`
- Checks SCOPE files exist
- Returns exit code 0 on success, 1 on failure

**`status-to-active.sh`**
- Updates mission status from `planned` to `active`
- Atomic file update

**`init-execution-log.sh`**
- Creates `.mission/execution.log` if it doesn't exist
- Uses template `libraries/scripts/init-execution-log.md`

## 4. Implementation Status

### 🚧 Future Work
1. Four-step execution workflow with structured logging
2. Automatic commit message generation in Step 3
3. Shell scripts for state validation and transitions
4. Display templates for success and failure outcomes
5. Execution log archival with completed missions
6. Commit message regeneration rules

### ✅ Current State
- Prompt template exists (`internal/templates/prompts/m.apply.md`)
- AI-driven workflow (no CLI implementation yet)
- Design ready for CLI implementation

## 5. Success Metrics
- ✅ **Architectural Consistency**: `m.apply.md` follows structured step-by-step pattern
- ✅ **Responsibility Clarity**: Clear division between shell scripts (deterministic) and AI (creative) tasks
- ✅ **Commit Message Generation**: Automatic conventional commit messages with regeneration rules
- ✅ **Execution Logging**: All steps logged with consistent format for debugging
- ✅ **Safety**: Rollback mechanism (`git checkout .`) on verification failure
- ✅ **Display Consistency**: Template-based output for success and failure scenarios
- 🚧 **Code Quality**: Single-pass execution (polish pass planned for future)
- 🚧 **CLI Migration**: Shell scripts work but could be replaced with Go CLI commands
