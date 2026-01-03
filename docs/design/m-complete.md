# Design: `m.complete` - Git Integration & Mission Archival

**Status**: 🚧 Future Enhancement (design complete, not implemented)  
**Last Updated**: 2026-01-03  
**Implementation**: Prompt template exists, CLI not implemented

## 1. Problem Statement
1.  **Manual Archival**: Current `m.complete` relies on shell scripts for mission archival and metrics collection.
2.  **No Git Integration**: Completed missions are not automatically committed, requiring manual git operations.
3.  **Architectural Inconsistency**: Prompt contains workflow logic instead of using CLI tools ("Thick Agent").
4.  **Missing Commit Standards**: No enforcement of conventional commit format for mission completion.
5.  **Fragmented CLI Structure**: Lacks unified mission command pattern established in `m.plan`, `m.clarify`, and `m.apply`.

## 2. Proposed Solution
1.  **Automated Git Commit**: Use go-git library to create conventional commits automatically.
2.  **"Thick Client, Thin Agent" Refactoring**: Move archival and git logic to CLI.
3.  **Commit Message Storage**: Store commit message in mission.md during m.apply execution.
4.  **Robust Archival**: CLI handles file operations, metrics aggregation, and backlog updates.
5.  **Unified CLI Structure**: Use `m mission` parent command for consistency with other workflows.

## 3. Architecture Overview

### 3.1 Commit Message Flow

**Generation (m.apply)**:
1. AI generates conventional commit message after successful execution
2. AI updates mission.md `## COMMIT_MESSAGE` section
3. Commit message stored in mission.md for m.complete to use

**Consumption (m.complete)**:
1. AI reads `## COMMIT_MESSAGE` section from mission.md
2. CLI extracts commit message and creates git commit
3. Commit message archived with mission in completed/ directory

**Benefits**:
- Single source of truth (mission.md)
- Persistent between commands
- Archived with mission context
- No additional files needed
- Template-driven consistency

### 3.2 Responsibility Division

#### CLI Responsibilities (Thick Client)
**Deterministic operations, file management, git operations**
- Mission state validation and status transitions
- File archival (mission.md, metrics.md, execution.log)
- Metrics aggregation and calculation
- Git staging and commit operations (working tree validation)
- Backlog updates (append patterns)
- Cleanup of active mission files
- Template loading and rendering

#### AI Responsibilities (Thin Agent)
**Content generation, analysis, decision-making**
- Extract commit message from mission.md
- Analyze mission outcomes from execution.log
- Identify refactoring opportunities from code patterns
- Populate completion templates with metrics
- Interpret CLI JSON outputs and follow next_step instructions

### 3.2 The Workflow (AI-Orchestrated)

**Prerequisites**: Run `m mission check --context complete` to validate mission state before completion.

| Step | CLI Responsibilities | AI Responsibilities |
|------|---------------------|--------------------|
| **1. Pre-completion Validation** | • `m mission check --context complete` - Validate mission state<br>• Verify mission status is `completed`<br>• Validate git working tree is clean | • Interpret CLI JSON output<br>• Follow `next_step` instructions (proceed/stop) |
| **2. Outcome Analysis** | • `m log` - Record analysis step | • Read mission.md COMMIT_MESSAGE field<br>• Read execution.log for verification results<br>• Validate mission success criteria met<br>• Identify refactoring patterns for backlog |
| **3. Archival** | • `m mission archive --patterns <json>` - Move files to completed/<br>• Generate mission metrics<br>• Update aggregate metrics.md<br>• Append patterns to backlog.md | • Provide refactoring patterns from analysis<br>• Interpret archival JSON output<br>• Confirm success before proceeding |
| **4. Git Commit** | • `m mission commit` - Stage and commit<br>• Validate working tree state<br>• Stage SCOPE files + .mission/ changes<br>• Create commit using go-git<br>• Return commit hash | • Extract commit message from mission.md<br>• Handle commit failures (dirty tree, conflicts)<br>• Interpret commit JSON output |
| **5. Cleanup & Display** | • `m mission update --status archived`<br>• Remove active mission files<br>• `m log` - Record completion | • Load `libraries/displays/complete-success.md`<br>• Populate with commit hash and metrics<br>• Display next steps to user |

