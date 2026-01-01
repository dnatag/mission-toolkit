# DUPLICATION ANALYSIS TEMPLATE

## Purpose
Identify existing code patterns that are similar to the requested feature to prevent redundancy and promote DRY (Don't Repeat Yourself) principles.

## Analysis Steps

### 1. Semantic Search
- **File Names**: Are there existing files with similar names? (e.g., `auth.go` vs `authentication.go`)
- **Function Names**: Are there functions doing similar tasks? (e.g., `ValidateUser` vs `CheckUser`)
- **Business Logic**: Is there logic that solves the same problem? (e.g., two different email validation regexes)

### 2. Pattern Recognition
- **Boilerplate**: Is this a standard pattern (e.g., CRUD) that exists elsewhere?
- **Utilities**: Can existing utility functions be reused instead of rewritten?
- **Configuration**: Can existing config structs be reused?

## Decision Logic & Output

Based on your analysis, choose ONE of the following outputs.

### A. No Duplication Found
If no significant duplication or reusable patterns are found:
```json
{
  "status": "none",
  "confidence": "high",
  "recommendation": "Proceed with new implementation."
}
```

### B. Reusable Code Found
If existing functions or utilities can be reused:
```json
{
  "status": "reusable_code",
  "confidence": "high",
  "recommendation": "Incorporate existing utilities into the plan.",
  "details": [
    {
      "file": "utils/validation.go",
      "symbol": "ValidateEmail",
      "action": "Reuse this function for email validation."
    }
  ]
}
```

### C. Refactoring Opportunity Found
If similar but not identical logic exists, suggesting a refactor:
```json
{
  "status": "refactor_opportunity",
  "confidence": "medium",
  "recommendation": "Add a step to the plan to extract common logic.",
  "details": [
    {
      "file": "users/handler.go",
      "symbol": "CreateUser",
      "action": "Extract database connection logic into a shared utility before implementing CreateProduct."
    }
  ]
}
```

### D. Exact Match Found
If the requested feature already exists:
```json
{
  "status": "exact_match",
  "confidence": "high",
  "recommendation": "STOP. Inform the user that the feature already exists.",
  "details": [
    {
      "file": "auth/jwt.go",
      "symbol": "GenerateJWT",
      "action": "The requested JWT generation function already exists."
    }
  ]
}
```
