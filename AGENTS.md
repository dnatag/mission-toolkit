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
- **Tests**: `tests/` - AI-native workflow validation framework

### Key Files
- `main.go` - CLI entry point
- `internal/templates/templates.go` - Template embedding
- `internal/version/version.go` - Version management
- `tests/cases/*.md` - AI-native test scenarios
- `tests/framework/` - AI reasoning validation framework

## Tech Stack

**Language**: Go 1.21+  
**Framework**: Cobra CLI + Bubble Tea TUI  
**Templates**: Go templates (embedded)  
**Testing**: AI-native reasoning validation + Go test

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
# AI-native testing with reasoning validation
# 1. Use scenarios from tests/cases/
# 2. Validate against tests/framework/assertion-patterns.md
# 3. Follow tests/framework/execution-guide.md
# 4. All tests pass: 32/32 assertions (100% success rate)
```