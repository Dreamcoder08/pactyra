# Pactyra Platform Specification

**Status:** Reviewed  
**Date:** 2026-07-26  
**Change ID:** SDD-000

## Purpose

Define the vision, boundaries, actors, value proposition, trust model, product invariants, success metrics, and scope for **Pactyra** — an open-source, terminal-first offensive operations kernel developed under Arkelythex. This specification establishes WHAT Pactyra is, WHO it serves, WHAT it must guarantee, and WHAT it explicitly does NOT do. It does not specify concrete schemas, wire formats, agent contracts, or UI layout; those are owned by dedicated SDD changes.

---

## Requirements

### Requirement: Platform Identity

Pactyra MUST present itself as an **open offensive operations kernel** — not a scanner, not an exploit framework, not a vulnerability manager.

#### Scenario: Public identity is consistent

- GIVEN the Pactyra repository and documentation
- WHEN a visitor encounters any public-facing surface (README, website, CLI help, conference abstract)
- THEN the canonical descriptor "Open Offensive Operations Kernel" is used
- AND the tagline "Bound by scope. Proven by evidence." appears
- AND the parent organization "Arkelythex" is credited
- AND the CLI binary is named `pactyra`
- AND the name is never presented as "Pactyra OS," "Pactyra Offensive Security Platform," or any hyphenated compound

#### Scenario: Repository model is transparent

- GIVEN the project's public repository
- WHEN an external developer evaluates Pactyra
- THEN the README clearly states: open-source (exact license TBD), self-hostable, terminal-first, AI-optional, provider-agnostic
- AND no mandatory services, registration, or proprietary dependencies are required to build or run

---

### Requirement: Actors and Personas

Pactyra MUST recognize distinct actor roles with defined responsibilities. The platform MUST serve the **technical offensive security operator** as its primary user.

#### Scenario: Actor roles are defined

- GIVEN the Pactyra platform documentation
- WHEN an operator, PM, or contributor reads the actors section
- THEN the following roles are enumerated with responsibilities:
  - **Operator** (primary user): executes offensive operations, runs tools, discovers and validates findings, produces evidence
  - **Engagement Manager / PM**: manages scope documents, approves scope expansions within pre-authorized boundaries, receives exported reports
  - **Client**: receives the final deliverable (report + evidence pack); never directly interacts with the platform in v0.1
  - **Reviewer / Junior Operator**: operates within enforced scope; work is reviewed via the shared evidence chain
  - **Auditor / Third Party**: verifies evidence packs independently without platform access
  - **Contributor**: extends tool adapters, agents, or integrations via the SDK

#### Scenario: Primary user is the technical operator

- GIVEN the Pactyra platform
- WHEN any design decision involves a tradeoff between operator experience and another role
- THEN the operator's needs take priority
- AND the platform NEVER requires operator actions that serve only administrative, compliance, or reporting purposes without direct operational value

#### Scenario: Client never touches the platform in v0.1

- GIVEN v0.1 of Pactyra
- WHEN an engagement produces deliverables
- THEN the client receives exported artifacts (report, evidence pack) through the PM
- AND there is NO client login, client dashboard, or client-facing interface in the platform

---

### Requirement: Primary User Profile — Alex Chen

The primary user of Pactyra MUST be a **technical offensive security operator** with the profile, pains, and needs of Alex Chen as defined in this requirement.

#### Scenario: Alex Chen profile is the design anchor

- GIVEN the Pactyra product requirements
- WHEN any feature, workflow, or interaction is designed
- THEN it MUST address one or more of the following pains:
  - Integration labor consuming 30–40% of engagement time
  - Scope violation risk from manual IP/domain entry
  - Evidence challenges where screenshots fail to prove findings
  - Handoff friction requiring knowledge-transfer meetings
  - Retesting archaeology requiring 4–6 hours of reconstruction
- OR it MUST enable one or more of the following needs:
  - Runtime scope enforcement that blocks out-of-scope actions
  - Automatic evidence recording without operator effort
  - A shared, queryable operational state
  - One-command retesting from evidence replay
  - Reports generated as views over evidence, not separate artifacts

