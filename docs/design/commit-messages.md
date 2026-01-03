# Design: Commit Message Flow

**Status**: 🚧 Future Enhancement (design complete, not implemented)  
**Last Updated**: 2026-01-03  
**Implementation**: Part of m.apply and m.complete (not yet implemented)

## Overview

Conventional commit messages are generated during `m.apply` execution and stored in `mission.md`, then consumed by `m.complete` for git commit creation.

## Flow Diagram

```
                    ┌─────────────┐
                    │   m.plan    │
                    │             │
                    │ Analyzes    │
                    │ intent      │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Ambiguous? │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │ YES                     │ NO
              ▼                         ▼
     ┌─────────────┐           ┌─────────────┐
     │  m.clarify  │           │   m.plan    │
     │  (optional) │           │             │
     │             │           │ Creates     │
     │ Asks        │           │ mission.md  │
     │ questions   │           └──────┬──────┘
     └──────┬──────┘                  │
            │                         │
            │ Re-runs m.plan          │
            └────────┬────────────────┘
                     │
                     ▼
            ┌─────────────┐
            │ Review      │
            │ mission.md  │
            │             │
            │ • INTENT    │
            │ • SCOPE     │
            │ • PLAN      │
            │ • VERIFY    │
            └──────┬──────┘
                   │
                   │ Approve?
                   ▼
            ┌─────────────┐
            │   m.apply   │
            │             │
            │ Executes    │
            │ + Polish    │
            │ + Generates │
            │   commit    │
            └──────┬──────┘
                   │
            ┌──────▼──────┐
            │ Review code │
            │ Adjustments?│
            └──────┬──────┘
                   │
          ┌────────┼────────┐
          │ YES             │ NO
          ▼                 ▼
   ┌─────────────┐   ┌─────────────┐
   │ User        │   │ m.complete  │
   │ requests    │   │             │
   │ changes     │   │ Archives    │
   │             │   │ Creates     │
   │ AI fixes +  │   │ git commit  │
   │ regenerates │   └─────────────┘
   │ commit msg  │
   └──────┬──────┘
          │
          └──────────────────────────┐
                                     │
                                     ▼
                            ┌─────────────┐
                            │ m.complete  │
                            │             │
                            │ Archives    │
                            │ Creates     │
                            │ git commit  │
                            └─────────────┘
```

## Implementation Details

### 1. Mission Template Structure

**File**: `.mission/libraries/missions/wet.md` and `dry.md`

```markdown
# MISSION

id: {{MISSION_ID}}
type: WET
track: {{TRACK}}
status: planned

## INTENT
{{REFINED_INTENT}}

## SCOPE
{{FILE_LIST}}

## PLAN
{{PLAN_STEPS}}

## VERIFICATION
{{VERIFICATION_COMMAND}}

## COMMIT_MESSAGE
(Generated during m.apply execution)

## EXECUTION INSTRUCTIONS
...
```

### 2. m.apply Step 3: Generate Commit Message

**Trigger**: After successful execution and verification

**Process**:
1. Extract context from mission.md (MISSION_ID, TYPE, TRACK, INTENT, SCOPE)
2. Generate conventional commit message:
   - **Type**: WET → `feat`, DRY → `refactor`
   - **Scope**: Primary module from SCOPE (e.g., `auth`, `api`, `cli`)
   - **Title**: Imperative mood, max 72 chars
   - **Description**: What and why (not how)
   - **Footer**: Mission metadata
3. Update mission.md `## COMMIT_MESSAGE` section

**Example Output**:
```
## COMMIT_MESSAGE
feat(auth): Add JWT authentication middleware

Implement token-based authentication for API endpoints using JWT.
Includes middleware for token validation and user context injection.

Mission-ID: M-20250115-001
Track: 2
Type: WET
```

### 3. m.complete Step 2: Read Commit Message

**Trigger**: During outcome analysis

