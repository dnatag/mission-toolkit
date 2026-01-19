# TEMPLATE VARIABLES REFERENCE

## ALL VARIABLES USED IN TEMPLATES

### Core Mission Variables
```
{{TRACK}}               # 1, 2, 3, 4 (complexity level)
{{MISSION_TYPE}}        # WET, DRY
{{MISSION_ID}}          # Track-Type-Timestamp (e.g., "2-WET-2024-01-15-14-30")
{{REFINED_INTENT}}      # Clarified user intent
{{MISSION_CONTENT}}     # Complete mission markdown content
{{TIMESTAMP}}           # 2024-01-15-14-30
```

### Scope & Planning Variables
```
{{FILE_COUNT}}          # Number of files in scope (numeric)
{{SUGGESTED_EDIT}}      # Direct edit suggestion for atomic tasks (Track 1)
{{SUB_INTENTS}}         # Formatted list of decomposed sub-intents (Track 4)
```

### Execution Variables
```
{{IMPLEMENTATION_DETAIL}} # What was implemented
{{REASONING}}           # Why this approach was chosen
{{KEY_FILES_CHANGED}}   # Which files were modified
{{FILE_NECESSITY}}      # Why these files were needed
{{TECHNICAL_APPROACH}}  # How it was implemented
{{APPROACH_RATIONALE}}  # Why this technical approach
{{ADDITIONAL_CHANGES}}  # Other changes made
{{CHANGE_NECESSITY}}    # Why additional changes were needed
{{CHANGE_DETAILS}}      # 3-4 bullet points with reasoning
{{CONTRACT_CHANGES}}    # Bullet points of modified public interfaces/structs
{{CRITICAL_SNIPPETS}}   # Code block showing core logic change
```

### Completion Variables
```
{{DURATION}}            # "45 minutes" or "1 hour 15 minutes"
{{FINAL_COMMIT_HASH}}   # Git commit hash (if applicable)
```

### Checkpoint Variables
```
{{CHECKPOINT_0}}        # First checkpoint name
{{CHECKPOINT_1}}        # Second checkpoint name
{{CHECKPOINT_2}}        # Third checkpoint name
{{COMPLETED_STEPS}}     # Number of completed steps (numeric)
{{TOTAL_STEPS}}         # Total number of steps (numeric)
```

### Failure Variables
```
{{FAILURE_REASON}}      # Why the mission failed
{{RETRY_ADVICE}}        # Guidance for retry attempt
```

## VARIABLE USAGE BY TEMPLATE

### plan-success.md
- {{TRACK}}
- {{MISSION_TYPE}}
- {{FILE_COUNT}}
- {{MISSION_CONTENT}}

### plan-atomic.md
- {{REFINED_INTENT}}
- {{SUGGESTED_EDIT}}

### plan-epic.md
- {{SUB_INTENTS}}

### apply-success.md
- {{CHANGE_DETAILS}}
- {{CONTRACT_CHANGES}}
- {{CRITICAL_SNIPPETS}}
- {{CHECKPOINT_0}}
- {{CHECKPOINT_1}}
- {{CHECKPOINT_2}}

### apply-failure.md
- {{FAILURE_REASON}}
- {{RETRY_ADVICE}}
- {{CHECKPOINT_0}}
- {{CHECKPOINT_1}}
- {{CHECKPOINT_2}}
- {{COMPLETED_STEPS}}
- {{TOTAL_STEPS}}

### complete-success.md
- {{MISSION_ID}}
- {{DURATION}}
- {{FINAL_COMMIT_HASH}}

### complete-failure.md
- {{FAILURE_REASON}}

## VARIABLE INITIALIZATION BY PROMPT

### m.plan.md Must Provide
```
TRACK              # From m analyze complexity
MISSION_TYPE       # From m analyze duplication + backlog check
FILE_COUNT         # Count of files in scope section
MISSION_CONTENT    # Content of .mission/mission.md
REFINED_INTENT     # From intent analysis
SUGGESTED_EDIT     # For Track 1 only
SUB_INTENTS        # For Track 4 only
```

### m.apply.md Must Provide
```
CHANGE_DETAILS     # Narrative summary
CONTRACT_CHANGES   # Interface/struct changes
CRITICAL_SNIPPETS  # Core logic code block
FAILURE_REASON     # If failure
RETRY_ADVICE       # If failure
CHECKPOINT_*       # If checkpoints used
COMPLETED_STEPS    # If partial completion
TOTAL_STEPS        # If partial completion
```

### m.complete.md Must Provide
```
MISSION_ID         # From mission frontmatter
DURATION           # Calculated from timestamps
FINAL_COMMIT_HASH  # From git log (if applicable)
FAILURE_REASON     # If completion fails
```

## VARIABLE TYPES

### String Variables
All variables except those listed below are strings.

### Numeric Variables
```
{{FILE_COUNT}}      # Integer: number of files
{{COMPLETED_STEPS}} # Integer: steps completed
{{TOTAL_STEPS}}     # Integer: total steps
```

### Multi-line Variables
```
{{MISSION_CONTENT}}    # Full markdown content
{{SUB_INTENTS}}        # Formatted list with bullets
{{CHANGE_DETAILS}}     # Formatted list with bullets
{{CONTRACT_CHANGES}}   # Formatted list with bullets
{{CRITICAL_SNIPPETS}}  # Markdown code block
```

## FORMATTING GUIDELINES

### Lists (SUB_INTENTS, CHANGE_DETAILS)
```markdown
- First item
- Second item
- Third item
```

### File Lists (KEY_FILES_CHANGED)
```markdown
- `path/to/file1.go`
- `path/to/file2.go`
```

### Duration Format
```
"45 minutes"
"1 hour 15 minutes"
"2 hours 30 minutes"
```

### Timestamp Format
```
2024-01-15-14-30
```

### Mission ID Format
```
{TRACK}-{TYPE}-{TIMESTAMP}
Example: 2-WET-2024-01-15-14-30
```

## CRITICAL NOTES

1. **Missing Variables**: Templates will break if required variables are undefined
2. **Prompt Responsibility**: Each prompt must initialize its required variables before loading templates
3. **Type Consistency**: Numeric variables must be numbers, not strings
4. **Format Requirements**: TIMESTAMP and DURATION formats must be consistent
5. **Multi-line Handling**: Use proper markdown formatting for lists
6. **File Paths**: Always use backticks for file paths in lists
7. **Variable Extraction**: Extract variables from mission.md frontmatter or content as needed