### 3.3 Step-by-Step Workflow Details

#### Step 1: Pre-completion Validation
1. **CLI**: Execute `m mission check --context complete`
   - Validate mission status is `completed`
   - Check git working tree is clean (no uncommitted changes outside mission scope)
   - Return JSON with `status`, `mission_id`, `next_step`, `warnings`
2. **AI**: Interpret CLI JSON output
   - If `next_step` is STOP: Load error template and halt
   - If `next_step` is PROCEED: Continue to analysis
   - If warnings present: Display to user
3. **CLI**: Log validation completion via `m log`

#### Step 2: Outcome Analysis
1. **AI**: Read `.mission/mission.md` to extract:
   - MISSION_ID, INTENT, TYPE (WET/DRY), TRACK, SCOPE
   - **COMMIT_MESSAGE** (stored during m.apply execution in `## COMMIT_MESSAGE` section)
2. **AI**: Read `.mission/execution.log` to validate:
   - Verification step completed successfully
   - No error entries after implementation
   - Polish pass outcome (if applicable)
3. **AI**: Analyze for refactoring patterns:
   - Code duplication across files
   - Abstraction opportunities
   - Security improvements needed
4. **AI**: Prepare backlog updates (if patterns found):
   - Format as markdown list items
   - Include file references and pattern descriptions
5. **CLI**: Log analysis completion via `m log`

#### Step 3: Archival
1. **CLI**: Execute `m mission archive --patterns <json>`
   - Move `.mission/mission.md` → `.mission/completed/MISSION-ID-mission.md`
   - Move `.mission/execution.log` → `.mission/completed/MISSION-ID-execution.log`
   - Generate `.mission/completed/MISSION-ID-metrics.md` with:
     - Mission duration (start to completion timestamp)
     - File count and track
     - Verification success/failure
     - Polish applied (yes/no)
   - Update `.mission/metrics.md` with aggregated statistics
   - Append patterns to `.mission/backlog.md` (if provided)
   - Return JSON with archived file paths
2. **AI**: Interpret archival JSON output
   - Verify `success: true`
   - Extract archived file paths for display
3. **CLI**: Log archival completion via `m log`

#### Step 4: Git Commit
1. **AI**: Extract commit message from mission.md `## COMMIT_MESSAGE` section
2. **CLI**: Execute `m mission commit`
   - Read commit message from `.mission/mission.md` `## COMMIT_MESSAGE` section
   - Validate working tree state (no conflicts, no untracked files in SCOPE)
   - Stage tracked files in SCOPE that have modifications
   - Stage `.mission/completed/` archived files
   - Create commit using go-git library
   - Return JSON with commit hash and stats
3. **AI**: Interpret commit JSON output
   - If `success: false`: Display error and stop
   - If `success: true`: Extract commit hash for display
4. **CLI**: Log commit completion via `m log`

#### Step 5: Cleanup & Display
1. **CLI**: Execute `m mission update --status archived`
   - Update mission status in archived file
   - Return JSON confirmation
2. **CLI**: Remove `.mission/mission.md` and `.mission/execution.log` (active files only)
3. **AI**: Load display template `libraries/displays/complete-success.md`
4. **AI**: Populate template with:
   - MISSION_ID and INTENT
   - Commit hash and message
   - Key metrics (duration, files modified, track)
   - Next steps (start new mission)
5. **CLI**: Log mission completion via `m log`

### 3.4 Conventional Commit Format

**Generation**: Commit message generated during `m.apply` execution (Step 3).

**Storage**: `.mission/mission.md` contains:
```markdown
## COMMIT_MESSAGE
<type>(<scope>): <title>

<description>

Mission-ID: <MISSION_ID>
Track: <TRACK>
Type: <WET|DRY>
```

**Generation Rules (Applied during m.apply Step 3):**

**Type Mapping:**
- WET missions → `feat` (new feature implementation)
- DRY missions → `refactor` (code improvement)
- Override if change is clearly a fix/docs/test/chore

