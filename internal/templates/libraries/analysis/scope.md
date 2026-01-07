# SCOPE ANALYSIS TEMPLATE

## Purpose
Determine which files need to be modified or created based on the user's intent, including whether test files should be included.

## Analysis Steps

### 1. File Discovery
- **Keyword Matching**: Use file search tool to find files matching intent keywords (e.g., "auth" → `auth.go`, `authentication.js`)
- **Dependency Analysis**: Check import/dependency graphs for affected modules
- **Related Files**: Identify files that interact with the target (e.g., handlers, models, services)

### 2. File Classification
Categorize each identified file:
- **Primary**: Files that directly implement the intent
- **Secondary**: Files that need updates due to dependencies (e.g., interfaces, types)
- **Test**: Test files for the primary and secondary files
- **Config**: Configuration files (only if explicitly mentioned in intent)

### 3. Test Inclusion Decision
Based on governance.md testability rules:

**INCLUDE test files if:**
- New logic or business rules are being added
- Bug fixes that need regression prevention
- Critical paths (authentication, payments, data integrity)
- Changes to public APIs or interfaces

**EXCLUDE test files if:**
- Trivial changes (typo fixes, formatting, comments)
- Documentation-only updates
- Configuration changes without logic
- Refactoring with existing comprehensive tests
- Framework-heavy code (e.g., CLI definitions, UI wiring) lacking custom business logic
- Trivial boilerplate or wiring code

## Output Format

Produce a structured list of files grouped by category:

```json
{
  "primary": [
    "path/to/main/file.go",
    "path/to/handler.go"
  ],
  "secondary": [
    "path/to/interface.go"
  ],
  "tests": [
    "path/to/main/file_test.go",
    "path/to/handler_test.go"
  ],
  "config": [],
  "reasoning": "Brief explanation of scope decisions and test inclusion/exclusion"
}
```

## Examples

### Example 1: New Feature
**Intent**: "Add user registration endpoint"
```json
{
  "primary": [
    "handlers/user.go",
    "models/user.go"
  ],
  "secondary": [
    "routes/api.go"
  ],
  "tests": [
    "handlers/user_test.go",
    "models/user_test.go"
  ],
  "config": [],
  "reasoning": "New business logic requires tests for registration validation and database operations"
}
```

### Example 2: Bug Fix
**Intent**: "Fix null pointer error in login handler"
```json
{
  "primary": [
    "handlers/auth.go"
  ],
  "secondary": [],
  "tests": [
    "handlers/auth_test.go"
  ],
  "config": [],
  "reasoning": "Bug fix in critical path (authentication) requires regression test"
}
```

### Example 3: Trivial Change
**Intent**: "Fix typo in error message"
```json
{
  "primary": [
    "handlers/user.go"
  ],
  "secondary": [],
  "tests": [],
  "config": [],
  "reasoning": "Trivial change to error message text does not require test updates"
}
```

### Example 4: Refactoring
**Intent**: "Extract database connection logic into utility"
```json
{
  "primary": [
    "utils/database.go"
  ],
  "secondary": [
    "handlers/user.go",
    "handlers/product.go"
  ],
  "tests": [
    "utils/database_test.go"
  ],
  "config": [],
  "reasoning": "New utility requires tests; existing handler tests cover integration"
}
```

### Example 5: Framework Wiring
**Intent**: "Add new flag to CLI command"
```json
{
  "primary": [
    "cmd/init.go"
  ],
  "secondary": [],
  "tests": [],
  "config": [],
  "reasoning": "Change is limited to framework configuration/wiring; no custom business logic to test"
}
```
