# Design: Pactyra System Architecture (SDD-000)

**Status:** Reviewed  
**Date:** 2026-07-26  
**Phase:** sdd-design  
**Change ID:** SDD-000

---

## 1. Executive Summary

This document defines the **system-level architecture** for Pactyra — the high-level component model, trust boundaries, data flow, technology rationale, repository structure, and deployment model. It answers "how does the whole thing fit together?" without specifying concrete wire formats, schemas, or protocol-level details (those belong to dedicated SDD changes: SDD-002, SDD-003, SDD-005, SDD-007, SDD-015, SDD-0XX).

The architecture enforces a single hard rule derived from the proposal's invariants: **Rust owns every trust decision; Go orchestrates everything else.** The Rust side is a small, auditable core that holds the signing key, validates scope, and produces receipts. The Go side is the operator-facing runtime — TUI, Docker management, agent coordination — that must ask Rust for permission before any tool hits the network. Docker workers are the third trust domain: fully isolated, network-egress-filtered, with no access to keys or the host.

For v0.1, the system is a **two-binary deployment**: `pactyra` (Go) and `pactyra-core` (Rust), communicating over a local Unix domain socket. Single operator, single engagement, one tool (Nmap), two agents (Scope Guardian + Recon Planner).

---

## 2. System Architecture Diagram

```
                         ┌─────────────────────────────────────────┐
                         │            Operator Terminal             │
                         │  ┌───────────────────────────────────┐  │
                         │  │         pactyra (Go binary)        │  │
                         │  │                                   │  │
                         │  │  ┌─────────┐  ┌────────────────┐  │  │
                         │  │  │   TUI   │  │  Orchestrator  │  │  │
                         │  │  │ (Scope  │  │  (Session,     │  │  │
                         │  │  │  View,  │  │   State,       │  │  │
                         │  │  │ Receipt │  │   Agent Coord) │  │  │
                         │  │  │ Browser)│  │                │  │  │
                         │  │  └────┬────┘  └───────┬────────┘  │  │
                         │  │       │                │           │  │
                         │  │  ┌────┴────────────────┴───────┐   │  │
                         │  │  │         Agents (Go)          │   │  │
                         │  │  │  ┌──────────┐ ┌──────────┐  │   │  │
                         │  │  │  │ Scope    │ │ Recon    │  │   │  │
                         │  │  │  │ Guardian │ │ Planner  │  │   │  │
                         │  │  │  │ (Go side)│ │          │  │   │  │
                         │  │  │  └────┬─────┘ └────┬─────┘  │   │  │
                         │  │  └───────┼────────────┼────────┘   │  │
                         │  │          │            │            │  │
                         │  │  ┌───────┴────────────┴────────┐   │  │
                         │  │  │    Worker Manager (Go)       │   │  │
                         │  │  │  ┌──────────────────────┐   │   │  │
                         │  │  │  │   Nmap Adapter        │   │   │  │
                         │  │  │  └──────────────────────┘   │   │  │
                         │  │  └─────────────┬───────────────┘   │  │
                         │  └────────────────┼───────────────────┘  │
                         │                   │                      │
                         └───────────────────┼──────────────────────┘
                                             │
                    ┌────────────────────────┼────────────────────────┐
                    │   Local IPC (Unix Socket)                       │
                    └────────────────────────┼────────────────────────┘
                                             │
                         ┌───────────────────┼───────────────────────┐
                         │   pactyra-core (Rust binary)              │
                         │                   │                       │
                         │  ┌────────────────┼────────────────────┐  │
                         │  │     IPC Server / Request Router     │  │
                         │  └────┬──────┬─────┬──────┬───────────┘  │
                         │       │      │     │      │              │
                         │  ┌────┴──┐ ┌─┴───┐ ┌┴──────┴────────┐   │
                         │  │Scope │ │Scope│ │ Receipt Chain   │   │
                         │  │Compiler│ │Guard│ │                 │   │
                         │  │       │ │ian  │ │ ┌─────────────┐ │   │
                         │  │Parse  │ │     │ │ │ Signer      │ │   │
                         │  │CIDR   │ │Valid-│ │ │ (Ed25519)   │ │   │
                         │  │Domains│ │ate   │ │ │             │ │   │
                         │  │Hash   │ │Allow/│ │ │ Append-Only │ │   │
                         │  │       │ │Deny  │ │ │ Hash Chain  │ │   │
                         │  └───────┘ └─────┘ │ │ Verification │ │   │
                         │                    │ └─────────────┘ │   │
                         │                    │ ┌─────────────┐ │   │
                         │                    │ │ Evidence    │ │   │
                         │                    │ │ Exporter    │ │   │
                         │                    │ └─────────────┘ │   │
                         │                    └─────────────────┘   │
                         │                                          │
                         │  ┌────────────────────────────────────┐  │
                         │  │         Key Storage                │  │
                         │  │  (File-backed, ~/.pactyra/keys)   │  │
                         │  └────────────────────────────────────┘  │
                         └──────────────────────────────────────────┘
                                             │
                         ┌───────────────────┼───────────────────────┐
                         │        Docker Engine (local)              │
                         │                   │                       │
                         │  ┌────────────────┼────────────────────┐  │
                         │  │         pactyra-worker              │  │
                         │  │   (ephemeral container per task)    │  │
                         │  │  ┌──────────────────────────────┐  │  │
                         │  │  │  Nmap (or future tools)       │  │  │
                         │  │  │  Network: scope-only egress   │  │  │
                         │  │  │  No host fs, no key access    │  │  │
                         │  │  └──────────────────────────────┘  │  │
                         │  └────────────────────────────────────┘  │
                         └──────────────────────────────────────────┘
```

---

## 3. Component Map

### 3.1 Rust Components (Trust Boundary)

All trust decisions live here. The Rust codebase is deliberately small — its correctness is the foundation of every product invariant.

