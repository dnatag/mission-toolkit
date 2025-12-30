#!/bin/bash

# Check if mission exists and has correct status
if [ ! -f ".mission/mission.md" ]; then
  echo "❌ No active mission found. Use /m.plan to create a new mission first."
  exit 1
fi

STATUS=$(grep "status:" .mission/mission.md | cut -d' ' -f2)
if [ "$STATUS" != "planned" ]; then
  echo "❌ Mission status is '$STATUS', expected 'planned'"
  exit 1
fi

echo "✅ Mission validation passed"