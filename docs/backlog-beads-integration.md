# Backlog Integration Strategy

## Overview

Mission Toolkit manages a structured backlog in `.mission/backlog.json` and provides **JSON-based import/export** for integration with external tools.

**Key Principle:** JSON as source of truth. Perfect roundtrip, AI-friendly, no data loss.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Mission Toolkit (.mission/)                     │
│  ┌─────────────┐  ┌──────────────┐                          │
│  │ mission.md  │  │ backlog.json │                          │
│  │ (current)   │  │ (future)     │                          │
│  └─────────────┘  └──────────────┘                          │
└──────────┬────────────────────────────────────┬─────────────┘
           │                                    │
           │ JSON Export / Import               │
           │ (Perfect roundtrip)                │
           │                                    │
    ┌──────▼──────────┐              ┌─────────▼─────────┐
    │  Beads          │              │ GitHub Issues     │
    │  (bd CLI)       │              │ (gh CLI)          │
    │  bd → JSON      │              │ gh → JSON         │
    └─────────────────┘              └───────────────────┘
```

## Core Concepts

### Separation of Concerns

**Mission State (mission.md):**
- Current work in progress
- Single active mission at a time
- Managed by `m mission` commands
- Lifecycle: planned → active → completed

**Backlog State (backlog.json):**
- Future work queue
- Multiple items tracked
- Managed by `m backlog` commands
- Structured JSON format

### JSON as Source of Truth

**Why JSON:**
- Perfect roundtrip (no data loss)
- AI-friendly (structured data)
- Easy to query and transform
- Preserves all metadata (rationale, dependencies, estimates)

**Human Readability:**
- `m backlog list` - Pretty table view
- `m backlog show --id <id>` - Detailed item view
- JSON is readable enough for quick edits

### Import/Export

**Export:**
- Export backlog.json to external tools
- Transform as needed (jq, scripts, AI agents)

**Import:**
- Import from external tools to backlog.json
- Merge or replace mode
- No data loss on roundtrip

## Backlog Management

### Backlog Storage (backlog.json)

Mission Toolkit stores backlog in structured JSON format:

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
      "id": "decomposed-001",
      "title": "Create status detection service",
      "type": "decomposed",
      "status": "open",
      "description": "Determines mission status from filesystem state",
      "rationale": "Core status detection logic is foundation for status command",
      "estimated_files": 3,
      "pattern_id": null,
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
      "pattern_id": null,
      "dependencies": ["decomposed-001"],
      "created_at": "2026-02-17T10:11:00Z"
    }
  ]
}
```

### Item Types

- **feature** - New functionality
- **refactor** - Code improvement (tracked for Rule of Three)
- **bug** - Bug fix
- **decomposed** - Sub-intent from epic decomposition
- **docs** - Documentation
- **test** - Test improvements

### Basic Commands

```bash
# Add items
m backlog add "description" --type TYPE [--pattern-id ID] [--estimated-files N]

# List items (pretty table)
m backlog list [--type TYPE] [--status STATUS]

# Show item details
m backlog show --id <id>

# Complete items
m backlog complete --id <id>

# Clean up completed items
m backlog cleanup
```

### Pattern Tracking (Rule of Three)

Track code duplication patterns to trigger DRY refactoring:

```bash
# Add pattern occurrence
m backlog add "Extract validation logic" \
  --type refactor \
  --pattern-id validation-pattern \
  --location "users/handler.go:45"

# List pattern occurrences
m backlog list --type refactor --pattern-id validation-pattern

# Count pattern occurrences
m backlog count --pattern-id validation-pattern
# Output: 3

# When count >= 3, next mission becomes DRY type
```

**Pattern count logic:**
```bash
# Count all items with pattern_id (including completed)
jq '[.items[] | select(.pattern_id == "validation-pattern")] | length' .mission/backlog.json
```

## JSON Export

Export backlog to external tools.

### Export Command

```bash
# Export entire backlog
m backlog export > backlog.json

# Export specific types
m backlog export --type refactor > refactors.json

# Export specific status
m backlog export --status open > open-items.json
```

