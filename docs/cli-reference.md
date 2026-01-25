# CLI Reference

Complete reference for all `m` CLI commands.

## Core Commands (Human)

```bash
m version                          # Show version
m init --ai <q|claude|kiro|opencode>  # Initialize project
m dashboard                        # Interactive TUI dashboard
```

## AI Commands

The following commands are invoked by AI assistants via prompt templates (`m.plan`, `m.apply`, `m.complete`, `m.debug`).

### Mission Lifecycle

```bash
m mission create --intent "description"
m mission check --context <plan|apply|complete|debug>
m mission update --status <active|executed|completed|failed>
m mission finalize
m mission archive
```

## Diagnosis Lifecycle

```bash
m diagnosis create --symptom "description"
m diagnosis check --context debug
m diagnosis update --section <hypotheses|investigation|root-cause|affected-files|recommended-fix|reproduction>
m diagnosis update --status <confirmed|inconclusive> --confidence <high|medium|low>
m diagnosis finalize
```

## Analysis Commands

```bash
m analyze intent "description"     # Analyze user intent
m analyze scope                    # Analyze mission scope
m analyze complexity               # Analyze complexity track
m analyze clarify                  # Check for clarification needs
m analyze duplication              # Check for code duplication
m analyze decompose                # Decompose epic intents
m analyze test                     # Analyze test requirements
```

## Backlog Management

```bash
m backlog list                     # List backlog items
m backlog add "item" --type <decomposed|refactor>
m backlog complete --item "exact text"
m backlog resolve --item "pattern"
m backlog cleanup                  # Remove completed items
```

## Checkpoint Management

```bash
m checkpoint create                # Create checkpoint
m checkpoint restore <name>        # Restore checkpoint
m checkpoint commit -m "message"   # Create commit
```

## Logging and Validation

```bash
m log --step "name" "message"      # Log execution step
m check "intent"                   # Validate intent
```