| Component                 | Crate                          | Responsibility                                                                                                                                                                                   | v0.1                         | v0.5+                            |
| ------------------------- | ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------- | -------------------------------- |
| **Scope Compiler**        | `pactyra-scope`                | Parse CIDR notation, domain lists, excluded hosts, port ranges into a canonical, machine-enforceable representation. Reject malformed input. Produce deterministic `scope_hash`.                 | ✅ Full                      | Extended for expansion envelope  |
| **Scope Guardian**        | `pactyra-guardian`             | Given a proposed tool invocation (target IP/domain, ports, tool, parameters), evaluate against the active compiled scope. Return allow/deny. Every decision is logged. This is the runtime gate. | ✅ Full                      | Extended for O4–O6 lattice rules |
| **Receipt Chain**         | `pactyra-chain`                | Append-only hash chain. Each receipt links to its parent via hash. Provides: receipt creation (with signing), chain integrity verification (detects tampering), chain traversal (for export).    | ✅ Integrity + Provenance    | Reproducibility layer added      |
| **Signer**                | (inside `pactyra-chain`)       | Holds the Ed25519 signing key. Signs every receipt. Key is generated at first run, stored in `~/.pactyra/keys/`. Never exposed outside the Rust process.                                         | ✅                           | TPM-backed optional              |
| **Evidence Exporter**     | (inside `pactyra-chain`)       | Produces the full chain + scope document + verification metadata as a standalone, verifiable artifact.                                                                                           | ✅ Basic (JSONL + signature) | Structured formats               |
| **IPC Server**            | (inside `pactyra-core` binary) | Listens on a local Unix domain socket. Routes requests from the Go binary to the appropriate internal component. Serializes responses.                                                           | ✅                           | —                                |
| **`pactyra-core` binary** | Cargo workspace root binary    | Combines all Rust crates into a single sidecar process. The `pactyra` Go binary spawns it as a child process on startup.                                                                         | ✅                           | —                                |

### 3.2 Go Components (Orchestration Boundary)

The Go binary is the operator's interface. It does not make trust decisions — it asks Rust.

