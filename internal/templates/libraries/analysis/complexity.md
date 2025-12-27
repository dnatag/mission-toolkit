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

## Assessment Rules
- Take the higher track from either file count or line count metrics
- Test files don't count toward complexity
- Apply domain multipliers after base assessment
- Final track = min(base_track + multipliers, 3)

## Complexity Scoring
- **Base Score**: File count (0=1pt, 1-5=2pts, 6-9=3pts, 10+=4pts)
- **Multipliers**: +1 for each domain risk
- **Final Track**: min(score, 3)

## Examples
- "Add missing item to array" (0 new files, 1 line) = Track 1
- "Add user CRUD" (3 new files, ~50 lines) = Track 2
- "Add payment processing" (2 new files, ~30 lines + payment API) = Track 3
- "Add database logging" (2 new files, ~40 lines + database) = Track 2