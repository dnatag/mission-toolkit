# Generate Commit Message

**Purpose**: Generate conventional commit message after polish pass

**Usage**: Called in Step 4 of m.apply after polish succeeds or rolls back

**Template**:
```
{{COMMIT_TYPE}}({{COMMIT_SCOPE}}): {{COMMIT_TITLE}}

{{COMMIT_DESCRIPTION}}

Mission-ID: {{MISSION_ID}}
Track: {{TRACK}}
Type: {{MISSION_TYPE}}{{POLISH_FOOTER}}
```

**Variables**:
- {{COMMIT_TYPE}} = Commit type (feat, refactor, fix, docs, test, chore)
  - WET missions → `feat`
  - DRY missions → `refactor`
  - Override if clearly fix/docs/test/chore
- {{COMMIT_SCOPE}} = Primary module/file scope (e.g., auth, api, cli)
  - Extract from first file in SCOPE or primary component
- {{COMMIT_TITLE}} = Imperative mood, max 72 chars, capitalize first letter, no period
- {{COMMIT_DESCRIPTION}} = Multi-line description explaining what and why (not how)
  - Wrap at 72 characters per line
  - Reflect ALL changes made (implementation + polish or implementation only)
- {{MISSION_ID}} = Mission identifier from mission.md
- {{TRACK}} = Complexity track (1-4)
- {{MISSION_TYPE}} = WET or DRY
- {{POLISH_FOOTER}} = Optional footer line
  - If polish skipped: `\nPolish-Skipped: checkpoint-creation-failed`
  - If polish succeeded or rolled back: empty string

**Rules**:
1. Generate AFTER polish pass completes (Step 4)
2. Regenerate after ANY user-requested code changes post-/m.apply
3. Description must reflect ALL changes, not just latest
4. If polish rolled back, description reflects only first pass implementation
5. If polish succeeded, description reflects implementation + polish improvements
