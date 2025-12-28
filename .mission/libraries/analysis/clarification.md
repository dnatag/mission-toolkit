# CLARIFICATION ANALYSIS TEMPLATE

## Purpose
Scan user intent for ambiguous requirements that need clarification before mission planning.

## Analysis Categories

### Bug/Issue Reports
- **Vague Problem**: "bug", "issue", "broken", "doesn't work" without specific symptoms
- **Missing Reproduction**: No steps to reproduce the problem
- **No Expected Behavior**: What should happen vs. what actually happens
- **Examples**: "Fix display bug", "API is broken", "Login doesn't work"

### Feature Requests  
- **Ambiguous Scope**: "improve", "enhance", "make better" without specifics
- **Missing Requirements**: No acceptance criteria or success definition
- **Unclear Integration**: How it connects to existing functionality
- **Examples**: "Make UI better", "Add user management", "Improve performance"

### Technical Tasks
- **Unspecified Technology**: Framework, library, or approach not defined
- **Missing Context**: Why this change is needed
- **Unclear Dependencies**: What other systems are affected
- **Examples**: "Add authentication", "Refactor code", "Update database"

### General Requirements
- **Technology Stack**: Unspecified frameworks, databases, or libraries
- **Business Logic**: Unclear validation rules, data relationships, or workflows
- **Integration Points**: External APIs, services, or data sources without details
- **Performance Requirements**: Unspecified response times, throughput, or scalability needs
- **Security Requirements**: Authentication, authorization, or data protection specifics

## Decision Logic
1. **Step 1**: Identify category (Bug/Feature/Technical/General)
2. **Step 2**: Check for missing details in that category
3. **Step 3**: Generate specific questions for gaps found
4. **Step 4**: **MANDATORY**: If ANY [CRITICAL] questions exist, MUST create clarification mission
5. **Step 5**: **DISCRETIONARY**: For [IMPORTANT] or [HELPFUL] questions, decide whether clarification mission is needed

## Clarification Thresholds
**MUST CLARIFY** - Mandatory clarification mission if ANY of these are missing:
- **Bug Reports**: Specific error symptoms, reproduction steps, or expected behavior
- **Feature Requests**: Clear functional requirements or integration points
- **Technical Tasks**: Technology choice, implementation approach, or system dependencies
- **Security/Performance**: Specific requirements for authentication, authorization, or performance targets

**MAY CLARIFY** - Use judgment to decide if clarification mission needed:
- Nice-to-have details that don't affect core implementation
- Styling preferences that can be decided during implementation
- Optional features that can be added later

## Context Preservation
When generating clarification questions:
- **Maintain Intent**: Reference user's original goal in questions
- **Build Context**: Each question should help refine the mission scope
- **Avoid Scope Creep**: Focus only on details needed for current intent

## Output Format
If clarifications needed, generate numbered questions for each ambiguous area identified.

**Question Quality Standards:**
- Each question must be specific and actionable
- Avoid yes/no questions - ask for concrete details
- Prioritize questions by impact: critical details first, nice-to-have details last
- Preserve user's original intent while seeking clarification

**Example Output:**
```
1. [CRITICAL] What specific error message appears when the login fails?
2. [IMPORTANT] Which browsers have you tested this issue on?
3. [HELPFUL] What user role should have access to this feature?
```