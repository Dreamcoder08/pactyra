# Scope Expansion Envelope Review (T-008)

**Date:** 2026-07-26
**Review scope:** Proposal §6.1 (Scope Management), Spec "Scope Expansion Envelope Policy"
**Verdict:** ✅ YAML policy is enforceable. Edge cases documented for SDD-002.

---

## 1. Policy Structure Validation

The expansion envelope policy (proposal §6.1) defines a YAML structure:

```yaml
expansion_policy:
  approver_roles:
    - engagement_manager

  allowed:
    additional_subdomains: true
    discovered_ips_within_authorized_asn: true
    authenticated_test_accounts: true

  limits:
    maximum_new_assets: 25
    maximum_duration_extension_hours: 8
    maximum_rps: 10

  always_requires_customer_approval:
    - new_root_domain
    - new_cidr
    - production_disruption_risk
    - social_engineering
    - data_access_beyond_proof
```

### 1.1 Parseability

The YAML structure is parseable with standard YAML libraries. Each section has a well-defined type:

- `approver_roles`: list of strings
- `allowed`: map of string → boolean
- `limits`: map of string → integer
- `always_requires_customer_approval`: list of strings

**Enforceability assessment:** Each section can be enforced programmatically:

| Section                             | Enforcement Mechanism                                                            |
| ----------------------------------- | -------------------------------------------------------------------------------- |
| `approver_roles`                    | Gate expansion approval on operator role matching one of the listed roles        |
| `allowed`                           | Only enable expansion capabilities where the value is `true`                     |
| `limits`                            | Track counters during engagement; reject expansions that exceed limits           |
| `always_requires_customer_approval` | Hard block — these items cannot be expanded without external (customer) approval |

### 1.2 Enforceability in Code

The policy is enforceable at the Rust Scope Guardian level:

1. **At scope activation:** Parse `expansion_policy` from scope YAML. Store in compiled scope.
2. **At expansion request:** Validate:
   - Is the requesting operator in `approver_roles`?
   - Is the requested capability in `allowed` and set to `true`?
   - Have any `limits` been exceeded?
   - Is the requested capability in `always_requires_customer_approval`? (reject without external approval)
3. **After expansion:** Increment counters, record cryptographic receipt of the expansion.

---

## 2. Edge Cases Identified

### 2.1 Overlapping Expansions

**Scenario:** Two concurrent expansion requests. PM-1 approves expanding to include discovered subnet A. PM-2 simultaneously approves expanding to include discovered subnet B. Both subnets overlap (e.g., A = 10.0.2.0/24, B = 10.0.2.0/25).

**Risk:** Race condition. The compiled scope could end up in an inconsistent state if both expansions are applied concurrently.

**Recommendation for SDD-002:** The Scope Guardian should process expansion requests sequentially (single-writer model). The compiled scope is guarded by a write lock during expansion. The second request sees the updated scope (including the first expansion) and validates against the new state.

### 2.2 Concurrent Expansion Requests

**Scenario:** Operator requests expansion for additional subdomains. Before PM approves, operator requests expansion for discovered IPs. Both are within allowed capabilities.

**Risk:** The operator could batch multiple expansion requests and hope one gets approved. This is not a security issue (PM still approves each), but it could create approval fatigue.

**Recommendation for SDD-002:** Consider batching: the operator submits one expansion request with multiple items. The PM approves or rejects the entire batch. Individual item rejection is a v0.5+ UX refinement.

### 2.3 Expansion Rollback

**Scenario:** PM approves expansion to include a new CIDR range. Operator scans it. Later, the client says "that range was never authorized — remove it." The operator has already produced receipts against that range.

**Risk:** The evidence chain is append-only. Receipts against the expanded range cannot be removed. The engagement's evidence integrity is preserved, but the scope expansion was erroneous.

**Recommendation for SDD-002:** Support a "scope contraction" receipt type. This does not remove prior receipts — it records that the scope was reduced, and receipts after the contraction point reference the new (smaller) scope_hash. The evidence pack includes the full scope history. This preserves integrity while documenting the error.

