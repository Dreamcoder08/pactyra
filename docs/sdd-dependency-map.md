# SDD Dependency Map

**Date:** 2026-07-26
**Source:** SDD-000 proposal §11.2 (Deferred to Future SDD Changes)
**Purpose:** Complete dependency graph for all downstream SDDs.

---

## 1. Complete SDD Inventory

| SDD     | Topic                         | Depends On       | v0.1 Required?   | Parallelizable With |
| ------- | ----------------------------- | ---------------- | ---------------- | ------------------- |
| SDD-000 | Vision, boundaries, non-goals | —                | ✅ (this change) | —                   |
| SDD-002 | Scope model and compiler      | SDD-000          | ✅ **Critical**  | SDD-005, SDD-003    |
| SDD-003 | Permission lattice O0–O6      | SDD-000, SDD-002 | ✅ (O0–O3 only)  | SDD-005             |
| SDD-005 | Tool and agent contracts      | SDD-000          | ✅ **Critical**  | SDD-002, SDD-003    |
| SDD-007 | Evidence and receipts         | SDD-000, SDD-002 | ✅ **Critical**  | SDD-010             |
| SDD-010 | Finding lifecycle             | SDD-007          | — (v0.5)         | SDD-015             |
| SDD-015 | Attack graph                  | SDD-005, SDD-007 | — (v0.5)         | SDD-010             |
| SDD-0XX | Operator TUI                  | SDD-002, SDD-005 | ✅               | SDD-007             |

---

## 2. Dependency Graph

```
SDD-000 (Vision & Governance)
  │
  ├── SDD-002 (Scope Model & Compiler) ──────┐
  │   ├── SDD-003 (Permission Lattice O0–O6) │
  │   ├── SDD-007 (Evidence & Receipts) ────┐│
  │   │   └── SDD-010 (Finding Lifecycle)   ││
  │   │       └── SDD-015 (Attack Graph)    ││
  │   └── SDD-0XX (Operator TUI)            ││
  │                                         ││
  └── SDD-005 (Tool & Agent Contracts) ─────┘│
      ├── SDD-015 (Attack Graph)             │
      └── SDD-0XX (Operator TUI)             │
```

---

## 3. v0.1 Minimal SDD Set

To ship v0.1 "First Light," the following SDDs must be completed (in order):

| Order | SDD                              | Rationale                                                                     |
| ----- | -------------------------------- | ----------------------------------------------------------------------------- |
| 1     | SDD-002 (Scope Model)            | Foundation — everything depends on scope parsing, compilation, and hashing.   |
| 2     | SDD-005 (Tool & Agent Contracts) | Defines Agent + Adapter interfaces. Can run in parallel with SDD-002.         |
| 3     | SDD-003 (Permission Lattice)     | O0–O3 only for v0.1. Depends on SDD-002 for scope model.                      |
| 4     | SDD-007 (Evidence & Receipts)    | Receipt schema, hash chain, signing, export. Depends on SDD-002 (scope_hash). |
| 5     | SDD-0XX (Operator TUI)           | Terminal UI. Depends on SDD-002 and SDD-005.                                  |

### Execution Strategy

**Wave 1 (parallel):** SDD-002 + SDD-005 — these are independent of each other and form the foundation.

**Wave 2 (parallel after Wave 1):** SDD-003 + SDD-007 — SDD-003 depends on SDD-002; SDD-007 depends on SDD-002. They can run in parallel since they don't depend on each other.

**Wave 3 (after Wave 2):** SDD-0XX — needs scope model (SDD-002) for scope view in TUI, and tool contracts (SDD-005) for command input that routes to agents.

---

## 4. v0.5 SDDs

These are NOT required for v0.1 but are on the critical path for v0.5:

