---
template: error-mission-exists
version: 1.0.0
---

# ❌ ERROR: Active Mission Already Exists

**Cannot create new mission** - there is already an active mission in progress.

## Current Mission

```
{{MISSION_CONTENT}}
```

## What would you like to do?

**1. Continue working on current mission:**
```bash
/m.apply
```

**2. Work on something else (pause current):**
```bash
m mission pause && /m.plan "{{NEW_INTENT}}"
```

**3. Current mission is done:**
```bash
/m.complete
```

**4. Abandon current mission** (⚠️ discards progress):
```bash
m mission archive --force && /m.plan "{{NEW_INTENT}}"
```

---
💡 **Tip**: Run `m dashboard` to view full mission status