**Scope Extraction:**
- Use primary file or module from SCOPE
- Examples: `auth`, `api`, `ui`, `core`, `cli`
- If multiple files, use common parent directory or component name

**Title Format:**
- Imperative mood ("Add" not "Added")
- No period at end
- Max 72 characters
- Capitalize first letter

**Description Format:**
- Explain what and why, not how
- Wrap at 72 characters
- Separate from title with blank line

**Example Commit Message:**
```
feat(auth): Add JWT authentication middleware

Implement token-based authentication for API endpoints using JWT.
Includes middleware for token validation and user context injection.

Mission-ID: M-20250115-001
Track: 2
Type: WET
```

**Validation Rules:**
- Type must be one of: feat, fix, refactor, docs, test, chore, style, perf
- Scope is optional but recommended
- Title is required and must be ≤72 chars
- Description is optional but recommended for non-trivial changes

### 3.5 Component Changes

#### A. Refactored Prompt Template: `m.complete.md`
Restructured to follow the 5-step pattern established in `m.plan`, `m.clarify`, and `m.apply`:
- **Step 1: Pre-completion Validation** - Use `m mission check --context complete`
- **Step 2: Outcome Analysis** - Extract COMMIT_MESSAGE from mission.md, validate execution.log
- **Step 3: Archival** - Use `m mission archive --patterns <json>`
- **Step 4: Git Commit** - Use `m mission commit`
- **Step 5: Cleanup & Display** - Use `m mission update --status archived` and display templates

#### B. Updated Mission Template: `libraries/missions/wet.md` and `dry.md`
Add COMMIT_MESSAGE section to mission templates:
```markdown
# MISSION

id: {{MISSION_ID}}
type: WET
status: planned

## INTENT
...

## COMMIT_MESSAGE
(Generated during m.apply execution)
```

**Flow**:
1. m.plan creates mission.md with empty COMMIT_MESSAGE section
2. m.apply generates commit message and updates COMMIT_MESSAGE section
3. m.complete reads COMMIT_MESSAGE section and creates git commitION_ID}}
INTENT: {{INTENT}}
TRACK: {{TRACK}}
TYPE: {{TYPE}}
SCOPE:
  {{SCOPE_FILES}}
COMMIT_MESSAGE: |
  {{COMMIT_TYPE}}({{COMMIT_SCOPE}}): {{COMMIT_TITLE}}
  
  {{COMMIT_DESCRIPTION}}
  
  Mission-ID: {{MISSION_ID}}
  Track: {{TRACK}}
  Type: {{TYPE}}
```

**Note**: AI populates COMMIT_MESSAGE during `m.apply` Step 5 (Completion) after successful verification.

#### C. New Library Template: `libraries/displays/complete-success.md`
Template for displaying completion results to user:
```markdown
✅ MISSION COMPLETED: {{MISSION_ID}}

📋 SUMMARY:
{{INTENT}}

📊 METRICS:
• Track: {{TRACK}}
• Files Modified: {{FILE_COUNT}}
• Duration: {{DURATION}}
• Verification: {{VERIFICATION_STATUS}}

🔗 GIT COMMIT:
{{COMMIT_HASH}}
{{COMMIT_MESSAGE}}

