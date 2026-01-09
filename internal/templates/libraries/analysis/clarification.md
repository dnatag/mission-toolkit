# CLARIFICATION ANALYSIS TEMPLATE

## Purpose
Scan user intent for ambiguous requirements that need clarification before mission planning.

## Analysis Categories

### 1. Scope & Requirements
- **Vague Problem**: "bug", "issue", "broken" without symptoms?
- **Ambiguous Scope**: "improve", "enhance" without specifics?
- **Missing Tech**: Unspecified framework, library, or database?

### 2. Domain & Risk Context (Critical for Plan Spec)
- **Security**: Does this involve auth, crypto, or PII?
- **Performance**: Are there latency or throughput constraints?
- **Compliance**: Does this touch GDPR, audit logs, or financial data?
- **Complexity**: Does this involve complex algorithms or AI?

*If these are relevant but undefined, you may need to clarify.*

## Output Format

Produce a JSON object with action and details.

```json
{
  "action": "CLARIFY" | "ASSUMPTIONS" | "CLEAR",
  "questions": ["Which auth method?"],  // Only if CLARIFY
  "assumptions": ["Using JWT"]          // Only if ASSUMPTIONS
}
```

**Examples:**

### A. Clarification Needed
```json
{
  "action": "CLARIFY",
  "questions": [
    "Which authentication method? (JWT/OAuth/Session)",
    "Should we add rate limiting?"
  ]
}
```

### B. Proceed with Assumptions
```json
{
  "action": "ASSUMPTIONS",
  "assumptions": [
    "Using JWT based on existing auth.go file",
    "No rate limiting required (not mentioned)"
  ]
}
```

### C. Intent Clear
```json
{
  "action": "CLEAR"
}
```

## Question Quality Standards
- **Specific**: Ask for concrete details (e.g., "Which error code?" not "What's wrong?").
- **Actionable**: The answer should directly affect the implementation plan.
- **Contextual**: Reference the user's original goal.
