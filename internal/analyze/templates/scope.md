# SCOPE ANALYSIS TEMPLATE

## Current Intent
{{.CurrentIntent}}

## Purpose
Determine which implementation files need to be modified or created based on the user's intent.

## Analysis Steps

### 1. File Discovery
- **Keyword Matching**: Use file search tool to find files matching intent keywords (e.g., "auth" → `auth.go`, `authentication.js`)
- **Dependency Analysis**: Check import/dependency graphs for affected modules
- **Related Files**: Identify files that interact with the target (e.g., handlers, models, services)

### 2. File Classification
Categorize each identified file:
- **Primary**: Files that directly implement the intent
- **Secondary**: Files that need updates due to dependencies (e.g., interfaces, types)
- **Config**: Configuration files (only if explicitly mentioned in intent)

**Note**: Test files are determined by a separate test analysis step.

## Output Format

Produce a JSON object with action and scope details.

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

## Examples

### Example 1: New Feature
**Intent**: "Add user registration endpoint"
```json
{
  "action": "PROCEED",
  "scope": {
    "primary": [
      "handlers/user.go",
      "models/user.go"
    ],
    "secondary": [
      "routes/api.go"
    ],
    "config": []
  },
  "reasoning": "New endpoint requires handler and model, with route registration"
}
```

### Example 2: Bug Fix
**Intent**: "Fix null pointer error in login handler"
```json
{
  "action": "PROCEED",
  "scope": {
    "primary": [
      "handlers/auth.go"
    ],
    "secondary": [],
    "config": []
  },
  "reasoning": "Bug fix isolated to authentication handler"
}
```

### Example 3: Refactoring
**Intent**: "Extract database connection logic into utility"
```json
{
  "action": "PROCEED",
  "scope": {
    "primary": [
      "utils/database.go"
    ],
    "secondary": [
      "handlers/user.go",
      "handlers/product.go"
    ],
    "config": []
  },
  "reasoning": "New utility with updates to existing handlers that use database connections"
}
```
