# TEST: M.APPLY Standard Mission Execution

## SCENARIO
**Given**: Active Track 2 WET mission with planned status and clear implementation scope
**When**: User executes M.APPLY to implement the planned mission
**Then**: Mission should execute all plan steps, modify only scoped files, pass verification, and transition to active status
**Because**: Standard WET missions should execute atomically with scope enforcement and mandatory verification

## MOCK DATA
**Mission State**:
```markdown
# MISSION
type: WET
track: 2
iteration: 1
status: planned

## INTENT
Add user profile endpoint with CRUD operations for user data management

## SCOPE
routes/profile.js
controllers/profile-controller.js

## PLAN
- [ ] Create profile route handlers (GET, POST, PUT, DELETE)
- [ ] Implement profile controller with validation
- [ ] Add error handling and response formatting
- [ ] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "profile"
```

**System State**:
- Mission File: Contains planned WET mission ready for execution
- Project Structure: Standard Node.js API with routes/ and controllers/ directories
- Git Status: Clean working directory
- Dependencies: All required packages installed

**Expected Implementation**:
- routes/profile.js: Express route definitions with CRUD endpoints
- controllers/profile-controller.js: Business logic with validation and error handling

## ASSERTIONS
- Mission status should transition from planned to active during execution
- Only files in SCOPE should be modified (routes/profile.js, controllers/profile-controller.js)
- All PLAN steps should be executed in order
- Verification command should run and pass
- Implementation should follow WET principles (allow duplication)
- No files outside SCOPE should be created or modified
- Git working directory should remain clean (no untracked changes outside scope)
- Mission file should maintain active status after successful execution

## VALIDATION METHOD
AI should execute M.APPLY prompt logic through the following reasoning process:

1. **Pre-execution Validation**: Verify mission has planned status and update to active
2. **Scope Enforcement**: Confirm only SCOPE files will be modified during execution
3. **Plan Execution**: Implement each PLAN step sequentially with WET approach
4. **File Creation**: Create routes/profile.js with CRUD endpoints
5. **Controller Implementation**: Create controllers/profile-controller.js with business logic
6. **Verification Execution**: Run npm test command and validate success
7. **Status Management**: Maintain active status on success, handle failure appropriately

## SUCCESS CRITERIA
**Pass Conditions**:
- Mission status transitions correctly (planned → active)
- Both SCOPE files created with appropriate content
- All PLAN steps completed successfully
- Verification command passes
- No files outside SCOPE modified
- Implementation follows WET principles (duplication allowed)
- Mission maintains active status after execution

**Fail Conditions**:
- Mission status incorrect or doesn't transition
- SCOPE files missing or inappropriate content
- PLAN steps skipped or executed incorrectly
- Verification command fails
- Files outside SCOPE modified
- Implementation violates WET principles
- Mission status changes to failed inappropriately

## EXPECTED REASONING TRACE
1. **Mission Validation**: "Mission has status: planned, updating to status: active for execution tracking"

2. **Scope Analysis**: "SCOPE contains 2 files: routes/profile.js and controllers/profile-controller.js. Will only modify these files during execution"

3. **Plan Step 1**: "Creating profile route handlers with GET /profile, POST /profile, PUT /profile/:id, DELETE /profile/:id endpoints"

4. **Plan Step 2**: "Implementing profile controller with input validation, business logic, and database operations"

5. **Plan Step 3**: "Adding comprehensive error handling with appropriate HTTP status codes and response formatting"

6. **WET Principle Application**: "Allowing duplication in initial implementation - not extracting common patterns yet, focusing on working functionality"

7. **Verification Execution**: "Running 'npm test -- --grep profile' to validate implementation works correctly"

8. **Success Confirmation**: "All PLAN steps completed, verification passed, mission remains active status for completion workflow"

## CONFIDENCE ASSESSMENT
**Expected Confidence**: HIGH
- Clear mission structure with specific implementation requirements
- Well-defined scope with standard file types
- Straightforward CRUD implementation pattern
- Standard verification approach with existing test framework

## RELATED TEST SCENARIOS
This test validates the core M.APPLY execution workflow. Related scenarios to test:
- Track 3 mission with security requirements (robust execution)
- Mission execution failure and recovery (status: failed)
- DRY mission execution (refactoring existing code)
- Verification failure handling and error recovery
