# Archive completed mission with metrics
TIMESTAMP=$(date +%Y-%m-%d-%H-%M)
mkdir -p .mission/completed

# Move mission file
mv .mission/mission.md ".mission/completed/${TIMESTAMP}-mission.md"

# Create metrics file
cat > ".mission/completed/${TIMESTAMP}-metrics.md" << 'EOF'
{{METRICS_CONTENT}}
EOF

echo "✅ Mission archived: .mission/completed/${TIMESTAMP}-mission.md"
echo "✅ Metrics saved: .mission/completed/${TIMESTAMP}-metrics.md"