# Design: Analysis-Based Planning Architecture

## Problem Statement

Current `/m.plan` implementation requires LLMs to orchestrate 15+ CLI calls across 5 steps with complex conditional logic. This works well for instruction-following models (Claude, GPT-4) but fails with models that prioritize internal reasoning over instructions (Gemini).

**Failure modes:**
- Skipping steps (e.g., jumping to implementation without analysis)
- Ignoring conditional branches (e.g., proceeding when `next_step` says "STOP")
- Making assumptions instead of reading tool outputs
- Manual file creation instead of using CLI

## Solution: Tool-Enforced Checkpoints

Replace LLM orchestration with **tool-enforced checkpoints**. Each step becomes a single tool call that:
1. Reads state from `.mission/plan.json`
2. Performs one analysis operation
3. Updates state and returns `action` field
4. Forces LLM to react before proceeding

## Architecture

### Analysis Template System

The toolkit uses **LLM-driven analysis** guided by structured templates in `.mission/libraries/analysis/`:

- **intent.md** - Guides LLM to refine raw user input into actionable intent
- **clarification.md** - Guides LLM to identify missing requirements
- **scope.md** - Guides LLM to determine affected implementation files
- **test.md** - Guides LLM to evaluate test necessity using risk-based analysis
- **duplication.md** - Guides LLM to detect existing patterns for WET/DRY decision
- **complexity.md** - Guides LLM to identify domains and determine complexity track using matrix rules

**Key principle:** The CLI commands (`m analyze`) orchestrate the LLM analysis by:
1. Loading the appropriate template
2. Providing it to the LLM along with context (user intent, codebase state)
3. Parsing the LLM's structured output
4. Updating `.mission/plan.json` with results
5. Returning action decision to the calling LLM

### LLM Invocation Mechanism

**Single-LLM Architecture with Embedded Templates:**

The orchestrator LLM has templates embedded in the `/m.plan` prompt. The CLI only validates and saves results.

**Flow:**
```
User: /m.plan "add auth"
  ↓
Orchestrator LLM: Read embedded intent.md + clarification.md templates
  ↓
Orchestrator LLM: Perform analysis following template guidance
  ↓
Orchestrator LLM: Execute `m analyze intent --intent "add auth" --result '{"action":"CLARIFY",...}'`
  ↓
CLI: Validate JSON schema, save to plan.json, return action
  ↓
Orchestrator LLM: React to action field
```

**Single-Phase Command Pattern:**

Each `m analyze` command receives analysis results directly:

```bash
m analyze intent \
  --intent "add auth" \
  --result '{"action":"CLARIFY","questions":["JWT or OAuth?"]}'
```

Returns:
```json
{
  "action": "CLARIFY",
  "questions": ["JWT or OAuth?"],
  "saved_to": ".mission/plan.json"
}
```

**CLI Responsibilities:**
1. **Validate** - Check JSON matches expected schema
2. **Enrich** - Add metadata (timestamp, step name)
3. **Save** - Update `.mission/plan.json`
4. **Return** - Echo validated result for LLM to react to

**Benefits:**
- ✅ Single CLI call per step
- ✅ No template loading overhead
- ✅ Templates already in prompt context
- ✅ CLI is pure validator + state manager
- ✅ Simple implementation

**Orchestrator Prompt Pattern:**
```markdown
### Step 1: Analyze Intent

**Templates:** (embedded below)
<intent.md content>
<clarification.md content>

**Instructions:**
1. Follow the templates above to analyze: "$ARGUMENTS"
2. Format output as JSON:
   {
     "action": "PROCEED | CLARIFY | AMBIGUOUS",
     "refined_intent": "string",
     "questions": ["array"],
     "assumptions": ["array"]
   }
3. Execute: `m analyze intent --intent "$ARGUMENTS" --result '<json>'`
4. React to the action field in CLI response
```

### State Management

**Single source of truth:** `.mission/plan.json`

