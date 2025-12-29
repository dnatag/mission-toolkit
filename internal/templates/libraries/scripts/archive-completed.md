# ARCHIVE COMPLETED MISSION SCRIPT

## Purpose
Archive completed mission with metrics to .mission/completed/ directory.

## Operations
1. **Create Directory**: Ensure .mission/completed/ exists
2. **Archive Mission**: Move mission.md with mission ID prefix
3. **Create Metrics**: Generate metrics file using completion.md template

## Template Variables Required
- {{MISSION_ID}} = Mission ID from mission.md id field
- {{METRICS_CONTENT}} = Content from completion.md template

## Output Files
- `.mission/completed/{{MISSION_ID}}-mission.md`
- `.mission/completed/{{MISSION_ID}}-metrics.md`

## Script Template
```bash
# Use mission ID for consistent naming
MISSION_ID="{{MISSION_ID}}"
mkdir -p .mission/completed

# Archive mission file
mv .mission/mission.md ".mission/completed/${MISSION_ID}-mission.md"

# Create metrics file
cat > ".mission/completed/${MISSION_ID}-metrics.md" << 'EOF'
{{METRICS_CONTENT}}
EOF

echo "✅ Mission archived: .mission/completed/${MISSION_ID}-mission.md"
echo "✅ Metrics saved: .mission/completed/${MISSION_ID}-metrics.md"
```