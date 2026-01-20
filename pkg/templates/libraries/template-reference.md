# TEMPLATE REFERENCE GUIDE

## Current Structure

```
libraries/
├── displays/                      # User output templates
│   ├── apply-failure.md               # ❌ MISSION FAILED
│   ├── apply-success.md               # ✅ MISSION EXECUTED
│   ├── checkpoint-failure-unstaged.md # ⚠️ CHECKPOINT FAILED (unstaged changes)
│   ├── complete-failure.md            # ❌ COMPLETION FAILED
│   ├── complete-success.md            # 🎉 MISSION COMPLETED
│   ├── error-mission-exists.md        # ❌ ERROR: Mission Already Exists
│   ├── error-no-mission.md            # ❌ ERROR: No Active Mission
│   ├── plan-atomic.md                 # ⚛️ ATOMIC TASK DETECTED
│   ├── plan-epic.md                   # 📋 EPIC DECOMPOSED
│   ├── plan-paused.md                 # ⏸️ MISSION PAUSED
│   └── plan-success.md                # ✅ MISSION CREATED
├── cli-reference.md           # CLI command reference
├── template-reference.md      # This file
└── variables-reference.md     # Variable naming guide

prompts/
├── m.plan.md                  # Planning workflow prompt
├── m.apply.md                 # Execution workflow prompt
└── m.complete.md              # Completion workflow prompt
```

## Display Templates

### Planning Phase
- **plan-success.md**: Standard mission created (Track 2 or 3)
- **plan-atomic.md**: Track 1 - too simple for mission, direct edit suggested
- **plan-epic.md**: Track 4 - decomposed to backlog
- **plan-paused.md**: Mission paused for clarification

### Execution Phase
- **apply-success.md**: Mission executed successfully
- **apply-failure.md**: Mission execution failed

### Completion Phase
- **complete-success.md**: Mission completed and archived
- **complete-failure.md**: Completion process failed

### Checkpoint Phase
- **checkpoint-failure-unstaged.md**: Checkpoint creation failed due to unstaged changes

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
m analyze decompose

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
- Checkpoint Failed: Use `.mission/libraries/displays/checkpoint-failure-unstaged.md`
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

### Core Commands
- `m init` - Initialize project with AI-specific templates
- `m dashboard` - Interactive TUI with mission details and logs
- `m version` - Show version information
- `m check` - Validate input strings
- `m docs` - Generate CLI documentation schema

### Mission Management
- `m mission check` - Check mission state and validate artifacts
- `m mission id` - Get or create mission ID
- `m mission create` - Create mission.md with intent
- `m mission update` - Update mission status or sections
- `m mission finalize` - Validate and display mission for review
- `m mission archive` - Archive mission files to completed directory
- `m mission mark-complete` - Mark plan step as complete
- `m mission pause` - Pause current mission
- `m mission restore` - Restore paused mission

### Analysis Tools
- `m analyze intent` - Analyze user intent
- `m analyze clarify` - Check for clarification needs
- `m analyze scope` - Determine affected files
- `m analyze test` - Analyze test requirements
- `m analyze duplication` - Check for code duplication
- `m analyze complexity` - Calculate complexity track
- `m analyze decompose` - Decompose epic intents

### Backlog Management
- `m backlog list` - List backlog items with filters
- `m backlog add` - Add backlog items
- `m backlog complete` - Mark item as complete
- `m backlog cleanup` - Remove completed items

### Checkpoint Management
- `m checkpoint create` - Create checkpoint
- `m checkpoint restore` - Restore checkpoint
- `m checkpoint clear` - Clear checkpoints
- `m checkpoint commit` - Create final commit

### Logging
- `m log` - Log messages to execution log

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