**Process**:
1. AI reads mission.md
2. Extracts `## COMMIT_MESSAGE` section content
3. Validates format (type, scope, title present)
4. Passes to CLI for git commit

### 4. m.complete Step 4: Create Git Commit

**Trigger**: After archival

**Process**:
1. CLI reads commit message from mission.md
2. Validates working tree state
3. Stages SCOPE files + .mission/completed/ files
4. Creates commit using go-git library
5. Returns commit hash

## Benefits

### ✅ Single Source of Truth
- mission.md contains all mission context including commit message
- No separate files or state management needed

### ✅ Persistent Storage
- Commit message survives between m.apply and m.complete invocations
- User can review/edit commit message before completing

### ✅ Archived with Mission
- Commit message preserved in completed/ directory
- Full context available for historical analysis

### ✅ Template-Driven
- Consistent with Mission Toolkit pattern
- AI generates message following conventional commit rules

### ✅ User Control
- User can edit `## COMMIT_MESSAGE` section before running m.complete
- Allows customization while maintaining format

## Conventional Commit Rules

### Type Mapping
- **WET missions** → `feat` (new feature)
- **DRY missions** → `refactor` (code improvement)
- **Override** if clearly: `fix`, `docs`, `test`, `chore`, `style`, `perf`

### Scope Extraction
- Use primary file/module from SCOPE
- Examples: `auth`, `api`, `ui`, `core`, `cli`, `db`
- If multiple files, use common parent directory or component

### Title Format
- Imperative mood ("Add" not "Added")
- Capitalize first letter
- No period at end
- Max 72 characters

### Description Format
- Explain what and why (not how)
- Wrap at 72 characters
- Separate from title with blank line
- Optional but recommended

### Footer Format
- Mission metadata (required):
  - `Mission-ID: <MISSION_ID>`
  - `Track: <TRACK>`
  - `Type: <WET|DRY>`

## Example Scenarios

### Scenario 1: Standard WET Mission

**m.apply generates**:
```
feat(api): Add product CRUD endpoints

Implement create, read, update, and delete operations for products.
Includes validation, error handling, and database integration.

Mission-ID: M-20250115-002
Track: 2
Type: WET
```

### Scenario 2: DRY Refactoring Mission

**m.apply generates**:
```
refactor(auth): Extract authentication logic to separate package

Move auth-related functions from handlers to dedicated auth package.
Improves code organization and enables reuse across modules.

Mission-ID: M-20250115-003
Track: 2
Type: DRY
```

### Scenario 3: Security Fix

**m.apply generates** (overrides WET → feat):
```
fix(auth): Prevent SQL injection in login endpoint

Add parameterized queries to user authentication logic.
Addresses security vulnerability in password validation.

Mission-ID: M-20250115-004
Track: 3
Type: WET
```

### Scenario 4: Polish Pass After Review

**After m.apply (initial)**:
```
## COMMIT_MESSAGE
feat(api): Add product CRUD endpoints

Implement create, read, update, and delete operations for products.
Includes validation, error handling, and database integration.

Mission-ID: M-20250115-005
Track: 2
Type: WET
```

**User reviews code and requests changes**:
```
User: "The validation logic has a bug - it's not checking for empty strings.
       Also add input sanitization for XSS prevention."
```

**AI makes fixes and updates COMMIT_MESSAGE**:
```
## COMMIT_MESSAGE
feat(api): Add product CRUD endpoints with security hardening

Implement create, read, update, and delete operations for products.
Includes input validation with empty string checks, XSS sanitization,
error handling, and database integration.

Mission-ID: M-20250115-005
Track: 2
Type: WET
```

**Key**: AI must regenerate commit message after ANY code changes to keep it accurate.

### Scenario 5: User Manual Edit Before Completion

**After m.apply**:
```
## COMMIT_MESSAGE
feat(api): Add product endpoints

Implement CRUD operations for products.

Mission-ID: M-20250115-005
Track: 2
Type: WET
```

