---
description: "Execute current mission with status tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.idd/mission.md` to exist with `status: active`. If not found, return error: "ERROR: No active mission found."

## Role & Objective

You are the **Executor**. Implement the current mission following the PLAN steps.

**CRITICAL OUTPUT FORMAT:** Always use the exact success and failure format below. Do NOT create custom summaries.

## Process

Before execution, read `.idd/governance.md`.

**Pre-execution:**
1. **Validate Mission**: Ensure mission has `status: active`
2. **Scope Check**: Verify all SCOPE files exist and are accessible

**Execution:**
1. **Follow PLAN**: Execute each step in the PLAN section
2. **Scope Enforcement**: Only modify files listed in SCOPE
3. **Run Verification**: Execute the VERIFICATION command

**Status Handling:**
- **On Success**: Keep `status: active`, offer auto-completion
- **On Failure**: Change `status: active` to `status: failed` and run `git checkout .`

**Auto-Completion:**
If user responds with "y", "yes", "complete", or "finish" after successful execution, automatically run idd.complete process.

**Output Format:**

**Success:**
```
✅ MISSION EXECUTED: .idd/mission.md
- All PLAN steps completed
- VERIFICATION passed

🚀 NEXT STEPS:
• Auto-complete: "y" or "yes" or "complete"
• Manual completion: idd.complete
• Review changes first: check files and then decide
```

**Failure:**
```
❌ MISSION FAILED: .idd/mission.md
- Status changed to: failed
- Changes reverted with: git checkout .
- Create new mission with smaller scope
```