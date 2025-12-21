# AI-NATIVE MOCK SCENARIOS

## PURPOSE
Library of reusable mock scenarios for AI-native testing using natural language descriptions instead of programming objects.

## MOCK SCENARIO LIBRARY

### User Input Scenarios

#### Simple Feature Requests
```markdown
**Scenario**: simple-crud-feature
**User Input**: "add user profile endpoint"
**Expected Complexity**: Standard feature (Track 2)
**Characteristics**: Clear intent, single domain, basic CRUD operations
```

```markdown
**Scenario**: simple-component-request  
**User Input**: "create login form component"
**Expected Complexity**: Standard feature (Track 2)
**Characteristics**: UI component, straightforward implementation
```

#### Security-Sensitive Requests
```markdown
**Scenario**: authentication-feature
**User Input**: "add JWT authentication with role-based access control"
**Expected Complexity**: Robust feature (Track 3)
**Characteristics**: Security domain, access control, authentication
**Domain Multipliers**: ["high-risk-integration", "regulatory-security"]
```

```markdown
**Scenario**: payment-integration
**User Input**: "integrate Stripe payment processing with PCI compliance"
**Expected Complexity**: Robust feature (Track 3)  
**Characteristics**: Payment processing, compliance requirements, external API
**Domain Multipliers**: ["high-risk-integration", "regulatory-security"]
```

#### Performance-Critical Requests
```markdown
**Scenario**: real-time-feature
**User Input**: "add real-time chat with sub-second message delivery"
**Expected Complexity**: Robust feature (Track 3)
**Characteristics**: Real-time processing, performance requirements
**Domain Multipliers**: ["performance-critical"]
```

```markdown
**Scenario**: search-optimization
**User Input**: "optimize search with caching and sub-second response times"
**Expected Complexity**: Robust feature (Track 3)
**Characteristics**: Performance optimization, caching, response time requirements
**Domain Multipliers**: ["performance-critical", "complex-algorithms"]
```

#### Vague or Ambiguous Requests
```markdown
**Scenario**: vague-request
**User Input**: "make the app better"
**Expected Complexity**: Requires clarification
**Characteristics**: Unclear intent, no specific requirements, needs clarification
```

```markdown
**Scenario**: overly-broad-request
**User Input**: "build complete e-commerce platform with payments and analytics"
**Expected Complexity**: Epic (Track 4 - decompose)
**Characteristics**: Multiple systems, broad scope, requires decomposition
```

### AI Analysis Scenarios

#### Standard Analysis Response
```markdown
**Scenario**: standard-ai-analysis
**Intent**: "Add user profile endpoint with CRUD operations"
**Suggested Files**: ["routes/profile.js", "controllers/profile.js"]
**Complexity Assessment**: "Standard CRUD endpoint implementation"
**Security Concerns**: None identified
**Performance Impact**: Minimal
**Domain Multipliers**: None
```

#### Security-Enhanced Analysis
```markdown
**Scenario**: security-ai-analysis
**Intent**: "Add JWT authentication system with role-based access control"
**Suggested Files**: ["auth/jwt.js", "auth/rbac.js", "middleware/auth.js", "models/roles.js"]
**Complexity Assessment**: "Security-sensitive authentication system"
**Security Concerns**: Authentication, authorization, token management
**Performance Impact**: Moderate (token validation overhead)
**Domain Multipliers**: ["high-risk-integration", "regulatory-security"]
```

#### Performance-Focused Analysis
```markdown
**Scenario**: performance-ai-analysis
**Intent**: "Optimize search with caching and real-time updates"
**Suggested Files**: ["search/engine.js", "cache/redis.js", "realtime/updates.js"]
**Complexity Assessment**: "Performance-critical search optimization"
**Security Concerns**: Cache security, data exposure
**Performance Impact**: High (sub-second requirements)
**Domain Multipliers**: ["performance-critical", "complex-algorithms"]
```

### System State Scenarios

#### Clean Project State
```markdown
**Scenario**: clean-project
**Mission File**: Does not exist
**Project Structure**: Standard directories present
**Previous Missions**: None
**Backlog Items**: Empty
**Characteristics**: Fresh start, no existing work
```

#### Active Mission State
```markdown
**Scenario**: active-mission
**Mission File**: Exists with status "active"
**Current Mission**: Track 2 WET mission in progress
**Project Structure**: Some implementation files present
**Characteristics**: Work in progress, mission conflict possible
```

#### Completed Mission State
```markdown
**Scenario**: completed-mission-history
**Mission File**: Does not exist (previous mission completed)
**Completed Missions**: 15 previous missions in archive
**Project Structure**: Established codebase
**Characteristics**: Mature project, established patterns
```

### Clarification Scenarios

#### Security Clarification
```markdown
**Scenario**: security-clarification-needed
**Initial Request**: "add user authentication"
**Clarification Questions**: 
  - "Which authentication method (JWT, OAuth, session-based)?"
  - "What security compliance requirements apply?"
  - "Should it integrate with existing user management?"
**Expected Responses**: Detailed security requirements
**Track Escalation**: Likely Track 2 → Track 3
```

#### Performance Clarification  
```markdown
**Scenario**: performance-clarification-needed
**Initial Request**: "optimize the search feature"
**Clarification Questions**:
  - "What are the current performance issues?"
  - "What response time targets are required?"
  - "Should caching or indexing be implemented?"
**Expected Responses**: Specific performance requirements
**Track Escalation**: Likely Track 2 → Track 3
```

## USAGE GUIDELINES

### Selecting Mock Scenarios
- Choose scenarios that match the test's complexity level
- Use realistic combinations of user input and AI analysis
- Include appropriate system state for test context
- Consider domain multipliers and their impact

### Customizing Scenarios
- Modify mock data to fit specific test requirements
- Combine multiple scenarios for complex test cases
- Adjust complexity levels to test edge cases
- Add domain-specific details as needed

### Scenario Validation
- Ensure mock scenarios reflect realistic user behavior
- Verify AI analysis responses are plausible
- Check that system states are consistent
- Validate domain multiplier detection logic

This library provides consistent, reusable mock scenarios that enable comprehensive testing of prompt logic through natural language specifications rather than programming constructs.
