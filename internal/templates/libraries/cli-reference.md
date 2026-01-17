# Mission Toolkit CLI Reference

This document is the master reference for all `m` CLI commands available for use in prompts.

## Core Commands

### `m init`
- **Purpose**: Initialize Mission Toolkit project with templates for specified AI type
- **Usage**: `m init --ai <type>`
- **Flags**:
  - `--ai`: AI assistant type (q, claude, kiro, opencode)
- **Output**: Creates `.mission/` directory with governance files and AI-specific prompt templates

### `m dashboard`
- **Purpose**: Display comprehensive mission dashboard with execution logs in interactive TUI
- **Usage**: `m dashboard`
- **Output**: Interactive terminal UI with split-pane view showing mission details, execution logs, and commit history

### `m version`
- **Purpose**: Show version information
- **Usage**: `m version`
- **Output**: Displays CLI version and embedded template version

### `m check`
- **Purpose**: Check if input is empty or whitespace
- **Usage**: `m check "<string to check>"`
- **Output (JSON)**:
  - `valid`: `true` or `false`
  - `next_step`: "PROCEED" or "ASK_USER"
  - `message`: Guidance message

## Mission Management (`m mission`)

### `m mission check`
- **Purpose**: Check mission state and validate artifacts
- **Usage**: `m mission check [--context <context>]`
- **Flags**:
  - `--context`: Optional context (e.g., "plan", "apply")
- **Output (JSON)**:
  - `next_step`: "PROCEED" or "STOP"
  - `message`: Status message
  - `mission_exists`: Boolean

### `m mission id`
- **Purpose**: Get or create mission ID
- **Usage**: `m mission id`
- **Output**: Mission ID string

### `m mission create`
- **Purpose**: Create mission.md with intent
- **Usage**: `m mission create --intent "<intent>"`
- **Flags**:
  - `--intent`: User's intent description
- **Output**: Creates `.mission/mission.md`

### `m mission update`
- **Purpose**: Update mission status or sections
- **Usage**: `m mission update [flags]`
- **Flags**:
  - `--status <status>`: Update mission status
  - `--section <section>`: Section to update (scope, plan, verification)
  - `--item <item>`: Add item to section (can be repeated)
  - `--content <content>`: Set section content
  - `--frontmatter <key=value>`: Update frontmatter field
- **Output**: Updates `.mission/mission.md`

### `m mission finalize`
- **Purpose**: Validate and display mission.md for review
- **Usage**: `m mission finalize`
- **Output (JSON)**:
  - `action`: "PROCEED" or "INVALID"
  - `errors`: List of validation errors (if any)

### `m mission archive`
- **Purpose**: Archive mission files to completed directory
- **Usage**: `m mission archive`
- **Output**: Moves mission files to `.mission/completed/<MISSION_ID>-*`

### `m mission mark-complete`
- **Purpose**: Mark a plan step as complete and log progress
- **Usage**: `m mission mark-complete --step <N> --status <STATUS> --message "<message>"`
- **Flags**:
  - `--step`: Step number to mark as complete (1-indexed)
  - `--status`: Status level for logging (INFO, SUCCESS, FAILED)
  - `--message`: Message to log for this step
- **Output**: Updates plan step checkbox and logs message

### `m mission pause`
- **Purpose**: Pause current mission and save to paused folder
- **Usage**: `m mission pause`
- **Output**: Moves mission to `.mission/paused/` with timestamp

### `m mission restore`
- **Purpose**: Restore a paused mission
- **Usage**: `m mission restore`
- **Output**: Restores most recent paused mission to active state

## Analysis Tools (`m analyze`)

### `m analyze intent`
- **Purpose**: Provide intent analysis template with user input
- **Usage**: `m analyze intent "<user-input>"`
- **Output**: Intent analysis template with injected user input

### `m analyze clarify`
- **Purpose**: Provide clarification analysis template with current intent
- **Usage**: `m analyze clarify`
- **Output**: Clarification template with current mission intent

