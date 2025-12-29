# COMPLEXITY ANALYSIS TEMPLATE

## Purpose
Analyze refined intent to determine track complexity and mission type.

## Mission Type Detection
1. **DRY Mission**: User explicitly requests refactoring existing duplication
   - Keywords: "Extract", "Refactor", "DRY", "consolidate", "eliminate duplication"
   - Examples: "Extract common validation logic", "Refactor duplicate API patterns"
2. **WET Mission**: All other feature development requests

## Base Complexity Assessment

| **TRACK 1 (Atomic)** | **TRACK 2 (Standard)** | **TRACK 3 (Robust)** | **TRACK 4 (Epic)** |
|---|---|---|---|
| Single line/function changes | Single feature implementation | Cross-cutting concerns | Multiple systems |
| 0 new implementation files | 1-5 new implementation files | 6-9 new implementation files | 10+ new implementation files |
| 1-5 lines of code changed | 10-100 lines of code changed | 100-500 lines of code changed | 500+ lines of code changed |
| "Fix typo", "Add to array", "Rename var" | "Add endpoint", "Create component" | "Add authentication", "Refactor for security" | "Build payment system", "Rewrite architecture" |

## Domain Multipliers (+1 track, max Track 3)
- **High-risk integrations**: Payment processing, financial transactions, authentication systems
- **Complex algorithms**: ML models, cryptography, real-time optimization
- **Performance-critical**: Sub-second response requirements, high-throughput systems
- **Regulatory/Security**: GDPR, SOX, PCI compliance, security-sensitive features
- **Cross-cutting concerns**: Logging, monitoring, caching, error handling that spans multiple modules

## Assessment Rules
- Take the higher track from either file count or line count metrics
- Test files don't count toward complexity
- Apply domain multipliers after base assessment
- Final track = min(base_track + multipliers, 3)

## CRITICAL: Assessment Guidelines
1. **Track 2 = Single features** (isolated functionality)
2. **Track 3 = Cross-cutting concerns** (spans multiple modules/infrastructure)
3. **Track 4 = Multiple systems** (10+ files OR requires new architecture/databases/services)
4. **"Spans multiple modules" + "infrastructure" = domain multiplier applies**
5. **When in doubt between Track 2/3: if it affects multiple parts of the system, it's Track 3**
6. **Track 4 triggers: new databases, external services, microservices, or architectural changes**

## Track 3 Scope Management
**Cross-cutting missions must remain atomic - enforce these limits:**
1. **Single concern only**: "Add logging" NOT "Add logging + monitoring + metrics"
2. **Minimal viable implementation**: Basic functionality first, enhancements later
3. **Max 6-9 files**: If exceeding 9 files, escalate to Track 4 for decomposition
4. **Clear boundaries**: Define exactly which modules are affected, exclude others
5. **Incremental approach**: "Add auth to API endpoints" NOT "Complete authentication system"

## Complexity Scoring
- **Base Score**: File count (0=1pt, 1-5=2pts, 6-9=3pts, 10+=4pts)
- **Multipliers**: +1 for each domain risk
- **Final Track**: min(score, 3)

## Edge Case Handling
- **File vs Line Conflict**: If file count suggests Track 2 but line count suggests Track 3, use higher track
- **Borderline Cases**: If assessment falls between tracks, consider implementation complexity over file count
- **Unknown Scope**: If file count cannot be estimated, default to Track 2 and note uncertainty

## Output Format
**Required Output:**
```
TRACK: [1-4]
CONFIDENCE: [High/Medium/Low]
REASONING: [Base complexity + domain multipliers + edge cases]
REASSESS_TRIGGERS: [Conditions that would change this assessment]
```

**Confidence Levels:**
- **High**: Clear file count, obvious domain classification
- **Medium**: Estimated file count, some domain ambiguity  
- **Low**: Vague requirements, multiple possible interpretations

**Reassessment Triggers:**
- New technical requirements discovered
- Scope significantly expanded or reduced
- Domain risk factors change (security, performance, compliance)
- File count estimates proven wrong during implementation

## Examples
- "Add missing item to array" (0 new files, 1 line) = Track 1
- "Add user CRUD" (3 new files, ~50 lines) = Track 2
- "Add payment processing" (2 new files, ~30 lines + payment API) = Track 3
- "Add database logging" (2 new files, ~40 lines + database) = Track 2