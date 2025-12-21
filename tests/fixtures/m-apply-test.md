# M.APPLY TEST FIXTURE

## TEST_MODE
dry_run: true
validate_only: true

## SETUP_MISSION
```markdown
# MISSION

type: WET
track: 2
iteration: 1
status: planned

## INTENT
Add user profile endpoint to API

## SCOPE
routes/profile.js
controllers/profile-controller.js

## PLAN
- [ ] Create profile controller with CRUD operations
- [ ] Add profile routes with authentication middleware
- [ ] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "profile"
```

## EXPECTED_EXECUTION_STEPS
1. **Pre-execution Validation**
   - Mission status: planned → active
   - Scope files validated

2. **Mission Execution**
   - Create routes/profile.js with GET/PUT endpoints
   - Create controllers/profile-controller.js with business logic
   - Run verification command

3. **Success Output**
   - Status remains active
   - Offer auto-completion

## EXPECTED_CHANGE_SUMMARY
```
feat: add user profile endpoint to API

- Created profile controller with CRUD operations → separates business logic from routing for maintainability
- Added profile routes with authentication middleware → ensures secure access to user data
- Implemented GET/PUT endpoints for profile management → provides complete profile functionality
- Added comprehensive test coverage → ensures endpoint reliability and security
```

## STEP_ASSERTIONS
- [ ] Mission status changed from planned to active
- [ ] All PLAN steps executed in sequence
- [ ] Files created according to SCOPE (in dry-run: validated only)
- [ ] VERIFICATION command executed successfully
- [ ] Change summary generated with proper format

## VALIDATION_INSTRUCTIONS
```
Execute @m.apply in validation mode:
- Validate mission execution logic
- Simulate file creation and verification
- DO NOT create actual files
- DO NOT modify .mission/mission.md status
- Return execution validation results only
```