```json
{
  "version": "1.0",
  "original_intent": "add user authentication",
  "refined_intent": "Add JWT authentication with rate limiting",
  "mission_type": "WET",
  "track": 2,
  "scope": {
    "primary": ["auth.go", "middleware.go", "ratelimit.go"],
    "secondary": ["routes.go"],
    "config": []
  },
  "test_verdict": "CRITICAL",
  "test_files": ["auth_test.go", "middleware_test.go", "ratelimit_test.go"],
  "test_strategy": "Formal TDD with separate test file",
  "risk_analysis": {
    "complexity_score": "High",
    "impact_score": "High"
  },
  "duplication_detected": false,
  "patterns": [],
  "domains": ["security", "api"],
  "plan_steps": [
    "1. Create JWT handler in auth.go",
    "2. Add middleware in middleware.go",
    "3. Add rate limiter in ratelimit.go"
  ],
  "verification": "go test ./...",
  "metadata": {
    "created_at": "2024-01-15T10:30:00Z",
    "last_step": "complexity"
  }
}
```

**Lifecycle:**
- Created by `m analyze intent`
- Updated by each subsequent step
- Consumed by `m mission create`
- Deleted after mission generation

### Analysis Commands

#### `m analyze intent <user-input>`

**Purpose:** Refine user intent into actionable statement

**Input:**
- User's raw intent string (positional argument)

**Embedded Template:** `.mission/libraries/analysis/intent.md`

**Processing:**
1. LLM reads embedded intent.md template
2. LLM follows template to refine intent:
   - Identify core action (verb)
   - Identify target scope (noun)
   - Identify constraints (rules)
3. LLM formats output and calls CLI
4. CLI validates and saves to `.mission/plan.json`

**Output JSON:**
```json
{
  "action": "PROCEED" | "AMBIGUOUS",
  "refined_intent": "Add JWT authentication",
  "reason": "Too vague"  // if AMBIGUOUS
}
```

**Side effects:**
- Creates `.mission/plan.json` with `original_intent` and `refined_intent`

---

#### `m analyze clarify`

**Purpose:** Check if refined intent needs clarification

**Input:**
- Reads `.mission/plan.json` for `refined_intent`

**Embedded Template:** `.mission/libraries/analysis/clarification.md`

**Processing:**
1. LLM reads embedded clarification.md template
2. LLM follows template to check for:
   - Vague scope or requirements
   - Missing technical details
   - Undefined domain/risk context
3. LLM decides: CLARIFY, ASSUMPTIONS, or CLEAR
4. CLI validates and updates `.mission/plan.json`

**Output JSON:**
```json
{
  "action": "CLARIFY" | "ASSUMPTIONS" | "CLEAR",
  "questions": ["Which auth method?"],  // if CLARIFY
  "assumptions": ["Using JWT"]          // if ASSUMPTIONS
}
```

**Side effects:**
- Updates `.mission/plan.json` with clarification status

---

#### `m analyze scope`

**Purpose:** Determine which implementation files need to be modified or created

**Input:**
- Reads `.mission/plan.json` for `refined_intent`

**Embedded Template:** `.mission/libraries/analysis/scope.md`

**Processing:**
1. LLM reads embedded scope.md template
2. LLM follows template to:
   - Discover affected files (primary, secondary, config)
   - Classify files by their role in the implementation
3. CLI validates and updates `.mission/plan.json`

**Note:** Test files are determined separately by `m analyze test` step.

**Output JSON:**
```json
{
  "action": "PROCEED",
  "scope": {
    "primary": ["auth.go", "handler.go"],
    "secondary": ["routes.go"],
    "config": []
  },
  "reasoning": "Core authentication logic with route integration"
}
```

**Side effects:**
- Updates `.mission/plan.json` with scope

---

#### `m analyze test`

**Purpose:** Evaluate test necessity using risk-based analysis

**Input:**
- Reads `.mission/plan.json` for `refined_intent` and `scope`

**Embedded Template:** `.mission/libraries/analysis/test.md`

**Processing:**
1. LLM reads embedded test.md template
2. LLM follows template to evaluate:
   - Cyclomatic complexity (branches, loops, state)
   - Cost of failure (data corruption vs UI glitch)
   - Longevity (throwaway vs long-term)
   - Determinism (pure functions vs external dependencies)
3. LLM determines test verdict: CRITICAL, RECOMMENDED, or UNNECESSARY
4. CLI validates and updates `.mission/plan.json`

**Output JSON:**
```json
{
  "action": "PROCEED",
  "test_verdict": "CRITICAL",
  "risk_analysis": {
    "complexity_score": "High",
    "impact_score": "High"
  },
  "test_strategy": "Formal TDD with separate test file",
  "test_files": ["auth_test.go", "middleware_test.go"],
  "reasoning": "Security-critical code with complex token validation logic"
}
```