### `m analyze scope`
- **Purpose**: Provide scope analysis template with current intent
- **Usage**: `m analyze scope`
- **Output**: Scope analysis template for determining affected files

### `m analyze test`
- **Purpose**: Provide test analysis template with current intent and scope
- **Usage**: `m analyze test`
- **Output**: Test requirement analysis template

### `m analyze duplication`
- **Purpose**: Provide duplication analysis template with current intent
- **Usage**: `m analyze duplication`
- **Output**: Duplication detection template for WET→DRY workflow

### `m analyze complexity`
- **Purpose**: Provide complexity analysis template with current intent and scope
- **Usage**: `m analyze complexity`
- **Output**: Complexity analysis template with track calculation

## Backlog Management (`m backlog`)

### `m backlog list`
- **Purpose**: List backlog items
- **Usage**: `m backlog list [--all] [--type <type>]`
- **Flags**:
  - `--all`: Include completed items
  - `--type`: Filter by type (decomposed, refactor, future)
- **Output**: List of backlog items

### `m backlog add`
- **Purpose**: Add one or more backlog items
- **Usage**: `m backlog add "<description>" [<description>...] --type <type>`
- **Flags**:
  - `--type`: Item type (decomposed, refactor, future) - REQUIRED
- **Output**: Confirmation message

### `m backlog complete`
- **Purpose**: Mark a backlog item as complete
- **Usage**: `m backlog complete --item "<exact text>"`
- **Flags**:
  - `--item`: Exact text of the item to complete - REQUIRED
- **Output**: Moves item to COMPLETED section with timestamp

### `m backlog resolve`
- **Purpose**: Mark a refactor opportunity as resolved via DRY conversion
- **Usage**: `m backlog resolve --item "<exact text>"`
- **Flags**:
  - `--item`: Exact text of the refactor item - REQUIRED
- **Output**: Marks item with [RESOLVED] prefix and DRY timestamp in-place

### `m backlog cleanup`
- **Purpose**: Remove completed items from the backlog
- **Usage**: `m backlog cleanup [--type <type>]`
- **Flags**:
  - `--type`: Filter by type (decomposed, refactor, future)
- **Output**: Number of items removed

## Checkpoint Management (`m checkpoint`)

### `m checkpoint create`
- **Purpose**: Create a checkpoint of current working directory state
- **Usage**: `m checkpoint create`
- **Output**: Creates checkpoint in `.mission/checkpoints/<MISSION_ID>-<timestamp>/`

### `m checkpoint restore`
- **Purpose**: Restore working directory to specified checkpoint
- **Usage**: `m checkpoint restore <checkpoint>`
- **Args**: Checkpoint name (e.g., "MISSION-123-20230101-120000")
- **Output**: Restores files from checkpoint

### `m checkpoint clear`
- **Purpose**: Clear all checkpoints for current mission
- **Usage**: `m checkpoint clear`
- **Output**: Removes checkpoint directory

### `m checkpoint commit`
- **Purpose**: Create final commit for the mission and clear checkpoints
- **Usage**: `m checkpoint commit -m "<message>"`
- **Flags**:
  - `-m`: Commit message (supports multi-line)
- **Output**: Creates consolidated commit and clears checkpoint tags

## Logging (`m log`)

### `m log`
- **Purpose**: Log messages to mission execution log
- **Usage**: `m log --step "<step>" "<message>"`
- **Flags**:
  - `--step`: Name of the current execution step - REQUIRED
  - `--level`: Log level (INFO, SUCCESS, ERROR) - Default: INFO
- **Output**: Appends formatted message to `.mission/execution.log`

## Notes

- All commands that output JSON should be parsed by the AI for decision-making
- Commands with `--item` flags require exact text matching
- The `m backlog resolve` command uses in-place marking with [RESOLVED] prefix
- Backlog items marked with [RESOLVED] remain visible in `m backlog list --type refactor`
