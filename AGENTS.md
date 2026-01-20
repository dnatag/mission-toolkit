# Mission Toolkit Development Guide

## Quick Reference

**Build**: `go build -o bin/m main.go`  
**Run**: `./bin/m --help`  
**Test**: `go test ./...`  
**Init Project**: `./bin/m init --ai q`

## Project Architecture

### Core Components
- **CLI Commands**: `cmd/` - Root, init, dashboard, version, analyze, mission, backlog, checkpoint, check, log
- **Templates**: `internal/templates/` - Embedded mission and prompt templates
- **TUI**: `internal/tui/` - Terminal user interface with dashboard
- **Mission Logic**: `internal/mission/` - Mission lifecycle management
- **Analysis**: `internal/analyze/` - Intent, scope, complexity analysis
- **Backlog**: `internal/backlog/` - Backlog management
- **Git Integration**: `internal/git/` - Git operations and checkpoints
- **Logger**: `internal/logger/` - Execution logging
- **Utils**: `internal/utils/` - File utilities and validation

### Key Files
- `main.go` - CLI entry point
- `internal/mission/mission.go` - Mission data structures and core logic
- `internal/templates/templates.go` - Template embedding and deployment
- `internal/version/version.go` - Version management
- `internal/templates/prompts/` - AI prompt templates (m.plan, m.apply, m.complete)
- `internal/templates/libraries/` - Display templates and references

## Tech Stack

**Language**: Go 1.25+  
**Framework**: Cobra CLI + Bubble Tea TUI  
**Templates**: Go templates (embedded)  
**Git**: go-git library for operations  
**Testing**: Go test + table-driven tests  
**Development**: Mission-driven development with AI assistance

## Essential Commands

### Development
```bash
go run main.go --help              # Run CLI
go build -o bin/m main.go          # Build binary
go test ./...                      # Run tests
go mod tidy                        # Clean dependencies
```

### CLI Usage
```bash
# Project initialization
m init --ai q                      # Initialize with Amazon Q
m init --ai claude                 # Initialize with Claude
m init --ai kiro                   # Initialize with Kiro
m init --ai opencode               # Initialize with OpenCode

# Mission management
m dashboard                        # Interactive TUI dashboard
m mission check --context plan    # Check mission state for planning
m mission create --intent "desc"   # Create new mission
m mission update --status active   # Update mission status
m mission archive                  # Archive completed mission

# Analysis tools
m analyze intent "user input"      # Analyze user intent
m analyze scope                    # Analyze mission scope
m analyze complexity               # Analyze complexity track
m analyze clarify                  # Check for clarification needs
m analyze duplication              # Check for code duplication
m analyze decompose                # Decompose epic intents
m analyze test                     # Analyze test requirements

# Backlog management
m backlog list                     # List backlog items
m backlog add "item" --type refactor # Add backlog item
m backlog complete --item "text"   # Mark item complete
m backlog resolve --item "text"    # Mark refactor resolved
m backlog cleanup                  # Remove completed items

# Git checkpoints
m checkpoint create                # Create checkpoint
m checkpoint restore <name>       # Restore checkpoint
m checkpoint commit -m "msg"       # Create commit

# Logging and validation
m log --step "name" "message"      # Log execution step
m check "intent"                   # Validate intent
m version                          # Show version
```

## AI Integration

### Supported AI Assistants
- **Amazon Q**: Uses `@m.plan`, `@m.apply`, `@m.complete` commands
- **Claude**: Uses `/m.plan`, `/m.apply`, `/m.complete` commands  
- **Kiro**: Uses `@m.plan`, `@m.apply`, `@m.complete` commands
- **OpenCode**: Uses `/m.plan`, `/m.apply`, `/m.complete` commands

### Prompt Templates
- `m.plan.md` - Planning phase with complexity analysis
- `m.apply.md` - Execution phase with two-pass implementation
- `m.complete.md` - Completion phase with commit generation

### Template Features
- CLI-exclusive state management
- JSON output parsing and conditional logic
- Template-driven analysis with embedded guidance
- Two-pass implementation with automatic rollback
- Rich commit message generation

## Development Workflow

### Mission-Driven Development
1. **Plan**: Use `/m.plan` to analyze intent and create structured mission
2. **Execute**: Use `/m.apply` for two-pass implementation with verification
3. **Complete**: Use `/m.complete` to generate commit and archive mission

### Key Principles
- **Atomic Scope**: Only modify files listed in mission SCOPE
- **WET→DRY Evolution**: Allow duplication first, refactor when patterns emerge
- **Mandatory Verification**: All missions must pass verification before completion
- **Template Consistency**: Use embedded templates for predictable outputs
- **Safety First**: Automatic rollback on polish failures

## Testing Strategy

### Unit Tests
- Table-driven tests for core logic
- Mock interfaces for external dependencies
- Filesystem abstraction with afero

### Integration Tests
- End-to-end CLI command testing
- Template rendering validation
- Git operations testing

### AI Workflow Testing
- Mission lifecycle validation
- Template output verification
- Error handling scenarios