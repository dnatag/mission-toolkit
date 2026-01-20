# Mission Toolkit CLI Reference

This document is the master reference for all `m` CLI commands available for use in prompts.

## Core Commands

### `m init`
- **Purpose**: Initialize Mission Toolkit project with templates for specified AI type
- **Usage**: `m init --ai <type>`
- **Flags**:
  - `--ai`: AI assistant type (required). Supported: q, claude, kiro, opencode
- **Output**: Creates `.mission/` directory with governance files and AI-specific prompt templates
- **Note**: Automatically initializes Git repository if not found

### `m dashboard`
- **Purpose**: Display comprehensive mission dashboard with execution logs in interactive TUI
- **Usage**: `m dashboard`
- **Features**:
  - Split-pane view: mission.md | execution.log | commit.msg (for completed)
  - Live refresh for active mission execution logs
  - Lazy loading for completed mission logs and commit messages
  - Keyboard navigation: ↑/↓ navigate, Enter view details, Tab switch panes, / search, q quit

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

### `m docs`
- **Purpose**: Generate CLI documentation schema
- **Usage**: `m docs --format json`
- **Flags**:
  - `--format`: Output format (json)
- **Output**: JSON schema of all CLI commands

## Mission Management (`m mission`)

### `m mission check`
- **Purpose**: Check mission state and validate artifacts
- **Usage**: `m mission check [--context <context>]`
- **Flags**:
  - `--context`: Context for validation (apply or complete)
- **Output (JSON)**:
  - `next_step`: "PROCEED" or "STOP"
  - `message`: Status message
  - `mission_exists`: Boolean

### `m mission id`
- **Purpose**: Get or create mission ID
- **Usage**: `m mission id`
- **Output**: Mission ID string (format: Track-Type-Timestamp)

### `m mission create`
- **Purpose**: Create mission.md with intent
- **Usage**: `m mission create --intent "<intent>"`
- **Flags**:
  - `--intent`: Intent text for initial mission creation
- **Output**: Creates `.mission/mission.md`

### `m mission update`
- **Purpose**: Update mission status or sections
- **Usage**: `m mission update [flags]`
- **Flags**:
  - `--status`: New mission status (planned, active, executed, completed, failed)
  - `--section`: Section to update (intent, verification, scope, plan)
  - `--item`: Items for list sections (can be repeated)
  - `--content`: Content for text sections
  - `--frontmatter`: Frontmatter key=value pairs
  - `--append`: Append items instead of replacing (default: false)
- **Output**: Updates `.mission/mission.md`

### `m mission finalize`
- **Purpose**: Validate and display mission.md for review
- **Usage**: `m mission finalize`
- **Output (JSON)**:
  - `action`: "PROCEED" or "INVALID"
  - `errors`: List of validation errors (if any)

### `m mission archive`
- **Purpose**: Archive mission files to completed directory and clean up obsolete files
- **Usage**: `m mission archive [--force]`
- **Flags**:
  - `--force`: Forcefully archive or no-op if no mission exists
- **Output**: Moves mission files to `.mission/completed/<MISSION_ID>-*`

### `m mission mark-complete`
- **Purpose**: Mark a plan step as complete and log progress
- **Usage**: `m mission mark-complete --step <N> --status <STATUS> --message "<message>"`
- **Flags**:
  - `--step`: Step number to mark as complete (required)
  - `--status`: Status level for logging (INFO, SUCCESS, FAILED, etc.) - default: INFO
  - `--message`: Message to log for this step
- **Output**: Updates plan step checkbox and logs message

### `m mission pause`
- **Purpose**: Pause current mission and save to .mission/paused/ folder
- **Usage**: `m mission pause`
- **Output**: Moves mission to `.mission/paused/<timestamp>-mission.md`

### `m mission restore`
- **Purpose**: Restore a paused mission from .mission/paused/ folder
- **Usage**: `m mission restore [mission-id]`
- **Output**: Restores specified or most recent paused mission to active state

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

### `m analyze decompose`
- **Purpose**: Provide decomposition analysis template for Track 4 epics
- **Usage**: `m analyze decompose`
- **Output**: Decomposition template for breaking down epic intents into sub-missions

## Backlog Management (`m backlog`)

### `m backlog list`
- **Purpose**: List backlog items
- **Usage**: `m backlog list [--include <types>] [--exclude <types>]`
- **Flags**:
  - `--include`: Include only these types (decomposed, refactor, future, completed)
  - `--exclude`: Exclude these types (decomposed, refactor, future, completed)
- **Output**: List of backlog items filtered by type

### `m backlog add`
- **Purpose**: Add one or more backlog items
- **Usage**: `m backlog add "<description>" [<description>...] --type <type>`
- **Flags**:
  - `--type`: Item type (decomposed, refactor, future) - REQUIRED
  - `--pattern-id`: Pattern ID for Rule-of-Three tracking (refactor type only)
- **Output**: Confirmation message

### `m backlog complete`
- **Purpose**: Mark a backlog item as complete
- **Usage**: `m backlog complete --item "<exact text>"`
- **Flags**:
  - `--item`: Exact text of the item to complete - REQUIRED
- **Output**: Moves item to COMPLETED section with timestamp

### `m backlog cleanup`
- **Purpose**: Remove completed items from the backlog
- **Usage**: `m backlog cleanup [--type <type>]`
- **Flags**:
  - `--type`: Filter by item type (decomposed, refactor, future)
- **Output**: Number of items removed

## Checkpoint Management (`m checkpoint`)

### `m checkpoint create`
- **Purpose**: Create a checkpoint of current working directory state
- **Usage**: `m checkpoint create`
- **Output**: Creates checkpoint in `.mission/checkpoints/<MISSION_ID>-<timestamp>/`

### `m checkpoint restore`
- **Purpose**: Restore working directory to specified checkpoint
- **Usage**: `m checkpoint restore <checkpoint> [--all]`
- **Args**: Checkpoint name (e.g., "MISSION-123-20230101-120000")
- **Flags**:
  - `--all`: Restore all mission changes (default: false)
- **Output**: Restores files from checkpoint

### `m checkpoint clear`
- **Purpose**: Clear all checkpoints for current mission
- **Usage**: `m checkpoint clear`
- **Output**: Removes checkpoint directory

### `m checkpoint commit`
- **Purpose**: Create final commit for the mission and clear checkpoints
- **Usage**: `m checkpoint commit -m "<message>"`
- **Flags**:
  - `-m, --message`: Commit message (supports multi-line) - REQUIRED
- **Output**: Creates consolidated commit and clears checkpoint tags

## Logging (`m log`)

### `m log`
- **Purpose**: Log messages to mission execution log
- **Usage**: `m log [message] [--step "<step>"] [--level <level>] [--file <path>]`
- **Flags**:
  - `--step, -s`: Mission step name (default: "General")
  - `--level, -l`: Log level (DEBUG, INFO, WARN, ERROR, SUCCESS) - default: INFO
  - `--file, -f`: Log file path (default: ".mission/execution.log", empty string for console only)
- **Output**: Appends formatted message to execution log

## Notes

- All commands that output JSON should be parsed by the AI for decision-making
- Commands with `--item` flags require exact text matching
- The `m backlog` commands use include/exclude filters instead of --all flag
- Mission status values: planned, active, executed, completed, failed
- Log levels: DEBUG, INFO, WARN, ERROR, SUCCESS
