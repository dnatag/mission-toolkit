# REFRESH METRICS SCRIPT

## Purpose
Update `.mission/metrics.md` with complete refresh of all sections based on completed missions.

## Execution Steps
1. **Scan Completed Missions**: Read all `*-mission.md` files in `.mission/completed/` directory
2. **Calculate Variables**: Generate all template variables from completed mission data
3. **Update Aggregate Sections**: Use template `.mission/libraries/metrics/aggregate.md` to replace:
   - AGGREGATE STATISTICS
   - TRACK DISTRIBUTION
   - WET-DRY EVOLUTION
   - QUALITY METRICS
   - RECENT COMPLETIONS
4. **Update Insights Sections**: Use template `.mission/libraries/metrics/insights.md` to replace:
   - PROCESS INSIGHTS
   - TECHNICAL LEARNINGS
   - RECOMMENDATIONS

## Template Variables to Calculate
**For aggregate.md:**
- {{TOTAL_MISSIONS}} = Count of completed missions
- {{SUCCESS_RATE}} = Percentage of successful missions
- {{AVG_DURATION}} = Average mission duration in minutes
- {{TRACK_1_COUNT}} = Count of Track 1 missions
- {{TRACK_2_COUNT}} = Count of Track 2 missions  
- {{TRACK_3_COUNT}} = Count of Track 3 missions
- {{TRACK_4_COUNT}} = Count of Track 4 missions
- {{WET_TO_DRY_RATIO}} = Ratio of WET to DRY missions
- {{COMMON_PATTERNS}} = Most frequent patterns identified
- {{BOTTLENECKS}} = Most common failure points

**For insights.md:**
- Analyze recent completed missions for workflow efficiency
- Identify failure patterns and technical learnings
- Generate actionable recommendations

## Error Handling
- **Missing completed/**: Skip if no completed missions exist
- **Template Missing**: Report error and stop
- **Parse Errors**: Log warning and continue with available data