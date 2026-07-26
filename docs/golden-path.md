# Golden Path Walkthrough: v0.1 "First Light"

**Date:** 2026-07-26
**Purpose:** Step-by-step validation of the 8-step golden path against v0.1 architecture components.
**Verdict:** ✅ All 8 steps are testable with v0.1 architecture. No v0.5+ dependencies.

---

## Golden Path Overview

The golden path is the v0.1 acceptance test: an operator loads a scope, runs discovery, attempts an out-of-scope action, sees it blocked, and exports a verifiable evidence pack — all against a local cyber range. If a third party can clone the repo, run the golden path, and verify the evidence pack without platform access, v0.1 is complete.

```
1. LOAD    → 2. HASH   → 3. DISCOVER → 4. VIOLATE  → 5. BLOCK
                                               ↓
8. EXPORT  ← 7. RECEIPTS ← 6. RECORD
```

---

## Step 1: Load and Compile an Authorized Scope

**What happens:**
The operator provides a scope YAML file (CIDR ranges, domain list, excluded hosts, port ranges). Pactyra parses, validates, and compiles it into a machine-enforceable representation.

**Components involved:**

| Component                        | Role                                                |
| -------------------------------- | --------------------------------------------------- |
| `cmd/pactyra` (Go)               | CLI entry point: `pactyra scope compile scope.yaml` |
| `internal/orchestrator` (Go)     | Routes the compile request                          |
| `internal/guardian` (Go)         | IPC client: sends `CompileScope` request            |
| `pactyra-scope` (Rust)           | Parses CIDR, domains, exclusions, port ranges       |
| `pactyra-core` IPC Server (Rust) | Routes request to Scope Compiler                    |

**Data flow:**

```
scope.yaml → Go CLI → IPC: CompileScope(scope_yaml) → Rust Scope Compiler
  → Parse CIDR notation, domain lists, excluded hosts, port ranges
  → Validate (reject malformed entries with clear errors)
  → Compile into canonical representation
  → Generate deterministic scope_hash
  → Return {scope_hash, compiled_scope_id, summary} to Go
```

**Acceptance test:** A valid `scope.yaml` produces a deterministic `scope_hash`. An invalid scope (overlapping CIDR, malformed domain) is rejected with a clear error message pointing to the specific line.

**v0.1 feasibility:** ✅ `pactyra-scope` crate implements CIDR parsing, domain validation, and hashing. All are v0.1 components.

---

## Step 2: Generate a scope_hash

**What happens:**
The compiled scope is cryptographically hashed. This `scope_hash` becomes the root of the evidence chain — every subsequent action references it.

**Components involved:**

| Component                      | Role                                                 |
| ------------------------------ | ---------------------------------------------------- |
| `pactyra-scope/hash.rs` (Rust) | Deterministic hash of canonical scope representation |
| `pactyra-chain` (Rust)         | Genesis receipt references `scope_hash`              |

**Data flow:**

```
Compiled scope (canonical bytes) → SHA-256 → scope_hash (hex string)
  → Genesis receipt created: {type: "scope_activation", scope_hash, timestamp, ...}
  → Subsequent receipts link to genesis via parent_hash chain
```

**Acceptance test:**

- Same scope YAML produces identical `scope_hash` across runs and machines.
- Different scope produces different `scope_hash`.
- Genesis receipt is verifiable and references the `scope_hash`.
- Every subsequent action receipt contains the `scope_hash` field.

**v0.1 feasibility:** ✅ Deterministic hashing and genesis receipt creation are v0.1 features. Hash chain is implemented in `pactyra-chain`.

---

## Step 3: Execute Permitted Discovery

**What happens:**
The operator expresses intent ("scan web ports on scope targets"). The Recon Planner translates this into Nmap invocations. The Scope Guardian validates targets against the active scope. The Worker Manager spawns a Docker container running Nmap. Output is captured, parsed, and receipted.

**Components involved:**

