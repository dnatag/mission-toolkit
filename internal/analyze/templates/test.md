# Test Analysis Template

## Current Mission Context

**Intent:** {{.CurrentIntent}}

**Files in Scope:**
{{.CurrentScope}}

## Test Analysis Instructions

### Step 1: Analyze Intent and Existing Code

**Read the Intent:**
- What functionality is being added or modified?
- What is the expected behavior?
- What are the inputs and outputs?

**Read Existing Scope Files (if they exist):**
- Use file read tool to examine each scope file
- What existing functions/types are present?
- What patterns and conventions are established?
- What dependencies and error handling exist?

**For New Files:**
- What functionality will this file provide based on intent?
- What similar files exist that show the expected pattern?

### Step 2: Predict Implementation Characteristics

Based on intent and existing code patterns:
- What logic complexity is expected? (simple | moderate | complex)
- What external dependencies will be involved? (filesystem, network, database, other packages)
- What error conditions are likely? (file not found, invalid input, network failure)
- What edge cases should be handled? (empty input, nil values, boundary conditions)

### Step 3: Categorize Test Necessity

**REQUIRES TESTS (High Priority):**
- Business logic with conditional branches or loops
- Functions that handle errors or validate input
- Code that transforms, parses, or processes data
- Integration with external dependencies (use mocks/fakes)
- Bug fixes (regression tests to prevent recurrence)
- Critical paths that affect system behavior

**SKIP TESTS (Low Value):**
- Simple getters/setters with no logic
- Trivial constructors that only assign fields
- Pass-through functions with no transformation
- Pure delegation to other tested functions

### Step 4: Design Test Cases

For each scope file requiring tests:

**File: [filename]**
- **Expected Function/Method:** [name based on intent]
- **Why Test?** [business logic | error handling | data transformation | integration | bug fix]
- **Test Cases:**
  1. [Happy path: expected input → expected output]
  2. [Edge case: boundary condition → expected behavior]
  3. [Error case: invalid input → expected error]
- **Test Data/Mocks:** [fixtures, mock dependencies, sample inputs]
- **Assertions:** [what to verify for correctness]

### Step 5: Test File Strategy

- Which test files need creation? (e.g., `filename_test.go`)
- What test helpers or fixtures are needed?
- What mocking strategy? (interfaces for dependencies, in-memory filesystems, test doubles)
- How to structure tests? (table-driven, subtests, setup/teardown)

### Step 6: Verification Approach

- What verification command will prove tests work? (e.g., `go test ./...`)
- What is minimum acceptable coverage for new logic?
- Are integration tests needed beyond unit tests?
