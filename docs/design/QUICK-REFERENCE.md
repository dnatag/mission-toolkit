# Quick Reference: m plan & m mission Commands

**AI Rule**: Never edit JSON files. Use these commands instead.

---

## Step 1: Clarification Mission

```bash
# Create plan.json with questions
m plan init --intent "Add user authentication" \
  --question "Which authentication method?" \
  --question "Should we support social login?"

# Generate clarification mission.md
m mission create --type clarification --file .mission/plan.json
```

---

## Step 3: Draft Spec

```bash
# Create plan.json with scope
m plan init --intent "Add JWT authentication" \
  --type WET \
  --scope internal/auth/jwt.go \
  --scope internal/auth/jwt_test.go \
  --domain security
```

---

## Step 4: Complexity Analysis

```bash
# Analyze and auto-update track
m plan analyze --file .mission/plan.json --update
```

**Output**: JSON with track, next_step, warnings

---

## Step 5: Validation

```bash
# Validate scope and verification
m plan validate --file .mission/plan.json
```

**Output**: JSON with valid status, errors, warnings

---

## Step 6: Final Mission

```bash
# Add plan steps and verification
m plan update --plan "Create JWT generation function" \
  --plan "Add validation middleware" \
  --plan "Write comprehensive tests" \
  --plan "Note: Allow duplication (WET principle)" \
  --verification "go test ./internal/auth/..."

# Generate final mission.md
m mission create --type final --file .mission/plan.json
```

---

## Command Summary

| Command | Purpose | Key Flags |
|---------|---------|-----------|
| `m plan init` | Create plan.json | `--intent`, `--type`, `--scope`, `--domain`, `--question` |
| `m plan update` | Update plan.json | `--plan`, `--verification` |
| `m plan analyze` | Determine track | `--file`, `--update` |
| `m plan validate` | Validate plan | `--file` |
| `m mission create` | Generate mission.md | `--type`, `--file` |

---

## Repeatable Flags

These flags can be used multiple times:
- `--scope`: Add multiple files
- `--domain`: Add multiple domains
- `--question`: Add multiple questions
- `--plan`: Add multiple plan steps

**Example**:
```bash
m plan init --intent "..." \
  --scope file1.go \
  --scope file2.go \
  --scope file3.go \
  --domain security \
  --domain performance
```

---

## plan.json Evolution

### After Step 1 (Clarification)
```json
{
  "intent": "...",
  "clarification_questions": ["Q1?", "Q2?"]
}
```

### After Step 3 (Draft)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"]
}
```

### After Step 4 (Analysis)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"],
  "track": "TRACK 2"
}
```

### After Step 6 (Finalize)
```json
{
  "intent": "...",
  "type": "WET",
  "scope": ["file1.go", "file2.go"],
  "domain": ["security"],
  "track": "TRACK 2",
  "plan": ["Step 1", "Step 2"],
  "verification": "go test ./..."
}
```

---

## Error Handling

```bash
# Missing required flag
$ m plan init
Error: required flag "intent" not set

# Invalid mission type
$ m mission create --type invalid --file plan.json
Error: invalid type: invalid (use clarification or final)

# File not found
$ m plan analyze --file missing.json
Error: open missing.json: no such file or directory
```

---

## Tips

1. **Always use `--update`** in Step 4 to auto-add track
2. **Use backslash `\`** for multi-line commands
3. **Quote strings** with spaces: `--intent "Add user auth"`
4. **Check JSON output** from analyze/validate commands
5. **Inspect plan.json** at any point for debugging
