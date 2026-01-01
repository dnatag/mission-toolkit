# DOMAIN ANALYSIS TEMPLATE

## Purpose
Identify specific technical or business domains that apply to the user's request. This information is used by the CLI to calculate complexity and risk.

## Valid Domains
You must select from this **strict list** of domains. Do not invent new ones.

### 1. Security (`security`)
- **Triggers**: Authentication, authorization, cryptography, PII handling, secrets management, input sanitization.
- **Example**: "Add login", "Encrypt password", "Fix SQL injection".

### 2. Performance (`performance`)
- **Triggers**: Latency requirements, throughput optimization, memory management, caching, database indexing, concurrency.
- **Example**: "Speed up API", "Reduce memory usage", "Add Redis cache".

### 3. Complex Algorithms (`complex-algo`)
- **Triggers**: Mathematical models, AI/ML, recursion, state machines, graph algorithms, custom data structures.
- **Example**: "Implement recommendation engine", "Pathfinding logic", "Parser implementation".

### 4. High Risk (`high-risk`)
- **Triggers**: Financial transactions, payments, data deletion (bulk), critical infrastructure, public API changes.
- **Example**: "Process refund", "Delete user account", "Change API signature".

### 5. Cross-Cutting (`cross-cutting`)
- **Triggers**: Changes affecting multiple distinct modules, logging infrastructure, configuration management, error handling strategies.
- **Example**: "Update logging format everywhere", "Refactor config loading", "Global error handler".

### 6. Real-Time (`real-time`)
- **Triggers**: WebSockets, streaming, event-driven architecture, polling.
- **Example**: "Live chat", "Stock ticker", "Notification stream".

### 7. Compliance (`compliance`)
- **Triggers**: GDPR, audit logs, legal requirements, accessibility (WCAG).
- **Example**: "Add consent banner", "Export user data", "Audit trail".

## Decision Logic
- **Default**: If none apply, the domain list is empty `[]`.
- **Multiple**: Select ALL that apply (e.g., `["security", "high-risk"]`).
- **Threshold**: If unsure, err on the side of caution and include the domain.

## Output Format
List of selected domain strings.

**Example:**
`["security", "performance"]`
