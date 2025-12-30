#!/bin/bash

# Update mission status to active
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/status: planned/status: active/' .mission/mission.md
else
    sed -i 's/status: planned/status: active/' .mission/mission.md
fi
echo "✅ Mission status updated to active"