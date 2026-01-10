# Design: Analysis-Based Planning Architecture v5

**Status**: ✅ IMPLEMENTED - Template-providing CLI + LLM-driven updates  
**Supersedes**: v4 (CLI validates JSON)  
**Implementation Date**: January 2026

## Problem Statement

Current `/m.plan` implementation requires LLMs to orchestrate 15+ CLI calls across 5 steps with complex conditional logic. This works well for instruction-following models (Claude, GPT-4) but fails with models that prioritize internal reasoning over instructions (Gemini).

**Failure modes:**
- Skipping steps (e.g., jumping to implementation without analysis)
- Ignoring conditional branches (e.g., proceeding when `next_step` says "STOP")
- Making assumptions instead of reading tool outputs
- Manual file creation instead of using CLI

**v3 Problem**: Embedded templates can be ignored by LLMs - no enforcement.  
**v4 Problem**: CLI parsing JSON is over-engineering - LLM should extract fields.

## Solution: Template-Providing CLI + LLM-Driven Analysis & Updates

**Key Insight**: CLI's job is to provide templates and update mission.md. LLM's job is to analyze and extract fields.

Each analysis step:
1. LLM calls `m analyze <step>` to get template
2. CLI returns template content
3. LLM performs analysis following template
4. LLM extracts relevant fields from its analysis
5. LLM calls `m mission update` to save results
6. CLI updates mission.md and returns confirmation

