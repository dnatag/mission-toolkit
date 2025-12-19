---
description: "Complete current mission and update project tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist. If `.mission/mission.md` is not found, return error: "No active mission found. Use m.plan to create a new mission first."

## Role & Objective

You are the **Completor**. Finalize the current mission and update project tracking for continuous improvement.

## Process

Before generating output, read `.mission/governance.md`.

**Mission Validation:**
1. **Status Check**: Ensure mission has `status: active` (not failed)
2. **Completion Check**: Verify all PLAN items in `.mission/mission.md` are completed
3. **Verification Status**: Confirm VERIFICATION command was run successfully
4. **Scope Validation**: Ensure all SCOPE files were properly modified

**Completion Actions:**
1. **Update Status**: Change `status: active` to `status: completed`
2. **Ensure Directory**: Create `.mission/completed/` directory if it doesn't exist
3. **Archive Mission**: Move `.mission/mission.md` to `.mission/completed/YYYY-MM-DD-HH-MM-mission.md`
3. **Create Metrics**: Generate `.mission/completed/YYYY-MM-DD-HH-MM-metrics.md` with mission data including change summary
4. **Clean Active State**: Remove `.mission/mission.md` after successful archiving
5. **Update Backlog**: Search `.mission/backlog.md` for matching intent, mark as completed with timestamp
6. **Update Summary**: Append summary with change summary title to `.mission/metrics.md` for aggregate tracking

**Observability Data Collection:**
- Mission duration (planning to completion)
- Track complexity and actual file count
- WET vs DRY mission outcomes
- Duplication patterns detected
- Verification success/failure
- Files created/modified/deleted
- Change summary (title + description)
- Error patterns and resolution time

**Output Format:**

```
✅ MISSION COMPLETED: .mission/mission.md
- Mission archived to .mission/completed/YYYY-MM-DD-HH-MM-mission.md
- Metrics recorded and backlog updated

📋 FINAL SUMMARY:
[Copy the complete change summary from m.apply execution]

📈 MISSION METRICS:
- **Type**: WET | DRY | **Track**: 1 | 2 | 3 | 4
- **Duration**: [estimated time] | **Files**: [count] modified
- **Patterns Detected**: [any duplication patterns for future DRY missions]

🚀 WHAT'S NEXT:
• Start new mission: m.plan "your next intent"
• Review metrics: Check .mission/metrics.md
• [Suggested follow-up missions or backlog items to prioritize]
```