#### Scenario: Alex's firm context constrains scope

- GIVEN Alex's firm (Aegis Security Partners: 7-person boutique, 3–4 engagements/month, mid-market SaaS/fintech clients)
- WHEN platform capabilities are prioritized
- THEN multi-tenant enterprise features, large-team coordination, and regulatory-audit workflows SHALL be deferred
- AND single-operator, single-engagement workflows SHALL ship first

---

### Requirement: Jobs to Be Done (JTBD)

The platform MUST address the following jobs to be done, framed from the operator's perspective.

#### Scenario: JTBD 1 — Load and enforce scope

- GIVEN an operator receives an authorized scope document from the PM
- WHEN the operator loads the scope into Pactyra
- THEN the platform MUST parse the scope into a machine-enforceable form
- AND every subsequent action MUST be validated against that scope before execution
- AND out-of-scope actions MUST be blocked, not warned

#### Scenario: JTBD 2 — Execute tooling within enforced boundaries

- GIVEN an active scope is loaded
- WHEN the operator invokes a tool (Nmap, Nuclei, custom script)
- THEN the platform MUST execute that tool in an isolated environment
- AND the tool MUST receive only scope-authorized parameters
- AND the operator MUST NOT be able to bypass scope by running tools directly

#### Scenario: JTBD 3 — Produce verifiable evidence automatically

- GIVEN the operator is conducting an engagement
- WHEN any tool invocation completes
- THEN the platform MUST produce an immutable evidence receipt
- AND the receipt MUST be created automatically, without operator action
- AND the receipt MUST be append-only — prior receipts cannot be modified

#### Scenario: JTBD 4 — Export evidence for external verification

- GIVEN an engagement has produced evidence receipts
- WHEN the operator or PM exports the engagement
- THEN the platform MUST produce a standalone evidence pack
- AND the pack MUST be verifiable without access to the Pactyra platform
- AND the pack MUST contain all receipts, the scope document, and the verification instructions

#### Scenario: JTBD 5 — Reproduce methodology for retesting

- GIVEN a previous engagement's evidence pack
- WHEN a retest is requested (v0.5+)
- THEN the platform MUST be capable of replaying the exact tool invocations from the evidence chain
- AND new evidence receipts MUST be generated for comparison

#### Scenario: JTBD 6 — Maintain operational continuity across sessions

- GIVEN the operator disconnects or switches machines
- WHEN the operator resumes the engagement
- THEN all prior actions, receipts, and operational state MUST be recoverable
- AND there MUST be no loss of evidence or scope binding

---

### Requirement: Core Journey — v0.1 Golden Path

Pactyra v0.1 MUST demonstrate the complete golden path: an 8-step sequence from loading scope through exporting a verifiable evidence pack. This journey is the acceptance test for v0.1 "First Light."

#### Scenario: Golden Path Step 1 — Load and compile an authorized scope

- GIVEN a scope document in a supported format (CIDR notation, domain list, excluded hosts)
- WHEN the operator loads the scope into Pactyra
- THEN the platform MUST parse and validate the scope
- AND the scope MUST be compiled into an internal, machine-enforceable representation
- AND invalid or malformed scope entries MUST be rejected with clear error messages

#### Scenario: Golden Path Step 2 — Generate a scope_hash

- GIVEN a validated, compiled scope
- WHEN the scope is bound to the engagement
- THEN the platform MUST produce a cryptographic `scope_hash`
- AND the `scope_hash` MUST be the root of the evidence chain
- AND every subsequent action MUST reference this `scope_hash`

#### Scenario: Golden Path Step 3 — Execute permitted discovery

- GIVEN an active scope and its `scope_hash`
- WHEN the operator invokes a permitted discovery action against a local cyber range
- THEN the platform MUST validate the target against the active scope
- AND the action MUST execute in an isolated worker
- AND an evidence receipt MUST be generated

#### Scenario: Golden Path Step 4 — Attempt an out-of-scope action

- GIVEN an active scope that authorizes a limited target set
- WHEN the operator (or a test script) deliberately attempts an action against a target NOT in the scope
- THEN the platform MUST detect the scope violation

#### Scenario: Golden Path Step 5 — Block before execution

