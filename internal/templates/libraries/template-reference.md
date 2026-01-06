# TEMPLATE REFERENCE GUIDE

## Current Structure

```
libraries/
├── analysis/           # Analysis guidance templates
│   ├── clarification.md    # Clarification analysis criteria
│   ├── domain.md           # Domain analysis criteria
│   ├── duplication.md      # Duplication analysis criteria
│   └── intent.md           # Intent analysis criteria
├── displays/           # User output templates
│   ├── apply-failure.md    # ❌ MISSION FAILED
│   ├── apply-success.md    # ✅ MISSION EXECUTED
│   ├── complete-success.md # 🎉 MISSION COMPLETED
│   ├── error-no-mission.md # ❌ ERROR: No Active Mission
│   ├── error-mission-exists.md # ⚠️ EXISTING MISSION DETECTED
│   ├── plan-atomic.md      # ⚛️ ATOMIC TASK DETECTED
│   ├── plan-epic.md        # 📋 EPIC DECOMPOSED
│   └── plan-success.md     # ✅ MISSION CREATED
```

## Usage in Prompts

Clear, specific references:
```markdown
# In m.plan.md
**Intent Analysis**: Use `libraries/analysis/intent.md`
**Clarification Analysis**: Use `libraries/analysis/clarification.md`
**Duplication Analysis**: Use `libraries/analysis/duplication.md`
**Domain Analysis**: Use `libraries/analysis/domain.md`
**On Success**: Use template `libraries/displays/plan-success.md`
**On Epic**: Use template `libraries/displays/plan-epic.md`
**On Atomic**: Use template `libraries/displays/plan-atomic.md`
**On Mission Exists**: Use template `libraries/displays/error-mission-exists.md`

# In m.apply.md
**On Success**: Use template `libraries/displays/apply-success.md`
**On Failure**: Use template `libraries/displays/apply-failure.md`
**On No Mission**: Use template `libraries/displays/error-no-mission.md`

# In m.complete.md
**On Success**: Use template `libraries/displays/complete-success.md`
**On No Mission**: Use template `libraries/displays/error-no-mission.md`
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
- Project files: `.mission/backlog.md`

## Benefits

1. **Clear References**: `libraries/displays/plan-success.md` is unambiguous
2. **Easy Maintenance**: One template per outcome
3. **Consistent Variables**: Same names across all templates
4. **Logical Organization**: Grouped by purpose
5. **LLM-Friendly**: Simple file references instead of complex instructions
6. **Path Consistency**: All use `.mission/` root for deployment