🚀 NEXT STEPS:
• Start new mission: /m.plan "<your intent>"
• Review metrics: m status
• Check backlog: cat .mission/backlog.md
```

#### D. CLI Commands (`cmd/mission.go`)

**`m mission check --context complete`**
- **Purpose**: Validate mission state before completion
- **Logic**:
  - Check `.mission/mission.md` exists
  - Validate status is `completed`
  - Check git working tree is clean (no uncommitted changes outside SCOPE)
  - Validate COMMIT_MESSAGE field exists in mission.md
  - Return JSON with validation results
- **Output Example**:
  ```json
  {
    "status": "completed",
    "mission_id": "M-20250115-001",
    "context": "complete",
    "next_step": "PROCEED",
    "warnings": [],
    "message": "Mission ready for completion"
  }
  ```

**`m mission archive --patterns <json>`**
- **Purpose**: Archive mission files and update metrics
- **Flags**:
  - `--patterns` (optional): JSON array of refactoring patterns to append to backlog
- **Logic**:
  - Read `.mission/mission.md` and `.mission/execution.log`
  - Generate MISSION-ID-based filenames
  - Move files to `.mission/completed/`
  - Calculate mission metrics (duration, file count, verification status)
  - Generate `.mission/completed/MISSION-ID-metrics.md`
  - Update aggregate `.mission/metrics.md`
  - Append patterns to `.mission/backlog.md` (if provided)
  - Return JSON with archived file paths
- **Output Example**:
  ```json
  {
    "success": true,
    "mission_id": "M-20250115-001",
    "archived_files": [
      ".mission/completed/M-20250115-001-mission.md",
      ".mission/completed/M-20250115-001-execution.log",
      ".mission/completed/M-20250115-001-metrics.md"
    ],
    "metrics": {
      "duration_minutes": 15,
      "files_modified": 3,
      "track": 2
    }
  }
  ```

**`m mission commit`**
- **Purpose**: Create git commit with message from mission.md
- **Library**: Use `github.com/go-git/go-git/v5`
- **Logic**:
  - Read COMMIT_MESSAGE from `.mission/mission.md`
  - Validate conventional commit format
  - Check git working tree state (no conflicts, no untracked files)
  - Read SCOPE from mission.md
  - Stage tracked files in SCOPE with modifications
  - Stage `.mission/completed/` directory
  - Create commit using go-git
  - Return JSON with commit hash and stats
- **Output Example**:
  ```json
  {
    "success": true,
    "commit_hash": "a1b2c3d4",
    "files_staged": 4,
    "message": "feat(auth): Add JWT authentication middleware"
  }
  ```

#### E. Git Integration (`internal/git/`)

**Package Structure:**
```go
package git

import "github.com/go-git/go-git/v5"

type CommitService struct {
    repo *git.Repository
}

// ValidateConventionalCommit checks commit message format
func (s *CommitService) ValidateConventionalCommit(msg string) error

// StageFiles stages specified files for commit
func (s *CommitService) StageFiles(files []string) error

// Commit creates a commit with the given message
func (s *CommitService) Commit(msg string) (string, error)

// GetStatus returns current git status
func (s *CommitService) GetStatus() (*Status, error)
```

**Why go-git over CLI:**
- Better error handling and validation
- No dependency on git binary installation
- Programmatic control over git operations
- Easier to test and mock
- Cross-platform consistency

## 4. Implementation Plan

### Phase 1: Git Integration
1. Add `github.com/go-git/go-git/v5` dependency
2. Create `internal/git/commit.go` with CommitService
3. Implement conventional commit validation
4. Add unit tests for git operations

### Phase 2: Archive Command
1. Create `cmd/mission.go` with `archive` subcommand (if not exists)
2. Implement `m mission archive --patterns <json>` subcommand
3. Add metrics calculation logic
4. Add unit tests for archival operations

### Phase 3: Commit Command
1. Implement `m mission commit` subcommand
2. Integrate CommitService from internal/git
3. Add SCOPE file staging logic
4. Add unit tests for commit operations

### Phase 4: Prompt Refactoring
1. Update `libraries/missions/wet.md` and `dry.md` with COMMIT_MESSAGE field
2. Create `libraries/displays/complete-success.md` template
3. Rewrite `m.complete.md` to follow 5-step pattern
4. Add explicit CLI command calls at each step
5. Update `m.apply.md` to populate COMMIT_MESSAGE field in Step 5

### Phase 5: Testing & Validation
1. Test happy path (validation → analysis → archive → commit → cleanup)
2. Test failure scenarios (invalid commit message, git conflicts)
3. Test with WET and DRY missions
4. Validate conventional commit format enforcement

## 5. Success Metrics
- **Architectural Consistency**: `m.complete.md` follows same structured pattern as other commands
- **Automated Commits**: 100% of completed missions result in git commits
- **Commit Quality**: All commits follow conventional commit format
- **Responsibility Clarity**: Clear division between CLI (git/archival) and AI (analysis/content)
- **Unified Logging**: All steps logged via `m log` with consistent format
- **Robustness**: Git operations handled gracefully with proper error messages
- **Maintainability**: Go library provides better control than shell scripts