**Note:** Since backlog is already JSON, export is essentially a copy with optional filtering.

### Export Output

Same format as backlog.json:

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
```

**For GitHub:**
```bash
# Export
m backlog export > backlog.json

# Transform and import to GitHub
jq -r '.items[] | "gh issue create --title \"\(.title)\" --body \"\(.description // "")\" --label \(.type)"' backlog.json | sh
```

### Epic Decomposition Export

When Track 4 epic is decomposed, items include dependencies:

```bash
# After epic decomposition
m analyze decompose
  → Adds decomposed items to backlog.json

# Export decomposed items
m backlog export --type decomposed > epic.json
```

**Output includes dependency information:**
```json
{
  "version": "1.0",
  "items": [
    {
      "id": "decomposed-001",
      "title": "Create status detection service",
      "type": "decomposed",
      "dependencies": [],
      "rationale": "Core status detection logic is foundation"
    },
    {
      "id": "decomposed-002",
      "title": "Add m status CLI command",
      "type": "decomposed",
      "dependencies": ["decomposed-001"],
      "rationale": "CLI interface depends on detection logic"
    }
  ]
}
```

## Import Integration

Import issues from external tools to start missions.

### Import Commands

```bash
# Import from Beads
m backlog import --from beads --id catalyst-toolkit-jwq.9

# Import from GitHub
m backlog import --from github --id 123

# Import from JSON file
m backlog import --from json --file issue.json
```

### Import Workflow

**Example: Start mission from Beads issue**

```bash
# 1. Browse issues in Beads
$ bd list
○ catalyst-toolkit-jwq.9 [● P2] [task] - Create status detection service

# 2. Import to mission-toolkit
$ m backlog import --from beads --id catalyst-toolkit-jwq.9

# Mission-toolkit:
# - Fetches issue details via: bd show catalyst-toolkit-jwq.9
# - Creates mission.md with intent from issue title/description
# - Adds external_link metadata to track origin
# - Ready for @m.plan workflow

# 3. Continue normal workflow
@m.plan
@m.apply
@m.complete
```

### Mission Metadata (External Link)

When importing, mission.md includes external link:

```yaml
---
id: 20260217105549-6632
intent: Create status detection service that determines mission status from filesystem state
status: planned
track: 2
external_link:
  source: beads
  id: catalyst-toolkit-jwq.9
---
```

### Adapter Pattern

Each external tool has an adapter:

```go
type BacklogAdapter interface {
    Import(id string) (*BacklogItem, error)
    Export(items []BacklogItem) (string, error)
    Sync(id string, status string) error
}

type BeadsAdapter struct {}
type GitHubAdapter struct {}
type JSONAdapter struct {}
```

**Beads Adapter:**
```go
func (a *BeadsAdapter) Import(id string) (*BacklogItem, error) {
    // Execute: bd show <id> --format json
    // Parse output
    // Return BacklogItem
}
```

## JSON Import

Import backlog items from external tools.

### Import Command

```bash
# Import and merge (skip duplicates by ID)
m backlog import --file backlog.json

# Import and replace entire backlog
m backlog import --file backlog.json --replace
```

### Import Behavior

**Merge Mode (Default):**
- Adds new items to backlog.json
- Skips items with duplicate IDs
- Updates items if ID exists (optional: use --update flag)
- Preserves existing items

**Replace Mode:**
- Clears backlog.json
- Adds all items from import file
- Use with caution

### Import from External Tools

**From Beads:**
```bash
# Export from Beads to JSON
bd list --format json > issues.json

