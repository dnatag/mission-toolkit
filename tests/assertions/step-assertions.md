# STEP ASSERTION FRAMEWORK

## PURPOSE
Framework for validating Mission Toolkit workflow steps without file system changes.

## ASSERTION TYPES

### 1. Mission Structure Assertions
**Step 1: Mission Creation**
- `mission_file_structure`: Validate mission has required sections (INTENT, SCOPE, PLAN, VERIFICATION)
- `mission_status_valid`: Check status is one of: clarifying, planned, active, completed, failed
- `mission_type_valid`: Check type is one of: WET, DRY, CLARIFICATION
- `mission_track_valid`: Check track is 1, 2, 3, or TBD

**Step 2: Content Validation**
- `intent_not_empty`: INTENT section contains meaningful description
- `scope_files_listed`: SCOPE contains file paths (one per line)
- `plan_has_steps`: PLAN contains actionable checklist items
- `verification_command_safe`: VERIFICATION contains safe, non-destructive command

### 2. Workflow Transition Assertions
**Step 3: Status Transitions**
- `status_transition_valid`: Status changes follow valid workflow (planned → active → completed)
- `clarification_to_planned`: CLARIFICATION missions properly transition to WET/DRY
- `active_to_completed`: Active missions properly complete with timestamp

**Step 4: Track Analysis**
- `track_complexity_correct`: Track assignment matches complexity analysis rules
- `track_escalation_handled`: Track 4 missions properly decomposed to backlog
- `domain_multipliers_applied`: Security/performance concerns increase track appropriately

### 3. Execution Logic Assertions
**Step 5: Plan Execution**
- `plan_steps_sequential`: Plan steps execute in logical order
- `scope_enforcement`: Only files in SCOPE are modified (in dry-run: validated)
- `verification_executed`: Verification command runs successfully
- `change_summary_generated`: Proper change summary with 4-bullet format

**Step 6: Completion Process**
- `three_step_completion`: All three completion steps execute (validation, archival, tracking)
- `metrics_updated`: Aggregate metrics properly updated with new completion
- `archival_proper`: Mission and metrics properly archived with timestamps
- `mission_reset`: Mission file reset to "No active mission" state

### 4. Safety Assertions
**Step 7: Dry-Run Validation**
- `no_file_creation`: No actual files created during dry-run testing
- `no_mission_modification`: .mission/mission.md not modified during validation
- `no_destructive_commands`: No destructive operations executed
- `validation_only_mode`: Only validation logic executed, no side effects

### 5. Assertion Execution Framework

#### Assertion Check Function
```
For each assertion in test fixture:
1. Extract assertion type and expected value
2. Execute corresponding validation logic
3. Compare actual vs expected result
4. Record pass/fail with details
5. Continue to next assertion
```

#### Assertion Result Format
```markdown
## ASSERTION RESULTS

### Mission Structure
- [✅/❌] mission_file_structure: Expected sections present
- [✅/❌] mission_status_valid: Status = "planned"
- [✅/❌] mission_type_valid: Type = "WET"
- [✅/❌] mission_track_valid: Track = 2

### Content Validation
- [✅/❌] intent_not_empty: INTENT contains 15 words
- [✅/❌] scope_files_listed: SCOPE contains 3 files
- [✅/❌] plan_has_steps: PLAN contains 4 steps
- [✅/❌] verification_command_safe: Command is read-only

### Workflow Transitions
- [✅/❌] status_transition_valid: planned → active transition
- [✅/❌] track_complexity_correct: Track 2 matches 3 files

### Execution Logic
- [✅/❌] plan_steps_sequential: Steps execute in order
- [✅/❌] scope_enforcement: Only SCOPE files affected
- [✅/❌] verification_executed: Command returned exit code 0
- [✅/❌] change_summary_generated: 4-bullet format confirmed

### Safety Validation
- [✅/❌] no_file_creation: No files created on disk
- [✅/❌] no_mission_modification: .mission/mission.md unchanged
- [✅/❌] validation_only_mode: Dry-run mode confirmed
```

## USAGE INSTRUCTIONS

1. Load test fixture
2. Execute workflow in dry-run mode
3. Run all applicable assertions
4. Generate assertion results
5. Compile into test report
