# SDD-003 Brief: Permission Lattice O0–O6

**Parent:** SDD-000 (Vision & Governance)
**Status:** Brief
**Date:** 2026-07-26

---

## 1. What SDD-003 Owns

SDD-003 implements the **permission lattice** — a seven-level trust hierarchy (O0–O6) that gates what actions an operator can perform at each level. The lattice is the capability model layered on top of the scope model (SDD-002). Scope says "what targets are authorized." Permission lattice says "what actions are authorized."

### Core Deliverables

1. **O0–O6 level definitions**: Capabilities, restrictions, escalation rules.
2. **Enforcement rules**: Which levels can do what. Level gating at Scope Guardian.
3. **Escalation procedure**: How an operator moves up (and down) the lattice.
4. **v0.1 scope**: O0–O3 (basic IP/port gating).
5. **v0.5 scope**: O0–O6 (full lattice with exploitation capabilities).

---

## 2. Requirements Extracted from SDD-000

### From Proposal

- **"O0–O6 permission lattice"** (Proposal §6.1): Escalating trust levels. O0 = read-only recon from external perspective. O6 = full exploitation with credential-based lateral movement.
- **"The operator's current level gates what tools and targets are available"** (Proposal §6.1).
- **"No AI agent can elevate its own permissions"** (Permanent Non-Goal 6): Escalation requires explicit operator action.
- **"Agents operate within the operator's current permission level"** (Proposal §5.1).

### From Spec

- **"v0.1 only needs O0–O3"** (Spec: v0.1 Scope Boundaries): Basic IP/port gating.
- **"Permission lattice beyond O0–O3 deferred to v0.5"** (Spec: v0.1 Scope Boundaries).
- **"No agent, AI model, or automated system SHALL be capable of elevating its own permissions"** (Product Invariant 3).

### From Design

- **Scope Guardian** (Design §3.1): Makes allow/deny decisions. Adding O4–O6 rules is an extension of decision logic.
- **Future-proofing** (Design §13): v0.1 architecture enables permission lattice without redesign.

---

## 3. Permission Level Definitions

### O0 — Read-Only Reconnaissance (External)

**What it allows:**

- Passive DNS lookups
- WHOIS queries
- SSL certificate inspection
- Shodan/censys queries (external data sources)

**What it blocks:**

- Any active network traffic to targets
- Port scanning
- Service detection
- Any tool that sends packets

**Use case:** PM loads scope and validates targets before authorizing active scanning.

### O1 — Active Reconnaissance (Basic)

**What it allows (includes O0):**

- TCP connect() scans on authorized ports
- Service version detection (Nmap -sV)
- OS detection (Nmap -O)
- HTTP/HTTPS requests to authorized targets

**What it blocks:**

- SYN scans, UDP scans, aggressive timing
- Authentication attempts
- Exploitation of any kind
- Traffic to unauthorized ports

**Use case:** Initial discovery phase. Operator maps the attack surface.

### O2 — Active Reconnaissance (Full)

**What it allows (includes O0–O1):**

- SYN scans, UDP scans, all Nmap scan types
- Aggressive timing templates
- Script scanning (Nmap -sC, safe scripts only)
- Full port range scanning on authorized targets

**What it blocks:**

- Authentication attempts (brute force, credential stuffing)
- Exploitation scripts (Nmap vuln category)
- Any tool that modifies target state

**Use case:** Full surface mapping. Operator has complete visibility of authorized targets.

### O3 — Limited Interaction

**What it allows (includes O0–O2):**

- HTTP POST/PUT/DELETE to authorized targets
- Authentication attempts with test credentials
- Directory enumeration (Gobuster, dirb)
- Vulnerability scanning (Nuclei, safe templates)

**What it blocks:**

- Exploitation payloads
- Credential reuse across targets
- Lateral movement
- Any action that could crash or disrupt a service

**Use case:** Operator validates findings and tests for known vulnerabilities without risk of disruption.

### O4 — Exploitation (Controlled)

**What it allows (includes O0–O3):**

- Single-target exploitation with proof-of-exploitation recording
- Metasploit modules (non-destructive payloads)
- Custom exploit scripts within container isolation
- Credential validation (not reuse)

**What it blocks:**

- Multi-target exploitation
- Lateral movement
- Persistence mechanisms
- Data exfiltration beyond proof

**v0.5+**

### O5 — Exploitation (Full)

**What it allows (includes O0–O4):**

- Multi-target exploitation campaigns
- Credential-based lateral movement
- Pivot chaining
- Custom payload execution

**What it blocks:**

- Persistence mechanisms
- C2 infrastructure deployment
- Data exfiltration beyond proof

**v0.5+**

### O6 — Full Engagement

**What it allows (includes O0–O5):**

- Everything authorized in the engagement scope
- Persistence testing
- Data exfiltration testing (controlled)
- C2 simulation
- Social engineering (if authorized in scope)

