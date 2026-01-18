# COMPLEXITY ANALYSIS TEMPLATE (IMPROVED)

## Current Mission Context

**Intent:** {{.CurrentIntent}}

**Scope:**
{{.CurrentScope}}

## Purpose
Determine mission complexity using a multi-factor scoring system that considers file count, domain criticality, and change characteristics.

## Step 1: Count Files with Weights

Count ALL files in scope with different weights:

**Implementation Files (1.0x weight):**
- Source code files with business logic
- Examples: `*.go`, `*.js`, `*.py`, `*.java` (non-test)

**Test Files (0.5x weight):**
- Unit tests: `*_test.go`, `*.test.js`, `*.spec.ts`
- Integration tests: `*_integration_test.go`
- E2E tests: `*.e2e.js`

**Documentation/Config (0.25x weight):**
- Documentation: `*.md`, `*.txt`
- Config files: `*.yaml`, `*.json`, `*.toml` (if they contain logic/complexity)

**Excluded (0x weight):**
- Pure config with no logic (simple key-value pairs)
- Generated files
- Vendored dependencies

**Calculation:**
```
Weighted File Count = (Implementation Files × 1.0) + (Test Files × 0.5) + (Docs/Config × 0.25)
```

**File Count Score:**
- 0-0.9 weighted files → 0 points (Track 1 candidate)
- 1.0-3.9 weighted files → 1 point
- 4.0-6.9 weighted files → 2 points
- 7.0-12.4 weighted files → 3 points
- 12.5+ weighted files → 5 points

## Step 2: Identify Critical Domains

Select ALL applicable domains and sum their weights:

### High-Impact Domains (2 points each)
- **Security** (`security`) - Auth, crypto, PII, secrets, input sanitization
- **High-Risk** (`high-risk`) - Payments, data deletion, critical infrastructure
- **Complex-Algo** (`complex-algo`) - AI/ML, graph algorithms, custom data structures

### Medium-Impact Domains (1 point each)
- **Performance** (`performance`) - Latency, caching, concurrency
- **Cross-Cutting** (`cross-cutting`) - Multi-module changes, global infrastructure
- **Compliance** (`compliance`) - GDPR, audit logs, accessibility

### Low-Impact Domains (0.5 points each)
- **Real-Time** (`real-time`) - WebSockets, streaming, events
- **Standard** (`standard`) - CRUD, simple business logic (default if none apply)

## Step 3: Identify Change Characteristics

Add points for complexity indicators:

**+1 point each:**
- Semantic breaking changes (behavior changes, removed features, API contract changes)
- Data migration required (schema changes, data transformation)
- External integrations (third-party APIs, webhooks)
- State management changes (session, cache, database state)

**+0.5 points each:**
- Mechanical breaking changes (parameter reordering, removing redundant params, renaming)
- New dependencies added
- Configuration changes required
- Multiple environments affected (dev, staging, prod)

## Step 4: Calculate Track (MANDATORY ARITHMETIC)

**REQUIRED:** You MUST show this exact calculation format:

```
Files: [impl_count] × 1.0 + [test_count] × 0.5 + [doc_count] × 0.25 = [weighted_total] = [file_points] pts
Domains: [domain_list] = [domain_points] pts
Characteristics: [char_list] = [char_points] pts
TOTAL: [file_points] + [domain_points] + [char_points] = [final_score] pts → Track [N]
```

**Track Mapping (Preliminary):**
- 0 pts + single file → Track 1
- 1-2 pts → Track 2
- 3-4 pts → Track 3
- 5+ pts → Track 4
- Multiple files → Minimum Track 2

*Note: Apply Special Rules below to determine Final Track.*

**Special Rules (Overrides Preliminary Track):**
- **Trivial Single File**: If Total Score = 0 AND single file AND trivial intent → Track 1 with suggested edit
- **Test-Only Changes**: If scope is ONLY test files AND no breaking changes → Cap at Track 3
- **Simple Mechanical Refactor Cap**: If ALL criteria below are met → Cap at Track 3 (overrides score-based Track 4):
  - Simple mechanical pattern (constructor signature, field addition, renaming, method signature)
  - Consistent pattern applied across multiple files
  - NO business logic changes
  - NO new dependencies or external integrations
  - NO critical domains (security, high-risk, complex-algo)

## Output Format

```json
{
  "track": 1 | 2 | 3 | 4,
  "action": "ATOMIC_EDIT" | "PROCEED" | "DECOMPOSE",
  "reasoning": "Explanation of track determination",
  "scoring": {
    "file_count": {
      "files": 3,
      "score": 1
    },
    "domains": {
      "identified": ["security"],
      "score": 2
    },
    "characteristics": {
      "identified": ["breaking_changes"],
      "score": 1
    },
    "total_score": 4,
    "final_track": 3
  },
  "special_rule_applied": "trivial_single_file" | "test_only_changes" | null,  // Optional: indicate if special rule was applied
  "suggested_edit": "..."  // Only if action is ATOMIC_EDIT
}
```

