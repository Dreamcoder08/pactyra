# Grant Narrative Package

**Date:** 2026-07-26
**Source:** SDD-000 T-021
**Purpose:** Reusable grant narrative adaptable to specific applications with <1 hour of tailoring.

---

## 1. Problem Statement

### The Fragmentation Crisis in Offensive Security

Authorized offensive security assessments — penetration tests, red team engagements, security audits — are essential infrastructure for the digital economy. Every SaaS company, fintech, and healthcare provider relies on these assessments to validate their defenses. Yet the operational tooling that powers these assessments is fragmented, unverifiable, and unsafe.

A typical offensive security engagement today runs across at least six disconnected surfaces: a Kali VM for tools, a terminal multiplexer for session state, a note-taking tool for findings, chat for team coordination, a report template for deliverables, and shared drives for scripts. The operator spends **30–40% of their engagement time on integration labor** — copying output, correlating findings, assembling evidence — rather than finding vulnerabilities.

More critically, there is **no verifiable model of authorization, execution, and evidence.** When a client disputes a finding, the operator's defense is screenshots and memory. When a junior operator runs a scan, one mistyped IP can expand a test into unauthorized territory. When a retest is needed six months later, the methodology must be reconstructed from notes and hope.

This fragmentation is not a tooling gap — it is a structural gap. No existing tool provides an **operations kernel** that binds authorization, execution, and evidence into a single integrity chain.

### Why This Matters Now

Two trends make this gap urgent:

1. **AI agents are entering offensive security.** Without a verifiable control plane, AI-assisted findings are less credible than manual ones. An AI agent that misinterprets scope or hallucinates a target can cause real legal and operational damage — and there is no standard way to prevent it.

2. **Regulatory pressure is increasing.** SOC 2, ISO 27001, and emerging cybersecurity regulations increasingly require evidence of methodology, not just results. "Trust me, I'm a pentester" is no longer sufficient.

---

## 2. Solution Overview

### Pactyra: Open Offensive Operations Kernel

Pactyra is an open-source, terminal-first offensive operations kernel that provides the missing control plane. It is not a scanner, not an exploit framework, not a vulnerability manager. It is the **layer between authorization and execution** that ensures every action is scoped, recorded, and verifiable.

**Bound by scope. Proven by evidence.**

Pactyra's three core capabilities:

1. **Runtime scope enforcement.** Every tool invocation is validated against a cryptographically-bound scope document before execution. Out-of-scope actions are blocked — not warned — before any network traffic is sent.

2. **Immutable evidence receipts.** Every action (tool invocation, blocked attempt, operator note) produces a signed, hash-chained receipt. The chain is append-only. Prior receipts cannot be modified without detection.

3. **Standalone evidence export.** Evidence packs are exportable and verifiable without access to the Pactyra platform. A third party — client, auditor, regulator — can verify the full chain with standard tools.

---

## 3. Public-Interest Thesis

### Why This Is Public Infrastructure

Pactyra is not a commercial product for a specific market segment. It is a **reference implementation** that proves a technical thesis:

> Authorized offensive automation does not have to be opaque, unsafe, or unverifiable. A control plane that enforces scope at runtime, records every action in an append-only evidence chain, and exports independently-verifiable evidence packs is feasible, useful, and should exist as open infrastructure.

The project's public-interest commitments:

- **Open-source.** Self-hostable. No mandatory services, registration, or proprietary dependencies.
- **AI-optional.** Works fully without any AI integration. AI providers are interchangeable accelerators, never a requirement.
- **Provider-agnostic.** No vendor lock-in. The core is and will remain the reference implementation.
- **Commercial guardrails.** Future commercial services must not distort or weaken the open-source core.

### Who Benefits

- **Small offensive security firms** (3–10 operators) that cannot afford enterprise tooling but face the same scope and evidence challenges as large consultancies.
- **Independent operators and researchers** who need verifiable proof of authorized methodology.
- **Clients and regulators** who receive independently-verifiable evidence packs instead of screenshots.
- **The open-source security community** that gains a reference architecture for safe, auditable offensive automation.

---

## 4. Technical Approach

### Architecture

Pactyra uses a **two-binary model** with a hard trust boundary:

- **Rust trust core (`pactyra-core`):** Owns scope validation, receipt signing (Ed25519), and evidence chain integrity. Small, auditable, fuzzable. The signing key never leaves this process.
- **Go orchestrator (`pactyra`):** Operator interface (terminal UI), agent coordination, Docker worker management. Must ask Rust for permission before any tool executes.
- **Docker workers:** Ephemeral, single-use containers with network egress restricted to scope-authorized targets. Defense-in-depth.

### Technology Choices

| Component      | Technology | Rationale                                                 |
| -------------- | ---------- | --------------------------------------------------------- |
| Trust boundary | Rust       | Memory safety, zero-cost abstractions, auditable binaries |
| Orchestration  | Go         | Docker SDK, concurrency (goroutines), fast compile        |
| Tool isolation | Docker     | Universal, network egress control, ephemeral by default   |
| Signing        | Ed25519    | Small keys, deterministic, well-understood                |
| Evidence       | Hash chain | Append-only, O(n) verification, no blockchain association |

### v0.1 Scope — "First Light" (8–12 weeks)

The v0.1 wedge demonstrates the complete golden path:

