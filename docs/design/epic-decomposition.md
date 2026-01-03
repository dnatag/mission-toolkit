# Design: Epic Decomposition Pattern (Track 4 Handler)

**Status**: 🚧 Future Enhancement (design complete, not implemented)  
**Last Updated**: 2025-01-15  
**Implementation**: Not yet started

**Note**: This document describes a planned feature with complete design but no code implementation. It will be implemented in a future phase of the Mission Toolkit.

## 1. Problem Statement
1. **Epic Intents**: Users often have large feature ideas (Track 4) that need decomposition into atomic missions.
2. **Manual Decomposition**: Current workflow requires manual breakdown of epics into backlog items.
3. **Lack of Structure**: No systematic approach to decompose features while maintaining context and dependencies.
4. **Missing Interview Pattern**: No interactive workflow to probe deeper into requirements and constraints.

## 2. Proposed Solution
1. **Two Decomposition Modes**: Automatic (default) or Interview (opt-in) - user chooses based on clarity.
2. **Structured Backlog Storage**: Store decomposed missions in `.mission/backlog.md` with metadata.
3. **AI-Only Pattern**: No CLI commands needed - AI directly manages backlog.md.
4. **Integration with m.plan**: Detect Track 4 and offer decomposition options.
5. **Interview Constraints**: Maximum 5 questions, user can skip, AI works with incomplete answers.

## 3. Decomposition Modes

### Mode 1: Automatic (Default)
**When to use**: Well-defined epics where requirements are clear from the intent.

**Process**:
1. AI analyzes intent and codebase context
2. Identifies logical sub-missions
3. Generates backlog immediately
4. Presents to user for confirmation
5. User can accept, adjust, or switch to Interview mode

**Benefits**: Fastest, good for standard patterns (CRUD, auth, etc.), minimal user input.

**Constraints**: AI must decompose into 3-7 sub-missions (not too granular, not too coarse).

### Mode 2: Interview (Opt-in)
**When to use**: Ambiguous epics with unclear requirements or many unknowns.

**Process**:
1. AI asks maximum 5 targeted questions
2. Each question must be justified ("I need X to determine Y")
3. User can skip questions ("I don't know yet" or "Skip")
4. AI works with incomplete answers and makes reasonable assumptions
5. AI documents all responses and assumptions
6. Generates backlog with context

**Benefits**: Uncovers hidden requirements without overwhelming user.

**Constraints**:
- Maximum 5 questions total (no follow-ups unless critical)
- Questions must be non-obvious and actionable
- User can exit interview anytime ("Just decompose it")
- AI must work with partial information

## 4. Integration with m.plan

When Track 4 detected, AI presents options:

```markdown
⚠️ TRACK 4 DETECTED: Epic Scope

Intent: "{{USER_INTENT}}"

This intent is too large for a single mission (estimated {{FILE_COUNT}} files).
I'll decompose it into smaller missions.

Choose approach:

1. **Automatic** (recommended) - I'll analyze and decompose immediately
   Best for: Clear requirements, standard patterns
   
2. **Interview** - I'll ask up to 5 questions to clarify requirements
   Best for: Ambiguous scope, many unknowns

Which mode? (1/2 or just press Enter for Automatic)
```

## 5. Mode 2: Interview Workflow

### Interview Constraints (CRITICAL)

1. **Maximum 5 Questions**: AI must prioritize the most impactful questions
2. **No Follow-ups**: Each question stands alone (no "can you elaborate?")
3. **Skippable**: User can respond "skip" or "don't know" to any question
4. **Justified**: Each question must explain why it's needed
5. **Exit Anytime**: User can say "just decompose it" to skip remaining questions

### Question Selection Strategy

AI selects up to 5 questions from these categories (prioritize by impact):

**1. Scope Boundaries** (Highest Priority)
- What's explicitly excluded from this feature?
- What's MVP vs nice-to-have?

**2. Technical Constraints**
- Any existing systems this must integrate with?
- Any performance or security requirements?

**3. User Flow**
- What's the primary user action or workflow?
- Any edge cases or error scenarios to handle?

**4. Data Model**
- What data needs to be stored or transformed?
- Any migrations or schema changes needed?

**5. Tradeoffs**
- Speed vs robustness - which matters more?
- What can be deferred to later?

### Interview Principles

1. **Non-Obvious**: Only ask what can't be inferred from intent or codebase
2. **Actionable**: Question must directly impact decomposition
3. **Justified**: Explain why answer matters ("I need X to determine Y")
4. **Graceful Degradation**: Work with incomplete answers
5. **Respect Time**: 5 questions max, no exceptions

