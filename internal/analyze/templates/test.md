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

**Output:** List test files to add to scope or create.

---

**Optional: If test files are needed, consider:**
- Test approach: table-driven, mocks, fixtures
- Verification command: `go test ./...`, `npm test`, etc.
