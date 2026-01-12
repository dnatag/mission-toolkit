---
template: checkpoint-failure-unstaged
version: 1.0.0
---

# ❌ Checkpoint Failed: Unstaged Changes

**Files:** {{UNSTAGED_FILES}}

**Choose an action:**

1. Stash & retry: `git stash && /m.apply`
2. Commit & retry: `git add . && git commit -m "WIP" && /m.apply`
3. Discard & retry: `git checkout . && /m.apply`

💡 After mission: `git stash pop` (if stashed)