### Interview Process

1. AI selects 3-5 most impactful questions
2. AI asks one question at a time with justification
3. User answers, skips, or exits
4. AI documents responses and assumptions for skipped questions
5. After 5 questions (or user exit), AI generates backlog
6. AI notes assumptions made for skipped questions

## 6. Mode 1: Automatic Workflow

### Analysis Process

1. **AI**: Analyze user intent and extract:
   - Core feature description
   - Implied scope and boundaries
   - Technical domain (auth, API, UI, etc.)

2. **AI**: Examine codebase context:
   - Existing file structure
   - Similar patterns or features
   - Dependencies and integration points

3. **AI**: Identify logical sub-missions:
   - Core functionality (must-have)
   - Supporting features (should-have)
   - Testing and validation
   - Documentation

4. **AI**: Determine execution order:
   - Foundation first (data models, schemas)
   - Core logic second (business logic)
   - Integration third (APIs, UI)
   - Testing and polish last

5. **AI**: Generate backlog items with estimates

6. **AI**: Present to user for review and adjustment

### Automatic Decomposition Strategy

**Analysis Steps**:
1. Identify feature type (CRUD, Auth, Integration, Refactor, etc.)
2. Extract core components from intent
3. Determine logical execution order (foundation → logic → integration)
4. Estimate complexity per sub-mission (aim for Track 2)
5. Generate 3-7 sub-missions with clear boundaries

**Decomposition Principles**:
- Each sub-mission should be independently testable
- Minimize dependencies between sub-missions
- Foundation before features (models before endpoints)
- Core logic before polish (functionality before optimization)
- Aim for Track 2 complexity per sub-mission

**Common Patterns** (use as guidance, not rigid templates):
- **CRUD**: Model → Create → Read → Update → Delete → Validation
- **Auth**: User model → Hashing → Token gen → Token validation → Endpoints
- **Integration**: Client setup → Core integration → Error handling → Testing
- **Refactor**: Extract logic → Update callers → Add tests → Remove old code

## 7. User Review & Adjustment

After decomposition (Automatic or Interview), AI presents breakdown:

```markdown
✅ EPIC DECOMPOSED: {{EPIC_NAME}}

📋 CREATED {{ITEM_COUNT}} SUB-MISSIONS:

1. BACKLOG-001: {{INTENT_1}} (Track {{TRACK_1}}, ~{{FILE_COUNT_1}} files)
2. BACKLOG-002: {{INTENT_2}} (Track {{TRACK_2}}, ~{{FILE_COUNT_2}} files)
...

🎯 SUGGESTED EXECUTION ORDER:
1. BACKLOG-001 (foundation)
2. BACKLOG-002 (core logic)
3. BACKLOG-003 (integration)

Does this breakdown make sense? (Y/N/Adjust)
```

**User Options**:
- **Y (Yes)**: Accept decomposition, add to backlog.md
- **N (No)**: Reject and try Interview mode (if was Automatic) or re-decompose
- **Adjust**: User describes changes ("Merge 2 and 3", "Split 1 into two parts")

**AI Actions**:
- If Y: Write to `.mission/backlog.md` and show next steps
- If N: Offer Interview mode or ask for guidance
- If Adjust: Apply changes and re-present for confirmation

## 8. Backlog Item Format

**Storage in `.mission/backlog.md`:**
```markdown
## Backlog

### BACKLOG-001: [Epic Name] - Sub-mission Intent
- **Epic**: Parent epic description
- **Estimated Track**: 2
- **Estimated Files**: `path/to/file1.go`, `path/to/file2.go`
- **Dependencies**: BACKLOG-000 (if any)
- **Priority**: high
- **Status**: pending
- **Notes**: Additional context or constraints
- **Created**: 2025-01-15T10:30:00Z

---
```

**Status Values**:
- `pending` - Not started
- `in-progress` - Mission created from this item
- `completed` - Mission completed
- `cancelled` - No longer needed

## 9. Implementation in m.plan

### Update to `m.plan.md` Prompt

Add Track 4 detection and decomposition offer:

```markdown
## Step 5: Complexity Analysis

Run: m plan analyze --file plan.json

If Track 4 detected:
  1. Display decomposition options (Automatic/Interview)
  2. Default to Automatic if user presses Enter
  3. Execute chosen mode
  4. Present decomposition for user review (Y/N/Adjust)
  5. If approved, write to backlog.md
  6. STOP (do not create mission.md)
  
If Track 1-3:
  Continue to validation
```