**User edits mission.md**:
```
## COMMIT_MESSAGE
feat(api): Add product CRUD endpoints with validation

Implement create, read, update, and delete operations for products.
Includes input validation, error handling, and database integration.
Adds comprehensive test coverage for all endpoints.

Mission-ID: M-20250115-005
Track: 2
Type: WET
```

**m.complete uses edited version** for git commit.

## Polish Pass Pattern

### When to Update COMMIT_MESSAGE

AI **MUST** regenerate the commit message whenever:

1. **User requests code changes** after initial m.apply
   - Bug fixes discovered during review
   - Additional features requested
   - Refactoring or optimization changes

2. **AI makes polish pass improvements**
   - Code quality improvements
   - Performance optimizations
   - Security hardening

3. **Verification failures require fixes**
   - Test failures that need code changes
   - Linting errors that change implementation

### Polish Pass Workflow

```
m.apply (initial)
  ↓ Generates commit message
  
User reviews code
  ↓ "Fix the validation bug"
  
AI makes changes
  ↓ MUST regenerate commit message
  ↓ Updates ## COMMIT_MESSAGE section
  
User: "Looks good"
  ↓
  
m.complete
  ↓ Uses updated commit message
```

### Implementation Rules

**Rule 1: Always Regenerate After Code Changes**
```
IF (AI modifies any file in SCOPE after initial m.apply):
  THEN regenerate ## COMMIT_MESSAGE section
  AND update description to reflect all changes
```

**Rule 2: Preserve Mission Metadata**
```
WHEN regenerating commit message:
  KEEP Mission-ID, Track, Type footer unchanged
  UPDATE title and description to match current code state
```

**Rule 3: Cumulative Description**
```
IF multiple polish passes:
  THEN description should reflect ALL changes made
  NOT just the latest change
```

### Example: Multiple Polish Passes

**Initial m.apply**:
```
feat(api): Add product endpoints

Implement basic CRUD operations for products.
```

**After polish pass 1** (fix validation bug):
```
feat(api): Add product endpoints with validation

Implement CRUD operations for products.
Includes input validation and empty string checks.
```

**After polish pass 2** (add XSS protection):
```
feat(api): Add product endpoints with security hardening

Implement CRUD operations for products.
Includes input validation, empty string checks, and XSS sanitization.
```

**Final commit message reflects complete implementation**.

### AI Prompt Guidance

Add to m.apply prompt:

```markdown
### Step 3: Generate Commit Message

**CRITICAL**: Generate commit message that reflects ACTUAL implementation.

**When to regenerate**:
- After initial implementation (always)
- After ANY code changes (polish pass, bug fixes, improvements)
- After verification fixes that modify code

**How to regenerate**:
1. Review ALL changes made in current session
2. Generate title that summarizes complete implementation
3. Write description covering all significant changes
4. Replace entire ## COMMIT_MESSAGE section
5. Preserve Mission-ID, Track, Type footer

**Example**:
IF user says "fix the validation bug":
  AND you modify validation logic:
  THEN update commit message to mention validation improvements
```

## Migration Notes

### For Existing Missions
- Old missions without `## COMMIT_MESSAGE` section: CLI generates default message
- New missions: Always include `## COMMIT_MESSAGE` section from m.plan

### Backward Compatibility
- m.complete checks if `## COMMIT_MESSAGE` section exists
- If missing: Generate default message from INTENT
- If present: Use stored message

## Testing Strategy

### Unit Tests
- Commit message generation logic (type mapping, scope extraction)
- Conventional commit format validation
- mission.md parsing and updating

### Integration Tests
- Full flow: m.plan → m.apply → m.complete
- Verify commit message persists between commands
- Validate git commit creation with correct message

### Edge Cases
- Empty SCOPE (use "core" as default scope)
- Very long INTENT (truncate title to 72 chars)
- Special characters in SCOPE files (sanitize for scope)
- User deletes `## COMMIT_MESSAGE` section (regenerate)