| Component                          | Package                   | Responsibility                                                                                                                                                                                                                                              | v0.1                          | v0.5+                        |
| ---------------------------------- | ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- | ---------------------------- |
| **CLI Entry Point**                | `cmd/pactyra`             | Binary entry point. Parses subcommands (`scope compile`, `run`, `export`). Launches TUI or headless mode. Spawns `pactyra-core` child process.                                                                                                              | ✅                            | —                            |
| **TUI**                            | `internal/tui`            | Terminal interface. Scope view (what's authorized), command input, receipt browser, status line. Uses a Go TUI framework (Bubble Tea or tview). The TUI is the primary interface — no web UI in v0.1.                                                       | ✅ Basic views                | Team awareness, shared views |
| **Orchestrator**                   | `internal/orchestrator`   | Holds session state: active scope reference, current `scope_hash`, engagement metadata, agent registry. Coordinates agent lifecycle. Owns the IPC client to `pactyra-core`.                                                                                 | ✅ Single-engagement          | Multi-engagement             |
| **Scope Guardian Agent (Go side)** | `internal/agent/guardian` | Go-side proxy for the Rust Scope Guardian. Receives validation requests from the Recon Planner or operator, forwards to Rust via IPC, returns allow/deny. Records blocked attempts as local events.                                                         | ✅                            | Extended for dynamic scope   |
| **Recon Planner Agent**            | `internal/agent/recon`    | Translates operator intent ("scan web ports on in-scope targets") into concrete Nmap invocations. Proposes invocations to Scope Guardian before execution. Handles retry on transient failures (never on scope denials).                                    | ✅ Basic (single-scan intent) | Complex multi-stage plans    |
| **Worker Manager**                 | `internal/worker`         | Manages Docker worker lifecycle. Creates ephemeral containers per tool invocation. Injects tool parameters (validated by Scope Guardian). Captures stdout/stderr. Enforces network egress via Docker network policies. Destroys container after completion. | ✅ Nmap workers               | Multi-tool, multi-operator   |
| **Nmap Adapter**                   | `internal/adapter/nmap`   | Validates Nmap-specific parameters. Translates Recon Planner intent into Nmap CLI flags. Parses Nmap XML output into structured result. Extracts hosts, ports, services, OS guesses.                                                                        | ✅                            | —                            |
| **Evidence Exporter (Go side)**    | `internal/export`         | Requests full evidence chain from Rust. Packages with metadata (engagement name, operator, dates). Produces standalone export (directory with receipts, scope, README with verification instructions).                                                      | ✅                            | —                            |
| **Guardian Client**                | `internal/guardian`       | IPC client library. Manages the Unix socket connection to `pactyra-core`. Handles serialization, timeouts, reconnection. The only Go package that talks to Rust.                                                                                            | ✅                            | —                            |

### 3.3 v0.5+ Components (Architectural Placeholders)

These are not built in v0.1 but their place in the architecture is defined now to avoid redesign.

| Component                  | Language   | Responsibility                                                             |
| -------------------------- | ---------- | -------------------------------------------------------------------------- |
| **Cartographer Agent**     | Python     | Consumes tool output, constructs attack graph nodes/edges.                 |
| **Validator Agent**        | Python     | Proof-of-exploitation. Crafts and records exploit attempts.                |
| **Report Generator Agent** | Python/Go  | Renders evidence graph into structured report. Template-driven.            |
| **Audit Agent**            | Go         | Continuously verifies chain integrity. Detects anomalies.                  |
| **Tool Adapter SDK**       | Go         | Interface + documentation for community tool adapters.                     |
| **Web Dashboard**          | TypeScript | Read-only views for PMs and clients. Never an operational control surface. |

---

## 4. Data Flow

### 4.1 Scope Loading Flow

```
Operator            pactyra (Go)              pactyra-core (Rust)
   │                     │                          │
   │  pactyra scope      │                          │
   │  compile scope.yaml │                          │
   │────────────────────>│                          │
   │                     │  CompileScope(scope_yaml)│
   │                     │─────────────────────────>│
   │                     │                          │──── Scope Compiler parses,
   │                     │                          │     validates, hashes
   │                     │                          │
   │                     │  {scope_hash,            │
   │                     │   compiled_scope_id}     │
   │                     │<─────────────────────────│
   │                     │                          │
   │  "Scope compiled:   │                          │
   │   a1b2c3d4...       │                          │
   │   3 IP ranges,      │                          │
   │   2 domains"        │                          │
   │<────────────────────│                          │
   │                     │                          │
   │  pactyra run        │                          │
   │  --scope a1b2c3d4   │                          │
   │────────────────────>│                          │
   │                     │  ActivateScope(hash)     │
   │                     │─────────────────────────>│
   │                     │  {ok}                    │
   │                     │<─────────────────────────│
   │                     │                          │
   │  TUI launches,      │                          │
   │  scope view active  │                          │
```

### 4.2 Tool Invocation Flow (Happy Path)

```
Operator        pactyra (Go)                  pactyra-core (Rust)       Docker
   │                 │                                │                    │
   │ "scan web       │                                │                    │
   │  ports on       │                                │                    │
   │  scope targets" │                                │                    │
   │────────────────>│                                │                    │
   │                 │──── Recon Planner translates   │                    │
   │                 │     intent → nmap -sV -p 80,   │                    │
   │                 │     443,8080,8443 10.0.1.0/24 │                    │
   │                 │                                │                    │
   │                 │──── Scope Guardian Agent:      │                    │
   │                 │     Validate(targets, ports)   │                    │
   │                 │───────────────────────────────>│                    │
   │                 │                                │──── Scope Guardian │
   │                 │                                │     checks against │
   │                 │                                │     compiled scope │
   │                 │                                │                    │
   │                 │  {allowed: true,               │                    │
   │                 │   invocation_id,               │                    │
   │                 │   scope_hash}                  │                    │
   │                 │<───────────────────────────────│                    │
   │                 │                                │                    │
   │                 │──── Worker Manager:            │                    │
   │                 │     CreateWorker(              │                    │
   │                 │       image=nmap,              │                    │
   │                 │       args=[...],              │                    │
   │                 │       network=scope-only)      │                    │
   │                 │───────────────────────────────────────────────────>│
   │                 │                                │   Container runs   │
   │                 │                                │   nmap, output     │
   │                 │                                │   captured         │
   │                 │  {stdout, stderr, exit_code}   │                    │
   │                 │<───────────────────────────────────────────────────│
   │                 │                                │                    │
   │                 │──── Nmap Adapter parses XML    │                    │
   │                 │──── Scope Guardian Agent:      │                    │
   │                 │     RecordReceipt(             │                    │
   │                 │       invocation_id,           │                    │
   │                 │       output_hash,             │                    │
   │                 │       exit_code)               │                    │
   │                 │───────────────────────────────>│                    │
   │                 │                                │──── Receipt Chain  │
   │                 │                                │     creates receipt │
   │                 │                                │     signs, appends │
   │                 │  {receipt_id, receipt_hash}    │                    │
   │                 │<───────────────────────────────│                    │
   │                 │                                │                    │
   │ "Scan complete:  │                                │                    │
   │  12 hosts,       │                                │                    │
   │  47 ports open"  │                                │                    │
   │<────────────────│                                │                    │
```

### 4.3 Scope Violation Flow (Blocked)

```
Operator        pactyra (Go)                  pactyra-core (Rust)       Docker
   │                 │                                │                    │
   │ "scan           │                                │                    │
   │  10.0.0.0/8"    │                                │                    │
   │────────────────>│                                │                    │
   │                 │──── Scope Guardian Agent:      │                    │
   │                 │     Validate(10.0.0.0/8, ...)  │                    │
   │                 │───────────────────────────────>│                    │
   │                 │                                │──── Checks scope   │
   │                 │                                │     10.0.1.0/24   │
   │                 │                                │     authorized     │
   │                 │                                │     10.0.0.0/8    │
   │                 │                                │     NOT in scope   │
   │                 │  {allowed: false,              │                    │
   │                 │   reason: "10.0.0.0/8 exceeds  │                    │
   │                 │   authorized 10.0.1.0/24",     │                    │
   │                 │   violation_receipt}           │                    │
   │                 │<───────────────────────────────│                    │
   │                 │                                │                    │
   │                 │  (Docker worker is NEVER       │                    │
   │                 │   created. No network traffic  │                    │
   │                 │   has left the host.)          │                    │
   │                 │                                │                    │
   │ "BLOCKED:        │                                │                    │
   │  10.0.0.0/8      │                                │                    │
   │  out of scope"   │                                │                    │
   │<────────────────│                                │                    │
```

### 4.4 Evidence Export Flow

```
Operator        pactyra (Go)                  pactyra-core (Rust)
   │                 │                                │
   │  pactyra export │                                │
   │  --output ./    │                                │
   │  evidence/      │                                │
   │────────────────>│                                │
   │                 │──── ExportChain()              │
   │                 │───────────────────────────────>│
   │                 │                                │──── Reads full chain
   │                 │                                │     Verifies integrity
   │                 │                                │     Produces: receipts,
   │                 │                                │     scope_doc,
   │                 │                                │     public_key,
   │                 │                                │     verification_metadata
   │                 │  {chain_data,                  │
   │                 │   public_key,                  │
   │                 │   root_hash}                   │
   │                 │<───────────────────────────────│
   │                 │                                │
   │                 │──── Packages into directory:   │
   │                 │     evidence/                  │
   │                 │     ├── receipts.jsonl         │
   │                 │     ├── scope.yaml             │
   │                 │     ├── metadata.json          │
   │                 │     ├── public_key.pem         │
   │                 │     └── VERIFY.md              │
   │                 │                                │
   │  "Evidence pack  │                                │
   │   exported to    │                                │
   │   ./evidence/"   │                                │
   │<────────────────│                                │
```

---

## 5. Trust Boundary Model

Pactyra has three trust domains, ordered from most-trusted to least-trusted. Data flows from less-trusted to more-trusted for decisions; decisions flow from more-trusted to less-trusted.

### 5.1 Domain 1: Rust Trust Core (`pactyra-core`)

**What it trusts:** Nothing outside itself. It trusts its own code, its own key material, and the operating system's process isolation.

**What it owns:**

- The signing key (Ed25519). Generated at first run, stored in `~/.pactyra/keys/`. Never exposed via IPC — the Go binary never sees it.
- All scope decisions. The Go binary asks "may I do X?" and the Rust binary answers yes or no. There is no code path in Go that can bypass this.
- The receipt chain. Rust creates, signs, and appends every receipt. The chain's integrity is verified by Rust on every read.
- Evidence export. Rust produces the chain data and verification metadata. Go packages it but cannot alter receipt content.

**What it does NOT trust:**

- The Go binary. If the Go binary sends garbage, Rust validates and rejects. If the Go binary crashes, Rust's state is unaffected (chain is on disk).
- The Docker worker output. Worker output is hashed and receipted as-is. Rust does not interpret Nmap XML.
- The network. Rust makes no outbound connections. Its only IPC is the local Unix socket.

**Attack surface:** The Unix socket IPC parser, the scope grammar parser, the hash chain validator. These are the three Rust surfaces that face untrusted input. All three must be fuzzed.

### 5.2 Domain 2: Go Orchestrator (`pactyra`)

**What it trusts:** The Rust binary's responses. It must trust that "allowed" means allowed and "denied" means denied. It trusts Docker Engine's container isolation.

**What it owns:**

- The operator experience (TUI).
- Agent logic (Recon Planner translations, retry decisions).
- Docker worker lifecycle.
- Tool adapter output parsing.
- Evidence pack packaging (formatting, not content).

**What it does NOT trust:**

- Operator input. Every command is validated through Rust before reaching a Docker worker.
- Its own state. If the Go process restarts, it recovers state from Rust (chain is the source of truth).
- Docker worker output. It parses and displays but does not trust it for security decisions.

**What it CANNOT do (by architecture, not convention):**

- Execute a tool without Rust approval. The Worker Manager is downstream of the Scope Guardian Agent.
- Modify a receipt after creation. Receipts live in Rust's chain.
- Access the signing key. The key never leaves the Rust process.
- Suppress a receipt. Every validation (allow or deny) produces a receipt in Rust before Go proceeds.

### 5.3 Domain 3: Docker Workers

**What it trusts:** Nothing. Workers are ephemeral, single-use, and fully isolated.

**What it owns:** Nothing persistent. Workers receive tool parameters, execute, produce output, and are destroyed.

**What it CANNOT do:**

- Access the host filesystem (no volume mounts to host paths).
- Access the signing key (key is in Rust process memory, not on a mounted path).
- Communicate with anything outside the authorized target set (network egress filtering via Docker network policies).
- Communicate with the Go or Rust processes (no shared networks with host).
- Survive beyond the tool invocation (destroyed on completion or timeout).

### 5.4 Trust Domain Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     TRUST DOMAIN 1                           │
│                   pactyra-core (Rust)                        │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │  Signing │  │  Scope   │  │ Receipt  │  │  Evidence   │ │
│  │   Key    │  │ Guardian │  │  Chain   │  │  Exporter   │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
│                                                             │
│  Trusts: nothing outside itself                             │
│  Attack surface: IPC parser, scope parser, chain validator │
└───────────────────────┬─────────────────────────────────────┘
                        │  Unix Socket (local only)
                        │  Request: validate, record, export
                        │  Response: allow/deny, receipt, chain
┌───────────────────────┴─────────────────────────────────────┐
│                     TRUST DOMAIN 2                           │
│                     pactyra (Go)                             │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │   TUI    │  │  Agents  │  │  Worker  │  │  Nmap       │ │
│  │          │  │          │  │  Manager │  │  Adapter    │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
│                                                             │
│  Trusts: Rust responses, Docker Engine isolation            │
│  Cannot: bypass Rust, modify receipts, access key          │
└───────────────────────┬─────────────────────────────────────┘
                        │  Docker API (local socket)
                        │  Create/start/stop containers
┌───────────────────────┴─────────────────────────────────────┐
│                     TRUST DOMAIN 3                           │
│                 Docker Workers (ephemeral)                   │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Tool container (Nmap)                                │  │
│  │  - Network: scope-only egress                         │  │
│  │  - Filesystem: ephemeral, no host mounts              │  │
│  │  - Keys: none accessible                              │  │
│  │  - Lifetime: single invocation, destroyed after       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
│  Trusts: nothing                                            │
│  Cannot: access host, key, or internal IPC                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Component Interaction Contracts

These are architectural contracts — they describe what each boundary guarantees, not the wire format. Concrete schemas belong to SDD-005 (tool/agent contracts) and SDD-007 (evidence/receipts).

### 6.1 Go ↔ Rust IPC Contract

**Transport:** Local Unix domain socket. The Go binary spawns `pactyra-core` as a child process. The socket path is passed via command-line argument or environment variable. If the Rust process dies, Go detects it via socket closure and restarts it — Rust's state is on disk so recovery is safe.

**Interaction model:** Request-response. Go sends a request; Rust sends exactly one response. Requests are idempotent where possible. No streaming in v0.1.

**Operations (architectural):**

| Operation            | Direction | Input (conceptual)                                | Output (conceptual)                                  | Idempotent?           |
| -------------------- | --------- | ------------------------------------------------- | ---------------------------------------------------- | --------------------- |
| `CompileScope`       | Go → Rust | Raw scope YAML/text                               | `scope_hash`, compiled scope metadata                | Yes                   |
| `ActivateScope`      | Go → Rust | `scope_hash`                                      | Confirmation                                         | Yes                   |
| `ValidateInvocation` | Go → Rust | Target IPs/domains, ports, tool                   | `{allowed: bool, invocation_id, violation_receipt?}` | Yes                   |
| `RecordReceipt`      | Go → Rust | `invocation_id`, output hash, exit code, metadata | `receipt_id`, `receipt_hash`                         | No (appends to chain) |
| `RecordNote`         | Go → Rust | Operator note text                                | `receipt_id`                                         | No                    |
| `RecordBlocked`      | Go → Rust | Violation details                                 | `receipt_id`                                         | No                    |
| `VerifyChain`        | Go → Rust | (none — verifies entire chain)                    | `{valid: bool, tampered_receipts: [...]}`            | Yes                   |
| `ExportChain`        | Go → Rust | (none)                                            | Full chain data, public key, metadata                | Yes                   |
| `GetScopeStatus`     | Go → Rust | (none)                                            | Current scope summary                                | Yes                   |

**Error model:** Every response includes a status. Errors are either:

- **Rejected** (invalid input, out-of-scope, malformed) — normal operations, receipt generated.
- **Fault** (IPC error, internal Rust error, chain corruption) — Go retries or surfaces to operator.

**Trust note:** The Go binary trusts the Rust binary's responses. There is no "trust but verify" — Rust IS the trust boundary. If Rust says denied, Go does not execute. If Rust says the chain is valid, Go reports it as valid.

### 6.2 Go Orchestrator ↔ Go Agents Contract

**Interaction model:** In-process. Agents are Go interfaces, not separate processes.

**Agent interface (conceptual):**

- An agent receives a **task** (structured from operator intent or prior agent output).
- An agent produces **proposals** (things to do: tool invocations, notes, scope expansion requests).
- Every proposal must pass through the Scope Guardian Agent before execution.
- Agents do not execute tools directly — they ask the Worker Manager through the Scope Guardian.

**v0.1 agents:**

- **Recon Planner:** receives operator intent ("scan web ports on scope"), produces Nmap invocation proposals. One intent → one or more Nmap invocations.
- **Scope Guardian Agent:** the Go-side proxy. Every proposal from any agent flows through it. It calls Rust `ValidateInvocation` and returns allow/deny.

**Future agents (v0.5+) follow the same contract:** propose → validate → execute → record. This is the architectural invariant that guarantees no agent, AI or otherwise, can bypass scope.

### 6.3 Go Worker Manager ↔ Docker Contract

**Interaction model:** Go Docker SDK. One ephemeral container per tool invocation.

**Worker lifecycle:**

1. Worker Manager receives an approved invocation (target, tool, parameters, `invocation_id`).
2. Worker Manager creates a Docker network with egress restricted to the authorized target set.
3. Worker Manager creates and starts a container from the tool image with the injected parameters.
4. Worker Manager captures stdout and stderr streams.
5. On completion or timeout, Worker Manager destroys the container and network.
6. Worker Manager returns `{stdout, stderr, exit_code, duration}` to the calling agent.

**Network egress enforcement:**

- Each worker gets a dedicated Docker network.
- Egress rules are derived from the authorized target set in the scope.
- For v0.1: egress is restricted to the CIDR ranges in the active scope.
- The worker cannot reach the host, the internet (unless in scope), or other workers.
- This is defense-in-depth: even if Rust validation has a bug, the network won't pass the traffic.

**Timeout:** Every worker has a maximum execution time. Default: 5 minutes for Nmap scans. Timeout kills the container. Receipt records timeout as exit condition.

### 6.4 Operator ↔ TUI Contract

**Interaction model:** Terminal keystrokes and commands. The TUI is the primary interface.

**Views (v0.1):**

- **Scope View:** Shows the active scope summary (IP ranges, domains, excluded hosts, ports). Read-only except for initial loading.
- **Command Input:** Where the operator types intents ("scan web ports on scope targets") or direct commands. Not a shell — the operator expresses intent, and the Recon Planner translates.
- **Receipt Browser:** Scrollable list of recent receipts. Shows action type, target, timestamp, allow/deny, receipt hash prefix. Filterable by type.
- **Status Line:** Current `scope_hash` prefix, engagement name, agent status, chain height.

**Command model (conceptual):**

- `pactyra scope compile <file>` — compile and load a scope (headless or TUI)
- `pactyra run --scope <hash>` — launch the TUI with an active scope
- Within TUI: operator types intents, not raw tool commands. The Recon Planner translates "scan web ports" into `nmap -sV -p 80,443,...`. This indirection is what allows scope validation to work — the operator never constructs tool parameters directly.
- `pactyra export --output <dir>` — export evidence pack (can be run from TUI or headless)
- `pactyra verify <evidence-dir>` — standalone verification of an evidence pack (no running platform needed; this is a separate, stateless command)

**What the operator CANNOT do from the TUI:**

- Run a tool directly with arbitrary parameters (no raw shell access)
- Bypass scope validation
- Delete or modify a receipt
- Access the signing key

---

## 7. Technology Stack Rationale

### 7.1 Rust — Trust Boundary

**Why Rust:**

- **Memory safety without garbage collection.** Scope validation, receipt signing, and chain verification must not have buffer overflows, use-after-free, or data races. Rust's ownership model eliminates these at compile time. A GC'd language would introduce unpredictable latency at signing time and larger attack surface in the runtime.
- **Zero-cost abstractions.** The scope compiler parses untrusted input. The chain verifier walks a potentially long hash chain. Both need to be fast without sacrificing safety.
- **Small, auditable binaries.** The Rust trust core is deliberately small. A security reviewer can audit the scope compiler, guardian, and chain in a few hundred lines each. The language's type system encodes invariants (e.g., "a receipt always has a parent hash except the genesis receipt") that would be runtime checks in other languages.
- **FFI-friendly.** If the sidecar model ever changes to a linked library (cgo), Rust compiles to C-compatible static libraries cleanly.

**Why NOT Rust for everything:**

- The Go ecosystem's Docker SDK is mature and well-maintained. Rewriting Docker client logic in Rust would add weeks to v0.1.
- TUI frameworks in Rust (Ratatui) are solid, but Go's Bubble Tea is more approachable for the operator-facing interaction model v0.1 needs.
- The orchestrator needs concurrency for managing multiple workers, agents, and the TUI simultaneously. Go's goroutine model is a natural fit.

### 7.2 Go — Orchestration

**Why Go:**

- **Docker SDK.** The official Docker client library is Go. Using it directly avoids wrapping a REST API.
- **Concurrency model.** The orchestrator manages the TUI event loop, multiple agent goroutines, worker lifecycle, and IPC to Rust — all concurrently. Goroutines and channels make this natural.
- **Single binary deployment.** Go compiles to a static binary. The operator experience is "download one binary" (plus the Rust sidecar, which Go manages).
- **Fast compile times.** The orchestration layer will iterate rapidly during development. Go's compile speed keeps the feedback loop tight.
- **Standard library.** `net`, `os/exec`, `encoding/json`, `crypto/sha256` — everything needed for the Go side is in the standard library or the Docker SDK.

**Why NOT Go for the trust boundary:**

- Garbage collection pauses are unpredictable. Not a problem for a TUI, but undesirable in a cryptographic signing path that must be deterministic.
- Go's type system cannot express "this buffer has been validated" at compile time. Rust's ownership can.
- The Go runtime is a larger trusted computing base than a minimal Rust binary.

### 7.3 Docker — Tool Isolation

**Why Docker:**

- **Universal availability.** Every offensive operator already has Docker or Podman. No new virtualization layer to install.
- **Network egress control.** Docker network policies provide defense-in-depth: even if scope validation has a bug, the worker container cannot reach out-of-scope targets.
- **Ephemeral by default.** Containers are designed to be created, used, and destroyed. Perfect for single-invocation workers.
- **Tool packaging.** Tool images (Nmap, future tools) are self-contained. No dependency conflicts with the host.
- **Resource limits.** CPU and memory limits prevent a runaway scan from impacting the host.

**Why NOT VMs or sandboxes:**

- Full VMs have higher startup latency (seconds vs. milliseconds for containers).
- Language-specific sandboxes (e.g., `seccomp`, `pledge`) are platform-specific and harder to configure correctly than Docker network policies.
- Docker provides a uniform interface across Linux distributions — important for a self-hostable tool.

### 7.4 Unix Domain Sockets — IPC

**Why Unix sockets:**

- **Local only.** No network exposure. The Rust binary only accepts connections from the local Go binary.
- **File-system permissions.** The socket file is created with `0600` permissions. Only the same user can connect.
- **No dependency.** Unix sockets are in the kernel. No message broker, no shared memory library.
- **Simple request-response.** Matches the Go ↔ Rust interaction model: ask, receive answer, proceed.

### 7.5 Ed25519 — Signing

**Why Ed25519:**

- Small keys (32 bytes private, 32 bytes public).
- Fast signing and verification.
- Deterministic signatures (no entropy required at signing time — important for reproducibility).
- Widely supported, well-understood security properties.
- The public key is compact enough to include in every exported evidence pack.

### 7.6 Hash Chain (not a Blockchain)

**Why a simple hash chain, not a blockchain or Merkle tree:**

- **v0.1 simplicity.** A linked list where each receipt hashes its parent is sufficient for append-only integrity. No consensus, no distributed nodes.
- **Verification is O(n).** For v0.1 engagement sizes (hundreds to low thousands of receipts), linear verification is fast enough.
- **No cryptocurrency or distributed ledger association.** Blockchain terminology would create confusion about the project's scope and intent.
- **Future upgrade path.** If performance requires it, the chain can be upgraded to a Merkle tree structure in v0.5 without changing the receipt format. But v0.1 ships with a simple linked hash chain.

---

## 8. Repository Structure

```
ark-onyx-offensive/                  # Monorepo root
│
├── README.md                        # Project identity, thesis, quickstart
├── LICENSE                          # TBD (open-source, protects mission)
├── CONTRIBUTING.md                  # How to contribute, adapter SDK (v0.5+)
├── SECURITY.md                      # Security policy, reporting, trust model
│
├── docs/
│   ├── architecture.md              # This design document
│   ├── golden-path.md               # Step-by-step v0.1 demo walkthrough
│   ├── threat-model.md              # Formal threat model (v0.5+)
│   └── cyber-range.md               # Local cyber range setup for testing
│
├── crates/                          # Rust workspace
│   ├── pactyra-scope/               # Scope parser, compiler, hash
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── lib.rs
│   │       ├── parser.rs            # CIDR + domain parsing
│   │       ├── compiler.rs          # Canonical representation
│   │       └── hash.rs              # Deterministic scope_hash
│   │
│   ├── pactyra-guardian/            # Runtime scope enforcement
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── lib.rs
│   │       └── guardian.rs          # Validate invocation against scope
│   │
│   ├── pactyra-chain/               # Receipt chain + signing
│   │   ├── Cargo.toml
│   │   └── src/
│   │       ├── lib.rs
│   │       ├── receipt.rs           # Receipt structure, creation
│   │       ├── chain.rs             # Append-only chain, verification
│   │       ├── signer.rs            # Ed25519 signing, key management
│   │       └── export.rs            # Chain export for evidence pack
│   │
│   └── pactyra-core/                # Rust binary (combines all crates)
│       ├── Cargo.toml
│       └── src/
│           ├── main.rs              # Entry point, IPC server setup
│           ├── server.rs            # Unix socket listener, request routing
│           └── handlers.rs          # Operation handlers
│
├── cmd/                             # Go binaries
│   └── pactyra/
│       └── main.go                  # Go entry point
│
├── internal/                        # Go internal packages (no external import)
│   ├── tui/                         # Terminal UI
│   │   ├── tui.go                   # Bubble Tea model setup
│   │   ├── views/
│   │   │   ├── scope.go             # Scope view
│   │   │   ├── command.go           # Command input
│   │   │   ├── receipts.go          # Receipt browser
│   │   │   └── status.go            # Status line
│   │   └── styles.go                # Colors, formatting
│   │
│   ├── orchestrator/                # Session & state management
│   │   ├── orchestrator.go          # Main orchestrator, agent lifecycle
│   │   ├── session.go               # Engagement session state
│   │   └── state.go                 # In-memory state, persistence
│   │
│   ├── guardian/                    # Rust guardian IPC client
│   │   ├── client.go                # Unix socket client
│   │   ├── protocol.go              # Request/response types (not wire format)
│   │   └── client_test.go
│   │
│   ├── agent/                       # Agent implementations
│   │   ├── agent.go                 # Agent interface definition
│   │   ├── guardian_agent.go        # Scope Guardian (Go-side proxy)
│   │   └── recon/                   # Recon Planner
│   │       ├── planner.go           # Intent → invocation translation
│   │       └── planner_test.go
│   │
│   ├── worker/                      # Docker worker management
│   │   ├── manager.go               # Container lifecycle
│   │   ├── network.go               # Network egress policies
│   │   └── manager_test.go
│   │
│   ├── adapter/                     # Tool adapters
│   │   ├── adapter.go               # Adapter interface definition
│   │   └── nmap/
│   │       ├── nmap.go              # Nmap parameter validation + execution
│   │       ├── parser.go            # XML → structured output
│   │       └── nmap_test.go
│   │
│   └── export/                      # Evidence pack packaging
│       ├── export.go                # Pack creation, file writing
│       └── verify.go                # Standalone verification logic
│
├── tests/                           # Integration tests
│   ├── golden_path_test.go          # Full golden path: scope → run → export
│   ├── scope_violation_test.go      # Out-of-scope blocking tests
│   ├── chain_integrity_test.go      # Tampered chain detection
│   └── fixtures/                    # Test data
│       ├── scope_valid.yaml
│       ├── scope_invalid.yaml
│       └── expected_outputs/
│
├── cyber-range/                     # Local test infrastructure
│   ├── docker-compose.yml           # Vulnerable target services
│   ├── README.md                    # How to start the range
│   └── targets/                     # Target service configs
│
├── go.mod
├── go.sum
├── Cargo.toml                       # Rust workspace manifest
├── Makefile                         # Build, test, run, lint targets
│
└── openspec/                        # SDD artifacts (committed)
    ├── config.yaml
    └── changes/
        └── SDD-000/
            ├── proposal.md
            ├── spec.md
            ├── design.md            # This file
            └── tasks.md             # (future: sdd-tasks phase)
```

### 8.1 Build System

**Makefile targets (v0.1):**

- `make build` — builds `pactyra` (Go) and `pactyra-core` (Rust)
- `make test` — runs Go tests, Rust tests, and integration tests
- `make lint` — `golangci-lint` + `clippy`
- `make run` — builds and launches the TUI with a test scope
- `make golden-path` — runs the full golden path integration test against the cyber range
- `make clean` — removes build artifacts

No Bazel or Nix in v0.1. Those add complexity without accelerating the 8-12 week timeline. The Rust workspace and Go module are each built with their native toolchains (`cargo build`, `go build`). The Makefile orchestrates them.

---

## 9. Deployment & Runtime Model

### 9.1 Operator Experience

The operator downloads or clones the repository. They need:

- **Go toolchain** (to build `pactyra`) or a pre-built binary
- **Rust toolchain** (to build `pactyra-core`) or a pre-built binary
- **Docker Engine** (running, with the operator in the `docker` group)
- **Nmap** (or the `pactyra-worker-nmap` Docker image)

**First run:**

```bash
$ git clone https://github.com/arkelythex/pactyra.git
$ cd pactyra
$ make build
$ ./pactyra scope compile engagement-123.yaml
Scope compiled: a1b2c3d4e5f6...
  3 IP ranges, 2 domains, 0 excluded hosts
$ ./pactyra run --scope a1b2c3d4e5f6
[pactyra TUI launches]
```

**Within the TUI:**

```
┌─ Scope: a1b2c3d4 (engagement-123) ─────────────────────────────┐
│ 10.0.1.0/24 (256 hosts)  │  *.example.com (2 resolved)         │
│ Excluded: 10.0.1.1        │  Ports: 1-1000, 8080, 8443         │
├─ Command ──────────────────────────────────────────────────────┤
│ > scan web ports on all scope targets                          │
├─ Receipts (12) ────────────────────────────────────────────────┤
│ #0012  ALLOW   scan    10.0.1.0/24    80,443,8080,8443  0.3s   │
│ #0011  ALLOW   scan    10.0.1.5        1-1000            2.1s   │
│ #0010  BLOCKED scan    10.0.0.0/8     —                 0.0s    │
│ ...                                                             │
├─ Status ────────────────────────────────────────────────────────┤
│ Chain: 12 receipts  │  Scope: a1b2c3d4  │  Agents: SG ●  RP ●  │
└─────────────────────────────────────────────────────────────────┘
```

### 9.2 Binary Distribution

**v0.1:** Build from source. Pre-built binaries for Linux x86_64 and ARM64 may be provided as GitHub releases.

**The two-binary model:**

- `pactyra` (Go): ~15-25 MB static binary. The operator runs this directly.
- `pactyra-core` (Rust): ~5-10 MB static binary. `pactyra` spawns it as a child process. The operator never runs it directly.

The operator only interacts with `pactyra`. The Rust sidecar is an implementation detail.

### 9.3 State Storage

```
~/.pactyra/
├── keys/
│   ├── ed25519_private.key         # Signing key (0600 permissions)
│   └── ed25519_public.pem          # Public key (0644)
├── engagements/
│   └── <engagement-name>/
│       ├── scope.compiled          # Compiled scope (Rust binary format)
│       ├── scope.yaml              # Original scope document
│       ├── chain.dat               # Receipt chain (Rust binary format)
│       └── metadata.json           # Engagement metadata
└── config.yaml                     # User configuration (optional)
```

All trust-critical data (keys, chain, compiled scope) lives under `~/.pactyra/` and is managed by the Rust binary. The Go binary reads/writes only through the IPC boundary.

### 9.4 Evidence Pack (Export Format)

When the operator exports:

```bash
pactyra export --output ./evidence-2026-07-26/
```

The resulting directory:

```
evidence-2026-07-26/
├── receipts.jsonl                 # One JSON object per line, each a receipt
├── scope.yaml                     # Original scope document
├── metadata.json                  # Engagement name, operator, dates, chain height
├── public_key.pem                 # Ed25519 public key for verification
└── VERIFY.md                      # Human-readable verification instructions
```

**Standalone verification (no platform needed):**

```bash
$ pactyra verify ./evidence-2026-07-26/
Verifying 47 receipts...
  Chain integrity: VALID
  All receipts signed by: a1b2c3d4...
  Scope hash: e5f6a1b2...
  No tampering detected.
```

The `verify` subcommand is a stateless operation. It reads the evidence pack, checks the hash chain, verifies every signature against the public key, and reports results. It does not require `pactyra-core` to be running.

---

## 10. v0.1 Minimal Architecture

### 10.1 What Ships in v0.1

| Layer           | Component                                                     | Status                        |
| --------------- | ------------------------------------------------------------- | ----------------------------- |
| **Rust**        | Scope Compiler (CIDR, domain list, excluded hosts)            | Full                          |
| **Rust**        | Scope Guardian (allow/deny per invocation)                    | Full                          |
| **Rust**        | Receipt Chain (append-only hash chain, Ed25519 signing)       | Integrity + Provenance layers |
| **Rust**        | Evidence Exporter (chain → standalone data)                   | Full                          |
| **Rust**        | `pactyra-core` binary (IPC server)                            | Full                          |
| **Go**          | CLI entry point (`scope compile`, `run`, `export`, `verify`)  | Full                          |
| **Go**          | TUI (scope view, command input, receipt browser, status line) | Full                          |
| **Go**          | Orchestrator (single session, single engagement)              | Full                          |
| **Go**          | Scope Guardian Agent (Go-side IPC proxy)                      | Full                          |
| **Go**          | Recon Planner (intent → Nmap invocations)                     | Full                          |
| **Go**          | Worker Manager (Docker lifecycle, network egress)             | Full                          |
| **Go**          | Nmap Adapter (parameter validation, XML parsing)              | Full                          |
| **Go**          | Evidence Exporter (pack packaging, standalone verifier)       | Full                          |
| **Docker**      | `pactyra-worker-nmap` image                                   | Full                          |
| **Integration** | Golden path test (8 steps)                                    | Full                          |
| **Integration** | Scope violation blocking test                                 | Full                          |
| **Integration** | Chain tampering detection test                                | Full                          |

### 10.2 What is Explicitly NOT in v0.1

| Capability                                               | Where It Goes          |
| -------------------------------------------------------- | ---------------------- |
| Attack graph (nodes, edges, query language)              | v0.5 (SDD-015)         |
| Adversarial validation / proof-of-exploitation           | v0.5                   |
| Retesting from evidence replay                           | v0.5                   |
| Dynamic scope expansion workflow                         | v0.5                   |
| Permission lattice O4–O6                                 | v0.5 (SDD-003)         |
| Team collaboration / multi-operator                      | v0.5                   |
| Tool adapter SDK                                         | v0.5 (SDD-005)         |
| Tools beyond Nmap (Nuclei, HTTP toolkit, custom scripts) | v0.5                   |
| Cartographer, Validator, Report Generator, Audit Agent   | v0.5                   |
| Python or TypeScript code                                | v0.5 / v1.0            |
| Web dashboard                                            | v1.0                   |
| AI integration (any provider)                            | v0.5 (always optional) |
| TPM-backed key storage                                   | v0.5+                  |
| Reproducibility layer (evidence model)                   | v0.5 (SDD-007)         |

### 10.3 v0.1 Architecture Diagram (Reduced)

```
┌──────────────────────────────────────────┐
│              pactyra (Go)                 │
│                                          │
│  TUI ◄──► Orchestrator ◄──► Agents       │
│                │              │           │
│                │              ▼           │
│                │         Scope Guardian   │
│                │         (Go proxy)       │
│                │              │           │
│                ▼              │           │
│           Worker Manager      │           │
│                │              │           │
│                ▼              │           │
│           Nmap Adapter        │           │
│                │              │           │
└────────────────┼──────────────┼───────────┘
                 │              │
                 │    ┌─────────┘
                 │    │  Unix Socket
                 │    │
┌────────────────┼────┼─────────────────────┐
│  pactyra-core  │    │                     │
│  (Rust)        │    │                     │
│                │    ▼                     │
│           IPC Server                      │
│                │                          │
│     ┌──────────┼──────────┐              │
│     ▼          ▼          ▼              │
│  Scope      Scope      Receipt            │
│  Compiler   Guardian   Chain              │
│                        ┌─ Signer          │
│                        └─ Exporter        │
│                                           │
│  Key Storage: ~/.pactyra/keys/            │
└───────────────────────────────────────────┘
                 │
                 │  Docker API
                 ▼
┌───────────────────────────────────────────┐
│         Docker Worker (Nmap)               │
│  Network: scope-only egress               │
│  Ephemeral, no host access                │
└───────────────────────────────────────────┘
```

---

## 11. Risks & Mitigations (Architectural)

| Risk                                       | Severity | Architectural Mitigation                                                            |
| ------------------------------------------ | -------- | ----------------------------------------------------------------------------------- |
| **Scope bypass via IPC bug**               | Critical | Unix socket permissions (0600). Rust validates all input. Fuzz the IPC parser.      |
| **Signing key compromise**                 | Critical | Key never leaves Rust process. File permissions (0600). TPM-backed in v0.5.         |
| **Chain corruption on disk**               | High     | Chain integrity verified on every read. Corruption detected, not silently accepted. |
| **Docker escape by malicious tool output** | Medium   | Workers are ephemeral. No host mounts. Network egress restricted.                   |
| **Go binary crash loses in-flight state**  | Medium   | Rust chain is the source of truth. Go recovers by re-reading chain on restart.      |
| **Sidecar IPC latency**                    | Low      | Local Unix socket. Scope validation is O(1) against compiled scope.                 |
| **Multi-language build complexity**        | Medium   | Makefile orchestrates `cargo` + `go`. Pre-built binaries for releases.              |
| **TUI complexity delays v0.1**             | Medium   | TUI starts minimal: 4 views, no animations, no plugins. Function over form.         |

---

## 12. Design Decisions Log

| Decision                                     | Rationale                                                                               | Reversible?                                                                                              |
| -------------------------------------------- | --------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| **Two-binary model (Go + Rust sidecar)**     | Clean trust boundary. No cgo complexity. Easier to test independently.                  | Yes — could link Rust as cgo library later, but no reason to in v0.1–v0.5.                               |
| **Unix domain socket for IPC**               | Local only, no network exposure, simple request-response.                               | Yes — could use stdin/stdout of child process, but socket is more flexible for future multi-process.     |
| **Simple hash chain (not Merkle tree)**      | Sufficient for v0.1 engagement sizes. Simpler to implement and verify.                  | Yes — can upgrade to Merkle tree in v0.5 without changing receipt format.                                |
| **Ed25519 (not RSA, not ECDSA)**             | Small keys, fast, deterministic, well-understood.                                       | Mostly no — Ed25519 is the right choice. Changing would require key migration.                           |
| **Operator types intents, not raw commands** | Indirection allows scope validation to work. The Recon Planner can optimize parameters. | Partially — could expose raw commands within scope for power users, but the intent model is the default. |
| **No Bazel/Nix in v0.1**                     | Makefile + native toolchains is sufficient. Bazel/Nix add weeks of setup.               | Yes — can adopt Bazel if monorepo grows to justify it in v0.5+.                                          |
| **Docker, not Podman or VMs**                | Universal, well-known, good Go SDK. Podman compatibility is a v0.5 consideration.       | Yes — Podman support can be added later via abstraction.                                                 |
| **Nmap as the only v0.1 tool**               | Avoids tool adapter fragmentation. Proves the architecture works with one tool first.   | Yes — tool adapter SDK in v0.5 opens the ecosystem.                                                      |
| **No attack graph in v0.1**                  | The attack graph is valuable but complex. Scope enforcement + evidence is the wedge.    | Yes — attack graph architecture is defined for v0.5, doesn't require v0.1 redesign.                      |

---

## 13. Future-Proofing: What the v0.1 Architecture Enables

The v0.1 architecture is deliberately minimal, but it does not paint the project into a corner:

- **Tool adapter SDK (v0.5):** The `adapter.Adapter` Go interface is defined in v0.1 with a single Nmap implementation. Adding Nuclei, HTTP toolkit, or custom scripts means implementing the same interface.
- **Multi-agent (v0.5):** The `agent.Agent` Go interface is defined in v0.1 with Scope Guardian and Recon Planner. Adding Cartographer and Validator means implementing the same interface with the same "propose → validate → execute → record" contract.
- **Attack graph (v0.5):** The evidence chain already records every tool output. The Cartographer consumes these receipts to build graph nodes. No v0.1 redesign needed.
- **Retesting (v0.5):** The evidence chain already records exact tool invocations. The retest engine replays them. No v0.1 redesign needed.
- **Permission lattice (v0.5):** The Scope Guardian already makes allow/deny decisions. Adding O4–O6 rules is an extension of the decision logic, not a new architecture.
- **AI integration (v0.5):** AI providers are agents that implement the same `agent.Agent` interface. They propose, the operator approves, Scope Guardian validates, Worker Manager executes. AI is just another agent — the architecture doesn't need to know.
- **Web dashboard (v1.0):** Read-only views over the evidence chain. No operational control — the terminal remains primary. The chain export already provides the data.

---

## Appendix A: Component Hierarchy (v0.1 Go Dependencies)

```
cmd/pactyra/main.go
  └── internal/tui           (terminal UI)
        └── internal/orchestrator  (session state, agent coordination)
              ├── internal/agent/guardian_agent  (Go-side Scope Guardian)
              │     └── internal/guardian        (Rust IPC client)
              ├── internal/agent/recon           (Recon Planner)
              │     └── internal/agent/guardian_agent
              └── internal/worker               (Docker manager)
                    └── internal/adapter/nmap    (Nmap adapter)
```

No circular dependencies. Every package depends downward. The `internal/guardian` package (Rust IPC client) is a leaf — no other Go package talks to Rust directly.

---

## Appendix B: Rust Crate Dependency Graph (v0.1)

```
pactyra-core (binary)
  ├── pactyra-scope    (scope parsing, compiling, hashing)
  ├── pactyra-guardian (runtime scope enforcement)
  │     └── pactyra-scope
  └── pactyra-chain    (receipt chain, signing, export)
        └── pactyra-scope  (for scope_hash root reference)
```

`pactyra-scope` is the foundational crate — everything depends on the scope model. `pactyra-chain` depends on `pactyra-scope` because the genesis receipt references the `scope_hash`.

---

## Appendix C: Key Terminology

| Term               | Definition                                                                                                                                                                         |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Scope**          | The machine-readable authorization document: IP ranges, domains, excluded hosts, port ranges, authorized techniques.                                                               |
| **scope_hash**     | The cryptographic hash of the compiled scope. The root of the evidence chain.                                                                                                      |
| **Receipt**        | An immutable record of an action (tool invocation, blocked attempt, operator note). Signed, hash-chained, append-only.                                                             |
| **Evidence Chain** | The linked list of receipts. Each receipt hashes its parent. Tampering with any receipt breaks the chain.                                                                          |
| **Evidence Pack**  | A standalone, exportable directory containing the full chain, scope document, public key, and verification instructions. Verifiable without the platform.                          |
| **Trust Boundary** | The architectural line that separates components that make trust decisions (Rust) from components that request them (Go) and components that are fully untrusted (Docker workers). |
| **Agent**          | A Go component that receives tasks and produces proposals. Agents propose; operators dispose.                                                                                      |
| **Worker**         | An ephemeral Docker container that executes a single tool invocation. Network-egress-filtered.                                                                                     |
| **Scope Guardian** | The agent (Rust side) + proxy (Go side) that validates every tool invocation against scope before execution. The architectural gatekeeper.                                         |
| **Recon Planner**  | The agent that translates operator intent into concrete tool invocations.                                                                                                          |
| **Golden Path**    | The 8-step v0.1 demonstration: scope → hash → discover → attempt violation → block → record → receipts → export.                                                                   |