### Decomposition Output Template

```markdown
✅ EPIC DECOMPOSED: {{EPIC_NAME}}

📋 CREATED {{ITEM_COUNT}} SUB-MISSIONS:

1. BACKLOG-001: {{INTENT_1}} (Track {{TRACK_1}}, ~{{FILE_COUNT_1}} files)
2. BACKLOG-002: {{INTENT_2}} (Track {{TRACK_2}}, ~{{FILE_COUNT_2}} files)
...

{{#IF_INTERVIEW_MODE}}
📝 ASSUMPTIONS MADE:
- {{ASSUMPTION_1}} (Question skipped: {{SKIPPED_QUESTION_1}})
- {{ASSUMPTION_2}} (Question skipped: {{SKIPPED_QUESTION_2}})
{{/IF_INTERVIEW_MODE}}

🎯 SUGGESTED EXECUTION ORDER:
1. BACKLOG-001 (foundation) - Start here
2. BACKLOG-002 (core logic) - Depends on #1
3. BACKLOG-003 (integration) - Depends on #2

🚀 NEXT STEPS:
• Review backlog: cat .mission/backlog.md
• Start first mission: /m.plan "{{FIRST_INTENT}}"
• Track progress: grep "Status:" .mission/backlog.md
```

## 10. Backlog Lifecycle Integration

### When Creating Mission from Backlog

**User runs**: `/m.plan "Add JWT token generation"`

**AI checks**: Does this match a backlog item?
- Search `.mission/backlog.md` for matching intent
- If found, include epic context in mission
- Update backlog item status to `in-progress`
- Add BACKLOG_ID to mission.md metadata

**Mission.md includes**:
```yaml
BACKLOG_ID: BACKLOG-001
EPIC: Build user authentication system
```

### When Completing Mission from Backlog

**During `m.complete`**:
- AI checks if mission has BACKLOG_ID metadata
- If yes, update backlog item status to `completed`
- Add completion timestamp
- Check if all items in epic are completed
- If epic complete, display summary message
- AI checks if mission has BACKLOG_ID
- If yes, update backlog item status to `completed`
- Check if all items in epic are completed
- If yes, display epic completion message

### Backlog Status Tracking

**User can check progress**:
```bash
# View all backlog items
cat .mission/backlog.md

# Count pending items
grep "Status: pending" .mission/backlog.md | wc -l

# View specific epic progress
grep -A 8 "Epic: Build user authentication" .mission/backlog.md
```

## 11. Example Flows

### Example 1: Interview Mode

```
User: /m.plan "Build user authentication system"

AI: ⚠️ TRACK 4 DETECTED: Epic Scope
    Choose: (1) Automatic (2) Interview [Enter = 1]

User: 2

AI: I'll ask up to 5 questions to clarify requirements.
    
    Q1: What's explicitly excluded from this auth system?
    (I need this to set clear boundaries for decomposition)
    
User: No OAuth, no 2FA, just basic email/password

AI: Q2: Any existing user data or systems this must integrate with?
    (I need this to determine migration/integration sub-missions)

User: Fresh start, no existing users

AI: Q3: What's the primary security requirement - speed or robustness?
    (I need this to determine if we need Track 3 security missions)

User: Robustness - this is production

AI: Q4: Any specific token requirements (expiration, refresh, etc.)?
    (I need this to scope the JWT implementation)

User: Skip - just standard JWT

AI: Q5: What can be deferred to later (password reset, email verification, etc.)?
    (I need this to focus on MVP)

User: Defer password reset and email verification

AI: ✅ EPIC DECOMPOSED: Build user authentication system
    
    📋 CREATED 5 SUB-MISSIONS:
    1. BACKLOG-001: Create user model and database schema (Track 2, ~2 files)
    2. BACKLOG-002: Add password hashing with bcrypt (Track 3, ~3 files)
    3. BACKLOG-003: Implement JWT token generation (Track 3, ~2 files)
    4. BACKLOG-004: Add token validation middleware (Track 3, ~2 files)
    5. BACKLOG-005: Create login endpoint with validation (Track 2, ~2 files)
    
    📝 ASSUMPTIONS MADE:
    - Standard JWT format (Question skipped: Q4)
    
    Does this breakdown make sense? (Y/N/Adjust)

User: Y

AI: ✅ Added to .mission/backlog.md
    🚀 Start with: /m.plan "Create user model and database schema"
```