**Enforcement**: LLM must call `m analyze` to get template (can't skip), then must call `m mission update` to proceed (tool-enforced checkpoint).

## Architecture

### Mission File Structure

```markdown
---
id: 20260109083903-4668
iteration: 1
status: planning
track: 3
type: WET
domains: ["security"]
---

## INTENT
Add JWT authentication with rate limiting

## SCOPE
cmd/auth.go
middleware/jwt.go
handlers/login.go
handlers/login_test.go

## PLAN
- [ ] 1. Create JWT handler in cmd/auth.go
- [ ] 2. Add middleware in middleware/jwt.go
- [ ] 3. Wire up login handler

## VERIFICATION
go test ./...
```

### Two-Command Pattern

**Step 1: Get Template**
```bash
m analyze intent "$ARGUMENTS"  # Only intent takes argument
m analyze clarify              # Others read from mission.md
# Returns: template content + context
```

**Step 2: Update Mission**
```bash
m mission update --section intent --content "Add JWT authentication"
# Returns: confirmation with action field
```

## Analysis Commands

### `m analyze intent "<user-input>"`

**Purpose:** Provide intent analysis template

**Input:**
- User's raw intent string (positional argument)

**CLI Processing:**
1. Load `.mission/libraries/analysis/intent.md`
2. Include user input in output
3. Return template content as text

**Output:**
```
# INTENT ANALYSIS TEMPLATE

## User Input
add auth

## Purpose
Distill raw user input into a clear, actionable, and scoped intent statement.

## Analysis Steps
1. Identify Core Action (The Verb)
2. Identify Target Scope (The Noun)
3. Identify Constraints (The Rules)

## Refinement Rules
1. Be Specific
2. Be Concise
3. Be Technical
4. Be Atomic

## Output Format
Produce a JSON object with action and refined intent:
{
  "action": "PROCEED" | "AMBIGUOUS",
  "refined_intent": "string",
  "reason": "string"  // if AMBIGUOUS
}
```

**LLM Flow:**
1. Call `m analyze intent "add auth"` to get template + user input
2. Follow template to analyze
3. Produce: `{"action": "PROCEED", "refined_intent": "Add JWT authentication"}`
4. Extract `refined_intent` field
5. Call `m mission create --intent "Add JWT authentication"`

---

### `m analyze clarify`

**Purpose:** Provide clarification analysis template

**CLI Processing:**
1. Load `.mission/libraries/analysis/clarification.md`
2. Read current INTENT from mission.md
3. Return template + current intent

**Output:**
```
# CLARIFICATION ANALYSIS TEMPLATE

## Current Intent
Add JWT authentication

## Purpose
Scan user intent for ambiguous requirements that need clarification.

## Analysis Categories
1. Scope & Requirements
2. Domain & Risk Context

## Output Format
{
  "action": "CLARIFY" | "ASSUMPTIONS" | "CLEAR",
  "questions": ["array"],
  "assumptions": ["array"]
}
```

**LLM Flow:**
1. Call `m analyze clarify` to get template + current intent
2. Follow template to analyze
3. Produce: `{"action": "CLARIFY", "questions": ["JWT or OAuth?"]}`
4. Display questions to user, STOP

---

### `m analyze scope`

**Purpose:** Provide scope analysis template

**CLI Processing:**
1. Load `.mission/libraries/analysis/scope.md`
2. Read current INTENT from mission.md
3. Return template + intent

**Output:**
```
# SCOPE ANALYSIS TEMPLATE

## Current Intent
Add JWT authentication

## Purpose
Determine which implementation files need to be modified or created.

## Analysis Steps
1. File Discovery
2. File Classification (primary, secondary, config)

## Output Format
{
  "action": "PROCEED",
  "scope": {
    "primary": ["array"],
    "secondary": ["array"],
    "config": ["array"]
  }
}
```

**LLM Flow:**
1. Call `m analyze scope` to get template + intent
2. Follow template to discover files
3. Produce: `{"action": "PROCEED", "scope": {"primary": ["auth.go", "handler.go"], ...}}`
4. Extract files from scope object
5. Call `m mission update --section scope --item "auth.go" --item "handler.go" ...`

---

### `m analyze test`

**Purpose:** Provide test analysis template

**CLI Processing:**
1. Load `.mission/libraries/analysis/test.md`
2. Read INTENT and SCOPE from mission.md
3. Return template + context

**Output:**
```
# TEST ANALYSIS TEMPLATE

## Current Intent
Add JWT authentication

## Current Scope
auth.go, handler.go, routes.go

## Purpose
Determine if automated tests are necessary using risk-based analysis.

## Analysis Criteria
1. Cyclomatic Complexity
2. Cost of Failure
3. Longevity
4. Determinism

## Output Format
{
  "action": "PROCEED",
  "test_verdict": "CRITICAL" | "RECOMMENDED" | "UNNECESSARY",
  "test_files": ["array"]
}
```

**LLM Flow:**
1. Call `m analyze test` to get template + context
2. Follow template to evaluate risk
3. Produce: `{"action": "PROCEED", "test_verdict": "CRITICAL", "test_files": ["auth_test.go"]}`
4. Extract test_files
5. Call `m mission update --section scope --item "auth_test.go" ...`

---

### `m analyze duplication`

**Purpose:** Provide duplication analysis template

**CLI Processing:**
1. Load `.mission/libraries/analysis/duplication.md`
2. Read INTENT and SCOPE from mission.md
3. Return template + context

**Output:**
```
# DUPLICATION ANALYSIS TEMPLATE

## Current Intent
Add JWT authentication

## Current Scope
auth.go, handler.go, routes.go

## Purpose
Identify existing code patterns similar to the requested feature.

## Output Format
{
  "duplication_detected": true | false,
  "patterns": ["array"],
  "mission_type": "WET"
}
```

**LLM Flow:**
1. Call `m analyze duplication` to get template + context
2. Follow template to scan for patterns
3. Produce: `{"duplication_detected": false, "patterns": [], "mission_type": "WET"}`
4. Extract mission_type
5. Call `m mission update --frontmatter type=WET`

---

### `m analyze complexity`

**Purpose:** Provide complexity analysis template

**CLI Processing:**
1. Load `.mission/libraries/analysis/complexity.md`
2. Read INTENT and SCOPE from mission.md
3. Return template + context

**Output:**
```
# COMPLEXITY ANALYSIS TEMPLATE

## Current Intent
Add JWT authentication

## Current Scope
auth.go (3 files total)

## Purpose
Identify domains and determine complexity track.

## Step 1: Identify Domains
[domain list]

## Step 2: Determine Track
[track calculation logic]

## Output Format
{
  "track": 1 | 2 | 3 | 4,
  "action": "ATOMIC_EDIT" | "PROCEED" | "DECOMPOSE",
  "domains": ["array"],
  "factors": {...}
}
```

**LLM Flow:**
1. Call `m analyze complexity` to get template + context
2. Follow template to identify domains and calculate track
3. Produce: `{"track": 3, "action": "PROCEED", "domains": ["security"], ...}`
4. Extract track and domains
5. Call `m mission update --frontmatter track=3 domains="security"`

---

### `m analyze plan`

**Purpose:** Provide plan generation guidance

**CLI Processing:**
1. Read all context from mission.md (INTENT, SCOPE, track, type)
2. Return context + planning guidance

**Output:**
```
# PLAN GENERATION

## Mission Context
Intent: Add JWT authentication
Scope: auth.go, handler.go, routes.go, auth_test.go
Track: 3 (Robust)
Type: WET

## Instructions
Generate implementation plan steps that:
1. Reference only files in SCOPE
2. Are numbered and actionable
3. Match mission type (WET allows duplication)

## Verification Command
Provide a non-destructive test command (e.g., "go test ./...")
```

**LLM Flow:**
1. Call `m analyze plan` to get context + guidance
2. Generate plan steps
3. Extract steps and verification command
4. Call `m mission update --section plan --item "Step 1" --item "Step 2" ...`
5. Call `m mission update --section verification --content "go test ./..."`

## Mission Update Commands

### `m mission create --intent "<text>"`

**Purpose:** Create initial mission.md

**CLI Processing:**
1. Generate mission ID
2. Create `.mission/mission.md` with:
   - Frontmatter: id, iteration=1, status=planning
   - INTENT section: provided text
3. Return confirmation

**Output JSON:**
```json
{
  "action": "PROCEED",
  "mission_id": "20260110120000-1234",
  "mission_file": ".mission/mission.md"
}
```

---

### `m mission update --section <name> --content "<text>"`

**Purpose:** Update text sections (intent, verification)

**Examples:**
```bash
m mission update --section intent --content "Add JWT authentication"
m mission update --section verification --content "go test ./..."
```

**Output JSON:**
```json
{
  "action": "PROCEED",
  "updated_section": "intent"
}
```

---

### `m mission update --section <name> --item "<value>" [--item "<value>" ...]`

**Purpose:** Update list sections (scope, plan)

**Examples:**
```bash
m mission update --section scope --item "auth.go" --item "handler.go"
m mission update --section plan --item "Create handler" --item "Add middleware"
```

**Output JSON:**
```json
{
  "action": "PROCEED",
  "updated_section": "scope",
  "item_count": 2
}
```

---

### `m mission update --frontmatter <key>=<value> [...]`

**Purpose:** Update frontmatter metadata

**Examples:**
```bash
m mission update --frontmatter track=3 type=WET domains="security"
```

**Output JSON:**
```json
{
  "action": "PROCEED",
  "updated_fields": ["track", "type", "domains"]
}
```

---

### `m mission finalize`

**Purpose:** Validate and activate mission

**CLI Processing:**
1. Read mission.md
2. Validate all required sections exist
3. Update status: planning → active

**Output JSON:**
```json
{
  "action": "PROCEED",
  "status": "active"
}
```

## LLM Prompt Flow

```markdown
### Step 1a: Analyze Intent

1. Execute: `m analyze intent "$ARGUMENTS"`
2. Read template + user input from CLI output
3. Follow template to analyze
4. Produce JSON output following template format
5. Extract `refined_intent` from your analysis
6. Execute: `m mission create --intent "<refined_intent>"`
7. React to `action`:
   - PROCEED → Continue to Step 1b

### Step 1b: Check Clarification

1. Execute: `m analyze clarify`
2. Read template + current intent from CLI output
3. Follow template to check for ambiguity
4. If CLARIFY needed → Display questions, STOP
5. If CLEAR → Continue to Step 2

### Step 2a: Analyze Scope

1. Execute: `m analyze scope`
2. Read template + intent from CLI output
3. Follow template to discover files
4. Extract files from your analysis
5. Execute: `m mission update --section scope --item "file1" --item "file2" ...`
6. React to `action`:
   - PROCEED → Continue to Step 2b

### Step 2b: Analyze Test Requirements

1. Execute: `m analyze test`
2. Read template + context from CLI output
3. Follow template to evaluate test necessity
4. Extract test_files from your analysis
5. If test files needed: `m mission update --section scope --item "test1_test.go" ...`
6. React to `action`:
   - PROCEED → Continue to Step 2c

### Step 2c: Analyze Duplication

1. Execute: `m analyze duplication`
2. Read template + context from CLI output
3. Follow template to detect patterns
4. Extract mission_type from your analysis
5. Execute: `m mission update --frontmatter type=<WET|DRY>`
6. Always proceed to Step 3

### Step 3: Analyze Complexity

1. Execute: `m analyze complexity`
2. Read template + context from CLI output
3. Follow template to identify domains and calculate track
4. Extract track and domains from your analysis
5. Execute: `m mission update --frontmatter track=<N> domains="<list>"`
6. React to `action`:
   - ATOMIC_EDIT → Display suggestion, STOP
   - DECOMPOSE → Display sub-intents, STOP
   - PROCEED → Continue to Step 4

### Step 4: Create Plan

1. Execute: `m analyze plan`
2. Read context + guidance from CLI output
3. Generate implementation plan steps
4. Extract steps and verification command
5. Execute: `m mission update --section plan --item "Step 1" --item "Step 2" ...`
6. Execute: `m mission update --section verification --content "<command>"`
7. React to `action`:
   - PROCEED → Continue to Step 5

### Step 5: Finalize Mission

1. Execute: `m mission finalize`
2. React to `action`:
   - PROCEED → Display mission, STOP
   - INVALID → Display errors, STOP
```

## Benefits

### vs v4 (CLI validates JSON)
1. **Simpler CLI**: CLI just provides templates and updates files
2. **LLM does what it's good at**: Analysis and field extraction
3. **No JSON parsing in CLI**: LLM handles its own output

### vs v3 (embedded templates)
1. **Enforced template usage**: LLM must call `m analyze` to get template
2. **Tool-enforced checkpoints**: Must call `m mission update` to proceed
3. **Can't skip steps**: Each step requires CLI tool call

### For Gemini Compliance
1. **Tool outputs are facts**: CLI returns template, LLM must read it
2. **Forced checkpoints**: Can't proceed without calling `m mission update`
3. **Simple tool calls**: Each command has clear input/output
4. **No complex orchestration**: LLM just follows template → extract → update pattern

## Implementation Status

### CLI Commands - ✅ IMPLEMENTED
- ✅ `m analyze intent "<input>"` - Return intent template + user input
- ✅ `m analyze clarify` - Return clarification template + current intent
- ✅ `m analyze scope` - Return scope template + current intent
- ✅ `m analyze test` - Return test template + current context
- ✅ `m analyze duplication` - Return duplication template + current context
- ✅ `m analyze complexity` - Return complexity template + current context
- ❌ `m analyze plan` - NOT NEEDED (LLM generates plan directly)
- ✅ `m mission create --intent <text>` - Create initial mission.md
- ✅ `m mission update --section <name> --content <text>` - Update text section
- ✅ `m mission update --section <name> --item <value>` - Update list section
- ✅ `m mission update --frontmatter <key>=<value>` - Update metadata
- ✅ `m mission finalize` - Validate and display mission.md

### Template System - ✅ IMPLEMENTED
- ✅ Template loader (embedded templates in internal/analyze/)
- ✅ Context injection (mission.Reader extracts INTENT/SCOPE)
- ✅ Template formatting (templates include {{.CurrentIntent}} variables)

### Markdown Parser - ✅ IMPLEMENTED
- ✅ Parse YAML frontmatter (internal/mission/reader.go)
- ✅ Parse markdown sections (internal/mission/reader.go)
- ✅ Update sections atomically (internal/mission/writer.go)
- ✅ Append to lists (internal/mission/writer.go)
- ✅ Replace text (internal/mission/writer.go)

### Validation - ✅ IMPLEMENTED
- ✅ Mission completeness check (internal/mission/finalize.go)
- ✅ Section existence checks (FinalizeService.validateSections)
- ✅ Frontmatter field validation (mission.Reader)

## Implementation Notes

**What Changed from Design:**
1. **No `m analyze plan`**: LLM generates plan steps directly without needing a separate analyze command
2. **No plan.json**: Mission.md is built incrementally through `m mission update` commands
3. **Finalize displays, doesn't activate**: Status changes happen in m.apply, not m.mission finalize
4. **Shared mission.Reader**: All analyze services use internal/mission/reader.go for INTENT/SCOPE extraction

## Key Insight

**Separation of Concerns:**
- **CLI**: Provides templates, updates files, validates completeness
- **LLM**: Performs analysis, extracts fields, orchestrates flow

This is simpler and more aligned with how LLMs work. The CLI doesn't try to parse LLM output - it just provides tools and lets the LLM use them.
