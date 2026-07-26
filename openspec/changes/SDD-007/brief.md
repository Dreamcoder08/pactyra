# SDD-007 Brief: Evidence and Receipts

**Parent:** SDD-000 (Vision & Governance)
**Status:** Brief
**Date:** 2026-07-26

---

## 1. What SDD-007 Owns

SDD-007 implements the **evidence model** — the immutable, cryptographically-signed record of every action in a Pactyra engagement. This is the "Proven by evidence" half of Pactyra's core promise.

### Core Deliverables

1. **Receipt schema**: Fields, types, serialization format.
2. **Hash chain**: Append-only linked list with parent hashing.
3. **Ed25519 signing**: Key generation, signing, verification.
4. **Chain verification**: Full chain walk, tamper detection.
5. **Evidence export**: Standalone, verifiable evidence pack.
6. **Standalone verifier**: `pactyra verify` subcommand.

---

## 2. Requirements Extracted from SDD-000

### From Spec

- **"Produce verifiable evidence automatically"** (JTBD 3): Every tool invocation → receipt. No operator action required.
- **"Export evidence for external verification"** (JTBD 4): Standalone pack verifiable without platform.
- **"Record actions and decisions"** (Golden Path Step 6): Every action (success, block, note) → receipt.
- **"Generate verifiable receipts"** (Golden Path Step 7): Integrity + provenance layers.
- **"Export a reproducible evidence pack"** (Golden Path Step 8): All receipts, scope, public key, VERIFY.md.
- **"Trust Model — Three Layers"** (Spec requirement): Integrity, provenance, reproducibility.
- **"Every action must produce an auditable entry"** (Product Invariant 2): Append-only, no suppression.
- **"Evidence packs must be independently verifiable"** (Product Invariant 8): No platform access required.
- **"100% usable export without platform access"** (v0.1 Success Metric): Standard tools, no proprietary format.

### From Design

- **Receipt Chain component** (Design §3.1): `pactyra-chain` crate with `receipt.rs`, `chain.rs`, `signer.rs`, `export.rs`.
- **IPC Operations** (Design §6.1): `RecordReceipt`, `RecordNote`, `RecordBlocked`, `VerifyChain`, `ExportChain`.
- **Evidence Export Flow** (Design §4.4): Rust produces chain data → Go packages into directory.
- **Hash Chain rationale** (Design §7.6): Simple linked list, not blockchain. O(n) verification sufficient for v0.1.
- **Ed25519 rationale** (Design §7.5): Small keys, fast, deterministic.
- **Evidence Pack format** (Design §9.4): `receipts.jsonl`, `scope.yaml`, `metadata.json`, `public_key.pem`, `VERIFY.md`.

### From Evidence Model Review (T-007)

- Integrity layer: hash chain + Ed25519 signing. Achievable in v0.1.
- Provenance layer: receipt fields capture who/when/tool/authorization. Achievable in v0.1.
- Reproducibility layer: deferred to v0.5. Receipt schema captures data needed for replay.
- Recommendations:
  - Include `schema_version` in every receipt
  - Use RFC 3339 timestamps with sub-second precision
  - Document that timestamp is system-clock-based (not trusted)
  - Include hash algorithm in receipt metadata
  - Gracefully handle: missing files, malformed JSON, truncated chain, wrong public key, future schema versions

---

## 3. Receipt Schema (Architectural)

### Receipt Fields

| Field            | Type     | Required    | Description                                                                 |
| ---------------- | -------- | ----------- | --------------------------------------------------------------------------- |
| `id`             | string   | Yes         | Unique receipt identifier (UUID or monotonic)                               |
| `schema_version` | string   | Yes         | Receipt format version (e.g., "1.0")                                        |
| `type`           | enum     | Yes         | `scope_activation`, `tool_invocation`, `violation_blocked`, `operator_note` |
| `timestamp`      | RFC 3339 | Yes         | When the receipt was created (system clock)                                 |
| `operator`       | string   | Yes         | Operator identity from engagement metadata                                  |
| `scope_hash`     | string   | Yes         | Active scope hash at time of action                                         |
| `parent_hash`    | string   | Yes         | SHA-256 of previous receipt (genesis uses scope_hash)                       |
| `invocation_id`  | string   | Conditional | Present for `tool_invocation` and `violation_blocked`                       |
| `tool`           | object   | Conditional | `{name, version}` — present for tool invocations                            |
| `target`         | object   | Conditional | `{ips: [...], ports: [...]}` — what was targeted                            |
| `parameters`     | object   | Conditional | Tool parameters as validated by Scope Guardian                              |
| `output_hash`    | string   | Conditional | SHA-256 of tool stdout+stderr                                               |
| `exit_code`      | int      | Conditional | Tool exit code (0 = success)                                                |
| `note`           | string   | Conditional | Free-text for `operator_note` type                                          |
| `reason`         | string   | Conditional | Block reason for `violation_blocked` type                                   |
| `authorization`  | string   | Yes         | Permission level at time of action (O0–O6)                                  |
| `signature`      | string   | Yes         | Ed25519 signature over receipt data (excludes signature field)              |
| `hash_algorithm` | string   | Yes         | "SHA-256"                                                                   |

