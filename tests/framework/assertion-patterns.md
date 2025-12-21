# AI-NATIVE ASSERTION PATTERNS

## PURPOSE
Standard patterns for reasoning-based assertions that AI agents can validate through logical analysis rather than code execution.

## ASSERTION PATTERN LIBRARY

### Track Complexity Assertions

#### Basic Track Calculation
```markdown
**Pattern**: track-calculation-assertion
**Format**: "Track calculation should result in Track [X] based on [Y] files and [Z] domain multipliers"
**Reasoning**: "File count [Y] maps to base Track [base], domain multipliers [list] add [count], final Track [X]"
**Validation**: AI applies track complexity rules and verifies calculation logic
```

#### Domain Multiplier Detection
```markdown
**Pattern**: domain-multiplier-assertion
**Format**: "Domain multiplier '[multiplier-name]' should be detected due to [trigger-phrase] in user input"
**Reasoning**: "Input contains '[trigger-phrase]' which indicates [domain-concern] requiring [multiplier-name]"
**Validation**: AI scans input for domain-specific keywords and applies multiplier detection rules
```

#### Track Escalation
```markdown
**Pattern**: track-escalation-assertion
**Format**: "Track should escalate from [original] to [final] due to [reason]"
**Reasoning**: "Base complexity [original] + [multiplier-factors] = [final] (capped at Track 3)"
**Validation**: AI traces escalation logic and verifies governance rule application
```

### Mission Structure Assertions

#### Required Sections
```markdown
**Pattern**: mission-structure-assertion
**Format**: "Mission should contain all required sections: [INTENT, SCOPE, PLAN, VERIFICATION]"
**Reasoning**: "Governance requires complete mission structure for proper execution and tracking"
**Validation**: AI checks generated mission for presence and completeness of each section
```

#### Mission Type Assignment
```markdown
**Pattern**: mission-type-assertion
**Format**: "Mission type should be [WET/DRY] for [reason]"
**Reasoning**: "Initial implementation uses WET approach, refactoring uses DRY approach"
**Validation**: AI determines appropriate mission type based on governance principles
```

#### Status Progression
```markdown
**Pattern**: status-progression-assertion
**Format**: "Mission status should be '[status]' after [action]"
**Reasoning**: "Workflow progression: planned → active → completed/failed"
**Validation**: AI verifies status follows proper workflow transitions
```

### Scope and Planning Assertions

#### File Scope Validation
```markdown
**Pattern**: scope-validation-assertion
**Format**: "Scope should contain [count] files: [file-list]"
**Reasoning**: "Implementation requires [file-types] to achieve [intent] with [approach]"
**Validation**: AI evaluates if proposed files are necessary and sufficient for intent
```

#### Plan Step Appropriateness
```markdown
**Pattern**: plan-step-assertion
**Format**: "Plan should include steps for [requirement] due to [complexity-factor]"
**Reasoning**: "[Complexity-factor] requires [specific-steps] to ensure [quality-attribute]"
**Validation**: AI checks if plan steps address all complexity requirements
```

#### WET vs DRY Planning
```markdown
**Pattern**: wet-dry-planning-assertion
**Format**: "Plan should include 'Allow duplication' note for WET mission"
**Reasoning**: "WET missions prioritize exploration over abstraction in initial implementation"
**Validation**: AI verifies planning approach matches mission type
```

### Verification and Safety Assertions

#### Verification Command Safety
```markdown
**Pattern**: verification-safety-assertion
**Format**: "Verification command should be safe and non-destructive"
**Reasoning**: "Verification validates implementation without modifying system state"
**Validation**: AI analyzes command for destructive operations (rm, chmod 777, etc.)
```

#### Verification Appropriateness
```markdown
**Pattern**: verification-appropriateness-assertion
**Format**: "Verification should test [specific-aspect] relevant to [mission-intent]"
**Reasoning**: "Verification must validate the core functionality being implemented"
**Validation**: AI checks if verification command tests the right functionality
```

### Clarification Workflow Assertions

#### Clarification Trigger Detection
```markdown
**Pattern**: clarification-trigger-assertion
**Format**: "Clarification should be triggered due to [ambiguity-type] in user input"
**Reasoning**: "Input lacks [specific-information] needed for proper implementation planning"
**Validation**: AI identifies missing information that prevents clear mission planning
```

#### Clarification Question Quality
```markdown
**Pattern**: clarification-question-assertion
**Format**: "Clarification questions should address [specific-gaps] in requirements"
**Reasoning**: "Questions must elicit information needed to resolve [ambiguity-type]"
**Validation**: AI evaluates if questions would resolve identified ambiguities
```

#### Post-Clarification Updates
```markdown
**Pattern**: post-clarification-assertion
**Format**: "After clarification, [mission-aspect] should be updated to reflect [new-information]"
**Reasoning**: "Clarification responses provide [specific-details] that change [planning-decisions]"
**Validation**: AI traces how clarification responses should modify mission planning
```

### Error and Edge Case Assertions

#### Invalid Input Handling
```markdown
**Pattern**: invalid-input-assertion
**Format**: "Invalid input [input-example] should trigger [error-response]"
**Reasoning**: "Input validation prevents processing of malformed or malicious requests"
**Validation**: AI determines appropriate error handling for various invalid inputs
```

#### Boundary Condition Testing
```markdown
**Pattern**: boundary-condition-assertion
**Format**: "Edge case [condition] should result in [expected-behavior]"
**Reasoning**: "Boundary conditions test limits of [algorithm/rule] implementation"
**Validation**: AI applies logic to edge cases and verifies consistent behavior
```

#### Governance Rule Enforcement
```markdown
**Pattern**: governance-enforcement-assertion
**Format**: "Governance rule [rule-name] should prevent [violation-scenario]"
**Reasoning**: "Rule enforcement maintains [quality-attribute] across all missions"
**Validation**: AI checks if governance rules are properly applied in scenario
```

## ASSERTION COMPOSITION PATTERNS

### Compound Assertions
```markdown
**Pattern**: compound-assertion
**Format**: "Mission should have Track [X] AND include [security-steps] AND use [verification-type]"
**Reasoning**: Multiple related requirements that must all be satisfied
**Validation**: AI validates each component assertion independently and collectively
```

### Conditional Assertions
```markdown
**Pattern**: conditional-assertion
**Format**: "IF [condition] THEN [assertion] BECAUSE [reasoning]"
**Reasoning**: Context-dependent expectations based on specific conditions
**Validation**: AI evaluates condition first, then applies appropriate assertion
```

### Sequential Assertions
```markdown
**Pattern**: sequential-assertion
**Format**: "First [assertion-1], then [assertion-2], finally [assertion-3]"
**Reasoning**: Order-dependent validations for multi-step processes
**Validation**: AI validates assertions in specified sequence
```

## VALIDATION METHODOLOGY

### Reasoning Process
1. **Parse Assertion**: Extract expected outcome and reasoning
2. **Apply Logic**: Use prompt logic rules to determine actual outcome
3. **Compare Results**: Check if actual matches expected
4. **Explain Variance**: If mismatch, identify where logic diverges
5. **Assess Confidence**: Rate confidence in validation result

### Quality Criteria
- **Specificity**: Assertions should be precise and measurable
- **Testability**: Must be validatable through logical reasoning
- **Relevance**: Should test important aspects of prompt behavior
- **Clarity**: Reasoning should be transparent and auditable

This pattern library enables consistent, comprehensive validation of prompt logic through structured reasoning rather than code execution.
