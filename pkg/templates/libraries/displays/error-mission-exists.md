---
template: error-mission-exists
version: 1.0.0
---

# ❌ Active Mission Already Exists

**Current Mission:** {{MISSION_CONTENT}}

**Choose an action:**

1. Continue: `/m.apply`
2. Pause & start new: `m mission pause && /m.plan "{{NEW_INTENT}}"`
3. Complete: `/m.complete`
4. Abandon (⚠️ discards): `m mission archive --force && /m.plan "{{NEW_INTENT}}"`
