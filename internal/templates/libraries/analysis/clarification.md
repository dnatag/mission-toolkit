# CLARIFICATION ANALYSIS TEMPLATE

## Purpose
Scan user intent for ambiguous requirements that need clarification before mission planning.

## Analysis Categories

### 1. Scope & Requirements
- **Vague Problem**: "bug", "issue", "broken" without symptoms?
- **Ambiguous Scope**: "improve", "enhance" without specifics?
- **Missing Tech**: Unspecified framework, library, or database?

### 2. Domain & Risk Context (Critical for Plan Spec)
- **Security**: Does this involve auth, crypto, or PII?
- **Performance**: Are there latency or throughput constraints?
- **Compliance**: Does this touch GDPR, audit logs, or financial data?
- **Complexity**: Does this involve complex algorithms or AI?

*If these are relevant but undefined, you may need to clarify.*

## Decision Logic & Output

Based on your analysis, choose ONE of the following outputs.

### A. Clarification Needed (STOP)
If ANY [CRITICAL] details are missing that prevent you from defining the **Scope** or **Domain**, you MUST stop.

**Output Format:**
```
🛑 CLARIFICATION NEEDED

I need a few more details to plan this mission effectively:

1. [CRITICAL] [Specific Question]
2. [IMPORTANT] [Specific Question]

Please provide these details so I can proceed.
```

### B. Proceed with Assumptions (CAUTION)
If only [HELPFUL] details are missing, you may proceed but must state your assumptions.

**Output Format:**
```
⚠️ PROCEEDING WITH ASSUMPTIONS

- Assuming [Assumption 1] based on [Reason]
- Assuming [Assumption 2] based on [Reason]

(Proceed to next step)
```

### C. Intent Clear (PROCEED)
If the request is clear and actionable.

**Output Format:**
```
✅ INTENT CLEAR
(Proceed to next step)
```

## Question Quality Standards
- **Specific**: Ask for concrete details (e.g., "Which error code?" not "What's wrong?").
- **Actionable**: The answer should directly affect the implementation plan.
- **Contextual**: Reference the user's original goal.
