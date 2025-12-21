# TEST: M.COMPLETE Mission Archival and Metrics

## SCENARIO
**Given**: Successfully executed Track 2 WET mission with active status and verified implementation
**When**: User executes M.COMPLETE to finalize and archive the mission
**Then**: Mission should be archived with metrics, backlog updated, and system prepared for next mission
**Because**: Completion workflow captures learning, maintains project history, and enables continuous improvement

## MOCK DATA
**Mission State**:
```markdown
# MISSION
type: WET
track: 2
iteration: 1
status: active

## INTENT
Add user profile endpoint with CRUD operations for user data management

## SCOPE
routes/profile.js
controllers/profile-controller.js

## PLAN
- [x] Create profile route handlers (GET, POST, PUT, DELETE)
- [x] Implement profile controller with validation
- [x] Add error handling and response formatting
- [x] Note: Allow duplication for initial implementation

## VERIFICATION
npm test -- --grep "profile"
```

**System State**:
- Mission File: Contains active WET mission with completed implementation
- Implementation Files: routes/profile.js and controllers/profile-controller.js exist and functional
- Verification: Tests passing successfully
- Git Status: Changes committed and ready for archival

**Expected Completion Data**:
- Mission Duration: ~45 minutes (typical Track 2 timing)
- Files Modified: 2 (matching SCOPE exactly)
- Lines Added: ~150 (standard CRUD implementation)
- Verification Status: PASSED

## ASSERTIONS
- Mission should be archived to .mission/completed/ with timestamp
- Metrics should be collected and stored with mission archive
- .mission/backlog.md should be updated with any identified patterns
- .mission/metrics.md should be updated with aggregate statistics
- .mission/mission.md should be cleared or reset for next mission
- Archive should include complete mission history and implementation details
- Metrics should track duration, complexity, success rate, and pattern detection
- Backlog should identify any duplication opportunities for future DRY missions

## VALIDATION METHOD
AI should execute M.COMPLETE prompt logic through the following reasoning process:

1. **Mission Status Validation**: Verify mission has active status and successful verification
2. **Metrics Collection**: Calculate duration, file count, line count, and success indicators
3. **Pattern Detection**: Analyze implementation for duplication opportunities
4. **Archive Creation**: Move mission to .mission/completed/ with timestamp and metrics
5. **Backlog Update**: Add any identified refactoring opportunities to backlog
6. **Aggregate Metrics**: Update .mission/metrics.md with new data points
7. **Mission Reset**: Clear .mission/mission.md for next mission

## SUCCESS CRITERIA
**Pass Conditions**:
- Mission archived to .mission/completed/YYYY-MM-DD-HH-MM-mission.md
- Metrics file created at .mission/completed/YYYY-MM-DD-HH-MM-metrics.md
- .mission/backlog.md updated with any identified patterns
- .mission/metrics.md updated with aggregate statistics
- .mission/mission.md cleared for next mission
- Archive contains complete mission history and implementation details
- Metrics accurately reflect mission execution (duration, files, lines, success)

**Fail Conditions**:
- Mission not archived or archived incorrectly
- Metrics missing or inaccurate
- Backlog not updated appropriately
- Aggregate metrics not updated
- .mission/mission.md not cleared
- Archive missing implementation details
- Metrics don't match actual execution data

## EXPECTED REASONING TRACE
1. **Status Verification**: "Mission has status: active with all PLAN steps completed and verification passed - ready for completion"

2. **Metrics Calculation**: "Mission duration: 45 minutes, files modified: 2 (routes/profile.js, controllers/profile-controller.js), estimated lines: ~150, verification: PASSED"

3. **Pattern Analysis**: "Scanning implementation for duplication patterns - profile CRUD follows standard pattern, potential for abstraction after 2-3 similar implementations"

4. **Archive Generation**: "Creating timestamped archive in .mission/completed/ with complete mission history and calculated metrics"

5. **Backlog Update**: "Adding 'Extract CRUD pattern abstraction' to backlog for future DRY mission after similar implementations emerge"

6. **Aggregate Metrics Update**: "Updating .mission/metrics.md with Track 2 success, 45-minute duration, 2-file scope for trend analysis"

7. **Mission Reset**: "Clearing .mission/mission.md and displaying completion summary with next steps"

## CONFIDENCE ASSESSMENT
**Expected Confidence**: HIGH
- Clear mission completion with successful verification
- Standard metrics collection for Track 2 WET mission
- Well-defined archival process with established patterns
- Straightforward backlog and metrics update workflow

## RELATED TEST SCENARIOS
This test validates the core M.COMPLETE archival workflow. Related scenarios to test:
- Track 3 mission completion with security metrics
- Failed mission completion and recovery archival
- DRY mission completion with refactoring metrics
- Multiple mission completion and trend analysis