- GIVEN a detected out-of-scope action
- WHEN the platform evaluates the action against the active scope
- THEN the platform MUST block the action BEFORE any network traffic is sent
- AND the platform MUST record the blocked attempt as an evidence receipt
- AND the platform MUST NOT issue a warning and proceed — it MUST block

#### Scenario: Golden Path Step 6 — Record actions and decisions

- GIVEN the engagement is in progress
- WHEN any action (successful tool invocation, blocked violation, operator note) occurs
- THEN the platform MUST produce an evidence receipt
- AND the receipt MUST be appended to the evidence chain
- AND the chain's integrity MUST be verifiable at any point

#### Scenario: Golden Path Step 7 — Generate verifiable receipts

- GIVEN a completed set of actions in the evidence chain
- WHEN receipts are verified
- THEN any receipt in the chain MUST be independently verifiable for integrity (content unchanged since registration)
- AND any receipt MUST be traceable for provenance (who, when, tool, authorization)
- AND the full chain MUST be reproducible (another operator can follow the documented procedure)

#### Scenario: Golden Path Step 8 — Export a reproducible evidence pack

- GIVEN a completed engagement with an evidence chain
- WHEN the operator exports the engagement
- THEN the platform MUST produce a standalone evidence pack
- AND the pack MUST contain all receipts, the scope document, and verification metadata
- AND the pack MUST be readable and verifiable without access to the Pactyra platform
- AND a third party MUST be able to independently verify the evidence chain from the exported pack

---

### Requirement: Three-Tier Value Proposition

Pactyra's value proposition MUST be structured as an ordered three-tier architecture.

#### Scenario: Primary pain is evidence challenges

- GIVEN the Pactyra value proposition
- WHEN the platform is presented to a potential user
- THEN the primary pain MUST be stated as: "Offensive results are disputable because evidence is fragmented, incomplete, or does not verifiably demonstrate what occurred"
- AND this pain MUST be the lead message in any pitch, README, or conference abstract

#### Scenario: Trust mechanism binds every artifact

- GIVEN the Pactyra value proposition
- WHEN describing HOW Pactyra solves the evidence problem
- THEN the trust mechanism MUST be stated as: "Pactyra binds every action and piece of evidence to a verifiable scope, identity, tool, version, authorization, and timeline"
- AND this mechanism MUST be demonstrable in v0.1

#### Scenario: Aha moment is the first blocked violation

- GIVEN a new operator using Pactyra for the first time
- WHEN the platform blocks an out-of-scope action before execution
- THEN the operator MUST recognize this as protection from contractual, legal, or reputational risk
- AND this moment MUST be the target "aha moment" for onboarding and demonstration

#### Scenario: Value architecture leads to adoption and retention

- GIVEN the three-tier value architecture
- WHEN a team evaluates Pactyra
- THEN the adoption driver MUST be credibility defense (evidence challenges)
- AND the retention driver MUST be risk reduction (scope enforcement + evidence automation)

---

### Requirement: Trust Model — Three Layers

Pactyra's trust model MUST provide three independently useful layers: integrity, provenance, and reproducibility. The platform MUST NOT promise "absolute cryptographic proof that a vulnerability exists."

#### Scenario: Integrity layer guarantees content immutability

- GIVEN an evidence receipt in the chain
- WHEN the receipt is verified
- THEN the platform MUST guarantee that the content has not changed since registration
- AND verification MUST be possible without trusting the platform
- AND tampering with any prior receipt MUST be detectable via hash chain breakage
- AND the integrity guarantee MUST NOT extend to claims about whether the original capture was accurate or complete

#### Scenario: Provenance layer guarantees traceability

- GIVEN an evidence receipt in the chain
- WHEN the receipt is inspected
- THEN the receipt MUST record: who performed the action, when it occurred, against which target, with which tool, under which authorization
- AND the full lineage from scope → action → receipt MUST be traceable
- AND the provenance guarantee MUST NOT extend to claims about whether the operator's methodology was correct

#### Scenario: Reproducibility layer guarantees repeatability

