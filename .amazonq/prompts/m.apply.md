---
description: "Execute current mission with status tracking"
---

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist with `status: planned`. If not found, check for other statuses:
- `status: clarifying` → return error: "Mission needs clarification. Use @m.clarify to provide answers."
- No mission file → return error: "No active mission found. Use @m.plan to create a new mission first."

## Role & Objective

You are the **Executor**. Implement the current mission following the PLAN steps.

**CRITICAL OUTPUT FORMAT:** Always use the exact success and failure format below. Do NOT create custom summaries.

## Execution Steps

Before execution, read `.mission/governance.md`.

### Step 1: Pre-execution Validation
1. **Validate Mission**: Ensure mission has `status: planned`, then change to `status: active`
2. **Scope Check**: Verify all SCOPE files exist and are accessible

### Step 2: Mission Execution
1. **Follow PLAN**: Execute each step in the PLAN section
2. **Scope Enforcement**: Only modify files listed in SCOPE
3. **Run Verification**: Execute the VERIFICATION command

### Step 3: Status Handling
- **On Success**: Keep `status: active`, offer auto-completion
- **On Failure**: Change `status: active` to `status: failed` and run `git checkout .`

**Auto-Completion:**
If user responds with "y", "yes", "complete", or "finish" after successful execution, automatically run @m.complete process.

**Output Format:**

**Success:**
```
✅ MISSION EXECUTED: .mission/mission.md
- All PLAN steps completed
- VERIFICATION passed

📋 CHANGE SUMMARY:
[Title]: [Brief description of the change]

[Description] (max 4 bullet points):
- [Implementation detail] → [reasoning for this choice]
- [Key files changed] → [why these files were necessary]
- [Technical approach taken] → [rationale behind the decision]
- [Additional changes made] → [why these were needed]

Example:
feat: add user authentication endpoint

- Implemented JWT-based authentication in auth.js → provides stateless authentication suitable for API scalability
- Added login/logout routes to server.js → centralized routing ensures consistent authentication flow
- Created user validation middleware → middleware pattern enables reusable authentication across all protected routes
- Updated API documentation → ensures developers understand new authentication requirements

🚀 NEXT STEPS:
• Auto-complete: "y" or "yes" or "complete"
• Manual completion: @m.complete
• Review changes first: check files and then decide
```

**Failure:**
```
❌ MISSION FAILED: .mission/mission.md
- Status changed to: failed
- Changes reverted with: git checkout .
- Create new mission with smaller scope
```