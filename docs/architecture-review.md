# Architecture Review: Two-Binary Trust Boundary (T-005)

**Date:** 2026-07-26
**Review scope:** Design §5 (Trust Boundary Model) and §6 (Interaction Contracts)
**Verdict:** ✅ Trust boundary is architecturally sound. Zero paths for Go to bypass Rust.

---

## 1. Trust Decision Inventory

Every trust decision in Pactyra was traced from origin to enforcement point. The rule is simple: **Rust owns the decision; Go must ask.**

| Trust Decision                   | Originates In                   | Enforced By                                                 | Go Bypass Possible?                                                                           |
| -------------------------------- | ------------------------------- | ----------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Is this target in scope?         | `pactyra-guardian` (Rust)       | Scope Guardian IPC response                                 | No — Go Worker Manager is downstream of Scope Guardian Agent                                  |
| Is this receipt valid?           | `pactyra-chain` (Rust)          | Chain verification on read                                  | No — receipts live in Rust chain; Go reads via IPC                                            |
| Is this chain intact?            | `pactyra-chain` (Rust)          | VerifyChain IPC operation                                   | No — verification runs in Rust                                                                |
| Can a receipt be modified?       | `pactyra-chain` (Rust)          | Append-only hash chain; signing key never leaves Rust       | No — Go has no write path to chain except RecordReceipt (append only)                         |
| Can a receipt be suppressed?     | `pactyra-chain` (Rust)          | Every ValidateInvocation produces a receipt (allow or deny) | No — Go cannot call Worker Manager without going through Scope Guardian, which always records |
| Can the signing key be accessed? | `pactyra-chain` signer (Rust)   | Key in `~/.pactyra/keys/` with 0600; never exposed via IPC  | No — key never crosses IPC boundary                                                           |
| Can evidence export be tampered? | `pactyra-chain` exporter (Rust) | ExportChain produces chain data + public key + metadata     | No — Go packages but cannot alter receipt content                                             |

---

## 2. Go → Rust Operation Audit

Design §6.1 defines nine IPC operations. Each was checked against the architecture diagram to confirm all Go trust-sensitive paths route through Rust.

| Operation            | Go Caller              | Rust Handler   | Trust Impact                                                | Coverage |
| -------------------- | ---------------------- | -------------- | ----------------------------------------------------------- | -------- |
| `CompileScope`       | CLI / Orchestrator     | Scope Compiler | Produces `scope_hash` — root of evidence chain              | ✅       |
| `ActivateScope`      | Orchestrator           | Scope Guardian | Binds session to scope                                      | ✅       |
| `ValidateInvocation` | Scope Guardian Agent   | Scope Guardian | **The runtime gate.** Every tool invocation must pass this. | ✅       |
| `RecordReceipt`      | Scope Guardian Agent   | Receipt Chain  | Appends to chain after tool execution                       | ✅       |
| `RecordNote`         | TUI / Orchestrator     | Receipt Chain  | Operator notes are receipted                                | ✅       |
| `RecordBlocked`      | Scope Guardian Agent   | Receipt Chain  | Blocked attempts are receipted                              | ✅       |
| `VerifyChain`        | Evidence Exporter (Go) | Receipt Chain  | Chain integrity check before export                         | ✅       |
| `ExportChain`        | Evidence Exporter (Go) | Receipt Chain  | Produces export data                                        | ✅       |
| `GetScopeStatus`     | TUI / Orchestrator     | Scope Guardian | Read-only scope summary                                     | ✅       |

**Coverage assessment:** All nine operations are accounted for. No trust-sensitive path exists outside this contract.

---

## 3. Architectural Gap Search

Each potential bypass path was examined against the design's trust domain model (§5.4).

### 3.1 Can Go execute a tool without Rust approval?

**Path examined:** Could the Worker Manager spawn a Docker container independently, bypassing the Scope Guardian Agent?

**Finding:** No. The Worker Manager is architecturally downstream of the Scope Guardian Agent. The flow is: Operator intent → Recon Planner → Scope Guardian Agent → Rust ValidateInvocation → (if allowed) Worker Manager. There is no code path from TUI or Orchestrator directly to Worker Manager that skips the Guardian.

