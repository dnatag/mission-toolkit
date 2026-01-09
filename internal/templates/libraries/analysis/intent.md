# INTENT ANALYSIS TEMPLATE

## Purpose
Distill raw user input into a clear, actionable, and scoped intent statement.

## Analysis Steps

### 1. Identify Core Action (The Verb)
- **Create**: New feature, file, or component.
- **Refactor**: Changing structure without changing behavior.
- **Fix**: Correcting a bug or error.
- **Optimize**: Improving performance or efficiency.
- **Update**: Upgrading dependencies or content.
- **Delete**: Removing code or features.

### 2. Identify Target Scope (The Noun)
- What specific component, module, or file is being acted upon?
- *Preference*: Use specific file paths or package names if known (e.g., `auth/handler.go` instead of "auth").
- *Example*: "Login Handler", "User Database", "API Routes".

### 3. Identify Constraints (The Rules)
- **Technology**: Specific libraries (e.g., "using Gin", "with Zap logger").
- **Behavior**: Specific outcomes (e.g., "return 404 on error").
- **Standards**: Coding styles or patterns (e.g., "table-driven tests").

## Refinement Rules
1.  **Be Specific**: Change "Fix login" to "Fix 500 error in login handler on invalid password".
2.  **Be Concise**: Remove polite filler ("Please", "I want to").
3.  **Be Technical**: Use correct terminology (e.g., "Middleware" instead of "middle thing").
4.  **Be Atomic**: If the request contains multiple distinct goals (e.g., "Fix bug AND add feature"), split them or focus on the primary one.

## Handling Ambiguity
If the intent is too vague to refine (e.g., "It doesn't work"), do NOT guess.
- **Action**: Flag as AMBIGUOUS.
- **Output**: "AMBIGUOUS: [Reason]"

## Output Format

Produce a JSON object with action and refined intent.

```json
{
  "action": "PROCEED" | "AMBIGUOUS",
  "refined_intent": "Add JWT authentication",
  "reason": "Too vague"  // Only if AMBIGUOUS
}
```

**Examples:**
- Raw: "Can you make the login faster?"
  ```json
  {
    "action": "PROCEED",
    "refined_intent": "Optimize login handler performance to reduce latency < 200ms"
  }
  ```

- Raw: "I need a new api for users."
  ```json
  {
    "action": "PROCEED",
    "refined_intent": "Create User REST API endpoints with CRUD operations"
  }
  ```

- Raw: "The auth is broken."
  ```json
  {
    "action": "PROCEED",
    "refined_intent": "Fix authentication failure in JWT middleware"
  }
  ```

- Raw: "Make it better."
  ```json
  {
    "action": "AMBIGUOUS",
    "reason": "Missing specific target or goal"
  }
  ```