- GIVEN an evidence receipt with a documented procedure
- WHEN another authorized operator follows that procedure
- THEN the procedure MUST be sufficiently detailed to repeat the validation
- AND the results of the repetition MUST be comparable to the original
- AND the reproducibility guarantee MUST NOT extend to claims that the target will respond identically (environments change)

#### Scenario: Trust model does not overpromise

- GIVEN the Pactyra trust model documentation
- WHEN the evidence guarantees are stated
- THEN the platform MUST explicitly state what each layer does NOT guarantee
- AND the platform MUST NOT use language implying "absolute proof" or "mathematical certainty of exploitability"

---

### Requirement: Product Invariants

Every feature, change, or extension to Pactyra MUST satisfy these non-negotiable product invariants.

#### Scenario: Invariant 1 — Every action must be bound to an active scope

- GIVEN the Pactyra platform at any version
- WHEN any action (tool invocation, agent decision, operator command) is executed
- THEN that action MUST be validated against an active, cryptographically-bound scope
- AND there MUST be NO mode, configuration, or bypass that allows unscoped actions
- AND the action MUST reference the current `scope_hash`

#### Scenario: Invariant 2 — Every action must produce an auditable entry

- GIVEN the Pactyra platform at any version
- WHEN any action completes (success, failure, or block)
- THEN the platform MUST produce an evidence receipt
- AND the receipt MUST be append-only — prior receipts cannot be modified or deleted
- AND there MUST be NO mode, configuration, or bypass that suppresses receipt generation

#### Scenario: Invariant 3 — The operator is in control

- GIVEN the Pactyra platform at any version
- WHEN an agent or automated system proposes an action
- THEN the operator MUST explicitly authorize that action before execution
- AND no agent, AI model, or automated system SHALL be capable of elevating its own permissions
- AND no finding SHALL be published, exported, or reported without operator review

#### Scenario: Invariant 4 — The terminal is the primary interface

- GIVEN the Pactyra platform at any version
- WHEN an operator interacts with the platform
- THEN the terminal (CLI/TUI) MUST be the primary control surface
- AND any web dashboard, GUI, or external interface MUST be a secondary view, not a replacement for terminal operations
- AND the terminal MUST support the full operational workflow without requiring a browser

#### Scenario: Invariant 5 — AI is optional, never required

- GIVEN the Pactyra platform at any version
- WHEN the platform is built, configured, or run
- THEN the platform MUST operate fully without any AI integration
- AND AI provider integrations (OpenAI, Anthropic, others) MUST be interchangeable
- AND no AI provider SHALL be a hard dependency for any core capability

#### Scenario: Invariant 6 — Scope enforcement is blocking, not advisory

- GIVEN the Pactyra platform at any version
- WHEN a tool invocation targets an out-of-scope resource
- THEN the platform MUST block execution before any network traffic is sent
- AND the platform MUST record the blocked attempt
- AND the platform MUST NOT warn and proceed

#### Scenario: Invariant 7 — Commercial services must not distort the open-source core

- GIVEN the Pactyra platform at any version
- WHEN commercial services, hosted offerings, or enterprise features are developed
- THEN the self-hostable, AI-optional, provider-agnostic core MUST remain fully functional
- AND no feature SHALL be removed from the core to create commercial differentiation
- AND the core MUST remain the reference implementation

#### Scenario: Invariant 8 — Evidence packs must be independently verifiable

- GIVEN the Pactyra platform at any version that produces evidence
- WHEN an evidence pack is exported
- THEN the pack MUST be verifiable without access to the Pactyra platform
- AND the verification procedure MUST be documented and executable with standard tools
- AND the pack MUST contain everything needed for verification (receipts, scope, metadata)

---

### Requirement: v0.1 Success Metrics (Hard Exit Criteria)

Pactyra v0.1 MUST satisfy all seven non-negotiable metrics before shipping. Any metric below threshold SHALL block the v0.1 release.

#### Scenario: Zero actions executed without valid scope

- GIVEN the v0.1 release candidate
- WHEN the golden path and all test scenarios are executed
- THEN there MUST be zero actions that reached execution without passing scope validation
- AND this MUST be verified by automated tests, not manual inspection

#### Scenario: Zero known violations not blocked