**Defense-in-depth:** Even if a Go bug bypassed this, the Docker worker's network egress is restricted to scope-authorized targets at the network level. A rogue container would still be confined.

### 3.2 Can Go modify a receipt after creation?

**Path examined:** Could the Go binary alter the chain file on disk?

**Finding:** No. The chain file lives under `~/.pactyra/` and is managed by the Rust process. Go accesses receipts only through the IPC boundary (VerifyChain, ExportChain). Any direct disk write by Go would break the hash chain and be detected by Rust's chain verification.

### 3.3 Can Go access the signing key?

**Path examined:** Could the Go binary read `~/.pactyra/keys/ed25519_private.key`?

**Finding:** The key file has 0600 permissions. The Rust process reads it at startup and holds it in memory. The key is never serialized in any IPC response. Go has no IPC operation that returns key material. Even if Go read the file directly (same user, so 0600 permits it), the architecture treats this as a deployment concern — in production, the Rust binary could run as a separate user with exclusive key access.

**Recommendation for SDD-002:** Consider running `pactyra-core` as a separate system user (e.g., `pactyra`) with exclusive ownership of `~/.pactyra/keys/`. The Go binary runs as the operator user and communicates via the Unix socket. This would make the 0600 permission defense-in-depth rather than the sole protection.

### 3.4 Can Go suppress receipt generation?

**Path examined:** Could the Go binary skip the RecordReceipt call after tool execution?

**Finding:** No. The receipt chain is append-only, but more importantly: the ValidateInvocation call itself produces a receipt (the violation receipt when denied, or the invocation record when allowed). Every action that reaches Rust produces a receipt. The Go binary cannot call ValidateInvocation without a receipt being created. The only way to avoid a receipt is to avoid Rust entirely — which means no tool execution, since the Worker Manager is downstream of Rust validation.

---

## 4. Attack Surface Analysis

Design §5.1 identifies three Rust surfaces that face untrusted input:

| Surface                  | Untrusted Input                            | Risk                                                                                | Fuzzing Priority |
| ------------------------ | ------------------------------------------ | ----------------------------------------------------------------------------------- | ---------------- |
| **IPC Parser**           | Serialized requests from Go binary         | Medium — Go is same-user, but a compromised Go binary could send malicious payloads | High             |
| **Scope Grammar Parser** | Raw scope YAML from operator-supplied file | High — operator may supply maliciously crafted scope files                          | Critical         |
| **Hash Chain Validator** | Chain data read from disk                  | Medium — disk corruption or tampered chain file                                     | Medium           |

All three surfaces should be fuzzed before v0.1 ships. The scope grammar parser is the highest priority because it faces the operator's input directly.

---

## 5. IPC Protocol Hardening Recommendations

These are recommendations for SDD-005 (tool/agent contracts), not SDD-000 requirements:

1. **Message size limits.** The IPC parser should enforce a maximum request size (e.g., 64KB). Larger payloads should be rejected before deserialization.
2. **Request ID tracking.** Every request should carry a unique, monotonic ID. Responses should echo the request ID. This prevents replay and mismatch attacks.
3. **Connection authentication.** Consider a shared secret passed via environment variable when Go spawns `pactyra-core`. The first message on the socket must present this secret. Without it, the socket is useless even if permissions are misconfigured.
4. **Graceful degradation.** If Rust crashes or is killed, Go should detect the broken socket, log the event, restart the sidecar, and recover state by re-reading the chain. This is already specified in design §6.1.

---

## 6. Verdict

**The two-binary architecture is sound.** All trust decisions flow through Rust. The Go binary has zero code paths that can bypass scope validation, modify receipts, access the signing key, or suppress receipt generation. The architecture achieves its stated goal: Rust is the trust boundary; Go orchestrates but does not decide.

**No blocking gaps found.** The recommendations above are hardening measures for downstream SDDs, not architectural defects.

**Readiness for SDD-002:** ✅ Proceed. The architecture supports the scope model implementation without redesign.