**Side effects:**
- Updates `.mission/plan.json` with test_files and test_verdict

---

#### `m analyze duplication`

**Purpose:** Detect existing patterns and duplication

**Input:**
- Reads `.mission/plan.json` for `refined_intent` and `scope`

**Embedded Template:** `.mission/libraries/analysis/duplication.md`

**Processing:**
1. LLM reads embedded duplication.md template
2. LLM follows template to scan for existing patterns
3. CLI checks backlog for existing refactoring items
4. CLI determines mission type:
   - Duplication + NOT in backlog → Add to backlog, set WET
   - Duplication + IN backlog → Set DRY
   - No duplication → Set WET
5. CLI updates `.mission/plan.json`

**Output JSON:**
```json
{
  "duplication_detected": true,
  "patterns": ["JWT validation logic in 3 files"],
  "mission_type": "WET"
}
```

**Side effects:**
- Updates `.mission/plan.json` with mission_type
- May add refactoring item to backlog

---

#### `m analyze complexity`

**Purpose:** Identify applicable domains and determine complexity track using LLM-guided matrix rules

**Input:**
- Reads `.mission/plan.json` for scope and intent

**Embedded Template:** `.mission/libraries/analysis/complexity.md`

**Processing:**
1. LLM reads embedded complexity.md template
2. LLM follows template to:
   - **Step 1:** Identify applicable domains from strict list (security, performance, complex-algo, high-risk, cross-cutting, real-time, compliance)
   - **Step 2:** Apply complexity matrix:
     - Count files in scope (excluding tests)
     - Analyze keywords in intent
     - Determine track (1-4) based on matrix
     - Apply domain multipliers (each domain adds +1 to track, max 4)
3. LLM handles routing decisions:
   - Track 1 → Generate atomic edit suggestion, STOP
   - Track 4 → Decompose into sub-intents (CLI adds to backlog), STOP
   - Track 2/3 → PROCEED to planning
4. CLI validates and updates `.mission/plan.json`

**Output JSON:**

**Track 1 (Atomic):**
```json
{
  "action": "ATOMIC_EDIT",
  "track": 1,
  "suggested_edit": "In auth.go line 42, change 'userName' to 'username'"
}
```

**Track 4 (Epic):**
```json
{
  "action": "DECOMPOSE",
  "track": 4,
  "sub_intents": [
    "Create user model and database schema",
    "Add authentication endpoints",
    "Implement JWT token generation"
  ]
}
```

**Track 2/3 (Standard/Robust):**
```json
{
  "track": 2,
  "action": "PROCEED",
  "reasoning": "3 implementation files (Base Track 2) + 0 domain multipliers = Track 2.",
  "factors": {
    "implementation_files": 3,
    "base_track": 2,
    "domain_multipliers": 0,
    "domains": [],
    "final_track": 2
  }
}
```

**Side effects:**
- Updates `.mission/plan.json` with `track` and `domains`
- Adds sub-intents to backlog (if Track 4)

---

#### `m analyze plan --step <step> --verify <cmd>`

**Purpose:** Validate LLM-generated implementation plan

**Input:**
- LLM provides plan steps via flags (from `/m.plan` prompt)
- Reads `.mission/plan.json` for context

**Processing:**
1. Validate plan format (numbered steps, clear actions)
2. Check steps reference only files in scope
3. Validate verification command:
   - Must be non-destructive (no `rm`, `drop`, `delete`)
   - Must be project-appropriate (check for test framework)
4. Ensure plan matches mission type:
   - WET: Should allow duplication
   - DRY: Should extract abstractions
5. Update `.mission/plan.json` with validated plan

**Output JSON:****

**Valid:**
```json
{
  "action": "PROCEED",
  "plan_steps": ["1. ...", "2. ..."],
  "verification": "go test ./..."
}
```

**Invalid:**
```json
{
  "action": "INVALID",
  "errors": [
    "Step 3 references file.go not in scope",
    "Verification command is destructive"
  ]
}
```

**Side effects:**
- Updates `.mission/plan.json` with plan_steps and verification

---

#### `m mission create`