- GIVEN a set of deliberate out-of-scope attempts covering CIDR boundaries, domain exclusions, and port restrictions
- WHEN each attempt is made against the scope enforcement engine
- THEN every attempt MUST be blocked before execution
- AND there MUST be zero cases where a known violation was warned instead of blocked

#### Scenario: Zero artifacts without verifiable provenance

- GIVEN the v0.1 release candidate
- WHEN the evidence chain from any test engagement is inspected
- THEN every artifact (tool output, receipt, export) MUST trace to its producing action, operator, tool, timestamp, and authorization
- AND there MUST be zero artifacts with missing or incomplete provenance

#### Scenario: Zero confirmed findings without reproducible evidence

- GIVEN the v0.1 release candidate
- WHEN a finding is included in the exported report
- THEN that finding MUST be backed by at least one evidence receipt that captured the tool invocation and its output
- AND the evidence MUST be sufficient for another operator to reproduce the detection

#### Scenario: One hundred percent of actions bound to scope_hash

- GIVEN any engagement run under Pactyra v0.1
- WHEN the evidence chain is audited
- THEN every action receipt MUST contain the active `scope_hash`
- AND there MUST be zero actions that exist outside the scope binding

#### Scenario: One hundred percent usable export without platform access

- GIVEN an evidence pack exported from Pactyra v0.1
- WHEN a third party with no access to the Pactyra platform attempts to read and verify the pack
- THEN the pack MUST be fully readable with standard tools
- AND the verification procedure MUST succeed without contacting any Pactyra service
- AND no proprietary format, tool, or key MUST be required to read or verify the pack

#### Scenario: Zero expansions without valid approval

- GIVEN the v0.1 scope model (basic IP/port gating, no dynamic expansion in v0.1)
- WHEN any scope modification is applied
- THEN the modification MUST be cryptographically recorded
- AND there MUST be zero cases where scope was modified without a valid, recorded approval

---

### Requirement: v0.5 and v1.0 Success Metrics (Targets)

Beyond v0.1, Pactyra SHALL target specific efficiency, adoption, and community metrics as aspirational goals.

#### Scenario: Operator efficiency targets (v0.5)

- GIVEN operators using Pactyra v0.5 on real engagements
- WHEN efficiency is measured against the Alex Chen baseline
- THEN integration labor ratio SHOULD reduce from 30–40% to ≤15% of engagement time
- AND report assembly time SHOULD reduce from 4–8 hours to ≤1 hour
- AND retesting preparation time SHOULD reduce from 4–6 hours to ≤30 minutes
- AND scope violation incidents SHOULD reach 0 (prevention, not mitigation)

#### Scenario: Adoption and community targets (v1.0)

- GIVEN Pactyra v1.0 has been released
- WHEN community health is assessed
- THEN GitHub stars SHOULD exceed 500
- AND active contributors SHOULD exceed 15 unique contributors per month
- AND public benchmarks SHOULD have been independently reproduced by 3+ external parties
- AND design partner teams SHOULD exceed 5 teams using Pactyra on real engagements
- AND grant or accelerator applications SHOULD number 3+ submitted with at least 1 accepted

---

### Requirement: v0.1 Scope Boundaries

Pactyra v0.1 MUST deliver the wedge: scope enforcement, evidence receipts, and professional export. Capabilities beyond this wedge MUST be deferred to v0.5 or later.

#### Scenario: v0.1 includes scope ingestion, enforcement, and tool isolation

- GIVEN the v0.1 release
- WHEN the platform is built and run
- THEN the following capabilities MUST be present:
  - Scope ingestion: parse CIDR notation and domain lists
  - Runtime scope enforcement: basic IP/port gating (O0–O3)
  - Isolated tool execution: Docker workers for Nmap
  - Evidence receipts: append-only, hash-chained
  - Standalone evidence export: verifiable without platform access
  - Single-operator mode
  - Terminal-first operator interface (TUI)
  - Multi-agent architecture: Scope Guardian + Recon Planner
  - Basic report generation: structured evidence dump
  - Cryptographic receipt verification (standalone)

#### Scenario: v0.1 explicitly excludes attack graph, validation, retesting, and team features

