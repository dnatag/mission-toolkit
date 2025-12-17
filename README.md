# Intent-Driven Development (IDD) Toolkit

> **"Slow down the process to speed up the understanding"**  
> *Intent defines the scope and approach — let purpose drive process*

## 🧠 The Philosophy

IDD is a minimalist workflow designed to bridge the gap between "Vibe Coding" (Chaos) and "Spec-Driven Development" (Bureaucracy).

We believe that AI coding fails in two extremes:

**🌀 The Vibe Trap (System 1):** You let the AI drive. It moves fast, hallucinates, and paints you into a corner. You feel frustrated.

**📝 The Spec Trap (System 2 Overload):** You write exhaustive documentation before coding. It works, but it alienates you from the codebase. You feel like a contributor, not an owner.

*System 1 and System 2 terminologies are from Daniel Kahneman's "Thinking, Fast and Slow"*

**✨ IDD is the Golden Ratio.** It forces a "🤝 Handshake" before every coding task. You don't write the code, but you authorize the architecture and verify the implementation.

### 🤝 The Handshake Moment

This creates the psychological sweet spot where you maintain ownership while leveraging AI capabilities. The secret is **Tactical System 2** — "slow down the process to speed up the understanding."

IDD speeds up understanding by working only on atomic-sized intents. Large epics get decomposed into manageable pieces where your brain can keep pace with AI speed. Instead of trying to comprehend massive changes all at once, you understand each small mission completely before moving forward. This piecemeal approach lets human cognition stay in sync with AI execution.

## ⚙️ How It Works

IDD implements this philosophy through a systematic approach that converts user intents into structured, executable missions with built-in complexity management and continuous improvement.

### 🔄 The Three-Phase Handshake

**1. 📝 Intent Analysis (Human)** - You define what you want, IDD categorizes complexity using a 4-track system

**2. 🛠️ Mission Planning (AI + Human)** - AI proposes scope and plan, you authorize the architecture

**3. 🚀 Execution & Learning (AI + Human)** - AI implements, you verify, system learns patterns

### Core Capabilities

- **🔄 WET→DRY Evolution**: Write Everything Twice (WET) first, then Don't Repeat Yourself (DRY) — allows duplication initially, then extracts abstractions when patterns emerge
- **🎯 Mission-Based Execution**: Breaks work into atomic, verifiable missions  
- **📈 Continuous Improvement**: Tracks metrics and patterns for process optimization

## 🤝 The Slash Commands

*Note: Use `/` for most AI assistants, `@` for Amazon Q CLI and Kiro CLI*

### 📝 `/idd.plan` - The Planning Handshake
Converts your intent into a structured mission. You define what, AI proposes how, you authorize.

```bash
# Example usage
/idd.plan "Add user authentication to the API"
```

**Features:**
- 🎯 4-track complexity analysis (Atomic, Standard, Robust, Epic)
- 📁 Automatic scope estimation and file identification
- 🔒 Security validation and input sanitization
- 📋 Backlog management for complex intents

### 🚀 `/idd.apply` - The Execution Handshake  
Implements your authorized plan. AI handles execution while you maintain oversight.

```bash
# Example usage
/idd.apply
```

**Features:**
- 🎯 Focused scope enforcement (only modify listed files)
- 🔄 WET vs DRY mission differentiation
- ✅ Mandatory verification execution
- 🔍 Pattern detection for future refactoring

### 📈 `/idd.complete` - The Learning Handshake
Captures what was accomplished and learned. Builds organizational memory for future missions.

```bash
# Example usage
/idd.complete
```

**Features:**
- 📁 Mission archival with timestamps
- 📊 Metrics collection and analysis
- 📋 Backlog updates and pattern tracking
- 📆 Historical data preservation

## Project Structure

```
.idd/
├── governance.md          # Core principles and workflow rules
├── backlog.md            # Future work and refactoring opportunities
├── metrics.md            # Aggregate performance statistics
├── mission.md            # Current active mission (auto-generated)
└── completed/            # Archived missions and detailed metrics
    ├── YYYY-MM-DD-HH-MM-mission.md
    └── YYYY-MM-DD-HH-MM-metrics.md

prompts/
├── idd.plan.md           # Planning prompt and complexity matrix
├── idd.apply.md          # Execution prompt and safety checks
└── idd.complete.md       # Completion prompt and observability
```

## Complexity Matrix

| Track | Scope | Files | Keywords | Action |
|-------|-------|-------|----------|--------|
| **TRACK 1** (Atomic) | Single line/function | 0 new files | "Fix typo", "Rename var" | Skip IDD, direct edit |
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
User Intent → idd.plan → .idd/mission.md → idd.apply → Verification → idd.complete → Archive
                ↓                                                                    ↓
            .idd/backlog.md ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←← .idd/metrics.md
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
   mkdir .idd
   # Copy governance.md, backlog.md, metrics.md templates
   ```

2. **📝 Plan Your First Mission**
   ```bash
   /idd.plan "Your development intent here"
   ```

3. **⚙️ Execute the Mission**
   ```bash
   /idd.apply
   ```

4. **🏁 Complete and Track**
   ```bash
   /idd.complete
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

## License

This project is licensed under the terms specified in the LICENSE file.