# SDD-005 Brief: Tool and Agent Contracts

**Parent:** SDD-000 (Vision & Governance)
**Status:** Brief
**Date:** 2026-07-26

---

## 1. What SDD-005 Owns

SDD-005 defines the **contracts** between Pactyra's orchestration layer and its tools and agents. These interfaces are the architectural guarantee that no agent — AI or otherwise — can bypass scope enforcement.

### Core Deliverables

1. **Agent interface** (`agent.Agent`): propose → validate → execute → record contract.
2. **Adapter interface** (`adapter.Adapter`): tool parameter validation, intent translation, output parsing.
3. **Scope Guardian Agent** (Go-side proxy): all proposals flow through this gate.
4. **Recon Planner Agent**: translates operator intent into concrete tool invocations.
5. **Worker Manager contract**: Docker lifecycle, network egress, timeout.
6. **Nmap adapter specification**: parameter validation, CLI flag mapping, XML parsing.

---

## 2. Requirements Extracted from SDD-000

### From Spec

- **"Execute tooling within enforced boundaries"** (JTBD 2): Tools run in isolated workers, receive only scope-authorized parameters.
- **"The operator is in control"** (Product Invariant 3): Agents propose, operators dispose. No autonomous execution.
- **"Scope enforcement is blocking, not advisory"** (Product Invariant 6): Validated by Scope Guardian Agent before Worker Manager.
- **"v0.1 supports exactly one tool — Nmap"** (v0.1 Scope Boundaries): Nmap adapter only.
- **"Multi-agent architecture"** (v0.1 Scope Boundaries): Scope Guardian + Recon Planner.

### From Design

- **Agent interface** (Design §6.2): In-process. Agents receive tasks, produce proposals.
- **"Propose → validate → execute → record" contract** (Design §6.2): Architectural invariant for all agents.
- **Recon Planner** (Design §3.2): Intent → Nmap invocations. Basic single-scan intent in v0.1.
- **Scope Guardian Agent** (Design §3.2): Go-side IPC proxy. All proposals must pass through.
- **Worker Manager** (Design §3.2, §6.3): Docker lifecycle, network egress, timeout.
- **Nmap Adapter** (Design §3.2): Parameter validation, CLI flag mapping, XML → structured output.
- **Tool Invocation Flow** (Design §4.2): Full happy-path sequence.
- **Scope Violation Flow** (Design §4.3): Blocked before Docker container creation.
- **Worker lifecycle** (Design §6.3): Create network, create container, capture output, destroy.

### From Architecture Review (T-005)

- Worker Manager MUST NOT be callable without prior Scope Guardian approval.
- IPC protocol hardening: message size limits, request ID tracking, connection authentication (recommendations for SDD-005).

---

## 3. Agent Interface (Architectural)

### `agent.Agent` (Go interface)

```go
// Agent receives tasks and produces proposals.
// All agents follow the "propose → validate → execute → record" contract.
type Agent interface {
    // Name returns the agent's identifier.
    Name() string

    // ReceiveTask accepts a structured task (from operator intent or prior agent output).
    ReceiveTask(task Task) error

    // Proposals returns the set of proposed actions.
    // Every proposal must pass through Scope Guardian before execution.
    Proposals() []Proposal
}

type Task struct {
    Intent   string            // Operator's expressed intent (natural language or structured)
    Context  map[string]any    // Additional context from orchestrator
}

type Proposal struct {
    Type      string            // "tool_invocation", "operator_note", "scope_expansion"
    Tool      string            // Tool name (e.g., "nmap")
    Target    TargetSpec        // IPs, domains, ports
    Params    map[string]any    // Tool-specific parameters
    Priority  int               // Execution priority (lower = sooner)
}
```

### Contract

1. An agent receives a `Task` from the Orchestrator.
2. An agent produces `Proposal`s.
3. Every `Proposal` of type `"tool_invocation"` must be validated by the Scope Guardian Agent.
4. If Scope Guardian returns `denied`, the proposal is discarded (receipt already recorded).
5. If Scope Guardian returns `allowed`, the proposal is forwarded to the Worker Manager.
6. Agents do NOT execute tools directly. They do NOT talk to Docker. They do NOT access the signing key.

### v0.1 Agents

| Agent                | Package                            | Responsibility                                                        |
| -------------------- | ---------------------------------- | --------------------------------------------------------------------- |
| Scope Guardian Agent | `internal/agent/guardian_agent.go` | Go-side IPC proxy. Routes ValidateInvocation to Rust. The gatekeeper. |
| Recon Planner        | `internal/agent/recon/planner.go`  | Translates operator intent into Nmap invocations.                     |

