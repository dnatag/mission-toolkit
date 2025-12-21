---
description: "Complete current mission and update project tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist. If `.mission/mission.md` is not found, return error: "No active mission found. Use @m.plan to create a new mission first."

## Role & Objective

You are the **Completor**. Finalize the current mission and update project tracking for continuous improvement.

## Execution Steps

Before generating output, read `.mission/governance.md`.

### Step 1: Mission Validation
1. **Status Check**: Ensure mission has `status: active` (not failed)
2. **Completion Check**: Verify all PLAN items in `.mission/mission.md` are completed
3. **Verification Status**: Confirm VERIFICATION command was run successfully
4. **Scope Validation**: Ensure all SCOPE files were properly modified

### Step 2: Mission Completion and Archival
1. **Update Mission**: Change `status: active` to `status: completed` and add `completed_at: YYYY-MM-DDTHH:MM:SS.sssZ`
2. **Archive Mission**: Move `.mission/mission.md` to `.mission/completed/YYYY-MM-DD-HH-MM-mission.md`
3. **Create Metrics**: Generate `.mission/completed/YYYY-MM-DD-HH-MM-metrics.md` with detailed data
4. **Clean Up**: Remove `.mission/mission.md` after successful archiving

### Step 3: Project Tracking Updates
1. **Update Summary**: Update `.mission/metrics.md` by refreshing all aggregate statistics and adding new completion to RECENT COMPLETIONS
2. **Update Backlog**: Search `.mission/backlog.md` for matching intent, mark as completed with timestamp

**Output Format:**

```markdown
# MISSION COMPLETED

**Timestamp**: YYYY-MM-DD HH:MM:SS
**Mission Type**: WET | DRY
**Track**: 1 | 2 | 3 | 4
**Files Modified**: [count]
**Duration**: [estimated time]

## CHANGE SUMMARY
[Copy the complete change summary from @m.apply execution]

Title: [Brief description]

Description (max 4 bullet points):
- [Implementation detail] → [reasoning for this choice]
- [Key files changed] → [why these files were necessary]
- [Technical approach taken] → [rationale behind the decision]
- [Additional changes made] → [why these were needed]

## OUTCOMES
- [ ] All PLAN items completed
- [ ] VERIFICATION passed
- [ ] Files properly modified
- [ ] Backlog updated (matching items marked as ✅ COMPLETED YYYY-MM-DD)

## PATTERNS DETECTED
(List any duplication patterns for future DRY missions)

## NEXT STEPS
(Suggested follow-up missions or backlog items to prioritize)

## METRICS CREATED
- **Detailed Metrics**: `.mission/completed/YYYY-MM-DD-HH-MM-metrics.md` (includes change summary)
- **Summary Updated**: `.mission/metrics.md` aggregate statistics refreshed
- **Historical Preservation**: All mission data preserved with timestamps
```

🚀 WHAT'S NEXT:
• Start new mission: @m.plan "your next intent"
• Review metrics: Check .mission/metrics.md
• [Suggested follow-up missions or backlog items to prioritize]