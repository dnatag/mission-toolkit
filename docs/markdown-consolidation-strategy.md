# Markdown Consolidation Strategy

**Status**: Planned  
**Date**: 2026-01-31  
**Mission ID**: 20260131163853-0526

## Executive Summary

Consolidate markdown document parsing across mission.md, diagnosis.md, and backlog.md by creating a unified `pkg/md` abstraction layer. This eliminates duplication in frontmatter parsing and section handling while adding metadata tracking to backlog.md.

## Current State Analysis

### Implementation Overview

| File | Frontmatter Parsing | Body Parsing | YAML Library | Issues |
|------|-------------------|--------------|--------------|--------|
| mission.md | Custom regex | String parsing | yaml.v2 | Complex, maintenance burden |
| diagnosis.md | String split | String parsing | yaml.v3 | Fragile, inconsistent YAML version |
| backlog.md | None | Line-by-line | N/A | No metadata tracking |
| validator.go | N/A | Goldmark AST | N/A | Works well for validation |

### Identified Problems

1. **Duplication**: Two different frontmatter parsing approaches (regex vs string split)
2. **Inconsistency**: Mixed YAML library versions (v2 vs v3)
3. **Fragility**: String split doesn't handle whitespace variations or edge cases
4. **Missing Metadata**: backlog.md has no tracking for last_updated or modification history
5. **Maintenance Burden**: ~50 lines of custom parsing code to maintain and test

### Common Pattern

All markdown files use simple structure:
- YAML frontmatter (mission.md, diagnosis.md) or none (backlog.md)
- Sections with headers (`##`)
- Content as either lists (`- [ ]`, `- [x]`) or string paragraphs

## Recommended Solution

### Architecture: pkg/md Abstraction Layer

Create a new package `pkg/md` that provides:

1. **Document Type**: Unified markdown document representation
2. **Frontmatter Parsing**: Using `adrg/frontmatter` library
3. **Section Helpers**: Reusable utilities for headers, lists, and content
4. **YAML Standardization**: Migrate to `gopkg.in/yaml.v3`

### Why This Approach?

**Option A: adrg/frontmatter + Simple String Parsing (CHOSEN)**
- Matches our simple use case (headers + lists/strings)
- Frontmatter library handles edge cases
- Keep simple string parsing for body (already works)
- Goldmark stays in validator.go for validation only

**Option B: adrg/frontmatter + Goldmark (Rejected)**
- Overkill for our simple markdown structure
- AST traversal adds complexity without benefit
- Current string parsing works fine

### Key Benefits

1. **Single Source of Truth**: One frontmatter parsing approach for all files
2. **Robustness**: Dedicated library handles edge cases (whitespace, special chars, delimiters)
3. **Consistency**: Same YAML library (v3) across codebase
4. **Code Reduction**: Eliminate ~50 lines of custom parsing code
5. **Backlog Metadata**: Track modification history and timestamps
6. **Maintainability**: Clear separation between document structure and business logic
7. **Testability**: Centralized testing for markdown operations

## Implementation Plan

### Phase 1: Foundation

**Backlog Items:**
1. Add adrg/frontmatter dependency for YAML frontmatter parsing
2. Create pkg/md package for markdown document abstraction layer
3. Implement pkg/md.Document with frontmatter parsing using adrg/frontmatter
4. Implement pkg/md section parsing helpers (headers, lists, string content)

**Deliverables:**
- `pkg/md/document.go` - Document type with frontmatter support
- `pkg/md/section.go` - Section parsing utilities
- `pkg/md/document_test.go` - Comprehensive test coverage
- `pkg/md/section_test.go` - Section parsing tests

**Example API:**
```go
// Parse document with frontmatter
doc, err := md.Parse(content)

// Access frontmatter
var meta MissionMetadata
doc.Frontmatter(&meta)

// Parse sections
sections := doc.Sections()
section := doc.GetSection("INTENT")

// Update section
doc.UpdateSection("PLAN", content)
doc.UpdateList("SCOPE", items, appendMode)
```

### Phase 2: Mission Package Migration