| Component                            | Role                                                                   |
| ------------------------------------ | ---------------------------------------------------------------------- |
| `internal/tui` (Go)                  | Operator types intent                                                  |
| `internal/agent/recon` (Go)          | Translates "scan web ports" → `nmap -sV -p 80,443,8080,8443 <targets>` |
| `internal/agent/guardian_agent` (Go) | Proxies validation to Rust                                             |
| `pactyra-guardian` (Rust)            | Checks targets against compiled scope → allow/deny                     |
| `internal/worker` (Go)               | Creates Docker network (scope-only egress), starts Nmap container      |
| `internal/adapter/nmap` (Go)         | Validates Nmap parameters, parses XML output                           |
| `pactyra-chain` (Rust)               | Records receipt: tool invocation, output hash, exit code               |

**Data flow:**

```
Operator: "scan web ports on scope targets"
  → Recon Planner: intent → [nmap -sV -p 80,443,8080,8443 10.0.1.0/24]
  → Scope Guardian Agent → IPC: ValidateInvocation(targets, ports, tool)
  → Rust Scope Guardian: check targets ⊆ compiled scope → {allowed: true, invocation_id}
  → Worker Manager: create Docker network (egress: 10.0.1.0/24:80,443,8080,8443)
  → Worker Manager: docker run pactyra-worker-nmap [...args]
  → Container executes Nmap, output captured
  → Nmap Adapter: parse XML → structured results
  → Scope Guardian Agent → IPC: RecordReceipt(invocation_id, output_hash, exit_code)
  → Rust Receipt Chain: create, sign, append receipt
  → TUI displays: "12 hosts, 47 ports open"
```

**Acceptance test:**

- Nmap scan completes against cyber range targets within scope.
- Docker worker has network egress restricted to scope targets only.
- Output is parsed into structured results (hosts, ports, services).
- Evidence receipt is created with invocation_id, output_hash, and exit_code.

**v0.1 feasibility:** ✅ All components (Recon Planner, Scope Guardian, Worker Manager, Nmap Adapter, Receipt Chain) are v0.1. Only Nmap tooling is supported.

---

## Step 4: Attempt an Out-of-Scope Action

**What happens:**
The operator (or test script) deliberately attempts an action against a target NOT in the active scope. For example, the scope authorizes `10.0.1.0/24` but the operator requests a scan of `10.0.0.0/8`.

**Components involved:**

| Component                            | Role                                                              |
| ------------------------------------ | ----------------------------------------------------------------- |
| `internal/agent/recon` (Go)          | Translates intent into invocation (including out-of-scope target) |
| `internal/agent/guardian_agent` (Go) | Forwards validation request to Rust                               |
| `pactyra-guardian` (Rust)            | Detects scope violation                                           |

**Data flow:**

```
Operator/intent: "scan 10.0.0.0/8"
  → Recon Planner: intent → [nmap ... 10.0.0.0/8]
  → Scope Guardian Agent → IPC: ValidateInvocation(10.0.0.0/8, ...)
  → Rust Scope Guardian: compiled scope = 10.0.1.0/24
  → 10.0.0.0/8 ⊄ 10.0.1.0/24 → DETECTED VIOLATION
```

**Acceptance test:** The violation is detected by Rust before any Go code proceeds to Worker Manager. No Docker container is created.

**v0.1 feasibility:** ✅ Scope Guardian crate implements IP range containment checking. CIDR subnet math is deterministic and testable.

---

## Step 5: Block Before Execution

**What happens:**
Upon detecting the violation, the Rust Scope Guardian returns `{allowed: false, reason: "...", violation_receipt}`. The Go Scope Guardian Agent receives this denial and does NOT forward to Worker Manager. No network traffic leaves the host.

**Components involved:**

| Component                            | Role                                            |
| ------------------------------------ | ----------------------------------------------- |
| `pactyra-guardian` (Rust)            | Returns denial + violation receipt              |
| `internal/agent/guardian_agent` (Go) | Receives denial, does NOT call Worker Manager   |
| `pactyra-chain` (Rust)               | Violation receipt is already created and signed |

**Data flow:**

```
Rust → IPC: {allowed: false, reason: "10.0.0.0/8 exceeds authorized 10.0.1.0/24", violation_receipt}
  → Go Scope Guardian Agent: receives denied
  → Worker Manager: NEVER CALLED
  → Docker: NO CONTAINER CREATED
  → Network: NO TRAFFIC SENT
  → TUI displays: "BLOCKED: 10.0.0.0/8 out of scope"
```