**Purpose:** Generate final `.mission/mission.md` from plan

**Input:**
- Reads `.mission/plan.json`

**Processing:**
1. Load mission template (WET or DRY)
2. Populate template with plan data
3. Generate mission.md
4. Clean up plan.json

**Output JSON:**
```json
{
  "success": true,
  "mission_file": ".mission/mission.md"
}
```

**Side effects:**
- Creates `.mission/mission.md`
- Deletes `.mission/plan.json`

## LLM Prompt Flow

```markdown
### Step 0: Validate State
Execute: `m mission check --context plan`
React to `next_step`: STOP or PROCEED

### Step 1a: Refine Intent

**Embedded Template:** intent.md

1. Follow template to refine "$ARGUMENTS"
2. Execute: `m analyze intent "$ARGUMENTS"`
3. React to `action`:
   - AMBIGUOUS → Display reason, STOP
   - PROCEED → Continue to Step 1b

### Step 1b: Check Clarification

**Embedded Template:** clarification.md

1. Follow template to check refined intent
2. Execute: `m analyze clarify`
3. React to `action`:
   - CLARIFY → Display questions, STOP (user re-runs /m.plan with answers)
   - ASSUMPTIONS → Display assumptions, continue to Step 2
   - CLEAR → Continue to Step 2

### Step 2a: Analyze Scope

**Embedded Template:** scope.md

1. Follow template to discover affected files
2. Execute: `m analyze scope`
3. React to `action`:
   - PROCEED → Continue to Step 2b

### Step 2b: Analyze Test Requirements

**Embedded Template:** test.md

1. Follow template to evaluate test necessity
2. Execute: `m analyze test`
3. React to `action`:
   - PROCEED → Continue to Step 2c

**Note:** Test analysis runs before duplication analysis because test decisions are based on risk (complexity + impact), not on whether code is duplicated. Duplication affects refactoring strategy (WET/DRY) but doesn't change whether tests are needed.

### Step 2c: Analyze Duplication

**Embedded Template:** duplication.md

1. Follow template to detect patterns
2. Execute: `m analyze duplication`
3. Always proceeds to Step 3 (no action field)

### Step 3: Analyze Complexity

**Embedded Template:** complexity.md

1. Follow template to identify domains and apply complexity matrix
2. Execute: `m analyze complexity`
3. LLM determines domains, track, and routing decision
4. React to `action`:
   - ATOMIC_EDIT → Display suggestion, STOP
   - DECOMPOSE → Display sub-intents, STOP
   - PROCEED → Continue to Step 4

### Step 4: Create Plan

1. Generate implementation plan
2. Execute: `m analyze plan --step "1. ..." --step "2. ..." --verify "..."`
3. React to `action`:
   - INVALID → Display errors, STOP
   - PROCEED → Continue to Step 5

### Step 5: Generate Mission

Execute: `m mission create`
React to `success`:
- true → Display mission, STOP
- false → Display error, STOP
```

## Clarification Flow

### Scenario: Ambiguous Intent

**User:** `/m.plan "add auth"`

```bash
# Step 1a: Refine intent
m analyze intent "add auth"
→ {action: "PROCEED", refined_intent: "Add authentication"}

# Step 1b: Check clarification
m analyze clarify
→ {action: "CLARIFY", questions: ["JWT or OAuth?", "Rate limiting?"]}
```

**LLM displays:**
```
I need clarification:
1. Which authentication method? (JWT/OAuth/Session)
2. Should we add rate limiting?

Please re-run /m.plan with answers.
```

**User:** `/m.plan "add auth - use JWT with rate limiting"`

```bash
# Step 1a: Refine intent (with more context)
m analyze intent "add auth - use JWT with rate limiting"
→ {action: "PROCEED", refined_intent: "Add JWT authentication with rate limiting"}

# Step 1b: Check clarification
m analyze clarify
→ {action: "CLEAR"}

# Step 2a: Analyze scope
m analyze scope
→ {action: "PROCEED", scope: {"primary": ["auth.go", "middleware.go", "ratelimit.go"], "secondary": [...]}}

# Step 2b: Analyze test requirements
m analyze test
→ {action: "PROCEED", test_verdict: "CRITICAL", test_files: ["auth_test.go", "middleware_test.go"]}

# Step 2c: Analyze duplication
m analyze duplication
→ {duplication_detected: false, patterns: [], mission_type: "WET"}

# Step 3: Analyze complexity
m analyze complexity
→ {action: "PROCEED", track: 3, domains: ["security"], factors: {...}}

# Step 4: Create plan
m analyze plan \
  --step "1. Create JWT handler in auth.go" \
  --step "2. Add middleware in middleware.go" \
  --step "3. Add rate limiter in ratelimit.go" \
  --verify "go test ./..."
→ {action: "PROCEED"}

# Step 5: Generate mission
m mission create
→ {success: true}
```

