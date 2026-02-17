# Backlog Integration Strategy

## Overview

Mission Toolkit uses `.mission/backlog.json` as the **primary backlog management tool** for tracking features, bugs, refactoring tasks, and decomposed intents. It also provides **JSON export** for integration with external tools.

**Key Principle:** backlog.json is mission-toolkit's backlog management tool. Users can add items manually, and mission-toolkit generates items automatically (epic decomposition). Pattern tracking for WET→DRY is determined by code analysis, not backlog state.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Mission Toolkit (.mission/)                     │
│  ┌─────────────┐  ┌──────────────┐                          │
│  │ mission.md  │  │ backlog.json │                          │
│  │ (current)   │  │ (backlog)    │                          │
│  └─────────────┘  └──────────────┘                          │
│                         │                                    │
│                         │ Primary backlog:                   │
│                         │ - Features                         │
│                         │ - Bugs                             │
│                         │ - Refactoring (patterns)           │
│                         │ - Decomposed intents               │
└─────────────────────────┼────────────────────────────────────┘
                          │
                          │ Optional: Export to external tools
                          │
                          ▼
                ┌──────────────────────┐
                │  External Tools      │
                │  - Beads             │
                │  - GitHub Issues     │
                │  - Jira, etc.        │
                └──────────────────────┘
```

## Core Concepts

### Separation of Concerns

**Mission State (mission.md):**
- Current work in progress
- Single active mission at a time
- Managed by `m mission` commands
- Lifecycle: planned → active → completed

**Backlog State (backlog.json):**
- Primary backlog for mission-toolkit
- Features, bugs, refactoring tasks, decomposed intents
- Managed by `m backlog` commands
- Structured JSON format

### Backlog Management

**How Items Are Added:**
1. **Manually** - User adds features, bugs via `m backlog add`
2. **Automatically** - Mission-toolkit generates:
   - Epic decomposition (Track 4) - Sub-intents with dependencies

**Note:** Pattern tracking (WET→DRY) is determined by code analysis during `m analyze duplication`, not by backlog state. See Pattern Tracking section below.

**Pattern Tracking (Rule of Three):**
- During `m analyze duplication`, patterns are detected
- Refactor items added to backlog.json with pattern_id
- Count patterns in backlog.json
- When count >= 3 → DRY mission triggered

**Epic Decomposition:**
- During `@m.plan`, Track 4 epics are decomposed
- Sub-intents added to backlog.json with dependencies

### Pattern Tracking (WET→DRY)

**Pattern tracking is self-contained in code analysis, not backlog state.**

**How It Works:**
1. During `@m.plan`, `m analyze duplication` runs
2. AI uses duplication.md template to scan codebase
3. AI finds similar code patterns and counts occurrences in actual code
4. AI determines mission type based on count:
   - count < 3 → WET (allow duplication)
   - count >= 3 → DRY (refactor pattern)

**Example:**
```bash
# During @m.plan
$ m analyze duplication

# AI scans codebase, finds:
# - Email validation pattern in users/handler.go:45
# - Email validation pattern in products/handler.go:78
# - Email validation pattern in orders/handler.go:23
# Count: 3 → DRY threshold reached

# Mission type automatically set to DRY
# User can optionally add refactor item to backlog for tracking
```

**Key Point:** Pattern count comes from analyzing actual code, not from counting backlog items. The duplication.md template guides AI to perform self-contained code analysis.

### JSON Export (Optional)

**Export Flow:**
1. User exports backlog.json to external tool
2. Transform JSON for tool format (jq, scripts)
3. Import to external tool (Beads, GitHub, Jira, etc.)
4. User can work in either tool

**Why Export:**
- Share backlog with team using external tools
- Leverage external tool features (dependency graphs, etc.)
- Flexibility to use preferred tools

## Backlog Management

### Backlog Storage (backlog.json)

Mission Toolkit uses backlog.json as the primary backlog:

```json
{
  "version": "1.0",
  "items": [
    {
      "id": "feature-001",
      "title": "Add user authentication",
      "type": "feature",
      "status": "open",
      "description": "Implement JWT-based authentication system",
      "estimated_files": 5,
      "dependencies": [],
      "created_at": "2026-02-17T10:00:00Z"
    },
    {
      "id": "bug-001",
      "title": "Fix null pointer in login handler",
      "type": "bug",
      "status": "open",
      "description": "Null pointer exception when password is empty",
      "location": "auth/handler.go:45",
      "created_at": "2026-02-17T10:02:00Z"
    },
    {
      "id": "refactor-001",
      "title": "Extract validation logic",
      "type": "refactor",
      "status": "open",
      "description": "Extract email validation into shared utility",
      "estimated_files": 3,
      "location": "users/handler.go:45",
      "created_at": "2026-02-17T10:05:00Z"
    },
    {
      "id": "decomposed-001",
      "title": "Create status detection service",
      "type": "decomposed",
      "status": "open",
      "description": "Determines mission status from filesystem state",
      "rationale": "Core status detection logic is foundation for status command",
      "estimated_files": 3,
      "dependencies": [],
      "created_at": "2026-02-17T10:10:00Z"
    }
  ]
}
```

### Item Types

- **feature** - New functionality
- **bug** - Bug fix
- **refactor** - Code improvement (tracked for Rule of Three)
- **decomposed** - Sub-intent from epic decomposition

### Adding Items

**Manually:**
```bash
# Add feature
m backlog add "Add user authentication" --type feature --estimated-files 5

