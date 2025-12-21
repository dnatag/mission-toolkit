# AI-NATIVE TEST EXECUTION GUIDE

## PURPOSE
Instructions for AI agents to execute AI-native tests through reasoning and logical analysis rather than code execution.

## EXECUTION METHODOLOGY

### Pre-Execution Setup
1. **Load Test Specification**: Read and understand the complete test case
2. **Identify Mock Scenarios**: Extract mock data and system state information
3. **Review Assertions**: Understand what needs to be validated and why
4. **Prepare Reasoning Framework**: Set up logical analysis approach

### Core Execution Process

#### Step 1: Scenario Analysis
```markdown
**Process**: Analyze Given/When/Then scenario structure
**Actions**:
- Parse initial conditions (Given)
- Identify trigger action (When) 
- Extract expected outcomes (Then)
- Understand reasoning (Because)

**Example**:
Given: "Clean project with no existing mission"
When: "User requests 'add user profile endpoint'"
Then: "Mission should be assigned Track 2"
Because: "Simple CRUD endpoint falls into standard complexity"
```

#### Step 2: Mock Data Integration
```markdown
**Process**: Incorporate mock data into reasoning context
**Actions**:
- Load user input from mock data
- Apply simulated AI analysis responses
- Establish system state conditions
- Set up environmental context

**Example**:
User Input: "add user profile endpoint"
AI Analysis: Intent clarified, 2 files suggested, no security concerns
System State: Clean project, no existing missions
```

#### Step 3: Prompt Logic Simulation
```markdown
**Process**: Mentally execute prompt logic using governance rules
**Actions**:
- Apply intent parsing logic to user input
- Simulate complexity analysis algorithms
- Execute track calculation rules
- Generate mission structure following governance
- Apply safety and validation checks

**Reasoning Trace Example**:
1. "Parse 'add user profile endpoint' → clear CRUD intent"
2. "Apply complexity analysis → 2 files suggested"
3. "Calculate track: 2 files = Track 2 base, no multipliers = Track 2 final"
4. "Generate WET mission with appropriate scope and plan"
```

#### Step 4: Assertion Validation
```markdown
**Process**: Validate each assertion against simulated results
**Actions**:
- Compare expected vs actual outcomes
- Apply assertion reasoning logic
- Document validation results
- Identify any discrepancies

**Validation Example**:
Assertion: "Track should be 2 based on file count"
Validation: "2 files → Track 2 base, no multipliers → Track 2 final ✓"
Result: PASS - Logic correctly applied
```

#### Step 5: Result Documentation
```markdown
**Process**: Document validation results with reasoning
**Actions**:
- Record pass/fail status for each assertion
- Explain reasoning for validation decisions
- Identify any logic gaps or errors
- Provide confidence assessment

**Documentation Format**:
Test: [Test Name]
Status: PASS/FAIL
Assertions: [X/Y passed]
Reasoning: [Detailed explanation]
Confidence: HIGH/MEDIUM/LOW
```

## EXECUTION PATTERNS

### Simple Feature Test Execution
```markdown
**Scenario Type**: Standard CRUD feature request
**Execution Approach**:
1. Confirm user intent is clear and specific
2. Apply standard complexity analysis (file count, domain scan)
3. Calculate track using base rules (no multipliers expected)
4. Validate mission structure follows WET pattern
5. Check verification command safety

**Common Assertions**:
- Track calculation accuracy
- Mission structure completeness  
- Scope appropriateness
- Plan step relevance
```

### Security Feature Test Execution
```markdown
**Scenario Type**: Security-sensitive feature request
**Execution Approach**:
1. Identify security-related keywords in user input
2. Apply domain multiplier detection for security concerns
3. Calculate track with security multiplier escalation
4. Validate robust planning for security requirements
5. Check security-appropriate verification

**Common Assertions**:
- Domain multiplier detection (security, high-risk)
- Track escalation logic (Track 2 → Track 3)
- Security-specific plan steps
- Compliance considerations
```

### Clarification Workflow Test Execution
```markdown
**Scenario Type**: Ambiguous or incomplete request
**Execution Approach**:
1. Identify ambiguities or missing information
2. Evaluate need for clarification vs direct processing
3. Generate appropriate clarification questions
4. Process clarification responses
5. Validate mission updates after clarification

**Common Assertions**:
- Ambiguity detection accuracy
- Clarification question quality
- Response processing logic
- Mission update appropriateness
```

## REASONING VALIDATION TECHNIQUES

### Logic Chain Verification
```markdown
**Technique**: Trace logical steps from input to output
**Process**:
1. Document each reasoning step
2. Verify each step follows governance rules
3. Check for logical consistency
4. Identify any gaps or errors

**Example Chain**:
Input → Intent Parsing → Complexity Analysis → Track Calculation → Mission Generation
```

### Counterfactual Analysis
```markdown
**Technique**: Test alternative scenarios to validate logic
**Process**:
1. Modify one input variable
2. Predict how output should change
3. Verify logic handles variation correctly
4. Confirm robustness of reasoning

**Example**:
"If user input included 'with authentication', would security multiplier be detected?"
```

### Boundary Condition Testing
```markdown
**Technique**: Test edge cases and limits
**Process**:
1. Identify boundary conditions (file counts, complexity thresholds)
2. Apply logic to boundary cases
3. Verify consistent behavior at boundaries
4. Check for off-by-one errors or edge case failures

**Example**:
"Does 5 files result in Track 2 and 6 files result in Track 3?"
```

## ERROR DETECTION AND HANDLING

### Common Logic Errors
- **Incorrect Track Calculation**: Wrong file count or multiplier application
- **Missing Domain Multipliers**: Failed to detect security/performance concerns
- **Inappropriate Mission Structure**: Wrong type, status, or section content
- **Unsafe Verification**: Destructive or inappropriate verification commands

### Error Analysis Process
1. **Identify Discrepancy**: Compare expected vs actual results
2. **Trace Logic Path**: Find where reasoning diverged from expected
3. **Classify Error Type**: Logic error, missing rule, or incorrect assumption
4. **Assess Impact**: Determine severity and scope of error
5. **Document Findings**: Record error details for improvement

### Confidence Assessment
- **HIGH**: All assertions pass, logic is clear and consistent
- **MEDIUM**: Most assertions pass, minor logic questions remain
- **LOW**: Multiple assertion failures or unclear reasoning paths

## EXECUTION BEST PRACTICES

### Systematic Approach
- Follow execution steps in order
- Document reasoning at each step
- Validate assumptions explicitly
- Check work against governance rules

### Thorough Analysis
- Consider multiple interpretation paths
- Test edge cases and boundary conditions
- Verify logic consistency across scenarios
- Question assumptions and validate rules

### Clear Documentation
- Record detailed reasoning traces
- Explain validation decisions
- Note confidence levels and uncertainties
- Provide actionable feedback for improvements

This execution guide enables AI agents to systematically validate prompt logic through structured reasoning, ensuring comprehensive and reliable test results without requiring code execution environments.
