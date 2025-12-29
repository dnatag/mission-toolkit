# TEMPLATE VARIABLES REFERENCE

## ALL VARIABLES USED IN TEMPLATES

### Core Mission Variables
```
{{TRACK}}               # 1, 2, 3, 4 (complexity level)
{{MISSION_TYPE}}        # WET, DRY, CLARIFICATION
{{MISSION_ID}}          # Track-Type-Timestamp (e.g., "2-WET-2024-01-15-14-30")
{{REFINED_INTENT}}      # Clarified user intent
{{FILE_LIST}}           # Files in scope (one per line)
{{PLAN_STEPS}}          # Implementation steps (bullet points)
{{VERIFICATION_COMMAND}} # Safe test command
{{MISSION_CONTENT}}     # Complete mission markdown
```

### Display Variables
```
{{FILE_COUNT}}          # Number of files in scope
{{TIMESTAMP}}           # 2024-01-15-14-30
{{DURATION}}            # "45 minutes" or "1 hour 15 minutes"
{{COMPLETION_DATE}}     # "2024-01-15 14:30:00"
{{DURATION_MINUTES}}    # 45 (numeric)
{{NEXT_STEP}}           # Next action for user (e.g., "/m.apply to execute this mission")
```

### Change Summary Variables
```
{{CHANGE_TITLE}}        # "feat: add user authentication"
{{CHANGE_DESCRIPTION}}  # Brief summary of change
{{CHANGE_DETAILS}}      # 3-4 bullet points with reasoning
{{IMPLEMENTATION_DETAIL}} # What was implemented
{{REASONING}}           # Why this approach was chosen
{{KEY_FILES_CHANGED}}   # Which files were modified
{{FILE_NECESSITY}}      # Why these files were needed
{{TECHNICAL_APPROACH}}  # How it was implemented
{{APPROACH_RATIONALE}}  # Why this technical approach
{{ADDITIONAL_CHANGES}}  # Other changes made
{{CHANGE_NECESSITY}}    # Why additional changes were needed
```

### Metrics Variables
```
{{VERIFICATION_STATUS}} # PASSED, FAILED, SKIPPED
{{COMPLETED_STEPS}}     # 4 (number of completed steps)
{{TOTAL_STEPS}}         # 5 (total number of steps)
{{QUALITY_SCORE}}       # 85% (calculated quality score)
{{FILES_MODIFIED}}      # 3 (actual files changed)
{{LINES_ADDED}}         # 127 (lines of code added)
{{LINES_REMOVED}}       # 23 (lines of code removed)
{{ITERATION_COUNT}}     # 1, 2, 3... (mission iteration)
{{SUCCESS_STATUS}}      # "SUCCESS" or "FAILED"
{{FAILURE_REASON}}      # Reason for failure (if failed)
```

### Learning Variables
```
{{DUPLICATION_FOUND}}   # "Yes" or "No"
{{SECURITY_ISSUES}}     # "None" or list of issues
{{PERFORMANCE_IMPACT}}  # "Minimal", "Moderate", "High"
{{PATTERNS_FOUND}}      # List of code patterns identified
{{REFACTORING_OPPORTUNITIES}} # List of refactoring suggestions
{{NEXT_MISSIONS}}       # Suggested follow-up missions
```

### DRY Mission Variables
```
{{ITERATION}}           # 2, 3, 4... (for DRY missions)
{{PARENT_MISSION}}      # Reference to original WET mission
```

### Clarification Variables
```
{{INITIAL_INTENT}}      # Original user input
{{CLARIFICATION_QUESTIONS}} # List of questions needing answers
{{ESTIMATED_FILES}}     # Provisional file list
```

### Epic Variables
```
{{SUB_INTENTS}}         # Formatted list of decomposed sub-intents
{{BACKLOG_ITEMS}}       # Formatted list for backlog display
```

### Logging Variables
```
{{LOG_ENTRY}}           # Formatted log entry for execution.log
{{TIMESTAMP}}           # Current timestamp (YYYY-MM-DD HH:MM:SS)
{{MISSION_ID}}          # Mission identifier for logging
```

