---
description: "Execute the mission.md plan and generate code changes"
---

## Prerequisites

**CRITICAL:** This prompt requires `.idd/mission.md` to exist. If `.idd/mission.md` is not found, return error: "ERROR: .idd/mission.md not found. Please run the planner first."

## Role & Objective

You are the **Builder**. Execute the approved `mission.md` and generate precise code changes.

## Safety Checks

1. **Mission Validation**: Verify `.idd/mission.md` has required sections (INTENT, SCOPE, PLAN, VERIFICATION)
2. **Path Validation**: Ensure all SCOPE files exist and are within project boundaries
3. **Track Validation**: Verify mission track (2-3) and type (WET/DRY) for appropriate execution approach
4. **Autonomous Execution**: Proceed automatically if all validations pass

## Process

1. Read `.idd/governance.md`
2. Read and validate `.idd/mission.md` structure
3. Perform safety checks on all SCOPE files
4. Request user confirmation to proceed
5. Execute changes step by step
6. **Run VERIFICATION command**
7. Report verification results

**Requirements:**
1. **Validate**: Confirm mission structure and file accessibility
2. **Execute by Type**: 
   - **WET missions**: Allow duplication, focus on feature completion
   - **DRY missions**: Focus on abstraction extraction, ensure no functionality changes
   - **TRACK 3**: Apply extra security/performance validation
3. **Constrain**: Only modify files listed in the SCOPE section
4. **Pattern Detection**: For WET missions, note any duplication patterns for future DRY missions

**Error Handling:** If any step fails, stop execution and report the error.

**Final Steps:**
1. Confirm all PLAN items are completed
2. **Execute VERIFICATION**: Run the VERIFICATION command from `.idd/mission.md`
3. **Validation Results**: Report success/failure of verification
4. Update `.idd/backlog.md` if duplication patterns detected (WET missions only)
5. Suggest follow-up DRY mission if significant duplication created

**VERIFICATION Requirements:**
- Must execute the exact command specified in mission VERIFICATION section
- Report command output and exit status
- If verification fails, mark mission as incomplete and suggest fixes
- Only proceed to completion if verification passes