### Receipt Types

| Type                | When Created                   | Key Fields                                          |
| ------------------- | ------------------------------ | --------------------------------------------------- |
| `scope_activation`  | Scope activated for engagement | scope_hash, operator                                |
| `tool_invocation`   | Tool execution completed       | invocation_id, tool, target, output_hash, exit_code |
| `violation_blocked` | Scope Guardian denies action   | invocation_id, target, reason                       |
| `operator_note`     | Operator adds note             | note                                                |

---

## 4. Hash Chain Structure

```
Genesis Receipt (scope_activation)
  parent_hash = scope_hash
  hash = SHA-256(parent_hash || receipt_data_without_signature)
  signed with Ed25519

Receipt #2 (tool_invocation)
  parent_hash = hash(Genesis Receipt)
  hash = SHA-256(parent_hash || receipt_data_without_signature)
  signed with Ed25519

Receipt #3 (violation_blocked)
  parent_hash = hash(Receipt #2)
  ...

Chain verification:
  For each receipt i (starting from genesis):
    Verify signature against public key
    Verify parent_hash == hash(receipt[i-1])
    Verify hash == SHA-256(parent_hash || receipt_data)
  Return {valid, tampered_receipts}
```

---

## 5. Evidence Export Format

```
evidence-YYYY-MM-DD/
├── receipts.jsonl          # One JSON receipt per line (appendable)
├── scope.yaml              # Original scope document
├── metadata.json           # Engagement name, operator, dates, chain_height, root_hash
├── public_key.pem           # Ed25519 public key (PEM format)
└── VERIFY.md               # Human-readable verification instructions
```

**receipts.jsonl**: Each line is a complete JSON receipt object. JSONL chosen because:

- Human-readable with standard tools
- Line-by-line verifiable (each line is self-contained)
- Appendable (new receipts at end)
- No framing or delimiters needed

---

## 6. Standalone Verification (`pactyra verify`)

```bash
$ pactyra verify ./evidence-2026-07-26/
Verifying 47 receipts...
  Chain integrity: VALID
  All receipts signed by: a1b2c3d4...
  Scope hash: e5f6a1b2...
  No tampering detected.
```

The `verify` subcommand:

1. Reads `receipts.jsonl` line by line
2. Reads `public_key.pem`
3. For each receipt: verifies signature, verifies hash chain link
4. Reports results
5. Does NOT require `pactyra-core` to be running (stateless)
6. Works with standard crypto libraries (no proprietary code)

---

## 7. Boundary with SDD-010 (Finding Lifecycle)

**SDD-007 provides:** The raw evidence chain — every action, its parameters, its output hash, its signed receipt.

**SDD-010 provides:** The finding workflow on top — detection → validation → confirmation → remediation → retest. SDD-010 reads receipts, correlates them into findings, and manages finding state transitions.

SDD-007 does NOT define:

- Finding states or transitions
- Operator review workflow
- Report generation from findings
- Retesting orchestration

---

## 8. v0.1 vs. v0.5 Scope

### v0.1 Deliverables

- [x] Receipt creation and signing (Ed25519)
- [x] Append-only hash chain (integrity layer)
- [x] Receipt provenance fields (provenance layer)
- [x] Chain verification (full walk, tamper detection)
- [x] Evidence export (JSONL + metadata + public key + VERIFY.md)
- [x] Standalone verification (`pactyra verify`)

### v0.5 Deliverables

- [ ] Reproducibility layer (evidence replay engine)
- [ ] Receipt diff comparison (original vs. retest)
- [ ] Human-readable procedure documentation
- [ ] Environmental precondition capture
- [ ] Binary chain format (for large engagements)
- [ ] TPM-backed key storage (optional)

---

## 9. Dependencies

| Depends On | Why                                             |
| ---------- | ----------------------------------------------- |
| SDD-000    | Vision, trust model, evidence model, invariants |
| SDD-002    | `scope_hash` is genesis receipt parent_hash     |

| Required By | Why                                        |
| ----------- | ------------------------------------------ |
| SDD-010     | Finding lifecycle reads evidence chain     |
| SDD-015     | Attack graph consumes tool output receipts |

---

## 10. Acceptance Criteria (from SDD-000 Spec)

1. Every tool invocation produces a signed receipt within the chain.
2. Every blocked violation produces a receipt.
3. Chain verification returns `{valid: true}` on a clean chain.
4. Tampering with any receipt is detected (specific receipt identified).
5. Evidence pack exports with all five required files.
6. `pactyra verify` succeeds on a valid pack and fails on a tampered pack.
7. Verification works without platform access (just the pack).
8. No proprietary format, tool, or key required to read/verify the pack.
