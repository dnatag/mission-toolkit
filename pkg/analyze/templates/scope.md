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

### 2. Breaking Change Analysis (Call Site Analysis)
**CRITICAL**: If the intent involves modifying existing public APIs, methods, functions, or interfaces:
1. **Identify Symbols**: List the specific functions, methods, or types being changed.
2. **Find Usages**: Use `find_usages` (preferred) or `grep` to locate all call sites.
   - *Tip*: When using grep, verify matches to avoid false positives (e.g., same method name on different struct).
3. **Include Callers**: Add files containing valid call sites to the **Secondary** scope.
4. **Check Implementations**: If changing an interface, find all structs implementing it.

**Triggers for call site analysis:**
- Method signature changes (parameters, return types)
- Function renaming or removal
- Interface contract changes
- Behavior changes that affect callers

**Example search patterns:**
- Go: `\.MethodName\(` or `FunctionName\(`
- JavaScript: `\.methodName\(` or `functionName\(`
- Python: `\.method_name\(` or `function_name\(`

### 3. File Classification
Categorize each identified file:
- **Primary**: Files that directly implement the intent (the definition of the logic/API).
- **Secondary**: Files that need updates due to dependencies (callers, interface implementations, types).
- **Config**: Configuration files (only if explicitly mentioned in intent).

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

### Example 4: Method Signature Change (Breaking Change)
**Intent**: "Refactor Reader methods to remove path parameter and use the path from the Reader struct instead"
```json
{
  "action": "PROCEED",
  "scope": {
    "primary": [
      "internal/mission/reader.go"
    ],
    "secondary": [
      "internal/mission/writer.go",
      "internal/mission/finalize.go",
      "internal/mission/check.go",
      "internal/mission/archiver.go",
      "internal/analyze/test.go",
      "internal/analyze/complexity.go",
      "internal/tui/update.go"
    ],
    "config": []
  },
  "reasoning": "Method signature change requires updating all call sites found via grep search for .Read(, .GetMissionID(, etc."
}
```
