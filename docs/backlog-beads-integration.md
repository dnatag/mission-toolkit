# Backlog Integration Strategy

## Overview

Mission Toolkit uses `.mission/backlog.md` as the **primary backlog management tool** for tracking features, bugs, refactoring tasks, and decomposed intents. It also provides **export capability** for integration with external tools.

**Key Principle:** backlog.md is mission-toolkit's backlog with embedded template instructions for AI. Users interact naturally ("add X to backlog"), and AI maintains the structured format. Pattern tracking for WET→DRY is determined by code analysis, not backlog state.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Mission Toolkit (.mission/)                     │
│  ┌─────────────┐  ┌──────────────┐                          │
│  │ mission.md  │  │ backlog.md   │                          │
│  │ (current)   │  │ (backlog)    │                          │
│  └─────────────┘  └──────────────┘                          │
│                         │                                    │
│                         │ Embedded template guides AI        │
│                         │ User: "add X to backlog"           │
│                         │ AI: Maintains structured format    │
│                         │                                    │
│                         │ Contains:                          │
│                         │ - Features                         │
│                         │ - Bugs                             │
│                         │ - Refactoring tasks                │
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

**Backlog State (backlog.md):**
- Primary backlog for mission-toolkit
- Features, bugs, refactoring tasks, decomposed intents
- Embedded template guides AI to maintain structure
- Human-readable markdown format

### Backlog Management

**How Items Are Added:**
1. **AI-Assisted** - User says "add X to backlog", AI follows embedded template
2. **Automatically** - Mission-toolkit generates:
   - Epic decomposition (Track 4) - Sub-intents with dependencies

**Note:** Pattern tracking (WET→DRY) is determined by code analysis during `m analyze duplication`, not by backlog state. See Pattern Tracking section below.

**Epic Decomposition:**
- During `@m.plan`, Track 4 epics are decomposed
- Sub-intents added to backlog.md with dependencies

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

### Export to External Tools (Optional)

**Export Flow:**
1. AI parses backlog.md and converts to JSON
2. User transforms JSON for external tool (jq, scripts)
3. Import to external tool (Beads, GitHub, Jira, etc.)

**Why Export:**
- Share backlog with team using external tools
- Leverage external tool features (dependency graphs, etc.)
- Flexibility to use preferred tools

## Backlog Management

### Backlog Storage (backlog.md)

Mission Toolkit uses backlog.md with embedded template instructions:

```markdown
---
# Backlog Template Instructions (for AI)
# When user says "add X to backlog" or similar:
# 1. Analyze the intent
# 2. Determine type (feature/bug/refactor/decomposed)
# 3. Estimate files if possible
# 4. Add item below in the appropriate section
# 5. Generate unique ID (type-timestamp format)
# 6. Set created timestamp
---

# Mission Backlog

## Features
- [ ] feature-20260217-001: Add user authentication
  - **Description**: Implement JWT-based authentication system
  - **Estimated Files**: 5
  - **Created**: 2026-02-17T10:00:00Z

- [ ] feature-20260217-002: Implement search functionality
  - **Description**: Full-text search across all entities
  - **Estimated Files**: 3
  - **Created**: 2026-02-17T10:05:00Z

## Bugs
- [ ] bug-20260217-001: Fix null pointer in login handler
  - **Description**: Null pointer exception when password is empty
  - **Location**: auth/handler.go:45
  - **Created**: 2026-02-17T10:10:00Z

## Refactoring
- [ ] refactor-20260217-001: Extract validation logic
  - **Description**: Extract email validation into shared utility
  - **Location**: users/handler.go:45
  - **Estimated Files**: 3
  - **Created**: 2026-02-17T10:15:00Z

## Decomposed Intents (from Track 4 Epics)
- [ ] decomposed-20260217-001: Create status detection service
  - **Description**: Determines mission status from filesystem state
  - **Rationale**: Core status detection logic is foundation
  - **Estimated Files**: 3
  - **Dependencies**: None
  - **Created**: 2026-02-17T10:20:00Z

- [ ] decomposed-20260217-002: Add m status CLI command
  - **Description**: CLI command that outputs current mission status
  - **Rationale**: CLI interface depends on detection logic
  - **Estimated Files**: 2
  - **Dependencies**: decomposed-20260217-001
  - **Created**: 2026-02-17T10:21:00Z
```

### Item Types

- **feature** - New functionality
- **bug** - Bug fix
- **refactor** - Code improvement
- **decomposed** - Sub-intent from epic decomposition

### Adding Items (AI-Assisted)

