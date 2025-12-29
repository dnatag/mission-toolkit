# TEMPLATE REFERENCE GUIDE

## Current Structure

```
libraries/
├── analysis/           # Analysis guidance templates
│   ├── clarification.md    # Clarification analysis criteria
│   └── complexity.md       # Complexity assessment rules
├── displays/           # User output templates
│   ├── apply-failure.md    # ❌ MISSION FAILED
│   ├── apply-success.md    # ✅ MISSION EXECUTED
│   ├── clarify-escalation.md # 🔄 TRACK ESCALATION
│   ├── clarify-success.md  # ✅ CLARIFICATION COMPLETE
│   ├── complete-failure.md # ❌ MISSION FAILED (ARCHIVED)
│   ├── complete-success.md # 🎉 MISSION COMPLETED
│   ├── error-no-mission.md # ❌ ERROR: No Active Mission
│   ├── plan-atomic.md      # ⚛️ ATOMIC TASK DETECTED
│   ├── plan-clarification.md # ❓ CLARIFICATION NEEDED
│   ├── plan-epic.md        # 📋 EPIC DECOMPOSED
│   ├── plan-paused.md      # ⏸️ MISSION PAUSED
│   └── plan-success.md     # ✅ MISSION CREATED
├── logs/               # Execution logging templates
│   └── execution.md        # Log entry format
├── metrics/            # Metrics templates
│   ├── aggregate.md        # Project-wide metrics
│   ├── completion.md       # Individual mission metrics
│   └── insights.md         # Process insights format
├── missions/           # Mission file templates
│   ├── clarification.md    # Clarification mission template
│   ├── dry.md             # DRY mission template
│   └── wet.md             # WET mission template
├── scripts/            # Operation templates
│   ├── archive-completed.md # Archive to .mission/completed/
│   ├── archive-current.md  # Archive to .mission/paused/
│   ├── create-mission.md   # Create .mission/mission.md
│   ├── init-execution-log.md # Initialize execution log
│   ├── refresh-metrics.md  # Update metrics.md
│   ├── status-to-active.md # Update mission status
│   └── validate-planned.md # Check mission status
└── variables/          # Variable calculation rules
    ├── file-list.md        # File estimation rules
    ├── timestamps.md       # Date/time formatting
    └── track-calculation.md # Track complexity logic
```

## Usage in Prompts

Clear, specific references:
```markdown
# In m.plan.md
**Clarification Analysis**: Use `libraries/analysis/clarification.md`
**Complexity Analysis**: Use `libraries/analysis/complexity.md`
**On Success**: Use template `libraries/displays/plan-success.md`
**On Clarification**: Use template `libraries/displays/plan-clarification.md`
**On Epic**: Use template `libraries/displays/plan-epic.md`
**On Atomic**: Use template `libraries/displays/plan-atomic.md`
**Mission Template**: Use `libraries/missions/wet.md`
**Create Script**: Use `libraries/scripts/create-mission.md`
**Log Initialization**: Use `libraries/scripts/init-execution-log.md`

# In m.clarify.md
**Complexity Reassessment**: Use `libraries/analysis/complexity.md`
**Mission Update**: Use `libraries/missions/wet.md` or `libraries/missions/dry.md`
**Success Display**: Use `libraries/displays/clarify-success.md`
**Track 4 Escalation**: Use `libraries/displays/clarify-escalation.md`

# In m.apply.md
**On Success**: Use template `libraries/displays/apply-success.md`
**On Failure**: Use template `libraries/displays/apply-failure.md`
**Status Script**: Use `libraries/scripts/status-to-active.md`
**Validation Script**: Use `libraries/scripts/validate-planned.md`
**Logging**: Use `libraries/logs/execution.md`

# In m.complete.md
**On Success**: Use template `libraries/displays/complete-success.md`
**On Failure**: Use template `libraries/displays/complete-failure.md`
**Archive Script**: Use `libraries/scripts/archive-completed.md`
**Metrics Refresh**: Use `libraries/scripts/refresh-metrics.md`
**Metrics Template**: Use `libraries/metrics/completion.md`
**Logging**: Use `libraries/logs/execution.md`
```

## Variable Standardization

Consistent naming across all templates:
```
{{TRACK}}           # 1, 2, 3, 4 (complexity)
{{MISSION_TYPE}}    # WET, DRY, CLARIFICATION
{{TIMESTAMP}}       # 2024-01-15-14-30
{{DURATION}}        # "45 minutes"
{{FILE_COUNT}}      # 3
{{MISSION_CONTENT}} # Full mission markdown
{{MISSION_ID}}      # Track-Type-Timestamp
```

## Path Consistency

All templates use `.mission/` root path when deployed:
- Mission files: `.mission/mission.md`
- Paused missions: `.mission/paused/`
- Completed missions: `.mission/completed/`
- Project files: `.mission/backlog.md`, `.mission/metrics.md`

## Benefits

1. **Clear References**: `libraries/displays/plan-success.md` is unambiguous
2. **Easy Maintenance**: One template per outcome
3. **Consistent Variables**: Same names across all templates
4. **Logical Organization**: Grouped by purpose
5. **LLM-Friendly**: Simple file references instead of complex instructions
6. **Path Consistency**: All use `.mission/` root for deployment