# Add bug
m backlog add "Fix null pointer in login" --type bug --location "auth/handler.go:45"

# Add refactor
m backlog add "Extract validation logic" --type refactor --location "users/handler.go:45"
```

**Automatically Generated:**
- **Pattern tracking** - During `m analyze duplication`
- **Epic decomposition** - During `@m.plan` for Track 4

### Viewing Backlog

```bash
# List items (pretty table)
m backlog list [--type TYPE] [--status STATUS]

# Show item details
m backlog show --id <id>
```

## JSON Export

Export generated backlog items to external tools.

### Export Command

```bash
# Export entire backlog
m backlog export > backlog.json

# Export specific types
m backlog export --type refactor > refactors.json
m backlog export --type decomposed > epic.json

# Export specific status
m backlog export --status open > open-items.json
```

### Export Output

```json
{
  "version": "1.0",
  "items": [
    {
      "id": "decomposed-001",
      "title": "Create status detection service",
      "type": "decomposed",
      "status": "open",
      "description": "Determines mission status from filesystem state",
      "rationale": "Core status detection logic is foundation",
      "estimated_files": 3,
      "dependencies": [],
      "created_at": "2026-02-17T10:00:00Z"
    },
    {
      "id": "decomposed-002",
      "title": "Add m status CLI command",
      "type": "decomposed",
      "status": "open",
      "description": "CLI command that outputs current mission status",
      "rationale": "CLI interface depends on detection logic",
      "estimated_files": 2,
      "dependencies": ["decomposed-001"],
      "created_at": "2026-02-17T10:01:00Z"
    }
  ]
}
```

### Transform for External Tools

**For Beads:**
```bash
# Export
m backlog export > backlog.json

# Transform and import to Beads
jq -r '.items[] | "bd create \"\(.title)\" --label type:\(.type) --description \"\(.description // "")\""' backlog.json | sh

# With dependencies (for decomposed items)
jq -r '.items[] | 
  if .dependencies | length > 0 then
    "bd create \"\(.title)\" --label type:\(.type) --depends-on \(.dependencies | join(","))"
  else
    "bd create \"\(.title)\" --label type:\(.type)"
  end' backlog.json | sh
```

**For GitHub:**
```bash
# Export
m backlog export > backlog.json

# Transform and import to GitHub
jq -r '.items[] | "gh issue create --title \"\(.title)\" --body \"\(.description // "")\" --label \(.type)"' backlog.json | sh
```

**For Jira (via CSV):**
```bash
# Export
m backlog export > backlog.json

# Transform to CSV
jq -r '["Summary","Description","Issue Type"], 
  (.items[] | [.title, .description, .type]) | @csv' backlog.json > backlog.csv

# Import CSV to Jira
```

## Complete Workflows

### Workflow 1: Epic Decomposition with Export

```bash
# 1. Plan epic intent
$ @m.plan "Implement swimlane status flow with m status command"

# Mission-toolkit detects Track 4, decomposes:
📊 Track 4 (Epic) - 15 files
🔄 Decomposing into sub-intents...
✅ Added 5 sub-intents to backlog.json

# 2. View decomposed items
$ m backlog list --type decomposed

# 3. Export to Beads
$ m backlog export --type decomposed > epic.json
$ jq -r '.items[] | "bd create \"\(.title)\" --label type:decomposed"' epic.json | sh

# 4. View in Beads with dependency graph
$ bd graph

# 5. Copy first sub-intent and start working
$ @m.plan "Create status detection service"
$ @m.apply
$ @m.complete
```

### Workflow 2: Pattern Tracking (Rule of Three)

```bash
# 1. First occurrence - WET mission
$ @m.plan "Add email validation to products"