### 2.4 Expansion After Engagement Close

**Scenario:** The engagement is complete. The operator has exported the evidence pack. Two weeks later, the PM realizes they forgot to approve an expansion that the operator used (the operator worked from an informal approval, not the formal envelope).

**Risk:** The evidence pack shows receipts against targets that were not in the formal scope. A client audit would flag this.

**Recommendation for SDD-002:** After engagement close (explicit operator action), the scope is frozen. No expansions are permitted. If a retroactive expansion is needed, it must be documented as an addendum to the evidence pack, not a modification of the chain. The chain's integrity is preserved; the addendum explains the discrepancy.

### 2.5 Limits Counter Reset

**Scenario:** `maximum_new_assets: 25`. Operator expands scope, uses 10 assets. Later, the client authorizes a new scope document (new engagement phase) with `maximum_new_assets: 25` again.

**Risk:** Ambiguity — do limits reset per engagement phase, per scope document, or per expansion?

**Recommendation for SDD-002:** Limits apply per compiled scope activation. If a new scope document is loaded, it starts fresh counters. Within a single scope activation, counters are monotonic (only increase). Document this explicitly in the scope model.

### 2.6 Implicit Expansions from Tool Discovery

**Scenario:** Nmap discovers a host at 10.0.1.50 that was not explicitly listed in the scope but falls within the authorized CIDR 10.0.1.0/24. Is this an expansion?

**Resolution:** No. The CIDR range already authorizes all 256 hosts. Individual host discovery within an authorized range is not scope expansion — it's expected behavior. The `discovered_ips_within_authorized_asn` capability in the `allowed` section refers to IPs discovered in a _different_ ASN that belongs to the same organization, not individual hosts within already-authorized CIDR ranges. Clarify this distinction in SDD-002.

### 2.7 Multi-Phase Engagements

**Scenario:** A large engagement has Phase 1 (external recon) and Phase 2 (internal network after VPN access). Each phase has different scope. Does the expansion envelope carry over?

**Recommendation for SDD-002:** The expansion envelope is bound to the scope document, not the engagement. Each new scope document has its own expansion envelope. When Phase 2 begins, the operator loads a new scope with a new envelope. Receipts from Phase 1 are in the evidence chain; Phase 2 receipts link to the new scope_hash. The chain is continuous across phases.

---

## 3. v0.1 vs. v0.5 Scope

### v0.1 (Static Scope)

- No dynamic expansion. Scope is loaded once at engagement start.
- The expansion envelope policy structure is defined in the scope YAML but not enforced (no expansion workflow exists to enforce it against).
- Scope validation is basic IP/port gating (O0–O3).

### v0.5 (Dynamic Scope)

- Full expansion envelope enforcement.
- PM approval workflow with cryptographic receipts.
- Scope contraction recording.
- Multi-phase engagement support.

---

## 4. Recommendations for SDD-002

1. **Implement the expansion envelope parser** even in v0.1. The data model should support it, even if the workflow doesn't use it. This avoids a schema migration in v0.5.

2. **Define scope states explicitly:**
   - `draft` — being edited, not yet compiled
   - `compiled` — parsed, validated, has `scope_hash`
   - `active` — bound to an engagement, enforcement active
   - `expanded` — active + at least one expansion applied
   - `contracted` — active + at least one contraction applied
   - `closed` — engagement complete, frozen

3. **Document the distinction between:**
   - **Scope expansion** (adding new targets to the authorized set — requires PM approval within envelope)
   - **Scope refinement** (discovering individual hosts within already-authorized CIDR — no approval needed)
   - **Scope contraction** (removing targets — records a receipt but does not delete prior receipts)

4. **Limit counters should be monotonic within a scope activation.** Document this in the scope model spec.

---

## Verdict

**The expansion envelope policy is enforceable in code.** The YAML structure is well-typed and each section maps to a specific enforcement mechanism. Edge cases are identified and documented for SDD-002 to address.

**No architectural blocker for v0.5 implementation.** The policy structure defined in SDD-000 is sufficient for SDD-002 to implement the scope model with forward compatibility for dynamic expansion.
