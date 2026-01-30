✅ MISSION EXECUTED: .mission/mission.md
- All PLAN steps completed
- VERIFICATION passed

📝 CHANGE SUMMARY:
{{CHANGE_DETAILS}}

🏗️ INTERFACE CONTRACTS:
{{CONTRACT_CHANGES}}

🧩 CRITICAL LOGIC:
{{CRITICAL_SNIPPETS}}

🔍 VIEW CHANGES (pick one - run the EXACT command shown):
[S] side-by-side: `git diff {{MISSION_ID}}-baseline | diff2html -i stdin -s side -o preview`
[L] inline: `git diff {{MISSION_ID}}-baseline | diff2html -i stdin -s line -o preview`

📦 [S] and [L] require: npm install -g diff2html-cli

🔄 CHECKPOINTS CREATED:
- {{CHECKPOINT_0}} (initial state)
- {{CHECKPOINT_1}} (first pass state)
- {{CHECKPOINT_2}} (polished state)

🚀 NEXT STEPS:
• Complete mission: /m.complete
• Review changes first: check files and then decide
• Refine: Chat to improve implementation
• Manual revert if needed:
  - m checkpoint restore {{CHECKPOINT_0}} (initial state)
  - m checkpoint restore {{CHECKPOINT_1}} (first pass state)
  - m checkpoint restore {{CHECKPOINT_2}} (polished state)