**Natural language interaction:**
```
User: "Add user authentication to the backlog"

AI: 
- Reads backlog.md
- Follows embedded template instructions
- Analyzes intent
- Generates ID: feature-20260217-003
- Estimates files: 5
- Adds structured item to Features section
- Sets timestamp

Result: Item added to backlog.md
```

**Example interactions:**
```
User: "Add a bug for the null pointer in login"
→ AI adds to Bugs section with location

User: "Add refactoring task to extract database connection logic"
→ AI adds to Refactoring section with location

User: "Track this: implement caching layer"
→ AI adds to Features section
```

**Automatically Generated:**
- **Epic decomposition** - During `@m.plan` for Track 4

### Viewing Backlog

Simply read backlog.md or use AI:
```
User: "What's in the backlog?"
→ AI reads and summarizes backlog.md

User: "Show me all bugs"
→ AI filters and displays Bugs section

User: "What features are pending?"
→ AI lists items from Features section
```
m backlog list [--type TYPE] [--status STATUS]

# Show item details
m backlog show --id <id>
```

## Export to External Tools

Export backlog items to external tools for team collaboration.

### Export Workflow

```bash
# 1. AI parses backlog.md and converts to JSON
m backlog export > backlog.json

# 2. Transform JSON for external tool
# 3. Import to external tool
```

### Export Command

```bash
# Export entire backlog to JSON
m backlog export > backlog.json

# Export specific types
m backlog export --type feature > features.json
m backlog export --type decomposed > epic.json
```

### Transform for External Tools

**For Beads:**
```bash
# Export and transform
m backlog export > backlog.json
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
# Export and transform
m backlog export > backlog.json
jq -r '.items[] | "gh issue create --title \"\(.title)\" --body \"\(.description // "")\" --label \(.type)"' backlog.json | sh
```

**For Jira (via CSV):**
```bash
# Export and transform to CSV
m backlog export > backlog.json
jq -r '["Summary","Description","Issue Type"], 
  (.items[] | [.title, .description, .type]) | @csv' backlog.json > backlog.csv

# Import CSV to Jira via web interface
```

## Complete Workflows

### Workflow 1: Epic Decomposition with Export

```bash
# 1. Plan epic intent
$ @m.plan "Implement swimlane status flow with m status command"

# Mission-toolkit detects Track 4, decomposes:
📊 Track 4 (Epic) - 15 files
🔄 Decomposing into sub-intents...
✅ Added 5 sub-intents to backlog.md

# 2. View decomposed items
User: "Show me the decomposed intents"
AI: Reads and displays Decomposed Intents section from backlog.md

# 3. Export to Beads
$ m backlog export --type decomposed > epic.json
$ jq -r '.items[] | "bd create \"\(.title)\" --label type:decomposed"' epic.json | sh

# 4. View in Beads with dependency graph
$ bd graph

# 5. Start working on first sub-intent
User: "Let's work on status detection service"
AI: Reads item from backlog.md, starts mission
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

# 3. Optional: Add refactor item to backlog for tracking
User: "Add refactor task to extract email validation pattern"
AI: Adds to Refactoring section in backlog.md
```

### Workflow 3: AI Agent Managing Backlog

```bash
# AI agent reads backlog
User: "What's in the backlog?"
AI: Reads and summarizes backlog.md

# AI agent analyzes and prioritizes
User: "What's the smallest task?"
AI: Analyzes Estimated Files, suggests smallest item

# AI agent helps start mission
User: "Let's work on status detection"
AI: Reads item from backlog.md, starts mission with @m.plan

# After completion, mark as done
User: "Mark status detection as complete"
AI: Updates checkbox in backlog.md: - [x] decomposed-20260217-001
```

## File Structure

```
.mission/
├── mission.md          # Current mission state
├── backlog.md          # Future work queue (AI-maintained)
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

### backlog.md (Future Work)

```markdown
---
# Backlog Template Instructions (for AI)
# When user says "add X to backlog" or similar:
# 1. Analyze the intent
# 2. Determine type (feature/bug/refactor/decomposed)
# 3. Estimate files if possible
# 4. Add item below in the appropriate section
# 5. Generate unique ID (type-timestamp format)
# 6. Set created timestamp
---

# Mission Backlog

## Features
- [ ] feature-20260217-001: Add user authentication
  - **Description**: Implement JWT-based authentication system
  - **Estimated Files**: 5
  - **Created**: 2026-02-17T10:00:00Z

## Bugs
- [ ] bug-20260217-001: Fix null pointer in login handler
  - **Description**: Null pointer exception when password is empty
  - **Location**: auth/handler.go:45
  - **Created**: 2026-02-17T10:10:00Z

## Refactoring
- [ ] refactor-20260217-001: Extract validation logic
  - **Description**: Extract email validation into shared utility
  - **Location**: users/handler.go:45
  - **Estimated Files**: 3
  - **Created**: 2026-02-17T10:15:00Z

## Decomposed Intents (from Track 4 Epics)
- [ ] decomposed-20260217-001: Create status detection service
  - **Description**: Determines mission status from filesystem state
  - **Rationale**: Core status detection logic is foundation
  - **Estimated Files**: 3
  - **Dependencies**: None
  - **Created**: 2026-02-17T10:20:00Z
```
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