**Critical path check:** The Worker Manager must only be called when `allowed: true`. There must be no code path that calls Worker Manager with a denied invocation. The design ensures this: the Scope Guardian Agent is the only caller of Worker Manager, and it only calls on `allowed: true`.

**Acceptance test:**

- Blocked action produces zero network traffic.
- Docker container is never created.
- Violation receipt exists in chain.
- TUI displays block reason clearly.
- No "warn but proceed" code path exists.

**v0.1 feasibility:** ✅ The Scope Guardian's allow/deny is binary. No partial-allow or warn-and-proceed mode is implemented.

---

## Step 6: Record Actions and Decisions

**What happens:**
Every action — successful tool invocation, blocked violation, operator note — produces an evidence receipt. The receipt is appended to the hash chain. The chain's integrity is verifiable at any point.

**Components involved:**

| Component                             | Role                                                                                                  |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `pactyra-chain/receipt.rs` (Rust)     | Receipt structure: action type, parameters, output hash, timestamp, operator, scope_hash, parent_hash |
| `pactyra-chain/chain.rs` (Rust)       | Append-only linked list; each receipt hashes its parent                                               |
| `pactyra-chain/signer.rs` (Rust)      | Ed25519 signing of each receipt                                                                       |
| `internal/tui/views/receipts.go` (Go) | Receipt browser: scrollable list with filters                                                         |

**Data flow:**

```
Every action → Rust RecordReceipt / RecordBlocked / RecordNote
  → Receipt created: {id, type, timestamp, operator, scope_hash, parent_hash, payload}
  → Signer signs receipt with Ed25519
  → Receipt appended to chain
  → Chain integrity: hash(current) = SHA-256(parent_hash || receipt_data)
  → Any tampering with prior receipt breaks the hash chain → detectable
```

**Receipt types in v0.1:**

- `scope_activation` — genesis receipt
- `tool_invocation` — successful tool execution
- `violation_blocked` — blocked out-of-scope attempt
- `operator_note` — free-text operator annotation

**Acceptance test:**

- All four receipt types are produced and verifiable.
- Hash chain is continuous (no gaps, no broken links).
- VerifyChain returns `{valid: true}` after golden path run.
- Tampering with any receipt is detected (test: flip a bit in chain.dat → VerifyChain returns `{valid: false}`).

**v0.1 feasibility:** ✅ `pactyra-chain` implements all four receipt types and hash chain verification.

---

## Step 7: Generate Verifiable Receipts

**What happens:**
Any receipt in the chain can be independently verified: integrity (content unchanged since registration), provenance (who, when, tool, authorization), and basic reproducibility (procedure documented).

**Components involved:**

| Component                        | Role                                                                    |
| -------------------------------- | ----------------------------------------------------------------------- |
| `pactyra-chain/chain.rs` (Rust)  | VerifyChain: walks chain, checks hashes, verifies signatures            |
| `internal/export/verify.go` (Go) | Standalone verifier: reads evidence pack, checks chain, reports results |

**Data flow:**

```
VerifyChain request → Rust walks entire chain:
  For each receipt:
    1. Verify signature against public key
    2. Verify parent_hash matches previous receipt's hash
    3. Verify scope_hash matches genesis
  → Return {valid: bool, tampered_receipts: [...], chain_height, root_hash}
```

**Acceptance test:**

- Every receipt is individually verifiable for integrity.
- Every receipt traces to who, when, tool, authorization.
- The full chain verification succeeds on a clean chain.
- A tampered receipt is detected with specific identification (which receipt, what mismatch).

**v0.1 feasibility:** ✅ Integrity and provenance layers are v0.1. Reproducibility layer (full replay) is deferred to v0.5, but the receipt structure captures the necessary data (tool, parameters, target) for future replay.

---

## Step 8: Export a Reproducible Evidence Pack

**What happens:**
The operator exports the engagement. Rust produces the full chain data, scope document, public key, and verification metadata. Go packages it into a directory. The pack is verifiable without access to the Pactyra platform.

**Components involved:**

