# Update mission status to active
sed -i 's/status: planned/status: active/' .mission/mission.md
echo "✅ Mission status updated to active"