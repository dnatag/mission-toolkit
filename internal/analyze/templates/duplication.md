# DUPLICATION ANALYSIS TEMPLATE

## Current Intent

{{.CurrentIntent}}

## Purpose
Identify existing code patterns similar to the requested feature. This informs the WET→DRY workflow decision.

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

### 3. Decision Logic

- **No duplication** → Report `duplication_detected: false`, empty patterns
- **Duplication found** → Report `duplication_detected: true`, list specific patterns with file locations
- **Mission type** → Always set to `"WET"` (CLI will check backlog and may override to `"DRY"`)

## Output Format

Produce a JSON object with duplication analysis.

**Note**: This analysis does not include an `action` field because it always proceeds to the next step. The CLI handles mission type determination (WET/DRY) based on backlog state.

```json
{
  "duplication_detected": true | false,
  "patterns": ["Description with file locations"],
  "mission_type": "WET"
}
```

## Examples

### Example 1: No Duplication
```json
{
  "duplication_detected": false,
  "patterns": [],
  "mission_type": "WET"
}
```

### Example 2: Duplication Found
```json
{
  "duplication_detected": true,
  "patterns": [
    "Email validation regex exists in users/handler.go:45 and products/handler.go:78",
    "Database connection logic duplicated in auth/db.go, users/db.go, and orders/db.go"
  ],
  "mission_type": "WET"
}
```
