# COMPLEXITY ANALYSIS TEMPLATE

## Purpose
Identify applicable technical/business domains and determine the complexity track of the mission based on scope, intent, and domain context. This guides the planning strategy (Atomic vs. Standard vs. Epic).

## Step 1: Identify Domains

Analyze the intent and scope to identify applicable domains. Select from this **strict list** (do not invent new ones):

### Valid Domains

**Security** (`security`)
- Triggers: Authentication, authorization, cryptography, PII handling, secrets management, input sanitization
- Examples: "Add login", "Encrypt password", "Fix SQL injection"

**Performance** (`performance`)
- Triggers: Latency requirements, throughput optimization, memory management, caching, database indexing, concurrency
- Examples: "Speed up API", "Reduce memory usage", "Add Redis cache"

**Complex Algorithms** (`complex-algo`)
- Triggers: Mathematical models, AI/ML, recursion, state machines, graph algorithms, custom data structures
- Examples: "Implement recommendation engine", "Pathfinding logic", "Parser implementation"

**High Risk** (`high-risk`)
- Triggers: Financial transactions, payments, data deletion (bulk), critical infrastructure, public API changes
- Examples: "Process refund", "Delete user account", "Change API signature"

**Cross-Cutting** (`cross-cutting`)
- Triggers: Changes affecting multiple distinct modules, logging infrastructure, configuration management, error handling strategies
- Examples: "Update logging format everywhere", "Refactor config loading", "Global error handler"

**Real-Time** (`real-time`)
- Triggers: WebSockets, streaming, event-driven architecture, polling
- Examples: "Live chat", "Stock ticker", "Notification stream"

**Compliance** (`compliance`)
- Triggers: GDPR, audit logs, legal requirements, accessibility (WCAG)
- Examples: "Add consent banner", "Export user data", "Audit trail"

**Decision Logic:**
- Default: If none apply, domains list is empty `[]`
- Multiple: Select ALL that apply (e.g., `["security", "high-risk"]`)
- Threshold: If unsure, err on the side of caution and include the domain

## Step 2: Determine Complexity Track

### Complexity Tracks

### Track 1: Atomic (Trivial)
**Characteristics:**
- 0-1 implementation files (excluding tests/docs)
- Trivial changes: typos, formatting, comments, simple renames
- No logic changes or new functionality
- **Action**: Suggest direct edit, no mission needed

### Track 2: Standard (Routine)
**Characteristics:**
- 2-5 implementation files
- Standard patterns: CRUD operations, new endpoint, add field
- Single module/package scope
- Low-medium risk
- **Action**: Create standard mission with step-by-step plan

### Track 3: Robust (Complex)
**Characteristics:**
- 6-9 implementation files OR
- 2-5 files + critical domain (security, performance, etc.)
- Cross-module changes
- High risk or complex orchestration
- **Action**: Create robust mission with detailed verification

### Track 4: Epic (Architectural)
**Characteristics:**
- 10+ implementation files OR
- Massive scope (entire subsystem, major refactor)
- Too large for single mission
- **Action**: Decompose into sub-intents, add to backlog, STOP

## Analysis Factors

### 1. Implementation Files
Count files in scope, EXCLUDING:
- Test files (`*_test.go`, `*.test.js`, etc.)
- Documentation (`*.md`, `*.txt`)

### 2. Domain Multipliers
Add +1 to the track (max Track 4) for EACH domain identified in Step 1.
All domains are considered critical and act as multipliers.

### 3. Calculation Logic
1. **Count Implementation Files**:
   - Exclude: test files, docs, config-only files
   - Count: source code files that contain logic

2. **Determine Base Track**:
   - 0-1 files → Track 1
   - 2-5 files → Track 2
   - 6-9 files → Track 3
   - 10+ files → Track 4

3. **Apply Domain Multipliers**:
   - If Base Track is 4 → Final Track is 4 (no adjustment)
   - Else:
     - Apply special case: If Track 1 + any domain → Upgrade to Track 2 first
     - Then: Final Track = Base Track + count of domains from Step 1
     - Cap Final Track at 4

## Output Format

Produce a JSON object with track, action, and reasoning.

```json
{
  "track": 1 | 2 | 3 | 4,
  "action": "ATOMIC_EDIT" | "PROCEED" | "DECOMPOSE",
  "reasoning": "Explanation of track determination",
  "factors": {
    "implementation_files": 3,
    "base_track": 2,
    "domain_multipliers": 1,
    "domains": ["security"],
    "final_track": 3
  },
  "suggested_edit": "In auth.go line 42, change 'userName' to 'username'"  // Only if action is ATOMIC_EDIT
}
```

**Action Mapping:**
- Track 1 → `ATOMIC_EDIT` (provide suggested_edit for LLM to display)
- Track 2-3 → `PROCEED` (continue to planning)
- Track 4 → `DECOMPOSE` (CLI will handle decomposition)

**Examples:**

**Case 1:**
Intent: "Fix typo in README"
Scope: ["README.md"]
Domains: []
Result:
```json
{
  "track": 1,
  "action": "ATOMIC_EDIT",
  "reasoning": "0 implementation files (README.md excluded). Trivial documentation change.",
  "factors": {
    "implementation_files": 0,
    "base_track": 1,
    "domain_multipliers": 0,
    "domains": [],
    "final_track": 1
  },
  "suggested_edit": "In README.md line 15, change 'authentification' to 'authentication'"
}
```

**Case 2:**
Intent: "Add JWT authentication"
Scope: ["auth.go", "middleware.go", "main.go"]
Domains: ["security"]
Result:
```json
{
  "track": 3,
  "action": "PROCEED",
  "reasoning": "3 implementation files (Base Track 2) + 1 domain multiplier (security) = Track 3.",
  "factors": {
    "implementation_files": 3,
    "base_track": 2,
    "domain_multipliers": 1,
    "domains": ["security"],
    "final_track": 3
  }
}
```

**Case 3:**
Intent: "Rewrite entire payment system"
Scope: [12 files]
Domains: ["security", "compliance"]
Result:
```json
{
  "track": 4,
  "action": "DECOMPOSE",
  "reasoning": "12 implementation files exceeds Track 4 threshold (10+). Requires decomposition.",
  "factors": {
    "implementation_files": 12,
    "base_track": 4,
    "domain_multipliers": 2,
    "domains": ["security", "compliance"],
    "final_track": 4
  }
}
```