**What it blocks:**

- Nothing within scope. This is the operator's maximum authorization.

**v0.5+**

---

## 4. Enforcement Rules

### Level Gating at Scope Guardian

The Rust Scope Guardian evaluates two questions for every action:

1. **Scope check** (SDD-002): "Is target X:port Y within the authorized scope?"
2. **Level check** (SDD-003): "Is action type Z permitted at the current permission level?"

Both must return `allowed` for the action to proceed.

### Level → Action Matrix

| Action                    | O0  | O1  | O2  | O3  | O4  | O5  | O6  |
| ------------------------- | --- | --- | --- | --- | --- | --- | --- |
| DNS lookup                | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  |
| TCP connect scan          | ❌  | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  |
| SYN scan                  | ❌  | ❌  | ✅  | ✅  | ✅  | ✅  | ✅  |
| Service detection         | ❌  | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  |
| OS detection              | ❌  | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  |
| HTTP GET/HEAD             | ❌  | ✅  | ✅  | ✅  | ✅  | ✅  | ✅  |
| HTTP POST/PUT/DELETE      | ❌  | ❌  | ❌  | ✅  | ✅  | ✅  | ✅  |
| Directory enumeration     | ❌  | ❌  | ❌  | ✅  | ✅  | ✅  | ✅  |
| Vulnerability scanning    | ❌  | ❌  | ❌  | ✅  | ✅  | ✅  | ✅  |
| Auth attempt (test creds) | ❌  | ❌  | ❌  | ✅  | ✅  | ✅  | ✅  |
| Single exploitation       | ❌  | ❌  | ❌  | ❌  | ✅  | ✅  | ✅  |
| Multi exploitation        | ❌  | ❌  | ❌  | ❌  | ❌  | ✅  | ✅  |
| Lateral movement          | ❌  | ❌  | ❌  | ❌  | ❌  | ✅  | ✅  |
| Persistence testing       | ❌  | ❌  | ❌  | ❌  | ❌  | ❌  | ✅  |

---

## 5. Escalation Procedure

### Level Transition Rules

1. **Escalation:** Operator must explicitly request level escalation. The request is recorded as a receipt.
2. **Approval:** For levels O4+, escalation requires PM approval (within scope expansion envelope policy).
3. **De-escalation:** Operator can voluntarily drop to a lower level at any time. Recorded as a receipt.
4. **Automatic de-escalation:** Not implemented in v0.1. Consider for v0.5: time-based (session timeout drops to O0).

### v0.1 Behavior (Simplified)

- O0–O3 are the only levels in v0.1.
- O0: default after scope activation. Operator must explicitly escalate.
- O1–O3: operator escalates at will (no PM approval needed for O0–O3).
- Every escalation/de-escalation produces a receipt.

---

## 6. Boundary with SDD-002 (Scope Model)

**SDD-002 checks:** "Is `10.0.1.50:443` within the authorized scope?" → target validation.
**SDD-003 checks:** "Is the operator's current level (O2) allowed to run SYN scans?" → capability validation.

Both gates must pass. The Scope Guardian evaluates both in a single `ValidateInvocation` call.

SDD-002 provides:

- `ValidateTarget(ip, port) → bool`

SDD-003 adds:

- `ValidateAction(action_type, level) → bool`

The combined gate: `allowed = ValidateTarget(ip, port) && ValidateAction(action_type, level)`

---

## 7. v0.1 vs. v0.5 Scope

### v0.1 Deliverables

- [x] O0–O3 level definitions
- [x] Level gating at Scope Guardian (combined with scope check)
- [x] Escalation/de-escalation receipt recording
- [x] Operator-initiated level changes (no PM approval required for O0–O3)

### v0.5 Deliverables

- [ ] O4–O6 level definitions
- [ ] PM approval workflow for O4+ escalation
- [ ] Level-specific tool restrictions (not just port gating)
- [ ] Automatic de-escalation (timeout, session end)
- [ ] Level inheritance for multi-operator teams

---

## 8. Dependencies

| Depends On | Why                                               |
| ---------- | ------------------------------------------------- |
| SDD-000    | Vision, permission lattice concept, non-goals     |
| SDD-002    | Scope Guardian combines scope check + level check |

| Required By | Why                                       |
| ----------- | ----------------------------------------- |
| SDD-0XX     | TUI displays current permission level     |
| SDD-005     | Agent proposals gated by permission level |

---

## 9. Acceptance Criteria (from SDD-000 Spec)

1. O0 operator cannot execute port scans (blocked at Scope Guardian).
2. O2 operator can run SYN scans on authorized targets.
3. O3 operator can run vulnerability scans on authorized targets.
4. Level escalation produces a receipt.
5. Level de-escalation produces a receipt.
6. No agent, AI or otherwise, can elevate its own level.
7. Level check and scope check are independent but both must pass.
