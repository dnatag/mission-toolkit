# Workflows

## Mission Lifecycle

```
m.plan → 🤝 Review mission.md → m.apply → 🤝 Review code → [Adjustments] → m.complete
(Handshake #1)                  (Handshake #2)
```

### Steps

1. **m.plan** creates mission.md with INTENT, SCOPE, PLAN, VERIFICATION
2. **🤝 Review & approve** the mission before execution (authorize the architecture)
3. **m.apply** executes, polishes, and generates commit message
4. **🤝 Review code** and optionally request adjustments (verify the implementation)
5. **m.complete** archives mission and creates git commit

## Bugfix Workflow

```
m.debug → 🤝 Review diagnosis.md → m.plan → m.apply → m.complete
(Investigation)                    (Fix Planning)
```

### Steps

1. **m.debug** investigates the bug and creates diagnosis.md with root cause analysis
2. **🤝 Review diagnosis** to understand the problem (evidence-based findings)
3. **m.plan** automatically consumes diagnosis.md to create a targeted fix mission
4. **m.apply** implements the fix with verification
5. **m.complete** archives both diagnosis and fix mission together

## Project Structure

```
.mission/
├── governance.md          # Core principles and workflow rules
├── backlog.md            # Future work and refactoring opportunities
├── mission.md            # Current active mission (auto-generated)
├── diagnosis.md          # Current bug diagnosis (auto-generated)
├── execution.log         # Current mission execution log
├── completed/            # Archived missions and detailed metrics
├── paused/               # Temporarily paused missions
└── libraries/            # Template system (embedded)

# AI-specific prompt directories:
.amazonq/prompts/         # Amazon Q prompts
.claude/commands/         # Claude commands
.kiro/prompts/           # Kiro prompts
.opencode/command/       # OpenCode commands
```
