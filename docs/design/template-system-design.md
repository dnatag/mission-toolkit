# TEMPLATE SYSTEM DESIGN RATIONALE

## PROBLEM IDENTIFIED
Complex multi-step instructions (10 steps in m.plan, 3 steps in m.complete) fail across different LLM models due to:
- Instruction following inconsistency
- Step skipping or reordering
- Template format variations
- Context window limitations

## SOLUTION IMPLEMENTED: TEMPLATE + SCRIPT HYBRID

### 1. Template Libraries
Replaced complex instructions with exact templates in `libraries/`:

```
libraries/
├── displays/           # User output templates
├── missions/           # Mission file templates  
├── scripts/            # Bash operation templates
├── metrics/            # Metrics templates
└── variables/          # Variable calculation rules
```

### 2. Variable-Based Approach
Instead of 10 complex steps, use clear variable mapping:

```markdown
**TRACK 2-3**: Use template `.mission/libraries/missions/wet.md` with variables:
- {{TRACK}} = 2 or 3
- {{REFINED_INTENT}} = Refined summary of the goal
- {{FILE_LIST}} = List of file paths, one per line
- {{PLAN_STEPS}} = Implementation steps as bullet points
- {{VERIFICATION_COMMAND}} = Safe shell command
```

### 3. Script Templates
For file operations, provide exact commands in `libraries/scripts/`:

```bash
# From libraries/scripts/archive-current.md
if [ -f ".mission/mission.md" ]; then
  mkdir -p .mission/paused
  TIMESTAMP=$(date +%Y-%m-%d-%H-%M)
  mv .mission/mission.md ".mission/paused/${TIMESTAMP}-mission.md"
fi
```

## IMPLEMENTATION COMPLETED

### ✅ Phase 1: Template Library
Created exact templates for:
- Mission files: `libraries/missions/wet.md`, `dry.md`, `clarification.md`
- Display outputs: `libraries/displays/[command]-[outcome].md`
- File operations: `libraries/scripts/[operation].md`
- Metrics: `libraries/metrics/completion.md`, `aggregate.md`
- Variables: `libraries/variables-reference.md`

### ✅ Phase 2: Prompt Refactoring
Refactored all slash commands to use templates:
- `m.plan.md` - Uses mission and display templates
- `m.apply.md` - Uses display and script templates
- `m.clarify.md` - Uses mission and display templates
- `m.complete.md` - Uses display, script, and metrics templates

### ✅ Phase 3: Safety Mechanisms
- Template validation through variable reference
- EXECUTION INSTRUCTIONS prevent LLM bypass
- Clear variable requirements prevent template breakage

## RESULTS ACHIEVED
- **Consistency**: Exact templates across all LLMs
- **Reliability**: Simple template filling vs complex logic
- **Maintainability**: Change template once, affects all models
- **Robustness**: Script templates handle file operations
- **Logic Preservation**: Claude Sonnet's proven workflows maintained

This approach successfully traded flexibility for reliability - exactly what was needed for production use across different LLM models.