**Backlog Items:**
5. Migrate pkg/mission/reader.go to use pkg/md abstraction
6. Migrate pkg/mission/writer.go to use pkg/md abstraction

**Changes:**
- Replace custom regex parsing with `md.Parse()`
- Remove ~30 lines of frontmatter extraction code
- Migrate from yaml.v2 to yaml.v3
- Keep legacy format support for backward compatibility
- Update all tests in reader_test.go and writer_test.go

**Impact:**
- Simplified reader.go implementation
- Consistent error messages from library
- Better edge case handling

### Phase 3: Diagnosis Package Migration

**Backlog Items:**
7. Migrate pkg/diagnosis/diagnosis.go to use pkg/md abstraction

**Changes:**
- Replace string split with `md.Parse()`
- Already uses yaml.v3 (no YAML change needed)
- Update all tests in diagnosis_test.go

**Impact:**
- More robust frontmatter parsing
- ~10 lines of code simplified

### Phase 4: Backlog Package Enhancement

**Backlog Items:**
8. Add frontmatter support to backlog.md with last_updated and last_action fields
9. Migrate pkg/backlog/manager.go to use pkg/md abstraction

**New Frontmatter Structure:**
```yaml
---
last_updated: 2026-01-31T17:05:00Z
last_action: "Added refactoring pattern: frontmatter-yaml-parsing"
---
```

**Changes:**
- Add frontmatter parsing to manager.go
- Update frontmatter on all modifications (Add, Complete, Cleanup)
- Maintain backward compatibility (handle backlog.md without frontmatter)
- Update all tests

**Benefits:**
- Track when backlog was last modified
- Record what action was performed
- Audit trail for backlog changes

### Phase 5: Cleanup

**Backlog Items:**
10. Migrate all packages from gopkg.in/yaml.v2 to gopkg.in/yaml.v3
11. Remove gopkg.in/yaml.v2 dependency from go.mod after migration complete

**Tasks:**
- Search codebase for yaml.v2 imports
- Update to yaml.v3
- Run full test suite
- Remove yaml.v2 from go.mod
- Update documentation

## Risk Mitigation

### Backward Compatibility
- Keep legacy format support in mission/reader.go
- Handle backlog.md without frontmatter gracefully
- Comprehensive test coverage for edge cases

### Testing Strategy
- Unit tests for pkg/md package
- Integration tests for each migrated package
- Edge case testing: whitespace variations, --- in YAML, malformed files
- Regression testing with existing mission.md and diagnosis.md files

### Rollout Plan
- Implement in separate missions with focused scope
- Use git checkpoints for easy rollback
- Gradual migration: foundation → mission → diagnosis → backlog → cleanup

### Rollback Plan
- Git checkpoints enable easy rollback at each phase
- Each phase is independently reversible
- No breaking changes to file formats (backward compatible)

## Success Metrics

- **Code Reduction**: ~50 lines of custom parsing code eliminated
- **Consistency**: Single frontmatter parsing approach across 3 files
- **YAML Standardization**: 100% migration to yaml.v3
- **Test Coverage**: Maintain or improve current coverage
- **Backlog Metadata**: 100% of backlog modifications tracked
- **Zero Regressions**: All existing tests pass after migration

## Related Patterns

This strategy resolves the following duplication patterns tracked in backlog:
- `[PATTERN:frontmatter-yaml-parsing][COUNT:2]` - Consolidated into pkg/md
- `[PATTERN:section-content-helpers][COUNT:2]` - Consolidated into pkg/md
- `[PATTERN:list-section-update][COUNT:2]` - Consolidated into pkg/md

## References

- **adrg/frontmatter**: https://github.com/adrg/frontmatter
- **Goldmark**: https://github.com/yuin/goldmark (kept for validation only)
- **YAML v3**: https://github.com/go-yaml/yaml/tree/v3
- **Mission**: .mission/mission.md (20260131163853-0526)

## Next Steps

1. Review and approve this strategy document
2. Create first mission: "Implement pkg/md abstraction layer foundation"
3. Execute phases sequentially with separate missions
4. Update this document with lessons learned after each phase