### Backlog Management (AI-Assisted)

**Natural Language (Recommended)**
```
User: "Add user authentication to backlog"
User: "Add bug for null pointer in login handler"
User: "Show me all features in the backlog"
User: "Mark feature-001 as complete"
```

AI reads/writes backlog.md following embedded template.

### Export Commands

**m backlog export**
```bash
# Parse backlog.md and export to JSON
m backlog export > backlog.json

# Export specific types
m backlog export --type feature > features.json
m backlog export --type decomposed > epic.json
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
```
User: "Mark feature-001 as complete"
```
- AI updates checkbox in backlog.md: - [x] feature-001
- Separate from mission step tracking

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
## Implementation Notes

### Backlog Template

The backlog.md file contains embedded instructions for AI:

```markdown
---
# Backlog Template Instructions (for AI)
# When user says "add X to backlog" or similar:
# 1. Analyze the intent
# 2. Determine type (feature/bug/refactor/decomposed)
# 3. Estimate files if possible
# 4. Add item below in the appropriate section
# 5. Generate unique ID (type-timestamp format)
# 6. Set created timestamp
---
```

### Export Implementation

```go
func (s *BacklogService) Export(typeFilter string) ([]byte, error) {
    // 1. Read backlog.md
    content, err := os.ReadFile(".mission/backlog.md")
    if err != nil {
        return nil, err
    }
    
    // 2. Parse markdown to extract items
    items := parseBacklogMarkdown(content)
    
    // 3. Filter by type if specified
    if typeFilter != "" {
        items = filterByType(items, typeFilter)
    }
    
    // 4. Convert to JSON
    output := BacklogJSON{
        Version: "1.0",
        Items:   items,
    }
    
    return json.MarshalIndent(output, "", "  ")
}
```

### Markdown Parsing

```go
func parseBacklogMarkdown(content []byte) []BacklogItem {
    var items []BacklogItem
    lines := strings.Split(string(content), "\n")
    
    var currentItem *BacklogItem
    var currentSection string
    
    for _, line := range lines {
        // Detect section headers
        if strings.HasPrefix(line, "## ") {
            currentSection = strings.TrimPrefix(line, "## ")
            continue
        }
        
        // Parse checkbox items
        if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
            // Extract ID and title
            // Parse metadata fields (Description, Estimated Files, etc.)
            currentItem = &BacklogItem{
                Type:   sectionToType(currentSection),
                Status: checkboxToStatus(line),
            }
            items = append(items, *currentItem)
        }
    }
    
    return items
}
```

## Benefits of This Approach

### Simplicity
- ✅ AI-assisted (no manual CLI commands)
- ✅ Human-readable markdown
- ✅ Embedded template guides AI
- ✅ Git-friendly diffs

### Clear Value
- ✅ Natural language interaction
- ✅ Consistent structure enforced by template
- ✅ Export to any external tool
- ✅ No complex command syntax

### Flexibility
- ✅ Works with any AI assistant
- ✅ Users can manually edit if needed
- ✅ Easy to add custom transformations
- ✅ Export to JSON for external tools

### Standalone
- ✅ No external dependencies
- ✅ Works completely offline
- ✅ Simple markdown format
- ✅ AI maintains structure

## Design Principles

1. **AI-First**: Natural language interaction, not CLI commands
2. **Template-Guided**: Embedded instructions ensure consistency
3. **Export-Only**: Clear value without import complexity
4. **Markdown Native**: Human-readable, git-friendly
5. **Minimal Complexity**: Simplest possible implementation

## Open Questions

1. **Template Evolution**: How to update template instructions over time?
2. **Validation**: Should AI validate structure before writing?
3. **Conflict Resolution**: How to handle concurrent edits?

## Related Documentation

- [Core Concepts](concepts.md) - WET→DRY workflow and Rule of Three
- [Workflows](workflows.md) - Mission lifecycle and epic decomposition
- [CLI Reference](cli-reference.md) - Complete command documentation
