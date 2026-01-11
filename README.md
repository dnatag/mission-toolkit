# Mission Toolkit

> **"Slow down the process to speed up the understanding"**  
> *Intent defines the scope and approach — let purpose drive process*

## 🧠 The Philosophy

Intent-Driven Atomic Development is a minimalist workflow designed to bridge the gap between "Vibe Coding" (Chaos) and "Spec-Driven Development" (Bureaucracy).

We believe that AI coding fails in two extremes:

**🌀 The Vibe Trap:** You let the AI drive. It moves fast, generates massive code changes beyond human comprehension, and paints you into a corner. You feel frustrated and alienated from your own codebase.

**📝 The Spec Trap:** You write exhaustive documentation before coding. AI generates large implementations that work, but the sheer volume alienates you from the codebase. You feel like a contributor, not an owner.

**✨ Intent-Driven Atomic Development is the Golden Ratio.** It forces a "🤝 Handshake" before every coding task and keeps changes within human comprehension limits. You don't write the code, but you authorize the architecture and verify the implementation.

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

*Note: Amazon Q CLI and Kiro CLI differences:*
- *Use '@' commands instead of '/' (e.g., @m.plan, @m.clarify, @m.apply, @m.complete)*
- *Inline arguments are ignored - the AI will prompt for input*

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
/m.clarify
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
├── execution.log         # Current mission execution log
├── completed/            # Archived missions and detailed metrics
│   ├── MISSION-ID-mission.md
│   ├── MISSION-ID-metrics.md
│   └── MISSION-ID-execution.log
├── paused/               # Temporarily paused missions
│   └── TIMESTAMP-mission.md
└── libraries/            # Template system (embedded)
    ├── analysis/         # Analysis guidance templates
    ├── displays/         # User output templates
    ├── logs/             # Execution logging templates
    ├── metrics/          # Metrics templates
    ├── missions/         # Mission file templates
    ├── scripts/          # Operation templates
    └── variables/        # Variable calculation rules

