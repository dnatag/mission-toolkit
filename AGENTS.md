# Mission Toolkit Development Guide

## Quick Reference

**Build**: `go build -o bin/m main.go`  
**Run**: `./bin/m --help`  
**Test**: `go test ./...`  
**Init Project**: `./bin/m init --ai-type q`

## Project Architecture

### Core Components
- **CLI Commands**: `cmd/` - Root, init, status, version commands
- **Templates**: `internal/templates/` - Embedded mission and prompt templates
- **TUI**: `internal/tui/` - Terminal user interface
- **Mission Logic**: `internal/mission/` - Mission file reader
- **Tests**: `tests/` - AI workflow validation (manual)

### Key Files
- `main.go` - CLI entry point
- `internal/templates/templates.go` - Template embedding
- `internal/version/version.go` - Version management
- `tests/fixtures/*.md` - AI test scenarios

## Tech Stack

**Language**: Go 1.21+  
**Framework**: Cobra CLI + Bubble Tea TUI  
**Templates**: Go templates (embedded)  
**Testing**: Go test + manual AI validation

## Essential Commands

### Development
```bash
go run main.go --help              # Run CLI
go build -o bin/m main.go          # Build binary
go test ./...                      # Run tests
```

### CLI Usage
```bash
./bin/m init --ai-type q           # Initialize project
./bin/m status                     # Show mission status
./bin/m version                    # Show version
```

### Testing AI Workflows
```bash
# Manual testing with AI agents
# 1. Use scenarios from tests/fixtures/
# 2. Validate against tests/assertions/
# 3. Follow tests/validation/workflow-validator.md
```