# Duplication analysis scans codebase:
# - Found: users/handler.go:45 (existing)
# - Found: products/handler.go:78 (current intent)
# Count: 2 → WET (allow duplication)

# Mission proceeds as WET

# 2. Second occurrence - still WET
$ @m.plan "Add email validation to orders"

# Duplication analysis scans codebase:
# - Found: users/handler.go:45
# - Found: products/handler.go:78
# - Found: orders/handler.go:23 (current intent)
# Count: 3 → DRY threshold reached!

# Mission type automatically set to DRY
# AI creates refactoring plan to extract pattern

# 3. Optional: Add refactor items to backlog for tracking
$ m backlog add "Extract email validation pattern" --type refactor
```

### Workflow 3: AI Agent Managing Backlog

```bash
# AI agent reads backlog
$ cat .mission/backlog.json

# AI agent analyzes and prioritizes items
$ jq '.items | sort_by(.estimated_files) | .[0]' .mission/backlog.json
# Output: Smallest item to start with

# AI agent copies intent to start mission
$ @m.plan "Create status detection service"

# After completion, AI agent marks item as completed
$ jq '(.items[] | select(.id == "decomposed-001") | .status) = "completed"' .mission/backlog.json > tmp.json
$ mv tmp.json .mission/backlog.json
```

# Start mission from backlog
$ @m.plan "Add user authentication"
$ @m.apply
$ @m.complete

# Mark backlog item complete
$ m backlog complete --item "Add user authentication"
```

## File Structure

```
.mission/
├── mission.md          # Current mission state
├── backlog.json        # Future work queue (structured JSON)
├── execution.log       # Execution log
└── libraries/          # Templates and references
```

### mission.md (Current Work)

```yaml
---
id: 20260217105549-6632
intent: Create status detection service
status: planned
track: 2
type: WET
---

## SCOPE
- pkg/mission/status.go
- pkg/mission/status_test.go
- cmd/status.go

## PLAN
1. Create StatusDetector service
2. Implement filesystem state detection
3. Add unit tests
...
```

### backlog.json (Future Work)

```json
{
  "version": "1.0",
  "items": [
    {
      "id": "feature-001",
      "title": "Add user authentication",
      "type": "feature",
      "status": "open",
      "description": "Implement JWT-based authentication system",
      "estimated_files": 5,
      "pattern_id": null,
      "dependencies": [],
      "created_at": "2026-02-17T10:00:00Z"
    },
    {
      "id": "refactor-001",
      "title": "Extract validation logic",
      "type": "refactor",
      "status": "open",
      "description": "Extract email validation into shared utility",
      "estimated_files": 3,
      "pattern_id": "validation-pattern",
      "dependencies": [],
      "location": "users/handler.go:45",
      "created_at": "2026-02-17T10:05:00Z"
    },
    {
      "id": "refactor-002",
      "title": "Extract validation logic",
      "type": "refactor",
      "status": "open",
      "pattern_id": "validation-pattern",
      "location": "products/handler.go:78",
      "created_at": "2026-02-17T10:06:00Z"
    },
    {
      "id": "refactor-003",
      "title": "Extract validation logic",
      "type": "refactor",
      "status": "completed",
      "pattern_id": "validation-pattern",
      "location": "orders/handler.go:23",
      "created_at": "2026-02-17T10:07:00Z",
      "completed_at": "2026-02-17T11:00:00Z"
    },
    {
      "id": "decomposed-001",
      "title": "Create status detection service",
      "type": "decomposed",
      "status": "open",
      "description": "Determines mission status from filesystem state",
      "rationale": "Core status detection logic is foundation for status command",
      "estimated_files": 3,
      "dependencies": [],
      "created_at": "2026-02-17T10:10:00Z"
    },
    {
      "id": "decomposed-002",
      "title": "Add m status CLI command",
      "type": "decomposed",
      "status": "open",
      "description": "CLI command that outputs current mission status",
      "rationale": "CLI interface depends on detection logic",
      "estimated_files": 2,
      "dependencies": ["decomposed-001"],
      "created_at": "2026-02-17T10:11:00Z"
    }
  ]
}
```

## CLI Commands Reference

### Backlog Management

**m backlog add**
```bash
m backlog add "description" --type TYPE [OPTIONS]

# Options:
#   --estimated-files N      Estimated number of files
#   --location PATH:LINE     Source location
#   --description TEXT       Detailed description

# Examples
m backlog add "Add user authentication" --type feature --estimated-files 5
m backlog add "Extract validation logic" --type refactor --location "users/handler.go:45"
m backlog add "Fix null pointer" --type bug --location "auth/handler.go:45"
```

