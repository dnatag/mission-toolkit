# Test Analysis Template

## Current Mission Context

**Intent:** {{.CurrentIntent}}

**Files in Scope:**
{{.CurrentScope}}

## Test Analysis Instructions

### Step 1: Check for Existing Test Files (MANDATORY)

**CRITICAL:** You MUST check for test files corresponding to ALL implementation files in scope.

For EACH implementation file in scope (exclude test files, docs, configs):

1. **Derive test filename by language:**
   - **Go**: `X.go` → `X_test.go`
   - **Python**: `X.py` → `test_X.py` or `X_test.py`
   - **JavaScript/TypeScript**: `X.js` → `X.test.js` or `X.spec.js`
   - **Java**: `X.java` → `XTest.java` (in `src/test/java/`)
   - **Rust**: `#[cfg(test)]` inline modules or `tests/X.rs`
   - **Other**: Search for files containing "test" and the base filename

2. **MANDATORY CHECK:** Use file read tool to verify if test file exists
   - **Report result:** "✅ EXISTS" or "❌ NOT FOUND" for each file
   - **No skipping:** You must check ALL implementation files, even if you plan to exclude them later

**Example Output Format:**
```
cmd/mission.go → cmd/mission_test.go: ❌ NOT FOUND
internal/mission/writer.go → internal/mission/writer_test.go: ✅ EXISTS
```

### Step 2: Analyze the Changes

**Read existing scope files** to understand what's being modified:
- What functionality is being added/changed?
- What's the logic complexity? (simple | moderate | complex)
- What dependencies are involved? (filesystem, network, database)

### Step 3: Decide What Needs Testing

For EACH implementation file:

**Decision Matrix:**

| Change Type | Test File Exists? | Action |
|-------------|------------------|--------|
| **Substantive** | Yes | Add test file to scope |
| **Substantive** | No | Create new test file |
| **Trivial** | Yes or No | Skip test file |

**Substantive Changes:**
- Logic with branches/loops
- Error handling or validation  
- Data transformation/parsing
- External dependencies
- Bug fixes

**Trivial Changes:**
- Simple getters/setters
- Field additions
- Pass-through functions

## Output Format

**REQUIRED:** Produce a JSON object with your analysis results.

```json
{
  "action": "ADD_TESTS" | "SKIP_TESTS",
  "analysis": [
    {
      "file": "pkg/checkpoint/service.go",
      "test_file": "pkg/checkpoint/service_test.go",
      "exists": true,
      "change_type": "substantive",
      "decision": "add_to_scope",
      "reason": "Consolidate() signature change requires test updates"
    }
  ],
  "test_files_to_add": ["pkg/checkpoint/service_test.go"]
}
```

### Field Definitions

- **action**: `ADD_TESTS` if any test files need to be added to scope, `SKIP_TESTS` if none needed
- **analysis**: Array with one entry per implementation file in scope
  - **file**: The implementation file being analyzed
  - **test_file**: Derived test filename (or "N/A" for non-code files)
  - **exists**: Whether the test file exists (`true`/`false`/`"N/A"`)
  - **change_type**: `"substantive"` or `"trivial"` based on decision matrix
  - **decision**: `"add_to_scope"`, `"create_new"`, or `"skip"`
  - **reason**: Brief justification for the decision
- **test_files_to_add**: Array of test file paths to add to mission scope (empty if none)

### Examples

**Example 1: Substantive change with existing test**
```json
{
  "action": "ADD_TESTS",
  "analysis": [
    {
      "file": "pkg/git/client.go",
      "test_file": "pkg/git/client_test.go",
      "exists": false,
      "change_type": "substantive",
      "decision": "skip",
      "reason": "Interface-only file, implementations tested separately"
    },
    {
      "file": "pkg/checkpoint/service.go",
      "test_file": "pkg/checkpoint/service_test.go",
      "exists": true,
      "change_type": "substantive",
      "decision": "add_to_scope",
      "reason": "Return type change affects existing test assertions"
    }
  ],
  "test_files_to_add": ["pkg/checkpoint/service_test.go"]
}
```

**Example 2: No tests needed**
```json
{
  "action": "SKIP_TESTS",
  "analysis": [
    {
      "file": "pkg/templates/displays/success.md",
      "test_file": "N/A",
      "exists": "N/A",
      "change_type": "trivial",
      "decision": "skip",
      "reason": "Markdown template, not executable code"
    }
  ],
  "test_files_to_add": []
}
```
