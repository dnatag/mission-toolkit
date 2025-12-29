# REFRESH METRICS SCRIPT

## Purpose
Update `.mission/metrics.md` aggregate statistics and insights based on completed missions in `.mission/completed/` directory.

## Operations
1. **Scan Completed Missions**: Read all `*-mission.md` files in `.mission/completed/`
2. **Calculate Aggregates**: Count totals, success rates, track distribution, average durations
3. **Update Summary**: Replace AGGREGATE STATISTICS and TRACK DISTRIBUTION sections using template `.mission/libraries/metrics/aggregate.md`
4. **Update Recent**: Replace RECENT COMPLETIONS section with last 5 missions
5. **Refresh Insights**: Replace insights sections using template `.mission/libraries/metrics/insights.md`

## Input Data Sources
- **Completed Missions**: `.mission/completed/*-mission.md` files
- **Current Metrics**: `.mission/metrics.md` (preserve PROCESS INSIGHTS)
- **Templates**: `libraries/metrics/aggregate.md`, `libraries/metrics/insights.md`

## Output
- **Updated**: `.mission/metrics.md` with refreshed statistics and insights
- **Preserved**: Existing PROCESS INSIGHTS section content

## Error Handling
- **Missing completed/**: Skip if no completed missions exist
- **Template Missing**: Report error and stop
- **Parse Errors**: Log warning and continue with available data

## Variables for Templates
**Aggregate Template Variables:**
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

**Insights Template Variables:**
- Based on analysis of recent completed missions and trends
- Populate with actual insights from mission data analysis