**m backlog list**
```bash
m backlog list [--type TYPE] [--status STATUS]

# Examples
m backlog list                    # All items
m backlog list --type refactor    # Only refactor items
m backlog list --status open      # Only open items
```

**m backlog show**
```bash
m backlog show --id ID

# Example
m backlog show --id feature-001
```

**m backlog complete**
```bash
m backlog complete --id ID

# Example
m backlog complete --id feature-001
```

**m backlog cleanup**
```bash
m backlog cleanup

# Removes all completed items
```

### Export Commands

**m backlog export**
```bash
m backlog export [--type TYPE] [--status STATUS] > output.json

# Examples
m backlog export > backlog.json
m backlog export --type bug > bugs.json
m backlog export --type decomposed > epic.json
m backlog export --status open > open-items.json
```

### Mission vs Backlog Operations

**Mission Step Completion (Current Work):**
```bash
m mission mark-complete --step 1 --status success --message "Added validation"
```
- Marks a step in current mission plan
- Updates mission.md execution section
- Used during @m.apply execution

**Backlog Item Completion (Future Work):**
```bash
m backlog complete --id feature-001
```
- Marks backlog item as completed in backlog.json
- Separate from mission step tracking
```
- Marks backlog item as done
- Marks [x] in backlog.md
- Separate from mission step tracking

## Implementation Notes

### JSON Export/Import

```go
type BacklogItem struct {
    Title       string   `json:"title"`
    Type        string   `json:"type"`
    Status      string   `json:"status"`
    PatternID   *string  `json:"pattern_id"`
    Dependencies []string `json:"dependencies"`
}

type BacklogJSON struct {
    Version string        `json:"version"`
    Items   []BacklogItem `json:"items"`
}
```

### Export Implementation

```go
func (s *BacklogService) ExportJSON() (*BacklogJSON, error) {
    // Read backlog.md
    content, err := os.ReadFile(".mission/backlog.md")
    if err != nil {
        return nil, err
    }
    
    // Parse markdown items
    items := parseBacklogItems(content)
    
    return &BacklogJSON{
        Version: "1.0",
        Items:   items,
    }, nil
}
```

### Import Implementation

```go
func (s *BacklogService) ImportJSON(file string, replace bool) error {
    // Read JSON file
    data, err := os.ReadFile(file)
    if err != nil {
        return err
    }
    
    var backlog BacklogJSON
    if err := json.Unmarshal(data, &backlog); err != nil {
        return err
    }
    
    if replace {
        // Clear backlog.md and write all items
        return s.replaceBacklog(backlog.Items)
    }
    
    // Merge: skip duplicates
    return s.mergeBacklog(backlog.Items)
}
```

### Pattern Count Query

```go
func (s *BacklogService) GetPatternCount(patternID string) (int, error) {
    // Read backlog.md
    content, err := os.ReadFile(".mission/backlog.md")
    if err != nil {
        return 0, err
    }
    
    // Count all lines with pattern (including [x] completed)
    pattern := fmt.Sprintf("[pattern:%s]", patternID)
    count := strings.Count(string(content), pattern)
    
    return count, nil
}
```

## Benefits of This Approach

### Simplicity
- ✅ Export-only (no import complexity)
- ✅ Single JSON format (universal)
- ✅ No merge logic, no duplicate detection
- ✅ Easy to test and maintain

### Clear Value
- ✅ Get structured data out (epic decomposition, patterns)
- ✅ Transform for any external tool
- ✅ User copies intent back when ready to work
- ✅ No bidirectional sync complexity

### Flexibility
- ✅ Works with any tool that accepts JSON
- ✅ Users choose their preferred backlog tool
- ✅ Easy to add custom transformations
- ✅ AI agents can manipulate backlog.json directly

### Standalone
- ✅ No external dependencies
- ✅ Works completely offline
- ✅ Simple JSON format
- ✅ Git-friendly

## Design Principles

1. **Unix Philosophy**: Do one thing well, provide standard output, let users compose tools
2. **Export-Only**: Clear value without import complexity
3. **JSON as Output Format**: Universal, AI-friendly, transformable
4. **User Copies Intent**: Simple workflow, no sync needed
5. **Minimal Complexity**: Simplest possible implementation

## Open Questions

1. **JSON Schema Versioning**: How to handle schema evolution?
2. **Backlog Cleanup**: Should completed items be auto-removed or kept for history?
3. **AI Agent Integration**: Should we provide helper scripts for common transformations?

## Related Documentation

- [Core Concepts](concepts.md) - WET→DRY workflow and Rule of Three
- [Workflows](workflows.md) - Mission lifecycle and epic decomposition
- [CLI Reference](cli-reference.md) - Complete command documentation