# Transform to mission-toolkit format
jq '{
  version: "1.0",
  items: [.[] | {
    id: ("beads-" + .id),
    title: .title,
    type: (.labels[] | select(startswith("type:")) | ltrimstr("type:")),
    status: (if .status == "closed" then "completed" else "open" end),
    description: .description,
    estimated_files: null,
    pattern_id: (.labels[] | select(startswith("pattern:")) | ltrimstr("pattern:") // null),
    dependencies: (.depends_on // []),
    created_at: .created_at
  }]
}' issues.json > backlog.json

# Import
m backlog import --file backlog.json
```

**From GitHub:**
```bash
# Export from GitHub to JSON
gh issue list --json number,title,labels,body,createdAt --limit 100 > issues.json

# Transform to mission-toolkit format
jq '{
  version: "1.0",
  items: [.[] | {
    id: ("github-" + (.number | tostring)),
    title: .title,
    type: (.labels[0].name // "feature"),
    status: "open",
    description: .body,
    estimated_files: null,
    pattern_id: null,
    dependencies: [],
    created_at: .createdAt
  }]
}' issues.json > backlog.json

# Import
m backlog import --file backlog.json
```

### Roundtrip: Export → Edit → Import

```bash
# 1. Export current backlog
m backlog export > backlog.json

# 2. Edit with jq or text editor
jq '.items |= map(select(.type == "bug"))' backlog.json > bugs-only.json

# 3. Import back (replace mode)
m backlog import --file bugs-only.json --replace
```

### AI Agent Workflow

Perfect for AI agents managing backlog:

```bash
# AI agent reads backlog
cat .mission/backlog.json

# AI agent modifies (adds, updates, removes items)
jq '.items += [{
  id: "ai-generated-001",
  title: "Refactor authentication flow",
  type: "refactor",
  status: "open",
  description: "Detected duplication in auth handlers",
  pattern_id: "auth-pattern",
  dependencies: [],
  created_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}]' .mission/backlog.json > .mission/backlog.tmp.json

# AI agent writes back
mv .mission/backlog.tmp.json .mission/backlog.json
```

## Complete Workflows

### Workflow 1: Epic Decomposition with Export

```bash
# 1. Plan epic intent
$ @m.plan "Implement swimlane status flow with m status command"

# Mission-toolkit detects Track 4, decomposes:
📊 Track 4 (Epic) - 15 files
🔄 Decomposing into sub-intents...

# 2. Review decomposition
$ cat .mission/decomposed.md
1. Create status detection service
2. Add m status CLI command
3. Add status file management
...

# 3. Export to Beads
$ m backlog export --format beads --decomposed | sh

# Creates issues in Beads with dependencies

# 4. View in Beads
$ bd graph catalyst-toolkit-jwq
[Dependency graph visualization]

# 5. Import first sub-intent
$ m backlog import --from beads --id catalyst-toolkit-jwq.9
$ @m.plan
$ @m.apply
$ @m.complete
$ m backlog sync --complete
```

### Workflow 2: Bulk Import from External Tool

```bash
# 1. Export from Beads
$ bd list --format json > issues.json

# 2. Transform to mission-toolkit format (AI agent or script)
$ jq '{version: "1.0", items: [...]}' issues.json > backlog.json

# 3. Import to mission-toolkit
$ m backlog import --file backlog.json

# 4. View imported items
$ m backlog list

# 5. Start working on items
$ @m.plan "Add user authentication"
$ @m.apply
$ @m.complete
```

### Workflow 3: AI Agent Managing Backlog

```bash
# AI agent reads current backlog
$ cat .mission/backlog.json

# AI agent analyzes codebase, detects patterns
# AI agent adds refactor items programmatically
$ jq '.items += [{
  id: "refactor-'$(date +%s)'",
  title: "Extract database connection logic",
  type: "refactor",
  pattern_id: "db-connection",
  location: "users/db.go:15",
  estimated_files: 3,
  dependencies: [],
  created_at: "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}]' .mission/backlog.json > .mission/backlog.tmp.json && mv .mission/backlog.tmp.json .mission/backlog.json

# AI agent checks pattern count
$ m backlog count --pattern-id db-connection
# Output: 3

# AI agent triggers DRY mission
$ @m.plan "Refactor database connection pattern"
```

### Workflow 3: Pattern Tracking (Rule of Three)

```bash
# 1. During planning, duplication detected
$ @m.plan "Add email validation to products"

# Duplication analysis finds similar code in users/handler.go

# 2. Add pattern occurrence
$ m backlog add "Extract validation logic" --type refactor --pattern-id validation-pattern

