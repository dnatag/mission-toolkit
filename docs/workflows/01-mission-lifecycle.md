---
name: Mission Lifecycle
description: Complete workflow from planning to completion with handshakes and reviews
---

# Mission Lifecycle - Detailed Workflow

## Overview

The Mission Toolkit workflow ensures you maintain ownership while leveraging AI capabilities through a series of handshakes and reviews.

## Complete Flow Diagram

```
                    ┌─────────────┐
                    │   m.plan    │
                    │             │
                    │ Analyzes    │
                    │ intent      │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Ambiguous? │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │ YES                     │ NO
              ▼                         ▼
     ┌─────────────┐           ┌─────────────┐
     │  m.clarify  │           │   m.plan    │
     │  (optional) │           │             │
     │             │           │ Creates     │
     │ Asks        │           │ mission.md  │
     │ questions   │           └──────┬──────┘
     └──────┬──────┘                  │
            │                         │
            │ Re-runs m.plan          │
            └────────┬────────────────┘
                     │
                     ▼
            ┌─────────────┐
            │ Review      │
            │ mission.md  │
            │             │
            │ • INTENT    │
            │ • SCOPE     │
            │ • PLAN      │
            │ • VERIFY    │
            └──────┬──────┘
                   │
                   │ Approve?
                   ▼
            ┌─────────────┐
            │   m.apply   │
            │             │
            │ Executes    │
            │ + Polish    │
            │ + Generates │
            │   commit    │
            └──────┬──────┘
                   │
            ┌──────▼──────┐
            │ Review code │
            │ Adjustments?│
            └──────┬──────┘
                   │
          ┌────────┼────────┐
          │ YES             │ NO
          ▼                 ▼
   ┌─────────────┐   ┌─────────────┐
   │ User        │   │ m.complete  │
   │ requests    │   │             │
   │ changes     │   │ Archives    │
   │             │   │ Creates     │
   │ AI fixes +  │   │ git commit  │
   │ regenerates │   └─────────────┘
   │ commit msg  │
   └──────┬──────┘
          │
          └──────────────────────────┐
                                     │
                                     ▼
                            ┌─────────────┐
                            │ m.complete  │
                            │             │
                            │ Archives    │
                            │ Creates     │
                            │ git commit  │
                            └─────────────┘
```

## Step-by-Step Breakdown

### 1. m.plan - Intent Analysis

**Purpose**: Convert user intent into structured mission

**Process**:
- Analyze intent for clarity and complexity
- If ambiguous → route to m.clarify
- If clear → create mission.md with INTENT, SCOPE, PLAN, VERIFICATION
- **Test Strategy**: AI evaluates if test files should be included based on value (new logic, bugs) vs. noise (trivial changes).

**Output**: `.mission/mission.md` ready for review

### 2. m.clarify - Clarification (Optional)

**Purpose**: Refine ambiguous intents

**Process**:
- Ask targeted questions to clarify requirements
- Gather missing details
- Re-run m.plan with clarified intent

**Output**: Updated mission.md with refined details

### 3. Review mission.md - Human Authorization

**Purpose**: You authorize the architecture before execution

**What to Review**:
- **INTENT**: Does it match what you want?
- **SCOPE**: Are the right files included?
- **PLAN**: Does the approach make sense?
- **VERIFICATION**: Will this prove it works?

**Decision**: Approve or request changes

### 4. m.apply - Execution

**Purpose**: Implement the authorized plan

**Process** (automatic):
1. Execute implementation steps
2. Run polish pass (code quality improvements)
   - **Quality-Driven Testing**: Review newly created tests. Ensure they are high-value (meaningful happy paths, critical edge cases) and robust. Remove low-value tests (e.g., trivial getters/setters).
3. Generate conventional commit message
4. Update mission.md with commit message

**Output**: Working code + commit message ready for review

### 5. Review Code - Human Verification

**Purpose**: Verify the implementation meets requirements

**What to Review**:
- Does code match the PLAN?
- Does verification pass?
- Any bugs or improvements needed?

**Decision**: Accept or request adjustments

### 6. Adjustments - Iterative Refinement (Optional)

**Purpose**: Fix issues discovered during review

**Process**:
- User describes needed changes
- AI makes fixes
- AI regenerates commit message to reflect changes
- Loop back to m.complete

**Output**: Refined code + updated commit message

### 7. m.complete - Archival & Commit

**Purpose**: Capture learnings and create git commit

**Process**:
- Archive mission files to `.mission/completed/`
- Update metrics and backlog
- Create git commit using stored commit message
- Clean up active mission files

**Output**: Git commit + archived mission data

## Key Principles

### Human-in-the-Loop
- **Before execution**: Review and approve the plan
- **After execution**: Review and verify the code
- **Optional refinement**: Request changes if needed

### Atomic Scope
- Each mission is small enough to comprehend
- Changes stay within human understanding limits
- You maintain ownership, not just contribution

### Continuous Learning
- Every mission archived with full context
- Metrics tracked for process improvement
- Patterns detected for future optimization

### Quality-Driven Testing
- **Scope-Matched**: Test coverage must match mission complexity.
- **High-Value**: Focus on meaningful happy paths and critical edge cases.
- **Essential**: Avoid testing for the sake of testing (e.g., trivial getters/setters).
- **Quality**: Tests must be robust, readable, and maintainable.

## Common Paths

### Happy Path (No Issues)
```
m.plan → Review → m.apply → Review → m.complete
```

### With Clarification
```
m.plan → m.clarify → m.plan → Review → m.apply → Review → m.complete
```

### With Adjustments
```
m.plan → Review → m.apply → Review → Adjustments → m.complete
```

### Multiple Adjustment Cycles
```
m.plan → Review → m.apply → Review → Adjustments → Review → Adjustments → m.complete
```

## Related Documentation

- [Commit Message Flow](02-commit-messages.md) - How commit messages are generated and stored
- [Epic Decomposition](03-epic-decomposition.md) - Handling Track 4 (Epic) missions
- [Governance Rules](../../.mission/governance.md) - Core principles and constraints
