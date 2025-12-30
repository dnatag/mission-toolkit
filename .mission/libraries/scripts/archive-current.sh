#!/bin/bash

# Archive existing mission to paused directory
if [ -f ".mission/mission.md" ]; then
  mkdir -p .mission/paused
  TIMESTAMP=$(date +%Y-%m-%d-%H-%M)
  mv .mission/mission.md ".mission/paused/${TIMESTAMP}-mission.md"
  echo "✅ Current mission archived to .mission/paused/${TIMESTAMP}-mission.md"
fi