- GIVEN the v0.1 release
- WHEN the platform is evaluated
- THEN the following capabilities MUST NOT be present:
  - Attack graph (deferred to v0.5)
  - Adversarial validation / proof-of-exploitation (deferred to v0.5)
  - Retesting from evidence replay (deferred to v0.5)
  - Team collaboration / shared graph (deferred to v0.5)
  - Dynamic scope expansion workflow (deferred to v0.5)
  - Permission lattice beyond O0–O3 (deferred to v0.5)
  - Web dashboard (deferred to v1.0 or later)
  - AI integration of any kind (deferred to v0.5, always optional)
  - Tool adapter SDK (deferred to v0.5)
  - Cartographer, Validator, Report Generator, Audit Agent (deferred to v0.5/v1.0)

#### Scenario: v0.1 supports exactly one tool — Nmap

- GIVEN the v0.1 release
- WHEN tool execution is invoked
- THEN Nmap MUST be the only supported external tool
- AND the Nmap adapter MUST validate parameters against scope before execution
- AND the Nmap adapter MUST parse output into structured evidence receipts

---

### Requirement: Permanent Non-Goals

Pactyra MUST NOT pursue the following capabilities or characteristics. These are permanent boundaries, not version deferrals.

#### Scenario: Non-goal boundaries are enforced

- GIVEN any version of Pactyra
- WHEN a feature, integration, or capability is proposed
- THEN the following MUST be rejected if they conflict with any permanent non-goal:
  1. Pactyra is not a Linux distribution; it runs on existing Linux installations
  2. Pactyra does not reimplement Burp Suite, BloodHound, Nmap, or Nuclei functionality
  3. Pactyra requires a cryptographically-bound scope document; there is no "unscoped mode"
  4. Pactyra does not generate or host malware, C2, or persistence infrastructure
  5. Pactyra does not distribute, sell, or broker exploits
  6. No AI agent can elevate its own permissions; escalation requires explicit operator action
  7. No autonomous exploitation without supervision; agents propose, operators dispose
  8. Pactyra is not a defensive SIEM, EDR, or SOAR platform
  9. No customer portal in v0.1; client-facing portals are post-v1.0 considerations
  10. Pactyra is not multi-tenant enterprise from v1; single-operator, single-engagement focus
  11. No full attack graph in v0.1
  12. v0.1 supports exactly Nmap, not dozens of incomplete integrations
  13. Pactyra does not guarantee to find all vulnerabilities; it improves rigor, not skill
  14. Commercial services must not distort or weaken the open-source core

---

### Requirement: Future Expansion — v0.5 (Community Ready)

Pactyra v0.5 SHALL expand the platform to be usable by external operators on real authorized engagements, with community formation as a primary goal.

#### Scenario: v0.5 capabilities are defined

- GIVEN the Pactyra roadmap
- WHEN v0.5 is planned
- THEN the following capabilities SHALL be targeted:
  - Permission lattice O0–O6 (full lattice)
  - Receipt verification tool (standalone CLI)
  - Dynamic scope expansion envelope with approval workflow
  - Tool adapter SDK (Nuclei, HTTP toolkit, custom script adapter)
  - Multi-operator worker isolation
  - Cartographer agent (attack graph construction)
  - Validator agent (proof-of-exploitation)
  - Full agent suite (Cartographer, Validator, Report Generator, Audit Agent)
  - Attack graph: state machine with credential nodes, pivot edges, annotations, query language
  - Evidence: full receipt chain with integrity + provenance + reproducibility
  - Basic retesting: evidence replay against remediated target
  - AI integration: optional provider plugins (OpenAI, Anthropic)
  - Community: public benchmarks, contribution guidelines, design partner program (5+ teams)
  - Terminal: shared attack graph, team awareness, handoff workflows

---

### Requirement: Future Expansion — v1.0 (Reference Implementation)

Pactyra v1.0 SHALL establish the platform as the recognized open-source reference implementation for authorized offensive automation.

#### Scenario: v1.0 capabilities are defined