### Metrics Refresh Variables
```
{{TOTAL_MISSIONS}}      # Count of completed missions
{{SUCCESS_RATE}}        # Percentage of successful missions
{{AVG_DURATION}}        # Average mission duration in minutes
{{TRACK_1_COUNT}}       # Count of Track 1 missions
{{TRACK_2_COUNT}}       # Count of Track 2 missions
{{TRACK_3_COUNT}}       # Count of Track 3 missions
{{TRACK_4_COUNT}}       # Count of Track 4 missions
{{WET_TO_DRY_RATIO}}    # Ratio of WET to DRY missions
{{COMMON_PATTERNS}}     # Most frequent patterns identified
{{BOTTLENECKS}}         # Most common failure points
```

### Atomic Task Variables
```
{{SUGGESTED_EDIT}}      # Direct edit suggestion for atomic tasks
```

## INITIALIZATION REQUIREMENTS

### Required in ALL Prompts
```bash
# These MUST be calculated/provided:
TRACK="[1-4 based on complexity]"
MISSION_TYPE="[WET|DRY|CLARIFICATION]"
TIMESTAMP="$(date +%Y-%m-%d-%H-%M)"
```

### Required for Mission Creation (m.plan)
```bash
REFINED_INTENT="[processed user input]"
FILE_LIST="[estimated files, one per line]"
PLAN_STEPS="[3-5 implementation steps]"
VERIFICATION_COMMAND="[safe test command]"
FILE_COUNT="[number of files in scope]"
MISSION_CONTENT="[complete mission markdown]"
```

### Required for Execution (m.apply)
```bash
CHANGE_TITLE="[brief change description]"
CHANGE_DESCRIPTION="[one-line summary]"
CHANGE_DETAILS="[4 bullet points with reasoning]"
# Plus all the individual change detail variables
```

### Required for Completion (m.complete)
```bash
MISSION_ID="[Track-Type-Timestamp]"
DURATION="[human readable duration]"
DURATION_MINUTES="[numeric minutes]"
COMPLETION_DATE="[YYYY-MM-DD HH:MM:SS]"
VERIFICATION_STATUS="[PASSED|FAILED|SKIPPED]"
COMPLETED_STEPS="[number]"
TOTAL_STEPS="[number]"
QUALITY_SCORE="[percentage]"
FILES_MODIFIED="[actual count]"
LINES_ADDED="[actual count]"
LINES_REMOVED="[actual count]"
SUCCESS_STATUS="[SUCCESS|FAILED]"
FAILURE_REASON="[reason if failed]"
LOG_ENTRY="[formatted log entry]"
```

### Required for Metrics Refresh
```bash
TOTAL_MISSIONS="[count]"
SUCCESS_RATE="[percentage]"
AVG_DURATION="[minutes]"
TRACK_1_COUNT="[count]"
TRACK_2_COUNT="[count]"
TRACK_3_COUNT="[count]"
TRACK_4_COUNT="[count]"
WET_TO_DRY_RATIO="[ratio]"
COMMON_PATTERNS="[list]"
BOTTLENECKS="[list]"
```

## DEFAULT VALUES FOR MISSING VARIABLES

```bash
# Use these defaults if variable not available:
{{DURATION}} → "Unknown duration"
{{QUALITY_SCORE}} → "Not calculated"
{{VERIFICATION_STATUS}} → "Not run"
{{SECURITY_ISSUES}} → "Not assessed"
{{PERFORMANCE_IMPACT}} → "Not assessed"
{{DUPLICATION_FOUND}} → "Not analyzed"
{{PATTERNS_FOUND}} → "None identified"
{{REFACTORING_OPPORTUNITIES}} → "None identified"
{{NEXT_MISSIONS}} → "None suggested"
{{SUCCESS_STATUS}} → "SUCCESS"
{{FAILURE_REASON}} → "Not applicable"
```

## CRITICAL NOTES

1. **Missing Variables**: Templates will break if required variables are undefined
2. **Prompt Responsibility**: Each prompt must initialize its required variables
3. **Calculation Order**: Some variables depend on others (MISSION_ID needs TRACK, MISSION_TYPE, TIMESTAMP)
4. **Type Consistency**: Numeric variables (FILE_COUNT, DURATION_MINUTES) vs string variables
5. **Format Requirements**: TIMESTAMP format must match across all templates