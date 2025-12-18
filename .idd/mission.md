# MISSION

type: WET
track: 2
iteration: 1
status: completed

## INTENT
Add --ai flag to init command with validation for supported AI types (q, claude, gemini, cursor, codex, cline, kiro)

## SCOPE
cmd/init.go
internal/templates/templates.go

## PLAN
- [x] Add --ai string flag to init command
- [x] Create list of supported AI types for validation
- [x] Add validation logic to check AI type is supported
- [x] Update init command to use WriteTemplates with specified AI type
- [x] Add error handling for invalid AI types
- [x] Update command description and help text
- [x] Note: Allow duplication for initial implementation

## VERIFICATION
go run . init --help | grep -q "\-\-ai" && go run . init --ai invalid 2>&1 | grep -i "error\|invalid" && echo "Validation tests passed"