- GIVEN the Pactyra roadmap
- WHEN v1.0 is planned
- THEN the following capabilities SHALL be targeted:
  - Operator console: web dashboard for reporting, review, and engagement management
  - Retesting: diff engine with side-by-side finding comparison
  - Report: customizable templates, executive summary generation, client-ready PDF export
  - Evidence: third-party verifiable receipt format
  - Tool ecosystem: community adapter marketplace, 20+ supported tools
  - Active contributor base with published benchmarks and conference presence
  - Grant funding and international collaboration
  - Commercial services (if any): strictly additive, never replacing the self-hostable core

---

### Requirement: Scope Expansion Envelope Policy

When dynamic scope expansion is implemented (v0.5+), it MUST follow a formal policy structure that constrains what can be expanded without additional customer approval.

#### Scenario: Expansion envelope policy structure exists

- GIVEN a scope document with an expansion policy section
- WHEN scope expansion is requested mid-engagement
- THEN the policy MUST define:
  - `approver_roles`: which roles can approve expansions (e.g., engagement_manager)
  - `allowed`: pre-authorized expansion capabilities (e.g., additional_subdomains, discovered_ips_within_authorized_asn, authenticated_test_accounts)
  - `limits`: numeric or temporal ceilings (e.g., maximum_new_assets, maximum_duration_extension_hours, maximum_rps)
  - `always_requires_customer_approval`: items that can never be expanded without client authorization (e.g., new_root_domain, new_cidr, production_disruption_risk, social_engineering, data_access_beyond_proof)

#### Scenario: PM activates pre-authorized capabilities, not freely expands

- GIVEN an active engagement with a scope expansion envelope
- WHEN the PM approves a scope expansion
- THEN the PM MUST only activate capabilities that the customer pre-authorized in the envelope
- AND every expansion MUST produce a cryptographic record of the change and its justification
- AND the PM MUST NOT freely expand scope beyond the envelope without customer approval

#### Scenario: Expansion envelope is NOT in v0.1

- GIVEN the v0.1 release
- WHEN scope changes are considered
- THEN dynamic scope expansion SHALL NOT be implemented
- AND the scope SHALL be static for the duration of the engagement

---

### Requirement: Architectural Principles

All technical decisions across all SDD changes MUST be guided by these principles.

#### Scenario: Operator is in control

- GIVEN any decision point in architecture or implementation
- WHEN a tradeoff exists between automation and operator control
- THEN the principle "agents propose, operators dispose" MUST govern
- AND no automated system SHALL publish, exploit, or report without explicit operator authorization

#### Scenario: Evidence is a byproduct

- GIVEN any action in the platform
- WHEN that action completes
- THEN evidence generation MUST be automatic — the operator never explicitly "creates evidence"
- AND the absence of evidence for an action SHALL be treated as a platform defect

#### Scenario: The terminal is the interface

- GIVEN any user-facing feature
- WHEN the interaction model is designed
- THEN the terminal (CLI/TUI) MUST be the primary control surface
- AND web interfaces SHALL be views over existing data, not operational control surfaces

#### Scenario: Scope is enforced, not advised

- GIVEN any scope validation implementation
- WHEN a violation is detected
- THEN the action MUST be blocked
- AND the implementation MUST NOT offer a "warn but proceed" code path in production

#### Scenario: The moat is integration

- GIVEN any feature proposal
- WHEN its value is evaluated
- THEN the feature's contribution to the end-to-end integrity chain (scope → execution → evidence → export) MUST be assessed
- AND isolated features without integration value SHALL be deprioritized

#### Scenario: Start with Rust and Go

- GIVEN the technology stack
- WHEN new components are designed
- THEN trust-boundary code (scope enforcement, receipt verification, sandboxing) SHALL be in Rust
- AND orchestration code (worker management, networking, Docker integration) SHALL be in Go
- AND Python and TypeScript SHALL be introduced only when the core is proven (v0.5+)

#### Scenario: Ship v0.1 in 12 weeks or less

- GIVEN the v0.1 development timeline
- WHEN scope decisions are made
- THEN a working operations kernel with one tool for one operator SHALL be prioritized over comprehensive architecture
- AND features that extend the timeline beyond 12 weeks SHALL be deferred
