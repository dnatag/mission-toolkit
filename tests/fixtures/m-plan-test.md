# M.PLAN TEST FIXTURE

## TEST_MODE
dry_run: true
validate_only: true

## INPUT
```text
add user authentication to the API
```

## EXPECTED_MISSION_OUTPUT
```markdown
# MISSION

type: WET
track: 2
iteration: 1
status: planned

## INTENT
Add user authentication system to API with login/logout endpoints and JWT token validation

## SCOPE
auth/auth.js
routes/auth.js
middleware/authenticate.js

## PLAN
- [ ] Create authentication service with JWT token generation
- [ ] Add login/logout API endpoints
- [ ] Implement authentication middleware for protected routes
- [ ] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "auth"
```

## STEP_ASSERTIONS
- [ ] Mission structure validated (not written to disk)
- [ ] Status set to "planned"
- [ ] Track correctly identified as 2
- [ ] SCOPE contains 3 files
- [ ] VERIFICATION command is safe (no destructive operations)

## VALIDATION_INSTRUCTIONS
```
Execute @m.plan in validation mode:
- Parse intent and generate mission structure
- Validate against expected output
- DO NOT write .mission/mission.md to disk
- DO NOT create any files in SCOPE
- Return validation results only
```