| Component                        | Role                                                                    |
| -------------------------------- | ----------------------------------------------------------------------- |
| `pactyra-chain/export.rs` (Rust) | ExportChain: reads full chain, verifies integrity, produces export data |
| `internal/export/export.go` (Go) | Packages data into directory structure, writes VERIFY.md                |
| `internal/export/verify.go` (Go) | `pactyra verify <dir>` — standalone verification                        |

**Data flow:**

```
Operator: pactyra export --output ./evidence/
  → Go Evidence Exporter → IPC: ExportChain()
  → Rust: reads full chain, verifies integrity, produces:
      - chain_data (all receipts)
      - scope_document (original scope.yaml)
      - public_key (Ed25519 public key PEM)
      - verification_metadata (chain height, root hash, timestamps)
  → Go: writes evidence/ directory:
      evidence/
      ├── receipts.jsonl       # One JSON receipt per line
      ├── scope.yaml           # Original scope document
      ├── metadata.json        # Engagement name, operator, dates, chain height
      ├── public_key.pem       # Ed25519 public key
      └── VERIFY.md            # Human-readable verification instructions

Standalone verification:
  $ pactyra verify ./evidence/
  Verifying 47 receipts...
    Chain integrity: VALID
    All receipts signed by: a1b2c3d4...
    Scope hash: e5f6a1b2...
    No tampering detected.
```

**Acceptance test:**

- Evidence pack contains all required files.
- `pactyra verify ./evidence/` succeeds on a valid pack.
- A third party with no platform access can verify the pack.
- No proprietary format, tool, or key is required.
- VERIFY.md instructions are clear and executable.

**v0.1 feasibility:** ✅ All components (Rust exporter, Go packager, standalone verifier) are v0.1. The standalone `verify` subcommand is a stateless operation that reads the evidence pack and checks the hash chain.

---

## Step Dependency Summary

| Step                      | v0.1 Components Used                                                       | v0.5+ Dependencies | Testable Now? |
| ------------------------- | -------------------------------------------------------------------------- | ------------------ | ------------- |
| 1. Load scope             | Scope Compiler, IPC                                                        | None               | ✅            |
| 2. Generate scope_hash    | Scope Compiler, Receipt Chain                                              | None               | ✅            |
| 3. Execute discovery      | Recon Planner, Scope Guardian, Worker Manager, Nmap Adapter, Receipt Chain | None               | ✅            |
| 4. Attempt violation      | Recon Planner, Scope Guardian Agent                                        | None               | ✅            |
| 5. Block before execution | Scope Guardian, Worker Manager (NOT called)                                | None               | ✅            |
| 6. Record actions         | Receipt Chain, Signer                                                      | None               | ✅            |
| 7. Verify receipts        | Receipt Chain verification                                                 | None               | ✅            |
| 8. Export evidence        | Evidence Exporter (Rust + Go), Standalone verifier                         | None               | ✅            |

---

## Integration Test Plan (for SDD-002+)

The golden path should be an automated integration test:

```go
// tests/golden_path_test.go
func TestGoldenPath(t *testing.T) {
    // Step 1: Compile scope
    scopeHash := compileScope(t, "tests/fixtures/scope_valid.yaml")

    // Step 2: Verify scope_hash is deterministic
    scopeHash2 := compileScope(t, "tests/fixtures/scope_valid.yaml")
    assert.Equal(t, scopeHash, scopeHash2)

    // Step 3: Execute permitted discovery
    results := runDiscovery(t, scopeHash, "scan web ports on scope targets")
    assert.NotEmpty(t, results.Hosts)

    // Step 4-5: Attempt out-of-scope, verify blocked
    err := runDiscovery(t, scopeHash, "scan 10.0.0.0/8")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "out of scope")

    // Step 6-7: Verify chain integrity
    chain := exportChain(t)
    assert.True(t, chain.Valid)

    // Step 8: Export and verify standalone
    evidenceDir := exportEvidence(t)
    assert.True(t, verifyStandalone(t, evidenceDir))
}
```

---

## Verdict

**All 8 golden path steps are testable with v0.1 architecture.** No step requires v0.5+ components (attack graph, adversarial validation, retesting, AI integration). The golden path is the v0.1 acceptance test — when it passes, v0.1 ships.
