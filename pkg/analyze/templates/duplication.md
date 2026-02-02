# DUPLICATION ANALYSIS TEMPLATE

## Current Intent

{{.CurrentIntent}}

## Purpose
Identify existing code patterns similar to the requested feature. This informs the WET→DRY workflow decision using Rule-of-Three tracking.

## Analysis Steps

### 1. Semantic Search
Scan the codebase for similar implementations:

- **File Names**: Existing files with similar names (e.g., `auth.go` vs `authentication.go`)
- **Function Names**: Functions doing similar tasks (e.g., `ValidateUser` vs `CheckUser`)
- **Business Logic**: Logic solving the same problem (e.g., multiple email validation implementations)

### 2. Pattern Recognition
Identify reusable patterns:

- **Boilerplate**: Standard patterns (CRUD, handlers) that exist elsewhere
- **Utilities**: Existing utility functions that could be reused
- **Configuration**: Existing config structures that could be extended

### 3. Cross-Package Pattern Search
**CRITICAL**: Architectural patterns often span multiple packages. Don't limit your search to the same business domain.

Identify architectural patterns that appear across different packages or modules:

- **Function Signatures**: Similar function signatures with same parameters and behavior across packages
  - Go: `UpdateList(section string, items []string, appendMode bool)` in mission vs diagnosis packages
  - Python: `update_list(section: str, items: List[str], append_mode: bool)` in mission vs diagnosis modules
  - JavaScript/TypeScript: `updateList(section: string, items: string[], appendMode: boolean)` in mission vs diagnosis modules
  
- **Helper Function Groups**: Related helper functions that indicate pattern reuse
  - Go: `extractExistingItems()`, `skipSectionContent()`, `addFormattedItems()`
  - Python: `extract_existing_items()`, `skip_section_content()`, `add_formatted_items()`
  - JavaScript: `extractExistingItems()`, `skipSectionContent()`, `addFormattedItems()`
  
- **Behavioral Patterns**: Same algorithm implemented in different contexts
  - Markdown section parsing and updating
  - List management with append/replace modes
  - Configuration file manipulation

**When to Record Pattern Occurrence:**
- If implementing similar logic in a different package/module → Check if pattern already exists in backlog
- If pattern exists → Run `m backlog add "[description]" --pattern-id [existing-id] --type refactor` to increment count
- If pattern is new → Record as occurrence #1: `m backlog add "[description]" --pattern-id [new-id] --type refactor`
- **Critical**: Pattern occurrences count across packages/modules, not just within same file

### 4. Pattern ID Generation
When duplication is found, generate a stable pattern ID:
- Use lowercase kebab-case (e.g., `email-validation`, `db-connection`, `list-section-update`)
- Be specific enough to identify the pattern uniquely
- Be general enough to match future occurrences

### 5. Rule of Three Decision Tree

**For each pattern detected, follow this decision tree:**

```
Is this pattern already in backlog?
├─ YES → Run: m backlog add "[description]" --pattern-id [existing-id] --type refactor
│         (CLI auto-increments count)
│         └─ Is count now >= 3?
│            ├─ YES → Mission type should be DRY (refactor the pattern)
│            └─ NO → Mission type remains WET (allow duplication)
│
└─ NO → Run: m backlog add "[description]" --pattern-id [new-id] --type refactor
          (Creates pattern with count=1)
          └─ Mission type is WET (first occurrence)
```

**Command Examples:**
```bash
# First occurrence (new pattern)
m backlog add "List section update with append mode" --pattern-id list-section-update --type refactor

# Second occurrence (increment existing pattern)
m backlog add "List section update with append mode" --pattern-id list-section-update --type refactor

# Third occurrence (triggers DRY mission)
m backlog add "List section update with append mode" --pattern-id list-section-update --type refactor
# At this point, mission type should be DRY to refactor the pattern
```

**Check existing patterns before creating new ones:**
```bash
m backlog list --include refactor
```

## Output Format

Produce a JSON object with duplication analysis.

**Note**: The CLI tracks pattern counts in backlog. When count reaches 3, mission type becomes DRY.

```json
{
  "duplication_detected": true | false,
  "patterns": [
    {
      "id": "pattern-id",
      "description": "Human-readable description",
      "locations": ["file1.go:45", "file2.go:78"]
    }
  ]
}
```

## Examples

### Example 1: No Duplication
```json
{
  "duplication_detected": false,
  "patterns": []
}
```

### Example 2: Same-Package Duplication
```json
{
  "duplication_detected": true,
  "patterns": [
    {
      "id": "email-validation",
      "description": "Email validation regex duplicated across handlers",
      "locations": ["users/handler.go:45", "products/handler.go:78"]
    },
    {
      "id": "db-connection",
      "description": "Database connection logic duplicated",
      "locations": ["auth/db.go:12", "users/db.go:15", "orders/db.go:18"]
    }
  ]
}
```

### Example 3: Cross-Package Pattern (Architectural Duplication)
```json
{
  "duplication_detected": true,
  "patterns": [
    {
      "id": "list-section-update",
      "description": "Markdown list section update with append mode support - identical algorithm across packages",
      "locations": [
        "pkg/mission/writer.go:UpdateList()",
        "pkg/diagnosis/diagnosis.go:UpdateList()"
      ]
    }
  ]
}
```

**Note for Example 3**: This represents the second occurrence of the `list-section-update` pattern. The AI should run:
```bash
m backlog add "List section update with append mode" --pattern-id list-section-update --type refactor
```
This increments the pattern count to 2. When a third occurrence is detected, mission type becomes DRY.