1. Load and compile an authorized scope
2. Generate a cryptographic scope_hash
3. Execute permitted discovery (Nmap) against a local cyber range
4. Attempt an out-of-scope action
5. Block it before execution
6. Record every action as an evidence receipt
7. Generate verifiable receipts
8. Export a standalone evidence pack

---

## 5. Team Background

_[To be completed with specific maintainer information]_

Pactyra is developed under **Arkelythex**, an open-source collective focused on offensive security tooling and infrastructure.

---

## 6. Milestones

### Phase 1: Foundation (SDD-000) — COMPLETE

- Vision, boundaries, non-goals defined
- System architecture designed (two-binary Rust + Go model)
- Repository initialized, CI/CD configured
- Downstream SDD briefs prepared (scope, evidence, agents, permissions)
- Grant landscape researched

### Phase 2: v0.1 "First Light" (8–12 weeks) — PLANNED

- Scope model and compiler (SDD-002)
- Tool and agent contracts (SDD-005)
- Permission lattice O0–O3 (SDD-003)
- Evidence and receipts (SDD-007)
- Operator TUI (SDD-0XX)
- **Deliverable:** Public repository where anyone can clone, build, run the golden path, and verify the evidence pack independently.

### Phase 3: v0.5 "Community Ready" (6–9 months after v0.1)

- Attack graph, adversarial validation, retesting
- Tool adapter SDK (community contributions)
- Permission lattice O0–O6
- Design partner program (5+ teams using on real engagements)
- **Deliverable:** Pactyra usable by external operators on real authorized engagements.

### Phase 4: v1.0 "Reference Implementation" (12–18 months after v0.1)

- Web dashboard, diff engine, 20+ supported tools
- Active contributor base, published benchmarks
- Conference presence, grant funding
- **Deliverable:** Recognized open-source reference implementation for authorized offensive automation.

---

## 7. Budget Estimate

### v0.1 Development (8–12 weeks)

| Category                     | Estimate          | Notes                                      |
| ---------------------------- | ----------------- | ------------------------------------------ |
| Core development (Rust + Go) | 400–600 hours     | Two experienced systems developers         |
| Testing + cyber range        | 80–120 hours      | Integration tests, CI, fuzzing             |
| Documentation                | 40–60 hours       | README, golden path walkthrough, VERIFY.md |
| Design partner outreach      | 20–30 hours       | Recruitment, onboarding, feedback calls    |
| **Total**                    | **540–810 hours** |                                            |

### v0.5 Development (6–9 months)

| Category           | Estimate            | Notes                                            |
| ------------------ | ------------------- | ------------------------------------------------ |
| Core development   | 800–1,200 hours     | Attack graph, validation, retesting, SDK         |
| Community building | 100–200 hours       | Conferences, benchmarks, contribution guidelines |
| **Total**          | **900–1,400 hours** |                                                  |

---

## 8. Why Open Source

Pactyra is open-source because:

1. **Trust requires transparency.** You cannot verify a scope enforcement engine you cannot read. The Rust trust core must be auditable.

2. **Community multiplies value.** A reference implementation gains credibility when independent operators reproduce benchmarks and contribute adapters.

3. **No gatekeepers.** Offensive security teams should not need to trust a vendor's word that scope was respected. They should be able to verify it themselves.

4. **Public infrastructure.** Like OpenSSL, like Let's Encrypt, like the Linux kernel — some software is too important to the ecosystem to be proprietary. Authorized offensive security infrastructure is one of those categories.

---

## 9. Technical Appendix

### Architecture Diagram

_(Include the full system architecture diagram from design §2)_

### Trust Boundary Model

_(Include the trust domain diagram from design §5.4)_

### Golden Path Walkthrough

_(Reference docs/golden-path.md — 8-step walkthrough with component mapping and data flow)_

### Evidence Model

_(Reference docs/evidence-model-review.md — integrity, provenance, reproducibility layers)_

---

## 10. Adaptation Guide

This narrative can be tailored to specific grant applications:

| Grant Program      | Emphasize                                                     | De-emphasize                                             |
| ------------------ | ------------------------------------------------------------- | -------------------------------------------------------- |
| NLNet NGI Zero     | Open infrastructure, internet health, security as public good | Commercial services, team size                           |
| OTF                | Digital security, human rights defenders, safety by design    | Exploitation capabilities                                |
| STF                | Digital infrastructure, supply chain security, maintenance    | Early stage (emphasize architectural maturity)           |
| GitHub Accelerator | Community, open-source sustainability, contributor growth     | Grant dependency                                         |
| Mozilla MOSS       | Open internet, transparency, user empowerment                 | Offensive framing (use "authorized security assessment") |

---

## Appendix: Key Terminology

| Term                  | Definition                                                                                            |
| --------------------- | ----------------------------------------------------------------------------------------------------- |
| **Operations kernel** | The system that manages authorization, orchestration, and evidence for offensive operations           |
| **Scope document**    | Machine-readable specification of authorized targets, techniques, and constraints                     |
| **Evidence receipt**  | Immutable, cryptographically-signed record of an action                                               |
| **Evidence chain**    | Append-only linked list of receipts; each receipt hashes its parent                                   |
| **scope_hash**        | Cryptographic hash of the compiled scope; root of the evidence chain                                  |
| **Golden path**       | The 8-step v0.1 demonstration: scope → hash → discover → violate → block → record → receipts → export |
