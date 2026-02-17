# Backlog Integration Strategy

## Overview

Mission Toolkit uses `.mission/backlog.json` as the **primary backlog management tool** for tracking features, bugs, refactoring tasks, and decomposed intents. It also provides **JSON export** for integration with external tools.

**Key Principle:** backlog.json is the single source of truth for mission-toolkit's backlog. Users can add items manually, and mission-toolkit generates items automatically (epic decomposition, pattern tracking). Export capability allows integration with external tools.

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
   - Pattern tracking (duplication detection)
   - Epic decomposition (Track 4)

**Pattern Tracking (Rule of Three):**
- During `m analyze duplication`, patterns are detected
- Refactor items added to backlog.json with pattern_id
- Count patterns in backlog.json
- When count >= 3 → DRY mission triggered

**Epic Decomposition:**
- During `@m.plan`, Track 4 epics are decomposed
- Sub-intents added to backlog.json with dependencies

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
      "pattern_id": "validation-pattern",
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

# Add refactor (with pattern tracking)
m backlog add "Extract validation logic" --type refactor --pattern-id validation-pattern --location "users/handler.go:45"
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

# Count pattern occurrences
m backlog count --pattern-id <pattern-id>
```

### Pattern Tracking (Rule of Three)

```bash
# During @m.plan, duplication analysis runs
$ m analyze duplication

# Finds patterns, adds to backlog.json:
# - pattern_id: "validation-pattern"
# - count: 2

# View pattern occurrences
$ m backlog list --type refactor --pattern-id validation-pattern
Pattern: validation-pattern (count: 2)
- refactor-001: Extract validation logic [users/handler.go:45]
- refactor-002: Extract validation logic [products/handler.go:78]

# Third occurrence triggers DRY
$ @m.plan "Add email validation to orders"
# Mission type: DRY (refactor pattern)
# Automatically adds refactor-003 to backlog.json
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
# 1. During planning, duplication detected
$ @m.plan "Add email validation to products"

# Duplication analysis finds similar code in users/handler.go
# Automatically adds to backlog.json:
# - pattern_id: "validation-pattern"
# - count: 2 (not yet at threshold)

# 2. View pattern occurrences
$ m backlog list --type refactor --pattern-id validation-pattern
Pattern: validation-pattern (count: 2)
- refactor-001: Extract validation logic [users/handler.go:45]
- refactor-002: Extract validation logic [products/handler.go:78]

# 3. Third occurrence triggers DRY
$ @m.plan "Add email validation to orders"
# Mission type: DRY (refactor pattern)
# Automatically adds refactor-003 to backlog.json

# 4. Export refactor tasks
$ m backlog export --type refactor --pattern-id validation-pattern > refactors.json
$ jq -r '.items[] | "bd create \"\(.title)\" --label type:refactor --label pattern:\(.pattern_id)"' refactors.json | sh
```

### Workflow 3: AI Agent Managing Backlog

```bash
# AI agent reads generated backlog
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
m backlog add "description" [--type TYPE] [--pattern-id ID]

# Examples
m backlog add "Add user authentication" --type feature
m backlog add "Extract validation logic" --type refactor --pattern-id validation-pattern
m backlog add "Fix null pointer" --type bug
```

**m backlog list**
```bash
m backlog list [--include TYPE] [--exclude TYPE] [--pattern-id ID]

# Examples
m backlog list                                    # All items
m backlog list --include refactor                 # Only refactor items
m backlog list --pattern-id validation-pattern    # Specific pattern
```

**m backlog complete**
```bash
m backlog complete --item "description"

# Example
m backlog complete --item "Add user authentication"
```

**m backlog cleanup**
```bash
m backlog cleanup [--type TYPE]

# Examples
m backlog cleanup              # Remove all completed items
m backlog cleanup --type bug   # Remove completed bugs only
```

### Export Commands

**m backlog export**
```bash
m backlog export --format json [--include TYPE] [--decomposed]

# Examples
m backlog export --format json > backlog.json
m backlog export --format json --include bug > bugs.json
m backlog export --format json --decomposed > epic.json
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
m backlog complete --item "Add user authentication"
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
