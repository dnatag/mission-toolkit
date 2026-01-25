# Core Concepts

## The Philosophy

Atomic Intent-Driven Development bridges the gap between "Vibe Coding" (Chaos) and "Spec-Driven Development" (Bureaucracy).

AI coding fails in two extremes:

**🌀 The Vibe Trap:** You let the AI drive. It moves fast, generates massive code changes beyond human comprehension, and paints you into a corner.

**📝 The Spec Trap:** You write exhaustive documentation before coding. AI generates large implementations that work, but the sheer volume alienates you from the codebase.

**✨ Atomic Intent-Driven Development is the Golden Ratio.** It forces a "🤝 Handshake" before every coding task and keeps changes within human comprehension limits.

## Why Atomic?

We deliberately work only with atomic-sized intents to maintain small scope. This actually slows down the process — you can't tackle massive features in one go. But this constraint gives you better understanding and genuine ownership.

> "Slow down the process to speed up the understanding"

## Complexity Matrix

| Track | Scope | Files | Keywords | Action |
|-------|-------|-------|----------|--------|
| **TRACK 1** (Atomic) | Single line/function | 0 new files | "Fix typo", "Rename var" | Skip Mission, direct edit |
| **TRACK 2** (Standard) | Single feature | 1-5 files | "Add endpoint", "Create component" | Standard WET mission |
| **TRACK 3** (Robust) | Cross-cutting concerns | Security/Auth/Performance | "Add authentication" | Robust WET mission |
| **TRACK 4** (Epic) | Multiple systems | 10+ files | "Build payment system" | Decompose to backlog |

*Note: Test files don't count toward complexity*

## WET-then-DRY Workflow

### 💧 WET Phase (Write Everything Twice)
- **Purpose**: Understand the problem domain through implementation
- **Approach**: Allow duplication to explore solutions
- **Outcome**: Working features with identified patterns

### 🌵 DRY Phase (Don't Repeat Yourself)
- **Trigger**: User explicitly requests refactoring after patterns emerge
- **Approach**: Extract abstractions based on observed duplication
- **Outcome**: Clean, maintainable code with appropriate abstractions

## Key Principles

### 1. Focused Scope
Only modify files explicitly listed in mission SCOPE. Prevents scope creep and enables precise impact assessment.

### 2. Atomic Execution
All changes broken into verifiable steps. Each mission has clear success criteria with mandatory verification.

### 3. Complexity Management
Automatic complexity detection and routing. Epic decomposition into manageable sub-missions.

### 4. Template-Driven Consistency
Embedded template system ensures consistent outputs. LLM-agnostic design works with any AI assistant.

### 5. Continuous Improvement
Pattern detection for process optimization. Execution logging for debugging.
