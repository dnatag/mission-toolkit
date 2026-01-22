# Debugging & Bugfix Strategy Recommendations

The current `m.plan` -> `m.apply` -> `m.complete` workflow is optimized for linear feature development. Debugging is exploratory and iterative. Below are strategies to adapt the workflow for debugging.

## 1. The "Probe" Phase (Enhanced `m.plan`)
**Problem:** You can't plan a fix if you don't know the root cause.
**Solution:** Add a "Diagnosis Mode" to `m.plan`.
- If intent is "Fix bug X", do not jump to implementation plan.
- Propose a **"Probe Mission"**:
    - **Scope**: Read logs, add debug prints, write reproduction test.
    - **Goal**: Identify root cause.
- **Action**: Detect keywords (bug, fix, crash) and switch to Diagnosis Template.

## 2. The "Reproduction First" Rule (Strict `m.apply`)
**Problem:** AI guesses fixes without verification.
**Solution:** Enforce Test-Driven Development (TDD) for bugs.
- **Loop**:
    1. **Fail**: Write test case that reproduces bug.
    2. **Fix**: Modify code.
    3. **Pass**: Verify test passes.
- **Action**: `m.apply` asks "Is this a bug fix?". If yes, require `repro_test` in plan.

## 3. The "Quick Fix" Track (Track 1 Optimization)
**Problem:** Small fixes feel too heavy.
**Solution:** Streamline Track 1.
- Allow `m.plan` to immediately execute trivial changes.
- Or create `@m.fix` for one-shot Plan+Apply.

## 4. Dedicated `@m.debug` Workflow
**Problem:** Debugging requires investigation, not just coding.
**Solution:** New `m.debug.md` prompt.
- **Role**: Investigator.
- **Workflow**: Analyze -> Hypothesize -> Verify (Probe) -> Iterate -> RCA.
- **Output**: Root Cause Analysis (RCA) feeding into `m.plan`.

## 5. Iterative `m.apply` (Looping)
**Problem:** `m.apply` rolls back on first failure, which is bad for debugging.
**Solution:** Allow looping/retrying.
- If verification fails: Analyze error -> Adjust code -> Retry (Limit 3).
- Only rollback if the *loop* fails.

## Recommendation
Start with **Option 5 (Iterative Apply)** and **Option 2 (Repro First)** within existing commands.

1. **Modify `m.plan`**: Explicitly ask for a **Reproduction Step** if intent is a bug fix.
2. **Modify `m.apply`**: Allow an **"Iterative Fix Loop"** (up to 3 retries) before rolling back.
