# AI-NATIVE TEST SPECIFICATIONS

## PURPOSE
Standard format for AI-native test specifications that enable AI agents to validate prompt logic through natural language reasoning.

## TEST SPECIFICATION FORMAT

### Basic Structure
```markdown
# TEST: [Descriptive Test Name]

## SCENARIO
**Given**: [Initial conditions and context]
**When**: [Prompt execution trigger or action]
**Then**: [Expected outcomes and results]
**Because**: [Reasoning for why this outcome is expected]

## MOCK DATA
**User Input**: [Test input data]
**AI Analysis**: [Simulated AI response/analysis]
**System State**: [Simulated file system or environment state]

## ASSERTIONS
- [Natural language assertion 1 with expected reasoning]
- [Natural language assertion 2 with expected reasoning]
- [Natural language assertion 3 with expected reasoning]

## VALIDATION METHOD
[Instructions for AI on how to validate this test case through reasoning]

## SUCCESS CRITERIA
**Pass Conditions**: [Clear criteria for test success]
**Fail Conditions**: [Clear criteria for test failure]

## EXPECTED REASONING TRACE
[Step-by-step logical process AI should follow to validate]
```

### Example Test Specification
```markdown
# TEST: Simple Feature Track Assignment

## SCENARIO
**Given**: Clean project with no existing mission
**When**: User requests "add user profile endpoint"
**Then**: Mission should be assigned Track 2 with appropriate scope
**Because**: Simple CRUD endpoint with 2-3 files falls into standard feature complexity

## MOCK DATA
**User Input**: "add user profile endpoint"
**AI Analysis**: 
  - Intent: "Add user profile endpoint with CRUD operations"
  - Suggested Files: ["routes/profile.js", "controllers/profile.js"]
  - Complexity Indicators: ["endpoint", "CRUD", "profile"]
  - Domain Multipliers: None detected

**System State**: 
  - No existing .mission/mission.md file
  - Clean project structure

## ASSERTIONS
- Track calculation should result in Track 2 based on 2 files and no domain multipliers
- Mission type should be WET for initial implementation
- Scope should contain exactly 2 files: routes and controller
- Plan should include CRUD operations and testing steps
- Verification command should be safe and non-destructive

## VALIDATION METHOD
AI should mentally execute the M.PLAN prompt logic:
1. Parse user intent for complexity indicators
2. Simulate AI analysis of suggested implementation
3. Apply track complexity calculation rules
4. Generate mission structure following governance principles
5. Validate results against expected outcomes

## SUCCESS CRITERIA
**Pass Conditions**: 
- Track calculation logic produces Track 2
- Mission structure includes all required sections
- Scope and plan align with simple feature requirements

**Fail Conditions**:
- Track calculation produces wrong track number
- Mission structure missing required sections
- Scope or plan inappropriate for complexity level

## EXPECTED REASONING TRACE
1. "User input 'add user profile endpoint' indicates a simple CRUD feature"
2. "AI analysis suggests 2 files (routes + controller) with no security/performance concerns"
3. "Track calculation: 2 files = Track 2 base, no domain multipliers = Track 2 final"
4. "Mission structure generated with WET type, planned status, appropriate scope"
5. "Verification: all assertions satisfied, test passes"
```

## SPECIFICATION GUIDELINES

### Writing Effective Scenarios
- Use clear Given/When/Then structure
- Include sufficient context for AI understanding
- Specify expected reasoning, not just outcomes
- Make assertions testable through logical reasoning

### Mock Data Principles
- Provide realistic but simplified data
- Include all information needed for validation
- Avoid programming language constructs
- Use natural language descriptions

### Assertion Best Practices
- State expectations in natural language
- Include reasoning for why assertion should be true
- Make assertions specific and measurable
- Cover both positive and negative cases

### Validation Instructions
- Provide clear steps for AI reasoning process
- Specify what logic should be applied
- Include decision points and criteria
- Make validation process auditable

This format enables AI agents to execute tests through reasoning rather than code execution, aligning with the AI-native testing philosophy.
