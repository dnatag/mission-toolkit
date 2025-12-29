# Initialize Execution Log

**Purpose**: Create `.mission/execution.log` if it doesn't exist

**Usage**: Check file existence before any logging operation

**Template**:
```
MISSION EXECUTION LOG
Started: {{TIMESTAMP}}

=== EXECUTION STEPS ===
```

**Variables**:
- {{TIMESTAMP}} = Current timestamp (YYYY-MM-DDTHH:MM:SS.sssZ)