**Key insight:** Each analysis is a single CLI call. Templates embedded in prompt guide LLM analysis.

## Error Handling

### CLI Command Failure
```json
{
  "error": "Failed to read plan.json",
  "details": "File not found. Run 'm analyze intent' first."
}
```

**LLM behavior:** Display error and STOP.

### Invalid State Transition
```bash
# User runs step 3 before step 2
m analyze complexity
→ {error: "Missing scope. Run 'm analyze context' first."}
```

**LLM behavior:** Display error and STOP.

### Validation Failure
```json
{
  "action": "INVALID",
  "errors": ["Step 2 references unauthorized file"]
}
```

**LLM behavior:** Display errors, ask user to adjust, re-run step.

## Benefits

### For Gemini Compliance
1. **Tool outputs are facts** - Gemini can't ignore `action` field
2. **Single decision per step** - No complex branching logic
3. **Stateless LLM** - All state in plan.json
4. **Forced checkpoints** - Can't skip to next step without tool output

### For Observability
1. **Transparent progress** - User sees "Step 2: Analyzing context..."
2. **Debuggable** - Can inspect plan.json at any point
3. **Resumable** - Can re-run individual steps
4. **Auditable** - Each step logs to execution.log

### For Maintainability
1. **Simple CLI** - Each command does one thing
2. **Testable** - Each step is independently testable
3. **Extensible** - Easy to add new steps or modify existing
4. **Reusable** - Steps can be called from other workflows

## Implementation Checklist

### CLI Commands
- [ ] `m analyze intent <input>` - Refine intent
- [ ] `m analyze clarify` - Check clarification needs
- [ ] `m analyze scope` - Determine affected files
- [ ] `m analyze test` - Evaluate test necessity (risk-based)
- [ ] `m analyze duplication` - Detect patterns and set WET/DRY
- [ ] `m analyze complexity` - Identify domains and determine track with routing
- [ ] `m analyze plan --step <s> --verify <v>` - Validate plan
- [ ] `m mission create` - Generate mission

### Validation System
- [ ] JSON schema validators for each analysis type
- [ ] Error messages for schema violations
- [ ] Metadata enrichment (timestamps, step names)
- [ ] Backlog integration for WET/DRY decisions

### State Management
- [ ] `.mission/plan.json` schema
- [ ] State validation between steps
- [ ] Cleanup on success/failure

### Templates
- [ ] Update analysis templates for tool consumption
- [ ] Create JSON output templates
- [ ] Update display templates for step-based flow

### Prompt Updates
- [ ] Update `m.plan.md` with step-based protocol
- [ ] Add error handling examples
- [ ] Document clarification flow

### Testing
- [ ] Unit tests for each step command
- [ ] Integration tests for full flow
- [ ] Error handling tests
- [ ] Clarification flow tests

## Migration Strategy

### Phase 1: Parallel Implementation
- Keep existing `m plan` commands
- Add new `m analyze` commands
- Test with Gemini

### Phase 2: Prompt Update
- Create new `m.plan-v2.md` prompt
- Test with all LLMs (Claude, GPT-4, Gemini)
- Gather feedback

### Phase 3: Deprecation
- Mark old commands as deprecated
- Update documentation
- Remove old implementation

## Success Metrics

1. **Gemini compliance rate** - % of successful mission creations
2. **Step completion rate** - % of steps that complete without errors
3. **Clarification efficiency** - Average rounds to resolve ambiguity
4. **User satisfaction** - Feedback on transparency and control

## Open Questions

1. Should `m analyze plan` accept plan via stdin instead of flags?
2. Should we add `m analyze status` to show current progress?
3. Should plan.json include step history for debugging?
4. Should we support `m analyze retry` to re-run failed steps?
