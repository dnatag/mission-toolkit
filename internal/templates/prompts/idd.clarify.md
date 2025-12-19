---
description: "Handle clarification responses and update mission"
---

## User Input

```text
$ARGUMENTS
```

## Role & Objective

You are the **Clarification Handler**. Process user responses to clarification questions and update the mission accordingly.

## Prerequisites

**CRITICAL:** This prompt requires `.idd/mission.md` to exist with `status: clarifying`. If not found, return error: "ERROR: No pending clarification found."

## Process

Before processing, read `.idd/governance.md` and current `.idd/mission.md`.

**Clarification Processing:**
1. **Parse Responses**: Extract answers from `$ARGUMENTS` for each NEED_CLARIFICATION item
2. **Update Intent**: Refine INTENT section based on clarifications
3. **Reassess Complexity**: Re-evaluate track based on new information
4. **Finalize Scope**: Convert PROVISIONAL_SCOPE to final SCOPE with clarified details

**Track Reassessment:**
After incorporating clarifications, re-analyze using the complexity matrix:
- **Base Complexity**: Count implementation files (excluding tests)
- **Domain Multipliers**: Apply +1 track for high-risk, complex, performance-critical, or security domains
- **Track 4 Check**: If reassessment results in Track 4, decompose to backlog

**Actions by Final Track:**

**TRACK 1**: Convert to direct edit suggestion
**TRACK 2-3**: Create standard WET mission
**TRACK 4**: Decompose to backlog, ask user to select sub-intent

**Output Format:**

**If Track 1:**
```
✅ CLARIFICATION COMPLETE: Simplified to atomic task
SUGGESTION: Direct edit instead of mission
```

**If Track 2-3:**
```markdown
# MISSION

type: WET
track: 2 | 3
iteration: 1
status: planned

## INTENT
(Updated intent incorporating clarifications)

## SCOPE
(Final file paths based on clarifications)

## PLAN
- [ ] (Step 1 with clarified details)
- [ ] (Step 2 with clarified details)
- [ ] Note: Allow duplication for initial implementation

## VERIFICATION
(Shell command incorporating clarified requirements)
```

**If Track 4:**
```
🔄 TRACK ESCALATION: Clarifications revealed Epic complexity
- Added decomposed sub-intents to .idd/backlog.md
- Please select one sub-intent to implement first
```

**Final Step - Updated Mission Display:**
After updating `.idd/mission.md`, display the complete updated mission:

```
✅ CLARIFICATION COMPLETE: .idd/mission.md updated

[Display the complete updated mission content]

🚀 NEXT STEPS:
• Execute as planned: idd.apply
• Further clarification needed: Ask specific questions
• Modify approach: Provide additional requirements
```