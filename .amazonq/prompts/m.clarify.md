---
description: "Handle clarification responses and update mission"
---

## User Input

```text
$ARGUMENTS
```

## Interactive Prompt

**CRITICAL:** Always check if `$ARGUMENTS` is empty or contains only whitespace first.

If `$ARGUMENTS` is empty, blank, or contains only whitespace:
- Ask: "What clarifications can you provide for the current mission?"
- Wait for user response
- Use the response as `$ARGUMENTS` and continue

## Role & Objective

You are the **Clarification Handler**. Process user responses to clarification questions and update the mission accordingly.

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist with `status: clarifying`. If not found, use template `.mission/libraries/displays/error-no-mission.md` with message: "No mission awaiting clarification. Use @m.plan to create a new mission first."

## Execution Steps

Before processing, read `.mission/governance.md` and current `.mission/mission.md`.

### Step 1: Clarification Processing
1. **Parse Responses**: Extract answers from `$ARGUMENTS` for each NEED_CLARIFICATION item
2. **Update Intent**: Refine INTENT section based on clarifications
3. **Reassess Complexity**: Re-evaluate track based on new information
4. **Finalize Scope**: Convert PROVISIONAL_SCOPE to final SCOPE with clarified details

### Step 2: Track Reassessment
After incorporating clarifications, re-analyze using the complexity matrix:
- **Base Complexity**: Count implementation files (excluding tests)
- **Domain Multipliers**: Apply +1 track for high-risk, complex, performance-critical, or security domains
- **Track 4 Check**: If reassessment results in Track 4, decompose to backlog

### Step 3: Mission Update
**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Actions by Final Track:**
- **TRACK 1**: Convert to direct edit suggestion
- **TRACK 2-3**: Use template `.mission/libraries/missions/wet.md` with clarified variables
- **TRACK 4**: Use template `.mission/libraries/displays/clarify-escalation.md`

**If Track 1:**
```
✅ CLARIFICATION COMPLETE: Simplified to atomic task
SUGGESTION: Direct edit instead of mission
```

**If Track 2-3**: Use template `.mission/libraries/missions/wet.md` with variables:
- {{TRACK}} = 2 or 3 (reassessed)
- {{REFINED_INTENT}} = Updated intent incorporating clarifications
- {{FILE_LIST}} = Final file paths based on clarifications
- {{PLAN_STEPS}} = Steps with clarified details
- {{VERIFICATION_COMMAND}} = Shell command incorporating clarified requirements

**If Track 4**: Use template `.mission/libraries/displays/clarify-escalation.md` with variables:
- {{BACKLOG_ITEMS}} = Decomposed sub-intents

**Final Step**: Use template `.mission/libraries/displays/clarify-success.md` with variables:
- {{MISSION_CONTENT}} = Complete updated mission content