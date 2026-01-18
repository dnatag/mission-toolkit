# DUPLICATION ANALYSIS TEMPLATE

## Current Intent

{{.CurrentIntent}}

## Purpose
Identify existing code patterns similar to the requested feature. This informs the WET→DRY workflow decision using Rule-of-Three tracking.

## Analysis Steps

### 1. Semantic Search
Scan the codebase for similar implementations:

- **File Names**: Existing files with similar names (e.g., `auth.go` vs `authentication.go`)
- **Function Names**: Functions doing similar tasks (e.g., `ValidateUser` vs `CheckUser`)
- **Business Logic**: Logic solving the same problem (e.g., multiple email validation implementations)

### 2. Pattern Recognition
Identify reusable patterns:

- **Boilerplate**: Standard patterns (CRUD, handlers) that exist elsewhere
- **Utilities**: Existing utility functions that could be reused
- **Configuration**: Existing config structures that could be extended

### 3. Pattern ID Generation
When duplication is found, generate a stable pattern ID:
- Use lowercase kebab-case (e.g., `email-validation`, `db-connection`)
- Be specific enough to identify the pattern uniquely
- Be general enough to match future occurrences

## Output Format

Produce a JSON object with duplication analysis.

**Note**: The CLI tracks pattern counts in backlog. When count reaches 3, mission type becomes DRY.

```json
{
  "duplication_detected": true | false,
  "patterns": [
    {
      "id": "pattern-id",
      "description": "Human-readable description",
      "locations": ["file1.go:45", "file2.go:78"]
    }
  ]
}
```

## Examples

### Example 1: No Duplication
```json
{
  "duplication_detected": false,
  "patterns": []
}
```

### Example 2: Duplication Found
```json
{
  "duplication_detected": true,
  "patterns": [
    {
      "id": "email-validation",
      "description": "Email validation regex duplicated across handlers",
      "locations": ["users/handler.go:45", "products/handler.go:78"]
    },
    {
      "id": "db-connection",
      "description": "Database connection logic duplicated",
      "locations": ["auth/db.go:12", "users/db.go:15", "orders/db.go:18"]
    }
  ]
}
```