## Examples

**Case 1: Fix Typo (Single File)**
Intent: "Fix typo in user.go"
Scope: ["user.go"]
```json
{
  "track": 1,
  "action": "ATOMIC_EDIT",
  "reasoning": "1.0 weighted file (single impl), trivial change, 0 points → Track 1",
  "scoring": {
    "file_count": {
      "implementation": 1,
      "test": 0,
      "docs": 0,
      "weighted_total": 1.0,
      "score": 0
    },
    "domains": {"identified": [], "score": 0},
    "characteristics": {"identified": [], "score": 0},
    "total_score": 0,
    "final_track": 1
  },
  "suggested_edit": "In user.go line 42, change 'userName' to 'username'"
}
```

**Case 2: Fix Single Test**
Intent: "Fix flaky test in user_test.go"
Scope: ["user_test.go"]
```json
{
  "track": 1,
  "action": "ATOMIC_EDIT",
  "reasoning": "0.5 weighted file (single test), 0 points → Track 1",
  "scoring": {
    "file_count": {
      "implementation": 0,
      "test": 1,
      "docs": 0,
      "weighted_total": 0.5,
      "score": 0
    },
    "domains": {"identified": [], "score": 0},
    "characteristics": {"identified": [], "score": 0},
    "total_score": 0,
    "final_track": 1
  },
  "suggested_edit": "In user_test.go line 15, add time.Sleep(10ms) before assertion"
}
```

**Case 3: Add Validation + Test (Multiple Files)**
Intent: "Add email validation"
Scope: ["user.go", "user_test.go"]
```json
{
  "track": 2,
  "action": "PROCEED",
  "reasoning": "1.5 weighted files (1 impl + 1 test), multiple files → minimum Track 2. Score: 1pt → Track 2",
  "scoring": {
    "file_count": {
      "implementation": 1,
      "test": 1,
      "docs": 0,
      "weighted_total": 1.5,
      "score": 1
    },
    "domains": {"identified": ["standard"], "score": 0},
    "characteristics": {"identified": [], "score": 0},
    "total_score": 1,
    "final_track": 2
  }
}
```

**Case 4: Refactor 3 Tests**
Intent: "Refactor test helpers"
Scope: ["helper_test.go", "setup_test.go", "fixtures_test.go"]
```json
{
  "track": 2,
  "action": "PROCEED",
  "reasoning": "1.5 weighted files (3 tests × 0.5), multiple files → minimum Track 2. Score: 1pt → Track 2",
  "scoring": {
    "file_count": {
      "implementation": 0,
      "test": 3,
      "docs": 0,
      "weighted_total": 1.5,
      "score": 1
    },
    "domains": {"identified": ["standard"], "score": 0},
    "characteristics": {"identified": [], "score": 0},
    "total_score": 1,
    "final_track": 2
  }
}
```

**Case 5: Payment System Rewrite**
Intent: "Rewrite payment processing"
Scope: [8 implementation files, 4 test files, 1 config]
```json
{
  "track": 4,
  "action": "DECOMPOSE",
  "reasoning": "10.25 weighted files (8×1.0 + 4×0.5 + 1×0.25) = 3pts + security (2pts) + high-risk (2pts) + compliance (1pt) + data migration (1pt) = 9pts → Track 4",
  "scoring": {
    "file_count": {
      "implementation": 8,
      "test": 4,
      "docs": 1,
      "weighted_total": 10.25,
      "score": 3
    },
    "domains": {"identified": ["security", "high-risk", "compliance"], "score": 5},
    "characteristics": {"identified": ["data_migration", "external_integrations"], "score": 2},
    "total_score": 10,
    "final_track": 4
  }
}
```

**Case 6: Simple Mechanical Refactor (Many Files)**
Intent: "Add path field to Reader struct and update NewReader constructor to accept path parameter"
Scope: [13 implementation files, 4 test files]
```json
{
  "track": 3,
  "action": "PROCEED",
  "reasoning": "15.0 weighted files (13×1.0 + 4×0.5) would normally be Track 4, but this qualifies for SIMPLE MECHANICAL REFACTOR CAP: constructor signature update with consistent pattern, no logic changes, no new dependencies, no critical domains → Capped at Track 3",
  "scoring": {
    "file_count": {
      "implementation": 13,
      "test": 4,
      "docs": 0,
      "weighted_total": 15.0,
      "score": 5
    },
    "domains": {"identified": [], "score": 0},
    "characteristics": {"identified": ["breaking_changes"], "score": 1},
    "total_score": 6,
    "final_track": 3,
    "special_rule_applied": "simple_mechanical_refactor_cap"
  }
}
```
