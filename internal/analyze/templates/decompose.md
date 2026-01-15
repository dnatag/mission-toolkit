# DECOMPOSE ANALYSIS TEMPLATE

## Current Intent

{{.CurrentIntent}}

## Current Scope

{{.CurrentScope}}

## Purpose
Break down an Epic (Track 4) intent into smaller, atomic sub-intents that can each be executed as independent missions.

## Decomposition Guidelines

### 1. Identify Natural Boundaries
- **Feature boundaries**: Separate distinct features or capabilities
- **Layer boundaries**: Split by architectural layers (API, service, data)
- **Domain boundaries**: Separate by business domain or module

### 2. Sub-Intent Requirements
Each sub-intent MUST be:
- **Independent**: Can be implemented without other sub-intents
- **Atomic**: Fits within Track 2-3 complexity (1-9 files)
- **Verifiable**: Has clear success criteria
- **Valuable**: Delivers incremental value when completed

### 3. Ordering Considerations
- **Dependencies**: List sub-intents in dependency order (foundations first)
- **Risk**: Front-load high-risk items for early feedback
- **Value**: Prioritize high-value items when dependencies allow

## Output Format

Produce a JSON object with decomposed sub-intents.

```json
{
  "action": "DECOMPOSE",
  "sub_intents": [
    {
      "intent": "Clear, actionable description",
      "rationale": "Why this is a separate sub-intent",
      "estimated_files": 3,
      "dependencies": []
    }
  ],
  "decomposition_rationale": "Overall explanation of decomposition strategy"
}
```

## Examples

### Example 1: Payment System
**Intent**: "Build complete payment processing system"
```json
{
  "action": "DECOMPOSE",
  "sub_intents": [
    {
      "intent": "Create payment data models and database schema",
      "rationale": "Foundation layer needed by all other components",
      "estimated_files": 3,
      "dependencies": []
    },
    {
      "intent": "Implement payment validation service",
      "rationale": "Core business logic for payment validation",
      "estimated_files": 4,
      "dependencies": ["payment data models"]
    },
    {
      "intent": "Add Stripe payment gateway integration",
      "rationale": "External integration isolated for easier testing",
      "estimated_files": 3,
      "dependencies": ["payment validation service"]
    },
    {
      "intent": "Create payment API endpoints",
      "rationale": "API layer depends on service layer",
      "estimated_files": 4,
      "dependencies": ["payment validation service", "Stripe integration"]
    }
  ],
  "decomposition_rationale": "Split by architectural layers: data → service → integration → API"
}
```

### Example 2: Authentication System
**Intent**: "Add user authentication with OAuth and JWT"
```json
{
  "action": "DECOMPOSE",
  "sub_intents": [
    {
      "intent": "Create user model and authentication database schema",
      "rationale": "Data layer foundation",
      "estimated_files": 2,
      "dependencies": []
    },
    {
      "intent": "Implement JWT token generation and validation",
      "rationale": "Core auth mechanism, independent of OAuth",
      "estimated_files": 3,
      "dependencies": ["user model"]
    },
    {
      "intent": "Add OAuth provider integration (Google)",
      "rationale": "External integration, can be added incrementally",
      "estimated_files": 4,
      "dependencies": ["JWT implementation"]
    },
    {
      "intent": "Create authentication middleware and protected routes",
      "rationale": "Depends on JWT validation being in place",
      "estimated_files": 3,
      "dependencies": ["JWT implementation"]
    }
  ],
  "decomposition_rationale": "Split by auth mechanism: data → JWT core → OAuth extension → middleware"
}
```
