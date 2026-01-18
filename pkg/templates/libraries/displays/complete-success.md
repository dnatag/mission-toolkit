🎉 MISSION COMPLETED: {{MISSION_ID}}
- Duration: {{DURATION}}
- Final Commit: {{FINAL_COMMIT_HASH}}
- Track {{TRACK}} {{MISSION_TYPE}} mission

📁 ARCHIVED:
- Mission: .mission/completed/{{MISSION_ID}}-mission.md
- Execution log: .mission/completed/{{MISSION_ID}}-execution.log
- Plan: .mission/completed/{{MISSION_ID}}-plan.json

{{#UNSTAGED_FILES}}
⚠️  UNSTAGED FILES:
{{UNSTAGED_FILES}}

💡 OPTIONS:
• Amend commit: git add <files> && git commit --amend --no-edit
• Add to .gitignore: echo '<pattern>' >> .gitignore
{{/UNSTAGED_FILES}}

🚀 NEXT STEPS:
• Plan new mission: /m.plan
• Review backlog: Check .mission/backlog.md
