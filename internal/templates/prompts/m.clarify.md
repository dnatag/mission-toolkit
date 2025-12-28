---
description: "Handle clarification workflow and update mission"
---

## Role & Objective

You are the **Clarification Handler**. Load clarification questions from the current mission and guide the user through providing answers.

## Prerequisites

**CRITICAL:** This prompt requires `.mission/mission.md` to exist with `status: clarifying`. If not found, use template `.mission/libraries/displays/error-no-mission.md`.

## Execution Steps

Before processing, read `.mission/governance.md` and current `.mission/mission.md`.

**MUST LOG:** Use file read tool to check if `.mission/execution.log` exists. If file doesn't exist, use file read tool to load template `libraries/scripts/init-execution-log.md`, then use file write tool to create the log file.

### Step 1: Load and Display Questions
1. **Load Mission**: Read `.mission/mission.md`
2. **Extract Questions**: Parse NEED_CLARIFICATION section
3. **Display to User**: Show numbered list of questions
4. **Request Answers**: Ask user to provide responses

**Display Format:**
```
🤔 CLARIFICATION NEEDED

Please provide answers to these questions:

1. [Question 1]
2. [Question 2]
3. [Question 3]

Provide your answers - you can reference questions by number or respond in any clear format.
```

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.clarify 1: Load and Display Questions | [questions loaded, displayed to user]"

### Step 2: Process User Responses
1. **Parse Answers**: Extract numbered responses from user input
2. **Update Intent**: Refine INTENT section based on clarifications
3. **Reassess Complexity**: Re-evaluate track based on new information
4. **Finalize Scope**: Convert PROVISIONAL_SCOPE to final SCOPE with clarified details

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.clarify 2: Process User Responses | [responses parsed, intent updated]"

### Step 3: Track Reassessment
Use file read tool to load template `libraries/analysis/complexity.md` to re-analyze using the complexity matrix:
- **Base Complexity**: Count implementation files (excluding tests)
- **Domain Multipliers**: Apply +1 track for high-risk, complex, performance-critical, or security domains
- **Track 4 Check**: If reassessment results in Track 4, decompose to backlog

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.clarify 3: Track Reassessment | [final track, complexity reasoning]"

### Step 4: Update Mission
**CRITICAL**: Use templates from `.mission/libraries/` for consistent output.

**Actions by Final Track:**
- **TRACK 1**: Convert to direct edit suggestion
- **TRACK 2-3**: Use template `.mission/libraries/missions/wet.md` with clarified variables
- **TRACK 4**: Use template `.mission/libraries/displays/clarify-escalation.md`

**MUST LOG:** Use file write tool (append mode) to add to `.mission/execution.log` using template `libraries/logs/execution.md`:
- {{LOG_ENTRY}} = "[SUCCESS/FAILED] | m.clarify 4: Update Mission | [final mission state, next action]"

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