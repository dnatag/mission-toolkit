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
1. Scan for vague keywords → Generate specific questions
2. Check for missing technical details → Request specifications  
3. Validate completeness → Proceed or clarify

## Output Format
If clarifications needed, generate numbered questions for each ambiguous area identified.