# AI-specific prompt directories:
.amazonq/prompts/         # Amazon Q prompts
.claude/commands/         # Claude commands
.kiro/prompts/           # Kiro prompts
.opencode/command/       # OpenCode commands
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
m.plan → [m.clarify] → 🤝 Review mission.md → m.apply → 🤝 Review code → [Adjustments] → m.complete
                        (Handshake #1)                  (Handshake #2)
```

**How it works:**
1. **m.plan** creates mission.md with INTENT, SCOPE, PLAN, VERIFICATION
2. **m.clarify** (optional) refines ambiguous requirements
3. **🤝 Review & approve** the mission before execution (authorize the architecture)
4. **m.apply** executes, polishes, and generates commit message
5. **🤝 Review code** and optionally request adjustments (verify the implementation)
6. **m.complete** archives mission and creates git commit

[See detailed workflow diagram →](docs/workflows/01-mission-lifecycle.md)

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

### 4. Template-Driven Consistency
- Embedded template system ensures consistent outputs
- Standardized variable system across all operations
- LLM-agnostic design works with any AI assistant

### 5. Continuous Improvement
- Detailed metrics collection and analysis
- Pattern detection for process optimization
- Historical data preservation for trend analysis
- Execution logging for debugging and learning

## 🚀 Getting Started

### Installation

#### Option 1: Download Pre-built Binaries (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/dnatag/mission-toolkit/releases):

```bash
# macOS (Intel)
curl -L https://github.com/dnatag/mission-toolkit/releases/download/v2.0.0/mission-toolkit_Darwin_x86_64.zip -o m.zip
unzip m.zip && chmod +x m && sudo mv m /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/dnatag/mission-toolkit/releases/download/v2.0.0/mission-toolkit_Darwin_arm64.zip -o m.zip
unzip m.zip && chmod +x m && sudo mv m /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/dnatag/mission-toolkit/releases/download/v2.0.0/mission-toolkit_Linux_x86_64.zip -o m.zip
unzip m.zip && chmod +x m && sudo mv m /usr/local/bin/

# Linux (arm64)
curl -L https://github.com/dnatag/mission-toolkit/releases/download/v2.0.0/mission-toolkit_Linux_arm64.zip -o m.zip
unzip m.zip && chmod +x m && sudo mv m /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/dnatag/mission-toolkit/releases/download/v2.0.0/mission-toolkit_Windows_x86_64.zip -OutFile m.zip
Expand-Archive m.zip -DestinationPath .
# Add to PATH manually
```

#### Option 2: Homebrew (macOS/Linux)

```bash
brew tap dnatag/mission-toolkit
brew install mission-toolkit
```

#### Option 3: Build from Source

```bash
# Prerequisites: Go 1.21+
git clone https://github.com/dnatag/mission-toolkit.git
cd mission-toolkit
go build -o m main.go
sudo mv m /usr/local/bin/
```

### Verify Installation

```bash
m version
```

### Quick Start

1. **📁 Initialize Project**
   ```bash
   # Initialize Mission Toolkit project with AI-specific templates
   m init --ai q
   
   # Supported AI types: q, claude, kiro, opencode
   # Creates .mission/ directory with governance files and prompt templates
   ```

2. **📊 Check Project Status**
   ```bash
   # Display comprehensive mission dashboard with execution logs
   m dashboard
   
   # Use ↑/↓ to navigate missions, Enter to view details, Tab to switch panes
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

## Template System Features

### Embedded Templates
- **Analysis Templates**: Clarification and complexity assessment guidance
- **Display Templates**: Consistent user output for all command outcomes
- **Mission Templates**: WET, DRY, and clarification mission structures
- **Script Templates**: Standardized operations for status updates and archival
- **Metrics Templates**: Individual mission and aggregate project metrics
- **Logging Templates**: Execution step tracking and debugging

### Variable Standardization
- Consistent naming across all templates ({{TRACK}}, {{MISSION_TYPE}}, etc.)
- Type-safe variable handling (numeric vs string)
- Default values for missing variables
- Cross-template variable dependencies

### AI-Agnostic Design
- Automatic slash prefix adaptation (@m.plan vs /m.plan)
- AI-specific directory structure (Amazon Q, Claude, Kiro, OpenCode)
- Template deployment to appropriate AI prompt directories
- Unified versioning for CLI and templates

## Observability Features

### Execution Logging
- Step-by-step mission execution tracking
- Timestamped log entries with success/failure status
- Archived logs with completed missions
- Debugging support for failed missions

### Metrics Tracking
- Mission duration and complexity correlation
- Track distribution and success rates
- WET→DRY evolution effectiveness
- Verification success/failure patterns
- Template system usage analytics

### Pattern Detection
- Automatic duplication identification
- Abstraction opportunity recognition
- Common failure pattern analysis
- Process bottleneck identification

### Historical Analysis
- Timestamped mission archives with full context
- Performance trend analysis
- Process evolution tracking
- Evidence-based improvements

## ✨ Benefits

- **🧠 Reduced Cognitive Load**: Atomic missions eliminate decision paralysis — your brain stays in sync with AI speed
- **👑 Maintained Ownership**: You authorize architecture and verify implementation — never feel like just a contributor
- **✅ Quality Assurance**: Mandatory verification and scope constraints prevent the "Vibe Trap" chaos
- **🛠️ Technical Debt Management**: Systematic WET→DRY evolution avoids premature abstraction
- **📈 Scalability**: Handles projects from toy features to enterprise systems through complexity decomposition
- **🔧 Template Consistency**: Embedded template system ensures reliable, predictable outputs across all AI assistants
- **📊 Full Observability**: Execution logging and metrics provide complete visibility into mission progress and outcomes
- **🔄 AI-Agnostic**: Works seamlessly with Amazon Q, Claude, Kiro, OpenCode, and other AI assistants

## Versioning

```bash
# Check current version
m version

# Update version (for maintainers)
./scripts/sync-version.sh v1.0.0
```

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