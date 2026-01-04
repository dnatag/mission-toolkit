# Design: AI-Driven Backlog Management

**Status**: 🚧 Design complete, not implemented

## 1. Problem Statement

The mission-driven workflow is effective for executing discrete tasks, but it currently lacks a structured, persistent system for managing work that is identified but not immediately acted upon. This includes:

1.  **Decomposed Epics**: When a user provides a large-scale intent (an "epic"), the `m.plan` prompt correctly identifies it as too complex for a single mission. The resulting sub-tasks need to be captured so they can be executed as individual missions later.
2.  **Refactoring Opportunities**: During planning or implementation, the AI may identify technical debt, code duplication, or other refactoring opportunities. These valuable insights are currently lost if not immediately addressed.
3.  **Future Enhancements**: Users and developers need a simple way to jot down ideas for future features or improvements that are not part of the current mission.

Without a formal backlog, these valuable insights and planned tasks are either lost or must be managed manually outside the toolkit, creating friction and losing context.

## 2. Proposed Solution: An AI-Managed Backlog

We will introduce a structured `backlog.md` file within the `.mission` directory, managed by a new `m backlog` CLI command suite. This system is designed to be primarily orchestrated by the AI, which will use the CLI tools to intelligently add, complete, and refine backlog items.

The solution consists of three parts:
1.  A defined structure for the `.mission/backlog.md` file.
2.  A new `m backlog` command with subcommands for the AI to manage the backlog programmatically.
3.  Integration of these commands into the core AI prompts.

## 3. `backlog.md` File Structure

The `.mission/backlog.md` file will be a standard markdown file organized into managed sections.

```markdown
# Mission Backlog

## DECOMPOSED INTENTS
*This section lists atomic tasks that were broken down from a larger user-defined epic.*
- [ ] Sub-intent 1 from original epic.
- [ ] Sub-intent 2 from original epic.

## REFACTORING OPPORTUNITIES
*This section lists technical debt and refactoring opportunities identified by the AI during planning or execution.*
- [ ] Refactor the `NewService` function in `service.go` to use an interface.

## FUTURE ENHANCEMENTS
*This section is for user-defined ideas and future feature requests.*
- [ ] Add Prometheus metrics to the API server.

## COMPLETED
*History of completed backlog items.*
- [x] Refactor user validation logic (Completed: 2023-10-27)
```

## 4. CLI Command Design: `m backlog`

A new top-level command, `m backlog`, will be created to interact with the `backlog.md` file.

### `m backlog list`

Displays the contents of the backlog. This is the primary tool the AI uses to gain context before making changes.

**Usage:**
```bash
m backlog list [--all]
```

-   **Default**: Lists only *open* items (unchecked `[ ]`) from the active sections.
-   **`--all`**: Includes the `## COMPLETED` section.

### `m backlog add`

Adds a new item to the backlog. The AI is responsible for checking for duplicates before adding.

**Usage:**
```bash
m backlog add "Item description" --type <TYPE>
```

-   **`--type <TYPE>`**: (Required) Specifies the section. Valid types: `decomposed`, `refactor`, `future`.

### `m backlog complete`

Marks a specific item as completed by finding an exact string match.

**Usage:**
```bash
m backlog complete --item "The exact text of the backlog item to complete"
```

-   **`--item <string>`**: (Required) The exact text of the open backlog item to mark as complete.
-   **Action**:
    1.  Finds the line containing the exact item text.
    2.  Marks it with `[x]`.
    3.  Appends a timestamp: `(Completed: YYYY-MM-DD)`.
    4.  Moves the line to the `## COMPLETED` section.

## 5. AI-Driven Workflow Integration

The AI will be instructed to use a "list-then-act" pattern to manage the backlog.

### Adding an Item (with Duplication Check)

1.  **AI Identifies Item**: During `m.plan` or `m.apply`, the AI identifies a potential backlog item.
2.  **AI Lists Backlog**: The AI runs `m backlog list` to get all open items.
3.  **AI Checks for Duplicates**: The AI compares its new item against the existing list to check for semantic duplicates.
4.  **AI Adds Item**: If no duplicate is found, the AI runs `m backlog add "New item" --type <type>`.

### Completing an Item

1.  **AI Reviews Context**: After a mission (e.g., in `m.complete`), the AI is prompted to review the backlog.
2.  **AI Lists Backlog**: The AI runs `m backlog list`.
3.  **AI Identifies Completed Items**: Based on the work just done, the AI identifies which, if any, backlog items are now resolved.
4.  **AI Completes Items**: For each resolved item, the AI runs `m backlog complete --item "Exact text of the item"`.

### Editing an Item (via `complete` and `add`)

The AI can "edit" a backlog item by combining the `complete` and `add` commands.

1.  **AI Identifies Vague Item**: The AI runs `m backlog list` and decides an item is poorly worded.
2.  **AI Completes Vague Item**: It executes `m backlog complete --item "Vague item text"`.
3.  **AI Adds Specific Item**: It then executes `m backlog add "New, more specific item text" --type <type>`.

## 6. Success Metrics

-   **Intelligent Management**: The backlog is self-managing, with the AI handling duplication checks and completion sign-offs.
-   **Seamless Capture**: Work is captured automatically via AI or easily via CLI.
-   **Lifecycle Management**: Items have a clear state (Open -> Completed) and history.
-   **Visibility**: `m backlog list` provides a focused view of pending work.
-   **Simplicity**: The underlying file is standard Markdown, editable by hand if needed, and the CLI commands are simple and deterministic.
