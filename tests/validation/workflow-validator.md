# WORKFLOW VALIDATOR

## PURPOSE
AI agent-based validation of Mission Toolkit prompt workflows using dry-run test fixtures.

## HOW TO RUN

### Quick Start
```bash
# Run all workflow tests
cd tests/
# Execute each test fixture with AI agent in dry-run mode
```

### Manual Test Execution
1. **Choose a test fixture**: `tests/fixtures/m-plan-test.md`
2. **Load the fixture**: Read INPUT and EXPECTED_OUTPUT sections
3. **Execute prompt**: Run corresponding @m.* command with INPUT in dry-run mode
4. **Compare results**: Validate actual output matches EXPECTED_OUTPUT
5. **Check assertions**: Verify all STEP_ASSERTIONS pass

### Example: Running M.PLAN Test
```bash
# 1. Load test fixture
cat tests/fixtures/m-plan-test.md

# 2. Extract input: "add user authentication to the API"

# 3. Execute @m.plan in dry-run mode with the input

# 4. Compare generated mission with EXPECTED_MISSION_OUTPUT

# 5. Validate all STEP_ASSERTIONS pass
```

### Automated Validation
```bash
# Run validation script (if implemented)
./tests/run-workflow-tests.sh

# Or execute via AI agent with validation instructions
# See VALIDATION PROCESS section below
```

## VALIDATION PROCESS

### 1. Test Execution
For each test fixture in `tests/fixtures/`:
1. Load test fixture (m-plan-test.md, m-clarify-test.md, etc.)
2. Extract INPUT and EXPECTED_OUTPUT sections
3. Execute corresponding prompt in dry-run mode
4. Compare actual output with expected output
5. Validate all STEP_ASSERTIONS

### 2. Dry-Run Mode Instructions
**CRITICAL**: All test executions must use dry-run validation mode:
- Parse and validate logic only
- DO NOT write files to disk
- DO NOT modify .mission/ directory
- DO NOT execute destructive commands
- Return validation results only

### 3. AI Agent Execution Commands

#### M.PLAN Validation
```
Execute @m.plan with input from fixture in dry-run mode:
- Validate mission structure generation
- Check track complexity analysis
- Verify scope identification
- Confirm plan step logic
```

#### M.CLARIFY Validation
```
Execute @m.clarify with clarification responses in dry-run mode:
- Validate clarification processing
- Check track reassessment
- Verify mission updates
- Confirm status transitions
```

#### M.APPLY Validation
```
Execute @m.apply with active mission in dry-run mode:
- Validate execution step logic
- Check file creation simulation
- Verify verification command safety
- Confirm change summary generation
```

#### M.COMPLETE Validation
```
Execute @m.complete with completed mission in dry-run mode:
- Validate all three completion steps
- Check archival logic
- Verify metrics calculation
- Confirm project tracking updates
```

### 4. Assertion Validation
For each STEP_ASSERTION in fixtures:
- Execute assertion check
- Record pass/fail result
- Collect validation details
- Generate test report

### 5. Test Report Format
```markdown
# WORKFLOW VALIDATION REPORT

## SUMMARY
- **Total Tests**: 4
- **Passed**: X
- **Failed**: Y
- **Success Rate**: Z%

## DETAILED RESULTS

### M.PLAN Test
- **Status**: PASS/FAIL
- **Assertions**: X/Y passed
- **Issues**: [list any failures]

### M.CLARIFY Test
- **Status**: PASS/FAIL
- **Assertions**: X/Y passed
- **Issues**: [list any failures]

### M.APPLY Test
- **Status**: PASS/FAIL
- **Assertions**: X/Y passed
- **Issues**: [list any failures]

### M.COMPLETE Test
- **Status**: PASS/FAIL
- **Assertions**: X/Y passed
- **Issues**: [list any failures]
```
