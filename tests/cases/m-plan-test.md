# TEST: M.PLAN Simple Feature Processing

## SCENARIO
**Given**: Clean project with no existing mission and standard directory structure
**When**: User requests "add user profile endpoint" through M.PLAN prompt
**Then**: Mission should be generated with Track 2 complexity and appropriate WET structure
**Because**: Simple CRUD endpoint with 2-3 files represents standard feature complexity without domain multipliers

## MOCK DATA
**User Input**: "add user profile endpoint"

**AI Analysis**: 
- Intent: "Add user profile endpoint with CRUD operations for user data management"
- Suggested Files: ["routes/profile.js", "controllers/profile-controller.js"]
- Complexity Assessment: "Standard CRUD endpoint implementation"
- Security Concerns: None identified (basic data operations)
- Performance Impact: Minimal (standard database queries)
- Domain Multipliers: None detected

**System State**: 
- Mission File: Does not exist (.mission/mission.md absent)
- Project Structure: Standard directories present
- Previous Missions: None (clean project)
- Governance Rules: Standard Mission Toolkit governance active

## ASSERTIONS
- Track calculation should result in Track 2 based on 2 files and no domain multipliers
- Mission type should be WET for initial implementation approach
- Mission status should be "planned" after generation
- Scope should contain exactly 2 files: routes/profile.js and controllers/profile-controller.js
- Plan should include CRUD operations, testing, and duplication allowance note
- Verification command should be safe (npm test or similar) and non-destructive
- Intent should be refined from user input to include CRUD specifics
- Mission should include all required sections: INTENT, SCOPE, PLAN, VERIFICATION

## VALIDATION METHOD
AI should execute M.PLAN prompt logic through the following reasoning process:

1. **Intent Analysis**: Parse "add user profile endpoint" to identify clear CRUD intent
2. **Complexity Assessment**: Analyze suggested files (2 files) and scan for domain multipliers
3. **Track Calculation**: Apply rules - 2 files = Track 2 base, no multipliers = Track 2 final
4. **Mission Generation**: Create WET mission structure following governance principles
5. **Safety Validation**: Ensure verification command is safe and scope is appropriate
6. **Structure Validation**: Confirm all required mission sections are present and complete

## SUCCESS CRITERIA
**Pass Conditions**: 
- Track calculation produces Track 2 (not Track 1 or Track 3)
- Mission structure includes INTENT, SCOPE, PLAN, VERIFICATION sections
- Scope contains appropriate files for CRUD endpoint implementation
- Plan includes relevant steps for profile endpoint development
- Verification command is safe and tests relevant functionality
- Mission type is WET with "planned" status

**Fail Conditions**:
- Track calculation produces incorrect track (Track 1, 3, or 4)
- Mission structure missing required sections
- Scope inappropriate for stated intent (too many/few files, wrong file types)
- Plan steps irrelevant to profile endpoint development
- Verification command unsafe or destructive
- Mission type incorrect or status inappropriate

## EXPECTED REASONING TRACE
1. **Input Processing**: "User input 'add user profile endpoint' is clear and specific, indicating CRUD operations for user profile data"

2. **AI Analysis Integration**: "Simulated AI analysis suggests 2 implementation files (routes + controller) with standard CRUD complexity and no special concerns"

3. **Domain Multiplier Scan**: "Scanning input for domain multipliers: no security terms (auth, JWT, encryption), no performance terms (real-time, sub-second, cache), no high-risk terms (payment, compliance) → no multipliers detected"

4. **Track Calculation**: "Base complexity: 2 files maps to Track 2 according to governance rules (1-5 files = Track 2). Domain multipliers: 0 detected. Final track: Track 2 + 0 = Track 2"

5. **Mission Structure Generation**: "Generate WET mission (initial implementation) with planned status, refined intent including CRUD specifics, scope with 2 identified files, plan with CRUD steps plus testing"

6. **Safety Validation**: "Verification command 'npm test -- --grep profile' is safe (read-only testing), scope files are appropriate for intent, no destructive operations"

7. **Final Validation**: "Generated mission satisfies all assertions: correct track, complete structure, appropriate scope and plan, safe verification"

## CONFIDENCE ASSESSMENT
**Expected Confidence**: HIGH
- Clear user intent with unambiguous requirements
- Straightforward complexity analysis with no edge cases
- Standard track calculation with well-defined rules
- Typical WET mission structure following established patterns

## RELATED TEST SCENARIOS
This test validates the core M.PLAN logic for standard features. Related scenarios to test:
- Security feature with domain multipliers (Track escalation)
- Performance-critical feature (Track escalation)
- Vague input requiring clarification
- Complex feature requiring decomposition (Track 4)
