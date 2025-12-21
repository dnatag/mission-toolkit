# Mission Toolkit

> **"Slow down the process to speed up the understanding"**  
> *Intent defines the scope and approach — let purpose drive process*

## 🧠 The Philosophy

Intent-Driven Atomic Development is a minimalist workflow designed to bridge the gap between "Vibe Coding" (Chaos) and "Spec-Driven Development" (Bureaucracy).

We believe that AI coding fails in two extremes:

**🌀 The Vibe Trap:** You let the AI drive. It moves fast, hallucinates, and paints you into a corner. You feel frustrated.

**📝 The Spec Trap:** You write exhaustive documentation before coding. It works, but it alienates you from the codebase. You feel like a contributor, not an owner.

**✨ Intent-Driven Atomic Development is the Golden Ratio.** It forces a "🤝 Handshake" before every coding task. You don't write the code, but you authorize the architecture and verify the implementation.

## ⚛️ Why Atomic?

This creates the psychological sweet spot where you maintain ownership while leveraging AI capabilities. The secret is **deliberate pacing** — "slow down the process to speed up the understanding."

We deliberately work only with atomic-sized intents to maintain small scope. This actually slows down the process — you can't tackle massive features in one go. But this constraint gives you better understanding and genuine ownership. When your brain can fully comprehend each small mission, you maintain control instead of becoming a passenger to AI speed.

## ⚙️ How It Works

The Mission Toolkit implements this philosophy through a systematic approach that converts user intents into structured, executable missions with built-in complexity management and continuous improvement.

### 🔄 The Three-Phase Handshake

**1. 📝 Intent Analysis (Human)** - You define what you want, the Mission Toolkit categorizes complexity using a 4-track system

**2. 🛠️ Mission Planning (AI + Human)** - AI proposes scope and plan, you authorize the architecture

**3. 🚀 Execution & Learning (AI + Human)** - AI implements, you verify, system learns patterns

### Core Capabilities

- **🔄 WET→DRY Evolution**: Write Everything Twice (WET) first, then Don't Repeat Yourself (DRY) — allows duplication initially, then extracts abstractions when patterns emerge
- **🎯 Mission-Based Execution**: Breaks work into atomic, verifiable missions  
- **📈 Continuous Improvement**: Tracks metrics and patterns for process optimization

## 🤝 The Slash Commands

*Note: Amazon Q CLI and Kiro CLI don't support inline arguments. Run m.plan or m.clarify without user input, and the AI will prompt for input.*

### 📝 `/m.plan` - The Planning Handshake
Converts your intent into a structured mission. You define what, AI proposes how, you authorize.

```bash
# Example usage
/m.plan "Add user authentication to the API"
```

**Features:**
- 🎯 4-track complexity analysis (Atomic, Standard, Robust, Epic)
- 📁 Automatic scope estimation and file identification
- 🔒 Security validation and input sanitization
- 📋 Backlog management for complex intents

### 🔍 `/m.clarify` - The Clarification Handshake (Optional)
Refines vague or complex intents before planning. Helps break down ambiguous requirements into actionable missions.

```bash
# Example usage
/m.clarify "Make the app better"
```

**Features:**
- 🎯 Intent disambiguation and scope refinement
- 📋 Requirement extraction from vague descriptions
- 🔄 Interactive clarification process
- 📝 Structured output ready for m.apply

### 🚀 `/m.apply` - The Execution Handshake  
Implements your authorized plan. AI handles execution while you maintain oversight.

```bash
# Example usage
/m.apply
```

**Features:**
- 🎯 Focused scope enforcement (only modify listed files)
- 🔄 WET vs DRY mission differentiation
- ✅ Mandatory verification execution
- 🔍 Pattern detection for future refactoring

### 📈 `/m.complete` - The Learning Handshake
Captures what was accomplished and learned. Builds organizational memory for future missions.

```bash
# Example usage
/m.complete
```

**Features:**
- 📁 Mission archival with timestamps
- 📊 Metrics collection and analysis
- 📋 Backlog updates and pattern tracking
- 📆 Historical data preservation

## Project Structure

```
.mission/
├── governance.md          # Core principles and workflow rules
├── backlog.md            # Future work and refactoring opportunities
├── metrics.md            # Aggregate performance statistics
├── mission.md            # Current active mission (auto-generated)
└── completed/            # Archived missions and detailed metrics
    ├── YYYY-MM-DD-HH-MM-mission.md
    └── YYYY-MM-DD-HH-MM-metrics.md

prompts/
├── m.clarify.md        # Clarification prompt for vague intents
├── m.plan.md           # Planning prompt and complexity matrix
├── m.apply.md          # Execution prompt and safety checks
└── m.complete.md       # Completion prompt and observability
```

## Complexity Matrix

