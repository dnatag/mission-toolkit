# TEST ANALYSIS TEMPLATE

## Purpose
Determine if automated tests are necessary, valuable, or wasteful for the requested coding task.

## Analysis Criteria

Evaluate based on "Test Value" metrics:

### 1. Cyclomatic Complexity
- Does the logic contain multiple branches, loops, or state changes?
- Simple linear code vs complex decision trees

### 2. Cost of Failure
- Data corruption or system crash → High cost
- Minor UI glitch or cosmetic issue → Low cost
- Security vulnerability → Critical cost

### 3. Longevity
- Throwaway script or prototype → Low value
- Long-term infrastructure or library → High value
- Core business logic → High value

### 4. Determinism
- Pure functions (easy to test) → High value
- Heavy external dependencies (hard/flaky to test) → Low value
- UI-heavy code → Low value

## Risk Scoring

**Complexity Score:**
- **Low**: Linear logic, no branches, simple data transformation
- **Medium**: Some conditionals, basic error handling, moderate logic
- **High**: Multiple branches, loops, state management, complex algorithms

**Impact Score:**
- **Low**: Cosmetic changes, non-critical features, easily reversible
- **Medium**: User-facing features, data processing, API endpoints
- **High**: Authentication, payments, data integrity, security, core infrastructure

## Output Format

Produce a JSON object with test verdict and reasoning.

```json
{
  "action": "PROCEED",
  "test_verdict": "CRITICAL" | "RECOMMENDED" | "UNNECESSARY",
  "risk_analysis": {
    "complexity_score": "Low" | "Medium" | "High",
    "impact_score": "Low" | "Medium" | "High"
  },
  "test_strategy": "Formal TDD with separate test file" | "Inline verification" | "Manual/Visual testing",
  "test_files": ["auth_test.go"],  // Only if CRITICAL or RECOMMENDED
  "reasoning": "Explanation of test decision"
}
```

## Decision Matrix

**🔴 CRITICAL (Formal TDD):**
- Complexity: High OR Impact: High
- Requires separate test file (e.g., `_test.go`, `.test.ts`)
- Examples: Authentication logic, payment processing, data validation

**🟡 RECOMMENDED (Inline Verification):**
- Complexity: Medium AND Impact: Medium
- No separate file needed, but include inline assertions
- Examples: Utility functions, data transformers, API clients

**🟢 UNNECESSARY (Manual/Visual):**
- Complexity: Low AND Impact: Low
- Automated tests would be over-engineering
- Examples: Pure config, CSS, simple glue code, framework wiring

## Examples

### Example 1: Authentication Logic (CRITICAL)
**Intent**: "Add JWT authentication"
**Scope**: ["auth.go", "middleware.go"]
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
  "reasoning": "Security-critical code with complex token validation logic. Failure could expose system to unauthorized access."
}
```

### Example 2: Data Transformer (RECOMMENDED)
**Intent**: "Add function to convert user data to CSV"
**Scope**: ["utils/export.go"]
```json
{
  "action": "PROCEED",
  "test_verdict": "RECOMMENDED",
  "risk_analysis": {
    "complexity_score": "Medium",
    "impact_score": "Medium"
  },
  "test_strategy": "Inline verification with example data",
  "test_files": ["utils/export_test.go"],
  "reasoning": "Moderate complexity with data transformation logic. Tests ensure correct formatting but failure is non-critical."
}
```

### Example 3: CSS Styling (UNNECESSARY)
**Intent**: "Update button colors to match brand"
**Scope**: ["styles/button.css"]
```json
{
  "action": "PROCEED",
  "test_verdict": "UNNECESSARY",
  "risk_analysis": {
    "complexity_score": "Low",
    "impact_score": "Low"
  },
  "test_strategy": "Manual/Visual testing",
  "test_files": [],
  "reasoning": "Pure styling change with no logic. Visual inspection is sufficient. Automated tests would be over-engineering."
}
```

### Example 4: CLI Flag Addition (UNNECESSARY)
**Intent**: "Add --verbose flag to CLI"
**Scope**: ["cmd/root.go"]
```json
{
  "action": "PROCEED",
  "test_verdict": "UNNECESSARY",
  "risk_analysis": {
    "complexity_score": "Low",
    "impact_score": "Low"
  },
  "test_strategy": "Manual/Visual testing",
  "test_files": [],
  "reasoning": "Framework wiring with no custom logic. Cobra framework handles flag parsing. Manual testing sufficient."
}
```

### Example 5: Bug Fix (CRITICAL)
**Intent**: "Fix null pointer error in login handler"
**Scope**: ["handlers/auth.go"]
```json
{
  "action": "PROCEED",
  "test_verdict": "CRITICAL",
  "risk_analysis": {
    "complexity_score": "Medium",
    "impact_score": "High"
  },
  "test_strategy": "Formal TDD with regression test",
  "test_files": ["handlers/auth_test.go"],
  "reasoning": "Bug fix in critical authentication path. Regression test prevents reintroduction of the bug."
}
```