### Example 2: Automatic Mode

```
User: /m.plan "Add CRUD endpoints for products"

AI: ⚠️ TRACK 4 DETECTED: Epic Scope
    Choose: (1) Automatic (2) Interview [Enter = 1]

User: [presses Enter]

AI: Analyzing intent and codebase...
    
    ✅ EPIC DECOMPOSED: Add CRUD endpoints for products
    
    📋 CREATED 6 SUB-MISSIONS:
    1. BACKLOG-006: Create product model and schema (Track 2, ~2 files)
    2. BACKLOG-007: Add create product endpoint (Track 2, ~2 files)
    3. BACKLOG-008: Add list products endpoint (Track 2, ~2 files)
    4. BACKLOG-009: Add update product endpoint (Track 2, ~2 files)
    5. BACKLOG-010: Add delete product endpoint (Track 2, ~2 files)
    6. BACKLOG-011: Add product validation and tests (Track 2, ~3 files)
    
    🎯 SUGGESTED EXECUTION ORDER:
    1. BACKLOG-006 (foundation) - Start here
    2. BACKLOG-007-010 (endpoints) - Any order
    3. BACKLOG-011 (validation) - Last
    
    Does this breakdown make sense? (Y/N/Adjust)

User: Adjust - merge validation into each endpoint

AI: ✅ ADJUSTED DECOMPOSITION:
    
    📋 CREATED 5 SUB-MISSIONS:
    1. BACKLOG-006: Create product model and schema (Track 2, ~2 files)
    2. BACKLOG-007: Add create product endpoint with validation (Track 2, ~3 files)
    3. BACKLOG-008: Add list products endpoint with filtering (Track 2, ~2 files)
    4. BACKLOG-009: Add update product endpoint with validation (Track 2, ~3 files)
    5. BACKLOG-010: Add delete product endpoint (Track 2, ~2 files)
    
    Does this work? (Y/N)

User: Y

AI: ✅ Added to .mission/backlog.md
    🚀 Start with: /m.plan "Create product model and schema"
```

### Example 3: Interview Exit

```
User: /m.plan "Refactor authentication system"

AI: ⚠️ TRACK 4 DETECTED: Epic Scope
    Choose: (1) Automatic (2) Interview [Enter = 1]

User: 2

AI: I'll ask up to 5 questions to clarify requirements.
    
    Q1: What's the main problem with the current auth system?
    (I need this to focus the refactoring scope)
    
User: Just decompose it - I'll clarify as we go

AI: Understood. Analyzing current auth system...
    
    ✅ EPIC DECOMPOSED: Refactor authentication system
    
    📋 CREATED 4 SUB-MISSIONS:
    1. BACKLOG-012: Extract auth logic to separate package (Track 2, ~4 files)
    2. BACKLOG-013: Update all auth callers to use new package (Track 2, ~5 files)
    3. BACKLOG-014: Add comprehensive auth tests (Track 2, ~3 files)
    4. BACKLOG-015: Remove deprecated auth code (Track 2, ~3 files)
    
    📝 ASSUMPTIONS MADE:
    - Standard refactoring pattern (User exited interview early)
    
    Does this breakdown make sense? (Y/N/Adjust)

User: Y

AI: ✅ Added to .mission/backlog.md
    🚀 Start with: /m.plan "Extract auth logic to separate package"
```

## 12. Benefits

### Simplicity
- No CLI commands needed
- AI directly manages backlog.md
- Two clear modes (Automatic default, Interview opt-in)

### User Experience
- Interview limited to 5 questions max
- User can skip questions or exit anytime
- AI works with incomplete information
- User reviews and adjusts decomposition before committing

### Flexibility
- Automatic mode for clear intents
- Interview mode for ambiguous intents
- User can adjust AI's decomposition
- Can switch modes if first attempt unsatisfactory

### Integration
- Seamless with existing m.plan workflow
- Backlog items tracked through completion
- Epic context preserved in missions
- No breaking changes to current system

### Maintainability
- No additional CLI code
- Pattern documented in prompts
- Works with any AI assistant
- Clear constraints prevent scope creep

## 13. Documentation Location

Add to mission-toolkit:

1. **README.md**: Brief mention of Track 4 decomposition
2. **docs/patterns/epic-decomposition.md**: Full pattern documentation
3. **Update m.plan.md**: Add Track 4 detection and mode selection
4. **Update governance.md**: Add backlog lifecycle rules
