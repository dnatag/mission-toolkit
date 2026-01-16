# Catalyst Toolkit User Guide

This guide provides a comprehensive overview of the Catalyst Toolkit workflows and commands.

## Quick Reference

**AI Rule**: Never edit JSON files. Use these commands instead.

### Step 1: Clarification Mission

```bash
# Create plan.json with questions
m plan init --intent "Add user authentication" \
  --question "Which authentication method?" \
  --question "Should we support social login?"

# Generate clarification mission.md
m mission create --type clarification --file .mission/plan.json
```

### Step 3: Draft Spec

```bash
# Create plan.json with scope
m plan init --intent "Add JWT authentication" \
  --type WET \
  --scope internal/auth/jwt.go \
  --scope internal/auth/jwt_test.go \
  --domain security
```

### Step 4: Complexity Analysis

```bash
# Analyze and auto-update track
m plan analyze --file .mission/plan.json --update
```

**Output**: JSON with track, next_step, warnings

### Step 5: Validation

```bash
# Validate scope and verification
m plan validate --file .mission/plan.json
```

**Output**: JSON with valid status, errors, warnings

### Step 6: Final Mission

```bash
# Add plan steps and verification
m plan update --plan "Create JWT generation function" \
  --plan "Add validation middleware" \
  --plan "Write comprehensive tests" \
  --plan "Note: Allow duplication (WET principle)" \
  --verification "go test ./internal/auth/..."

# Generate final mission.md
m mission create --type final --file .mission/plan.json
```

### Command Summary

| Command | Purpose | Key Flags |
|---------|---------|-----------|
| `m plan init` | Create plan.json | `--intent`, `--type`, `--scope`, `--domain`, `--question` |
| `m plan update` | Update plan.json | `--plan`, `--verification` |
| `m plan analyze` | Determine track | `--file`, `--update` |
| `m plan validate` | Validate plan | `--file` |
| `m mission create` | Generate mission.md | `--type`, `--file` |

### Repeatable Flags

These flags can be used multiple times:
- `--scope`: Add multiple files
- `--domain`: Add multiple domains
- `--question`: Add multiple questions
- `--plan`: Add multiple plan steps

**Example**:
```bash
m plan init --intent "..." \
  --scope file1.go \
  --scope file2.go \
  --scope file3.go \
  --domain security \
  --domain performance
```

### plan.json Evolution

#### After Step 1 (Clarification)
```json
{
  "intent": "...",
  "clarification_questions": ["Q1?", "Q2?"]
}
```

#### After Step 3 (Draft)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"]
}
```

#### After Step 4 (Analysis)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"],
  "track": "TRACK 2"
}
```

#### After Step 6 (Finalize)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"],
  "track": "TRACK 2",
  "plan": ["Step 1", "Step 2"],
  "verification": "go test ./..."
}
```

### Error Handling

```bash
# Missing required flag
$ m plan init
Error: required flag "intent" not set

# Invalid mission type
$ m mission create --type invalid --file plan.json
Error: invalid type: invalid (use clarification or final)

# File not found
$ m plan analyze --file missing.json
Error: open missing.json: no such file or directory
```

### Tips

1. **Always use `--update`** in Step 4 to auto-add track
2. **Use backslash `\`** for multi-line commands
3. **Quote strings** with spaces: `--intent "Add user auth"`
4. **Check JSON output** from analyze/validate commands
5. **Inspect plan.json** at any point for debugging

---

## Mission Lifecycle

The Mission Toolkit workflow ensures you maintain ownership while leveraging AI capabilities through a series of handshakes and reviews.

### Complete Flow Diagram

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

### Step-by-Step Breakdown

#### 1. m.plan - Intent Analysis

**Purpose**: Convert user intent into structured mission

**Process**:
- Analyze intent for clarity and complexity
- If ambiguous → route to m.clarify
- If clear → create mission.md with INTENT, SCOPE, PLAN, VERIFICATION
- **Test Strategy**: AI evaluates if test files should be included based on value (new logic, bugs) vs. noise (trivial changes).

**Output**: `.mission/mission.md` ready for review

#### 2. m.clarify - Clarification (Optional)

**Purpose**: Refine ambiguous intents

**Process**:
- Ask targeted questions to clarify requirements
- Gather missing details
- Re-run m.plan with clarified intent

**Output**: Updated mission.md with refined details

#### 3. Review mission.md - Human Authorization

**Purpose**: You authorize the architecture before execution

**What to Review**:
- **INTENT**: Does it match what you want?
- **SCOPE**: Are the right files included?
- **PLAN**: Does the approach make sense?
- **VERIFICATION**: Will this prove it works?

**Decision**: Approve or request changes

#### 4. m.apply - Execution

**Purpose**: Implement the authorized plan

**Process** (automatic):
1. Execute implementation steps
2. Run polish pass (code quality improvements)
   - **Quality-Driven Testing**: Review newly created tests. Ensure they are high-value (meaningful happy paths, critical edge cases) and robust. Remove low-value tests (e.g., trivial getters/setters).
3. Generate conventional commit message
4. Update mission.md with commit message

**Output**: Working code + commit message ready for review

#### 5. Review Code - Human Verification

**Purpose**: Verify the implementation meets requirements

**What to Review**:
- Does code match the PLAN?
- Does verification pass?
- Any bugs or improvements needed?

**Decision**: Accept or request adjustments

#### 6. Adjustments - Iterative Refinement (Optional)

**Purpose**: Fix issues discovered during review

**Process**:
- User describes needed changes
- AI makes fixes
- AI regenerates commit message to reflect changes
- Loop back to m.complete

**Output**: Refined code + updated commit message

#### 7. m.complete - Archival & Commit

**Purpose**: Capture learnings and create git commit

**Process**:
- Archive mission files to `.mission/completed/`
- Update metrics and backlog
- Create git commit using stored commit message
- Clean up active mission files

**Output**: Git commit + archived mission data

### Key Principles

#### Human-in-the-Loop
- **Before execution**: Review and approve the plan
- **After execution**: Review and verify the code
- **Optional refinement**: Request changes if needed

#### Atomic Scope
- Each mission is small enough to comprehend
- Changes stay within human understanding limits
- You maintain ownership, not just contribution

#### Continuous Learning
- Every mission archived with full context
- Metrics tracked for process improvement
- Patterns detected for future optimization

#### Quality-Driven Testing
- **Scope-Matched**: Test coverage must match mission complexity.
- **High-Value**: Focus on meaningful happy paths and critical edge cases.
- **Essential**: Avoid testing for the sake of testing (e.g., trivial getters/setters).
- **Quality**: Tests must be robust, readable, and maintainable.

### Common Paths

#### Happy Path (No Issues)
```
m.plan → Review → m.apply → Review → m.complete
```

#### With Clarification
```
m.plan → m.clarify → m.plan → Review → m.apply → Review → m.complete
```

#### With Adjustments
```
m.plan → Review → m.apply → Review → Adjustments → m.complete
```

#### Multiple Adjustment Cycles
```
m.plan → Review → m.apply → Review → Adjustments → Review → Adjustments → m.complete
```
