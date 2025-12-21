# M.COMPLETE TEST FIXTURE

## TEST_MODE
dry_run: true
validate_only: true

## SETUP_MISSION
```markdown
# MISSION

type: WET
track: 2
iteration: 1
status: active

## INTENT
Add user profile endpoint to API

## SCOPE
routes/profile.js
controllers/profile-controller.js

## PLAN
- [x] Create profile controller with CRUD operations
- [x] Add profile routes with authentication middleware
- [x] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "profile"
```

## EXPECTED_COMPLETION_STEPS
1. **Mission Validation**
   - Status: active ✅
   - All PLAN items completed ✅
   - VERIFICATION passed ✅

2. **Mission Completion and Archival**
   - Status: active → completed
   - Add completed_at timestamp
   - Archive to .mission/completed/YYYY-MM-DD-HH-MM-mission.md
   - Create .mission/completed/YYYY-MM-DD-HH-MM-metrics.md

3. **Project Tracking Updates**
   - Update .mission/metrics.md aggregate statistics
   - Check .mission/backlog.md for matching items

## EXPECTED_METRICS_UPDATE
```markdown
## AGGREGATE STATISTICS
- **Total Missions**: 21 completed
- **Success Rate**: 100% (21/21 successful)
- **Average Duration**: ~4 minutes

## RECENT COMPLETIONS
- 2025-12-20 Track 2: feat: add user profile endpoint to API
- [previous completions...]
```

## STEP_ASSERTIONS
- [ ] Mission status changed to completed with timestamp
- [ ] Mission archived to completed/ directory
- [ ] Detailed metrics file created
- [ ] Aggregate metrics updated in .mission/metrics.md
- [ ] Mission file reset to "No active mission"
- [ ] All three completion steps executed

## VALIDATION_INSTRUCTIONS
```
Execute @m.complete in validation mode:
- Validate all three completion steps
- Simulate archival and metrics updates
- DO NOT create actual archive files
- DO NOT modify .mission/metrics.md
- Return completion validation results only
```
