# Mission Toolkit CLI Reference

This document is the master reference for all `m` CLI commands available for use in prompts.

## `m plan`

### `m plan check`
- **Purpose**: Checks if a new mission can be started. Cleans up stale artifacts.
- **Usage**: `m plan check`
- **Output (JSON)**:
  - `status`: "OK" or "ERROR"
  - `error`: Message if an active mission exists.
  - `mission_id`: The new ID for the session.

### `m plan analyze`
- **Purpose**: Calculates complexity and checks for missing tests.
- **Usage**: `m plan analyze --file .mission/plan.json`
- **Input**: `plan.json` with `scope` and `domain`.
- **Output (JSON)**:
  - `track`: Calculated track (1-4).
  - `recommendation`: "proceed" or "decompose".
  - `warnings`: List of issues (e.g., "Missing test file...").
- **Reference**: See `tools/complexity-reference.md` for logic.

### `m plan validate`
- **Purpose**: Performs security and safety checks.
- **Usage**: `m plan validate --file .mission/plan.json`
- **Input**: `plan.json` with `scope` and `verification`.
- **Output (JSON)**:
  - `valid`: `true` or `false`.
  - `errors`: List of critical issues.
  - `warnings`: List of non-critical issues.

### `m plan generate`
- **Purpose**: Generates the final `.mission/mission.md` file.
- **Usage**: `m plan generate --file .mission/plan.json`
- **Input**: A complete `plan.json`.
- **Output (JSON)**:
  - `success`: `true` or `false`.
  - `output_file`: Path to the generated mission.

## `m log`

### `m log`
- **Purpose**: Appends a formatted message to `.mission/execution.log`.
- **Usage**: `m log --step "Step Name" "Message content"`
- **Flags**:
  - `--level`: (Optional) "INFO", "SUCCESS", "ERROR". Default is "INFO".
  - `--step`: (Required) The name of the current execution step.
