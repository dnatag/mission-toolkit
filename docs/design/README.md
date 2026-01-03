# Mission Toolkit Design Documents

Technical design specifications for developers working on the Mission Toolkit.

## Core Command Designs

### [m-plan.md](m-plan.md)
CLI-driven planning workflow with complexity analysis and validation.

**Status**: ✅ Implemented

### [m-clarify.md](m-clarify.md)
Clarification workflow for refining ambiguous intents.

**Status**: ✅ Implemented (prompt-based)

### [m-apply.md](m-apply.md)
Mission execution with commit message generation and verification.

**Status**: 🚧 Future Enhancement

### [m-complete.md](m-complete.md)
Git integration and mission archival with metrics collection.

**Status**: 🚧 Future Enhancement

## Supporting Designs

### [commit-messages.md](commit-messages.md)
Conventional commit message lifecycle (generation → storage → consumption).

**Status**: 🚧 Future Enhancement

### [template-system.md](template-system.md)
Template architecture rationale and AI-agnostic design.

**Status**: ✅ Implemented

### [epic-decomposition.md](epic-decomposition.md)
Track 4 epic handler with automatic and interview decomposition modes.

**Status**: 🚧 Future Enhancement (design complete, not implemented)

## Meta Documentation

### [REVIEW.md](REVIEW.md)
Comprehensive review of all design docs with completeness assessment.

## Status Legend

- ✅ **Implemented** - Design complete and code exists
- 🚧 **Future Enhancement** - Design complete, awaiting implementation

## Related Documentation

- [Workflows](../workflows/) - User-facing workflow guides
- [Governance](../../.mission/governance.md) - Core principles and rules
- [Main README](../../README.md) - Project overview

## For Contributors

When creating new design documents:

1. Use clear, descriptive names (e.g., `feature-name.md`)
2. Include status badge at the top
3. Follow structure: Problem → Solution → Architecture → Implementation
4. Link to related code when implemented
5. Update this README with the new document