| SDD     | Topic                     | Depends On       | Notes                                                        |
| ------- | ------------------------- | ---------------- | ------------------------------------------------------------ |
| SDD-010 | Finding lifecycle         | SDD-007          | Detection → validation → confirmation → remediation → retest |
| SDD-015 | Attack graph              | SDD-005, SDD-007 | State machine with query language                            |
| SDD-XXX | Tool adapter SDK          | SDD-005          | Community adapter marketplace                                |
| SDD-XXX | AI integration (optional) | SDD-005          | Provider plugins — never required                            |
| SDD-XXX | Multi-operator            | SDD-003, SDD-007 | Team collaboration, shared graph                             |
| SDD-XXX | Retesting engine          | SDD-007, SDD-010 | Evidence replay + diff engine                                |

---

## 5. Cross-SDD Boundary Clarification

### What SDD-002 owns vs. SDD-003

**SDD-002 (Scope Model):**

- CIDR parsing, domain validation, port range parsing
- Excluded host handling
- Canonical representation format
- Deterministic `scope_hash` generation
- Scope activation and binding to engagement
- Expansion envelope policy parser (for v0.5 forward compatibility)
- Scope states: draft, compiled, active, expanded, contracted, closed

**SDD-003 (Permission Lattice):**

- O0–O6 level definitions and capabilities
- Level escalation/de-escalation rules
- Enforcement: which levels can do what
- v0.1 scope: O0–O3 (basic IP/port gating)
- v0.5 scope: O0–O6 (full lattice)

**Boundary:** SDD-003 checks "is this _action type_ permitted at this _permission level_?" SDD-002 checks "is this _target_ within the authorized _scope_?" They are complementary gates.

### What SDD-005 owns vs. SDD-007

**SDD-005 (Tool & Agent Contracts):**

- `agent.Agent` Go interface: propose → validate → execute → record
- `adapter.Adapter` Go interface: validate parameters, translate intent, parse output
- Nmap adapter specification
- Scope Guardian agent (Go-side proxy)
- Recon Planner agent
- Worker Manager contract
- Scope violation flow

**SDD-007 (Evidence & Receipts):**

- Receipt schema (fields, types, serialization)
- Hash chain structure (linked list, parent hash)
- Ed25519 signing (key generation, signing, verification)
- Chain verification (full walk, tamper detection)
- Evidence export format (JSONL, metadata, public key)
- Standalone verification tool (`pactyra verify`)
- Three-layer trust model (integrity, provenance, reproducibility)

**Boundary:** SDD-005 defines "how tools and agents interact." SDD-007 defines "how their actions are recorded and proven." SDD-005's contract says "every invocation produces a receipt." SDD-007 defines what that receipt looks like and how it's verified.

### What SDD-007 owns vs. SDD-010

**SDD-007 (Evidence & Receipts):**

- Low-level: receipt creation, signing, chaining, verification
- Export: standalone evidence pack
- Integrity + provenance layers

**SDD-010 (Finding Lifecycle):**

- High-level: detection → validation → confirmation → remediation → retest
- Finding states and transitions
- Operator review workflow (confirm/reject)
- Report generation from evidence chain
- v0.1: basic evidence dump
- v0.5: structured finding management

**Boundary:** SDD-007 provides the raw evidence chain. SDD-010 provides the finding workflow on top of that chain.

---

## 6. Parallel vs. Serial Summary

| Phase       | SDDs                     | Execution               | Rationale                                                         |
| ----------- | ------------------------ | ----------------------- | ----------------------------------------------------------------- |
| v0.1 Wave 1 | SDD-002, SDD-005         | **Parallel**            | No mutual dependency. Both depend only on SDD-000.                |
| v0.1 Wave 2 | SDD-003, SDD-007         | **Parallel**            | Both depend on SDD-002 but not on each other.                     |
| v0.1 Wave 3 | SDD-0XX                  | **Serial** after Wave 2 | Depends on SDD-002 (scope views) and SDD-005 (agent routing).     |
| v0.5        | SDD-010, SDD-015, others | **Mostly parallel**     | SDD-010 depends on SDD-007. SDD-015 depends on SDD-005 + SDD-007. |