---

## 4. Adapter Interface (Architectural)

### `adapter.Adapter` (Go interface)

```go
// Adapter integrates an external tool with Pactyra's orchestration layer.
type Adapter interface {
    // ToolName returns the tool identifier (e.g., "nmap").
    ToolName() string

    // ValidateParams checks tool-specific parameters for correctness.
    // Called BEFORE the invocation is forwarded to Scope Guardian.
    ValidateParams(params map[string]any) error

    // BuildCommand translates a validated proposal into a command + args.
    BuildCommand(proposal Proposal) (cmd string, args []string, err error)

    // ParseOutput parses tool stdout into structured results.
    ParseOutput(stdout []byte) (Result, error)
}

type Result struct {
    Hosts    []Host      // Discovered hosts
    Services []Service   // Discovered services
    Raw      any         // Tool-specific structured output
}
```

### v0.1 Adapters

| Adapter | Package                  | Responsibility                                          |
| ------- | ------------------------ | ------------------------------------------------------- |
| Nmap    | `internal/adapter/nmap/` | Validate Nmap params, build CLI flags, parse XML output |

---

## 5. Worker Manager Contract

### Lifecycle

1. **Receive:** Worker Manager receives an approved invocation (target, tool, parameters, `invocation_id`).
2. **Network:** Create Docker network with egress restricted to authorized targets.
3. **Container:** Create and start container from tool image with injected parameters.
4. **Capture:** Capture stdout and stderr streams.
5. **Destroy:** On completion or timeout, destroy container and network.
6. **Return:** `{stdout, stderr, exit_code, duration}` to calling agent.

### Network Egress

- Each worker gets a dedicated Docker network.
- Egress rules derived from authorized target set in scope.
- v0.1: egress restricted to CIDR ranges in active scope.
- Worker cannot reach host, internet (unless in scope), or other workers.

### Timeout

- Default: 5 minutes for Nmap scans.
- Timeout kills the container.
- Receipt records timeout as exit condition.

---

## 6. Scope Violation Flow

```
Operator intent → Recon Planner → Proposal
  → Scope Guardian Agent → ValidateInvocation IPC → Rust Scope Guardian
  → {allowed: false, reason, violation_receipt}
  → Scope Guardian Agent: does NOT forward to Worker Manager
  → Worker Manager: NEVER CALLED
  → Docker: NO CONTAINER CREATED
  → Network: NO TRAFFIC SENT
```

**Critical invariant:** There must be NO code path that creates a Docker worker without prior Scope Guardian approval. The Worker Manager is downstream of the Scope Guardian Agent and is ONLY called when `allowed: true`.

---

## 7. v0.1 vs. v0.5 Scope

### v0.1 Deliverables

- [x] `agent.Agent` interface with Scope Guardian Agent + Recon Planner implementations
- [x] `adapter.Adapter` interface with Nmap implementation
- [x] Worker Manager with Docker lifecycle
- [x] Network egress enforcement
- [x] Scope violation flow (block → receipt → display)
- [x] Recon Planner: basic single-scan intent translation

### v0.5 Deliverables

- [ ] Tool Adapter SDK (community adapters)
- [ ] Nuclei adapter
- [ ] HTTP toolkit adapter
- [ ] Custom script adapter
- [ ] Multi-operator worker isolation
- [ ] Cartographer, Validator, Report Generator, Audit Agent implementations
- [ ] AI agent integration (implements `agent.Agent` interface)

---

## 8. Dependencies

| Depends On | Why                                              |
| ---------- | ------------------------------------------------ |
| SDD-000    | Vision, agent architecture, tool isolation model |

| Required By | Why                                            |
| ----------- | ---------------------------------------------- |
| SDD-0XX     | TUI command input routes to Recon Planner      |
| SDD-015     | Attack graph consumes tool output via adapters |
| SDD-XXX     | Tool adapter SDK extends the Adapter interface |

---

## 9. Acceptance Criteria (from SDD-000 Spec)

1. Recon Planner translates "scan web ports on scope targets" into valid Nmap invocations.
2. Scope Guardian Agent routes all proposals to Rust. No bypass path exists.
3. Nmap adapter validates parameters before execution.
4. Worker Manager creates Docker container with scope-only egress.
5. Worker Manager captures stdout/stderr and returns to calling agent.
6. Out-of-scope proposal is blocked before Docker container creation.
7. Blocked attempt produces a receipt and displays in TUI.
8. Worker timeout kills container and records timeout receipt.