| Track | Scope | Files | Keywords | Action |
|-------|-------|-------|----------|--------|
| **TRACK 1** (Atomic) | Single line/function | 0 new files | "Fix typo", "Rename var" | Skip Mission, direct edit |
| **TRACK 2** (Standard) | Single feature | 1-5 files | "Add endpoint", "Create component" | Standard WET mission |
| **TRACK 3** (Robust) | Cross-cutting concerns | Security/Auth/Performance | "Add authentication", "Refactor for security" | Robust WET mission |
| **TRACK 4** (Epic) | Multiple systems | 10+ files | "Build payment system", "Rewrite architecture" | Decompose to backlog |

*Note: Test files don't count toward complexity*

## 🔄 WET-then-DRY Workflow

### 💧 WET Phase (Write Everything Twice)
- **Purpose**: Understand the problem domain through implementation
- **Approach**: Allow duplication to explore solutions
- **Outcome**: Working features with identified patterns

### 🌵 DRY Phase (Don't Repeat Yourself)
- **Trigger**: User explicitly requests refactoring after patterns emerge
- **Approach**: Extract abstractions based on observed duplication
- **Outcome**: Clean, maintainable code with appropriate abstractions

## Mission Lifecycle

```
User Intent → [m.clarify] → m.plan → .mission/mission.md → m.apply → Verification → m.complete → Archive
                              ↓                                                                    ↓
                          .mission/backlog.md ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←← .mission/metrics.md
```

## Key Principles

### 1. Focused Scope
- Only modify files explicitly listed in mission SCOPE
- Prevents scope creep and unintended changes
- Enables precise impact assessment

### 2. Atomic Execution
- All changes broken into verifiable steps
- Each mission has clear success criteria
- Mandatory verification before completion

### 3. Complexity Management
- Automatic complexity detection and routing
- Epic decomposition into manageable sub-missions
- Progressive disclosure of complexity

### 4. Continuous Improvement
- Detailed metrics collection and analysis
- Pattern detection for process optimization
- Historical data preservation for trend analysis

## 🚀 Getting Started

1. **📁 Initialize Project**
   ```bash
   # Initialize Mission Toolkit project with AI-specific templates
   m init --ai-type q
   
   # Supported AI types: q, claude, gemini, cursor, codex, kiro
   # Creates .mission/ directory with governance files and prompt templates
   ```

2. **📊 Check Project Status**
   ```bash
   # Display interactive TUI showing current and completed missions
   m status
   
   # Use ↑/↓ to navigate missions, Enter to view details, / to search
   # Shows mission progress and provides clear next steps
   ```

3. **📝 Plan Your First Mission**
   ```bash
   /m.plan "Your development intent here"
   ```

4. **⚙️ Execute the Mission**
   ```bash
   /m.apply
   ```

5. **🏁 Complete and Track**
   ```bash
   /m.complete
   ```

## Observability Features

### Metrics Tracking
- Mission duration and complexity correlation
- Track distribution and success rates
- WET→DRY evolution effectiveness
- Verification success/failure patterns

### Pattern Detection
- Automatic duplication identification
- Abstraction opportunity recognition
- Common failure pattern analysis
- Process bottleneck identification

### Historical Analysis
- Timestamped mission archives
- Performance trend analysis
- Process evolution tracking
- Evidence-based improvements

## ✨ Benefits

- **🧠 Reduced Cognitive Load**: Atomic missions eliminate decision paralysis — your brain stays in sync with AI speed
- **👑 Maintained Ownership**: You authorize architecture and verify implementation — never feel like just a contributor
- **✅ Quality Assurance**: Mandatory verification and scope constraints prevent the "Vibe Trap" chaos
- **🛠️ Technical Debt Management**: Systematic WET→DRY evolution avoids premature abstraction
- **📈 Scalability**: Handles projects from toy features to enterprise systems through complexity decomposition

## Versioning

Mission Toolkit uses **dual versioning** to enable independent evolution of CLI functionality and embedded templates.

### Version Structure

- **CLI Version**: Tracks command-line tool functionality (`internal/version/version.go`)
- **Template Version**: Tracks embedded workflow templates (`internal/templates/templates.go`)

```bash
$ m version
CLI: v1.2.3
Templates: v2.1.0
```

### Independent Evolution Benefits

- **Focused Updates**: CLI bugs don't require template version bumps
- **Template Innovation**: New AI features can evolve templates independently
- **Clear Compatibility**: Users understand which components changed
- **Selective Adoption**: Teams can update CLI or templates as needed

## Release Process

### Creating a Release

1. **Tag the release**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Automated build**: GitHub Actions will automatically:
   - Run tests and validation
   - Build cross-platform binaries (Linux, macOS, Windows)
   - Create zip archives for each platform
   - Generate checksums and changelog
   - Publish release with artifacts

3. **Download binaries**: Users can download platform-specific zips from the GitHub releases page.

### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

## License

This project is licensed under the terms specified in the LICENSE file.