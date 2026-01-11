# TEMPLATE REFERENCE GUIDE

## Current Structure

```
libraries/
├── displays/              # User output templates
│   ├── apply-failure.md       # ❌ MISSION FAILED
│   ├── apply-success.md       # ✅ MISSION EXECUTED
│   ├── complete-failure.md    # ❌ COMPLETION FAILED
│   ├── complete-success.md    # 🎉 MISSION COMPLETED
│   ├── error-mission-exists.md # ❌ ERROR: Mission Already Exists
│   ├── error-no-mission.md    # ❌ ERROR: No Active Mission
│   ├── plan-atomic.md         # ⚛️ ATOMIC TASK DETECTED
│   ├── plan-epic.md           # 📋 EPIC DECOMPOSED
│   ├── plan-paused.md         # ⏸️ MISSION PAUSED
│   └── plan-success.md        # ✅ MISSION CREATED
├── cli-reference.md       # CLI command reference
├── template-reference.md  # This file
└── variables-reference.md # Variable naming guide

prompts/
├── m.plan.md              # Planning workflow prompt
├── m.apply.md             # Execution workflow prompt
└── m.complete.md          # Completion workflow prompt
```

## Display Templates

### Planning Phase
- **plan-success.md**: Standard mission created (Track 2 or 3)
- **plan-atomic.md**: Track 1 - too simple for mission
- **plan-epic.md**: Track 4 - decomposed to backlog
- **plan-paused.md**: Mission paused for clarification

### Execution Phase
- **apply-success.md**: Mission executed successfully
- **apply-failure.md**: Mission execution failed

### Completion Phase
- **complete-success.md**: Mission completed and archived
- **complete-failure.md**: Completion process failed

### Error States
- **error-no-mission.md**: No active mission found
- **error-mission-exists.md**: Active mission already exists

## Usage in Prompts

### In m.plan.md
```markdown
# Analysis Commands
m analyze intent "<user-input>"
m analyze clarify
m analyze scope
m analyze test
m analyze duplication
m analyze complexity

# Display Templates
- Track 1 (Atomic): Use `.mission/libraries/displays/plan-atomic.md`
- Track 4 (Epic): Use `.mission/libraries/displays/plan-epic.md`
- Track 2/3 (Success): Use `.mission/libraries/displays/plan-success.md`
- Paused: Use `.mission/libraries/displays/plan-paused.md`
- Mission Exists: Use `.mission/libraries/displays/error-mission-exists.md`
```

### In m.apply.md
```markdown
# Display Templates
- Success: Use `.mission/libraries/displays/apply-success.md`
- Failure: Use `.mission/libraries/displays/apply-failure.md`
- No Mission: Use `.mission/libraries/displays/error-no-mission.md`
```

### In m.complete.md
```markdown
# Display Templates
- Success: Use `.mission/libraries/displays/complete-success.md`
- Failure: Use `.mission/libraries/displays/complete-failure.md`
- No Mission: Use `.mission/libraries/displays/error-no-mission.md`
```

## CLI Commands Reference

All CLI commands are documented in `libraries/cli-reference.md`:
- Core: `m init`, `m status`, `m version`, `m check`
- Mission: `m mission check|id|create|update|finalize|archive`
- Analysis: `m analyze intent|clarify|scope|test|duplication|complexity`
- Backlog: `m backlog list|add|complete|resolve|cleanup`
- Checkpoint: `m checkpoint create|restore|clear`
- Logging: `m log`

## Variable Standardization

See `libraries/variables-reference.md` for complete variable naming guide.

Common variables:
```
{{TRACK}}              # 1, 2, 3, 4 (complexity)
{{MISSION_TYPE}}       # WET, DRY
{{TIMESTAMP}}          # 2024-01-15-14-30
{{DURATION}}           # "45 minutes"
{{FILE_COUNT}}         # 3
{{MISSION_CONTENT}}    # Full mission markdown
{{MISSION_ID}}         # Track-Type-Timestamp
{{REFINED_INTENT}}     # Clarified user intent
{{NEW_INTENT}}         # New intent when mission exists
{{SUB_INTENTS}}        # List of decomposed sub-intents
{{SUGGESTED_EDIT}}     # Atomic edit suggestion
```

## Path Consistency

All templates use `.mission/` root path when deployed:
- Mission files: `.mission/mission.md`
- Execution log: `.mission/execution.log`
- Backlog: `.mission/backlog.md`
- Governance: `.mission/governance.md`
- Completed: `.mission/completed/<MISSION_ID>-*`
- Checkpoints: `.mission/checkpoints/<MISSION_ID>-<timestamp>/`

## Template Loading

Templates are loaded using the file read tool:
```markdown
Use file read tool to load template `.mission/libraries/displays/plan-success.md`
```

Fill variables and display to user:
```markdown
Fill template with:
- {{TRACK}}: From mission frontmatter
- {{MISSION_TYPE}}: From mission frontmatter
- {{FILE_COUNT}}: Count of files in scope
- {{MISSION_CONTENT}}: The content of `.mission/mission.md`
```

## Benefits

1. **Clear References**: `libraries/displays/plan-success.md` is unambiguous
2. **Easy Maintenance**: One template per outcome
3. **Consistent Variables**: Same names across all templates
4. **Logical Organization**: Grouped by purpose (displays, prompts)
5. **LLM-Friendly**: Simple file references instead of complex instructions
6. **Path Consistency**: All use `.mission/` root for deployment
7. **Embedded**: Templates are embedded in the CLI binary
8. **Versioned**: Templates version matches CLI version