# 3. Check pattern count
$ m backlog list --include refactor --pattern-id validation-pattern
Pattern: validation-pattern (count: 2)

# 4. Third occurrence triggers DRY
$ @m.plan "Add email validation to orders"
# Mission type: DRY (refactor pattern)

# 5. Export refactor tasks
$ m backlog export --format beads --include refactor
bd create "Extract validation logic" --label type:refactor --label pattern:validation-pattern
```

### Workflow 4: Standalone (No External Tool)

```bash
# Works completely standalone
$ m backlog add "Add user authentication" --type feature
$ m backlog add "Fix login bug" --type bug
$ m backlog list

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
#   --pattern-id ID          Pattern identifier for refactor items
#   --estimated-files N      Estimated number of files
#   --location PATH:LINE     Source location for refactor items
#   --description TEXT       Detailed description
#   --dependencies ID,ID     Comma-separated dependency IDs

# Examples
m backlog add "Add user authentication" --type feature --estimated-files 5
m backlog add "Extract validation logic" --type refactor --pattern-id validation-pattern --location "users/handler.go:45"
m backlog add "Fix null pointer" --type bug
```

**m backlog list**
```bash
m backlog list [OPTIONS]

# Options:
#   --type TYPE              Filter by type
#   --status STATUS          Filter by status (open, completed)
#   --pattern-id ID          Filter by pattern ID

# Examples
m backlog list                                    # All items (table view)
m backlog list --type refactor                    # Only refactor items
m backlog list --pattern-id validation-pattern    # Specific pattern
m backlog list --status open                      # Only open items
```

**m backlog show**
```bash
m backlog show --id ID

# Example
m backlog show --id refactor-001
# Output: Full item details including description, rationale, dependencies
```

**m backlog count**
```bash
m backlog count --pattern-id ID

# Example
m backlog count --pattern-id validation-pattern
# Output: 3
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

# Removes all completed items from backlog.json
```

### Export Commands

**m backlog export**
```bash
m backlog export [OPTIONS] > output.json

# Options:
#   --type TYPE              Export specific type only
#   --status STATUS          Export specific status only

# Examples
m backlog export > backlog.json
m backlog export --type refactor > refactors.json
m backlog export --status open > open-items.json
```

### Import Commands

**m backlog import**
```bash
m backlog import --file FILE [--replace]

# Examples
m backlog import --file backlog.json          # Merge (skip duplicate IDs)
m backlog import --file backlog.json --replace # Replace entire backlog
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
- ✅ Single JSON format (universal)
- ✅ No tool-specific adapters
- ✅ Simple export/import commands
- ✅ Easy to test and maintain

### User Control
- ✅ User handles transformations (jq, scripts)
- ✅ Explicit import/export operations
- ✅ No automatic synchronization
- ✅ Clear data ownership

### Flexibility
- ✅ Works with any tool that supports JSON
- ✅ Users choose their preferred backlog tool
- ✅ Easy to add custom transformations
- ✅ Supports bulk operations

### Standalone
- ✅ No external dependencies
- ✅ Works completely offline
- ✅ Simple markdown format
- ✅ Git-friendly

## Design Principles

1. **Unix Philosophy**: Do one thing well, provide standard formats, let users compose tools
2. **Explicit Over Implicit**: User triggers operations, no hidden behavior
3. **JSON as Universal Format**: Single format, user transforms as needed
4. **Graceful Degradation**: Full functionality without external tools
5. **Minimal Complexity**: Simplest possible implementation

## Open Questions

1. **JSON Schema Versioning**: How to handle schema evolution?
2. **Duplicate Detection**: Match by title only, or include other fields?
3. **Bulk Operations**: Should we support batch import/export?
4. **Validation**: Should we validate JSON schema on import?

## Related Documentation

- [Core Concepts](concepts.md) - WET→DRY workflow and Rule of Three
- [Workflows](workflows.md) - Mission lifecycle and epic decomposition
- [CLI Reference](cli-reference.md) - Complete command documentation
