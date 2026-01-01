# COMPLEXITY ANALYSIS REFERENCE

## Purpose
This document explains how the `m plan analyze` CLI tool calculates complexity. Use this to understand how your `scope` and `domain` inputs affect the mission Track.

## CLI Logic Overview

The CLI calculates the Track deterministically:
`Final Track = Base Track (File Count) + Domain Multipliers`

### 1. Base Track (File Count)
The CLI counts implementation files in your `scope` (ignoring tests/docs).

| **File Count** | **Base Track** | **Description** |
|---|---|---|
| 0 | **Track 1** | Atomic/Trivial change |
| 1-5 | **Track 2** | Standard feature |
| 6-9 | **Track 3** | Robust/Cross-cutting |
| 10+ | **Track 4** | Epic (Too large -> Decompose) |

### 2. Domain Multipliers (+1 Track)
The CLI adds +1 to the Track if you include any of these in the `domain` list:

- **`security`**: Auth, crypto, PII, permissions.
- **`performance`**: Latency, throughput, memory optimization.
- **`complex-algo`**: Math, AI, recursion, state machines.
- **`high-risk`**: Payments, data deletion, critical infrastructure.
- **`cross-cutting`**: Logging, config, error handling affecting multiple modules.
- **`real-time`**: WebSockets, streaming, concurrency.
- **`compliance`**: GDPR, audit logs, legal requirements.

*Note: Track is capped at 3 (unless file count forces Track 4).*

## Your Role
You do NOT calculate the Track. Your job is to:
1.  **Define Scope**: List all necessary files in `plan.json`.
2.  **Define Domain**: List all applicable risk domains in `plan.json`.
3.  **Run CLI**: Execute `m plan analyze` and respect its output.

## Handling Track 4
If the CLI returns `recommendation: "decompose"` (Track 4):
- **STOP**. Do not generate the mission.
- **Action**: Propose a decomposition plan to the user (break the Epic into smaller Track 2/3 missions).
