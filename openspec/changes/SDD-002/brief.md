# SDD-002 Brief: Scope Model and Compiler

**Parent:** SDD-000 (Vision & Governance)
**Status:** Brief
**Date:** 2026-07-26

---

## 1. What SDD-002 Owns

SDD-002 implements the **scope model** — the machine-enforceable representation of what an operator is authorized to target. This is the foundation of Pactyra's trust model: every action references a `scope_hash`, and the Scope Guardian evaluates every invocation against the compiled scope.

### Core Deliverables

1. **Scope parser** (`pactyra-scope` crate): Parse scope YAML into typed structures.
2. **Scope compiler**: Convert parsed scope into a canonical, machine-enforceable representation.
3. **Deterministic `scope_hash`**: SHA-256 hash of the canonical scope bytes.
4. **Scope activation**: Bind a compiled scope to an engagement session.
5. **Expansion envelope parser**: Parse (but not enforce in v0.1) the expansion policy.

---

## 2. Requirements Extracted from SDD-000

### From Spec

- **"Load and compile an authorized scope"** (Golden Path Step 1): Parse CIDR notation, domain lists, excluded hosts, port ranges. Reject malformed input with clear error messages.
- **"Generate a scope_hash"** (Golden Path Step 2): Deterministic hash. Same scope → same hash across runs and machines.
- **"Every action must be bound to an active scope"** (Product Invariant 1): Actions reference `scope_hash`.
- **"Scope enforcement is blocking, not advisory"** (Product Invariant 6): Out-of-scope actions are blocked.
- **"Runtime scope enforcement"** (v0.1 Scope Boundaries): Basic IP/port gating (O0–O3).
- **"Scope expansion envelope policy"** (Spec requirement): YAML structure for v0.5 forward compatibility.

### From Design

- **Scope Compiler component** (Design §3.1): `pactyra-scope` crate with `parser.rs`, `compiler.rs`, `hash.rs`.
- **Scope Loading Flow** (Design §4.1): `CompileScope` → `ActivateScope` IPC operations.
- **Scope Enforcement Flow** (Design §4.2-4.3): `ValidateInvocation` against compiled scope.
- **Expansion envelope policy structure** (Design reference to Proposal §6.1): `approver_roles`, `allowed`, `limits`, `always_requires_customer_approval`.
- **Scope state model** (Design §9.3): Compiled scope stored in `~/.pactyra/engagements/<name>/scope.compiled`.

### From Architecture Review (T-005)

- Scope grammar parser faces untrusted input (operator-supplied YAML). **Highest fuzzing priority.**
- Recommendation: define scope states explicitly — `draft`, `compiled`, `active`, `expanded`, `contracted`, `closed`.

### From Scope Expansion Review (T-008)

Edge cases for SDD-002 to address:

- Overlapping expansions (sequential processing with write lock)
- Concurrent expansion requests (batch approval model)
- Expansion rollback (scope contraction receipt, not deletion)
- Expansion after engagement close (frozen scope, addendum for retroactive corrections)
- Limits counter reset (per scope activation, monotonic within activation)
- Implicit expansions from tool discovery (within CIDR ≠ expansion)
- Multi-phase engagements (new scope document per phase, continuous chain)

---

## 3. Scope Grammar

### Supported Inputs (v0.1)

```yaml
targets:
  cidr:
    - '10.0.1.0/24'
    - '192.168.0.0/16'
  domains:
    - 'example.com'
    - '*.example.com' # Wildcard subdomain
  excluded:
    hosts:
      - '10.0.1.1'
      - '10.0.1.254'
  ports:
    include:
      - '1-1000'
      - '8080'
      - '8443'
    # exclude not needed in v0.1 — include is the allowlist
```

### Validation Rules

- CIDR must be valid IPv4 notation (no IPv6 in v0.1)
- Domains must be valid DNS names (RFC 1035)
- Excluded hosts must be within authorized CIDR ranges (warn if not)
- Port ranges: `start-end` or single port; start ≤ end; 1–65535
- Empty scope is invalid (must have at least one target)

### Canonical Representation

The compiled scope must be deterministically sortable:

- CIDR ranges sorted by network address
- Domains sorted alphabetically
- Excluded hosts sorted by IP
- Ports normalized (merge overlapping ranges, sort)

---

## 4. Boundary with SDD-003 (Permission Lattice)

**SDD-002 checks:** "Is `10.0.1.50:443` within the authorized scope?"
**SDD-003 checks:** "Is the operator's current permission level (O2) allowed to perform port scanning?"

These are complementary gates. SDD-002 provides the target validation; SDD-003 provides the capability validation. Both must pass for an action to proceed.

SDD-002 must provide:

- `ValidateTarget(ip, port) → bool` — is this target+port in scope?
- `ValidateCIDR(cidr) → bool` — is this CIDR fully contained in authorized CIDR?

SDD-003 consumes these to build permission-level rules on top.

---

## 5. Boundary with SDD-005 (Tool & Agent Contracts)

SDD-005 defines the `agent.Agent` and `adapter.Adapter` interfaces. SDD-002 provides the scope validation that the Scope Guardian Agent calls before any tool execution.

SDD-002 must provide:

- `scope_hash` as a stable reference
- `GetScopeStatus()` for TUI display
- `ValidateInvocation(targets, ports) → {allowed, invocation_id}` as the runtime gate

SDD-005's Scope Guardian Agent calls SDD-002's `ValidateInvocation` before forwarding to Worker Manager.

---

## 6. v0.1 vs. v0.5 Scope

### v0.1 Deliverables

- [x] CIDR parsing (IPv4 only)
- [x] Domain list parsing (exact and wildcard)
- [x] Excluded hosts
- [x] Port range parsing
- [x] Canonical compilation (deterministic ordering)
- [x] Deterministic `scope_hash` (SHA-256)
- [x] Scope activation/deactivation
- [x] ValidateInvocation (target vs. scope check)

### v0.5 Forward Compatibility

- [ ] Expansion envelope enforcement
- [ ] Dynamic scope expansion workflow
- [ ] Scope contraction recording
- [ ] Multi-phase engagement support
- [ ] IPv6 support

The expansion envelope policy structure must be parseable in v0.1 (forward compatibility) even though the workflow is not enforced until v0.5.

---

## 7. Dependencies

| Depends On | Why                                               |
| ---------- | ------------------------------------------------- |
| SDD-000    | Vision, invariants, golden path, scope boundaries |

| Required By | Why                                          |
| ----------- | -------------------------------------------- |
| SDD-003     | Permission lattice gates on scope validation |
| SDD-007     | Evidence chain root is `scope_hash`          |
| SDD-0XX     | TUI scope view displays compiled scope       |

---

## 8. Acceptance Criteria (from SDD-000 Spec)

1. Scope parser rejects malformed CIDR with clear error messages.
2. Same scope YAML produces identical `scope_hash` across runs and machines.
3. Different scope YAML produces different `scope_hash`.
4. `ValidateInvocation(10.0.0.0/8)` returns `denied` when scope is `10.0.1.0/24`.
5. `ValidateInvocation(10.0.1.50, 443)` returns `allowed` when `10.0.1.0/24` and port 443 are in scope.
6. Excluded host `10.0.1.1` returns `denied` even though it's within the authorized CIDR.
7. Empty scope, overlapping CIDR, malformed domain → rejected at parse time.
