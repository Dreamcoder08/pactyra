# Proposal: Pactyra — Open Offensive Operations Kernel (SDD-000)

**Status:** Reviewed  
**Date:** 2026-07-26  
**Author:** SDD Proposal Phase  
**Change ID:** SDD-000

---

## 1. Executive Summary

Pactyra is an **open-source, terminal-first offensive operations kernel for authorized, policy-constrained, and evidence-backed security assessments.** Developed under Arkelythex, it is not a commercial platform for a specific market segment — it is a **reference implementation** that demonstrates a technical thesis: authorized offensive automation can be safe, open, auditable, and reproducible.

**Bound by scope. Proven by evidence.**

The project solves a structural problem that no existing tool addresses: **offensive operations are fragmented across tools, scripts, terminals, and notes that share no verifiable model of authorization, execution, and evidence.** Pactyra provides that model. It binds every action to a cryptographically-verified scope, records every decision in an append-only evidence chain, and produces exportable, independently-verifiable evidence packs.

Pactyra is **provider-agnostic and AI-optional.** It must work without any AI integration. AI providers (OpenAI, Anthropic, and others) may be integrated as optional accelerators — agents that propose, suggest, or validate — but the core authorization, orchestration, and evidence systems operate independently of any model.

The project's public-interest thesis is simple: as AI agents increasingly participate in offensive security, the industry needs a reference implementation that proves safety, auditability, and reproducibility are achievable. Pactyra aims to become that reference — and a strong candidate for grants, credits, accelerators, and international collaboration.

Commercial services, hosted collaboration, and enterprise features may exist in the future. They must not distort or weaken the open-source core.

---

## 2. Public-Interest Thesis

### 2.1 Why This Matters Now

Offensive security is undergoing an automation wave. AI agents are being integrated into reconnaissance, exploitation, and reporting pipelines. Without a verifiable model of authorization, execution, and evidence, this trend produces two risks:

1. **Unsafe automation.** An AI agent that misinterprets a scope document or hallucinates a target can cause real legal and operational damage — and there is no standard control plane to prevent it.
2. **Unauditable results.** When an AI agent produces a finding, how does a human operator — or a client, or a regulator — verify what actually happened? Without an evidence chain, AI-assisted findings are less credible than manual ones.

Pactyra addresses both risks at the architectural level: scope is machine-enforced, not advisory; every action produces an immutable receipt; and evidence packs are independently verifiable without access to the platform.

### 2.2 The Reference Implementation Goal

Pactyra is not just a tool — it is an argument. The argument is:

> Authorized offensive automation does not have to be opaque, unsafe, or unverifiable. A control plane that enforces scope at runtime, records every action in an append-only evidence chain, and exports independently-verifiable evidence packs is feasible, useful, and should exist as open infrastructure.

The project aims to:

- **Demonstrate** the thesis with a working, self-hostable implementation.
- **Publish** reproducible benchmarks against public cyber ranges.
- **Build** a community of operators, researchers, and contributors.
- **Qualify** for grants, credits, accelerators, and international collaboration.

### 2.3 Repository Model

| Property                | Commitment                                                                                                                                                                                                    |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **License**             | AGPL-3.0 — strong copyleft. Protects the open-source core from commercial enclosure. Compatible with grants and accelerators. Commercial services must use the AGPL-3.0 core or negotiate separate licensing. |
| **Distribution**        | Self-hostable. Clone, build, run. No mandatory services.                                                                                                                                                      |
| **Interface**           | Terminal-first. A web dashboard may exist later for reporting and review, but operations happen in the terminal.                                                                                              |
| **AI dependency**       | None. AI is an optional accelerator, never a requirement.                                                                                                                                                     |
| **Provider lock-in**    | None. Provider-agnostic by design. OpenAI, Anthropic, and others are interchangeable integrations.                                                                                                            |
| **Commercial services** | Permitted in the future (hosted collaboration, enterprise features). Must not distort or weaken the open-source core.                                                                                         |

---

## 3. Problem Statement

### 2.1 The Fragmentation Tax

A typical specialized offensive security engagement today runs across at least six disconnected surfaces:

| Surface                                        | What It Holds                      | How It Fails                                                                |
| ---------------------------------------------- | ---------------------------------- | --------------------------------------------------------------------------- |
| Kali VM / local tools                          | Nmap, Burp, Nuclei, custom scripts | No runtime scope enforcement; `nmap` does not know your rules of engagement |
| Terminal multiplexer (tmux/Zellij)             | Session state, command history     | Scrollback is not evidence; sessions are lost on disconnect                 |
| Note-taking (Notion, Obsidian, markdown files) | Findings, hypotheses, methodology  | Manual; diverges from actual tool output; no cryptographic binding          |
| Chat (Slack, Signal)                           | Team coordination, handoffs        | Ephemeral; context lost when operator leaves                                |
| Report template (Word, LaTeX, custom tool)     | Final deliverable                  | Assembled retroactively from notes; no live link to evidence                |
| git / shared drive                             | Scripts, payloads, wordlists       | No version-to-engagement traceability                                       |

The operator spends 30–40% of their engagement time on **integration labor**: copying output, correlating findings across tools, verifying that nothing was tested outside scope, assembling evidence for the report, and reconstructing what happened when a client questions a finding two weeks later. This is time not spent finding vulnerabilities.

### 2.2 Four Hard Problems No Script Solves

1. **Runtime scope enforcement.** A mistyped IP or a forgotten `--exclude` flag can expand an Nmap scan into a /16 subnet the operator was never authorized to touch. No shell alias or wrapper script catches this consistently, and the consequences range from client trust erosion to legal liability.

2. **Evidence integrity under challenge.** When a client disputes a finding — "we fixed that," "that's out of scope," "prove you actually exploited this" — the operator's defense is screenshots, terminal scrollback, and memory. None of it is cryptographically verifiable. None of it proves the exact command was run at the exact time under the exact scope.

3. **Attack path continuity across operators.** A junior operator picks up where a senior left off. They inherit tmux panes, markdown notes, and Slack threads. The mental model of the attack graph — which services were discovered, which were exploited, which pivots were attempted and failed — lives in the senior operator's head. When the senior operator is on another engagement, that model is unavailable.

4. **Retesting as archaeology.** Six months after the engagement, the client has remediated. The team needs to retest. Which exact Nmap flags were used? Which Nuclei templates? Which custom exploit variant? The answer is usually "reconstruct from notes and hope." There is no replayable evidence chain.

### 2.3 Why Existing Tools Don't Solve This

- **Metasploit Pro / Cobalt Strike** are exploitation frameworks, not engagement management platforms. They don't enforce scope, build attack graphs across heterogeneous tools, or produce evidence receipts.
- **PlexTrac / Dradis** are reporting and collaboration tools. They manage findings but don't orchestrate operations or enforce permissions at runtime.
- **Custom scripts + Makefiles** are fragile, operator-specific, and provide no cryptographic guarantees.
- **Enterprise vulnerability management platforms (Qualys, Tenable, Rapid7)** are defensive, scan-based, and designed for vulnerability programs — not for authorized offensive operations with adversarial validation.

The gap is not "another tool." The gap is **an operations kernel that binds authorization, execution, evidence, and retesting into a single integrity chain**.

---

## 3. Beachhead User Profile

### Alex Chen — Lead Operator, Aegis Security Partners

**Firm profile:** 7-person specialized offensive security firm. Three senior operators, two mid-level, one junior, one client-facing PM. 3–4 engagements per month. Clients are mid-market SaaS companies and fintechs that need SOC 2 / ISO 27001 evidence and are starting to require penetration tests with traceable methodology.

**Alex's role:** Primary technical operator. Leads 70% of engagements. Decides tooling, methodology, and technical approach. Reviews junior operator work. Writes or reviews the technical sections of every report.

**Daily workflow (today):**

1. Receives scope document from PM (PDF, email, or spreadsheet). Manually parses IP ranges, domains, excluded hosts, and authorized techniques.
2. Opens Kali VM. Launches tmux with a standard pane layout: Recon, Exploitation, Notes, Logs.
3. Runs initial recon: `nmap`, `subfinder`, `httpx`. Copies output manually into notes.
4. Pivots to Burp for web targets, Nuclei for known CVEs, custom scripts for business logic.
5. Documents findings incrementally in a Notion page. Screenshots evidence.
6. When stuck or handing off, posts in Slack: "can someone check the auth bypass on app.target.com? details in Notion."
7. At engagement end: exports Notion page, manually assembles report in LaTeX, cross-references screenshots with findings, sends to PM.
8. Six months later: client asks for retesting. Alex searches Slack and Notion for what was actually done. Reconstructs methodology from memory and partial notes. Spends 4–6 hours just preparing.

**Alex's pain points (in their own words):**

- "I spend more time documenting what I did than actually doing it."
- "Every time I type an IP, I'm one typo away from scanning something I shouldn't. I've had two near-misses this year."
- "When a client pushes back on a finding, I have screenshots. But screenshots aren't proof — they're just pictures."
- "Training juniors is hard because I can't replay my sessions. I have to sit with them and narrate."
- "Retesting is the worst part of the job. It's archaeology on my own work from six months ago."

**What Alex wants (not what they'd say, but what they need):**

- A terminal where typing `scan 10.0.0.0/8` is impossible because the scope only authorizes `10.0.1.0/24`.
- A system that automatically records every command, its output, its timing, and its cryptographic signature — without Alex having to think about it.
- An attack graph that builds itself as Alex works, so the junior operator can see "we've exploited host A, pivoted to B via service X, and B's port 8443 is unexplored."
- One command to retest: replay the exact evidence chain from the original engagement against the remediated target.
- A report that is a view over the evidence graph, not a separate document to assemble.

---

## 4. Core Promise & Value Proposition

### 4.1 The Three-Tier Value Architecture

Pactyra's value proposition is ordered, not flat:

```text
Primary pain:
Offensive results are disputable because evidence is fragmented,
incomplete, or does not verifiably demonstrate what occurred.

Trust mechanism:
Pactyra binds every action and piece of evidence to a verifiable
scope, identity, tool, version, authorization, and timeline.

Aha moment:
The system blocks an out-of-scope action before it creates
contractual, legal, or reputational risk.
```

**They adopt Pactyra to defend the credibility of their work; they stay because Pactyra reduces their operational and legal risk.**

### 4.2 The Core Promise

> **Turn an approved scope into a reproducible, evidence-backed security operation — from discovery through retesting.**

This is not "automate pentesting." The operator is still the operator. Pactyra does not decide what to scan, what to exploit, or how to chain vulnerabilities. It provides the **scaffolding** that makes the operator faster, safer, and more credible: scope enforcement, orchestration, evidence, and replay.

### 4.3 The Evidence Model: Three Layers

Pactyra does not promise "absolute cryptographic proof that a vulnerability exists." That is impossible — an operator can manipulate the environment before capturing evidence, and interpretation of findings always involves human judgment. Instead, Pactyra provides three independently useful layers:

| Layer               | What It Guarantees                                                                        | What It Does NOT Guarantee                                      |
| ------------------- | ----------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| **Integrity**       | The content hasn't changed since registration (hash-chained).                             | That the original capture was accurate or complete.             |
| **Provenance**      | Who, when, where, with what tool, and under what authorization produced the artifact.     | That the operator's methodology was correct.                    |
| **Reproducibility** | Another authorized operator can repeat the validation following the documented procedure. | That the target will respond identically (environments change). |

**The promise: evidence with verifiable integrity, traceable provenance, and documented reproducibility.** Not "absolute proof."

### 4.4 The Value Chain

```
Scope Document (input)
    │
    ▼
[Scope Guardian] ── Parses, validates, cryptographically binds scope to engagement
    │
    ▼
[Operator Terminal] ── Every command is scoped, recorded, receipted
    │
    ▼
[Agent Orchestration] ── Recon Planner, Cartographer, Validator run tools in isolated workers
    │
    ▼
[Attack Graph] ── Live state of the engagement: discovered, exploited, pivoted, dead ends
    │
    ▼
[Evidence Receipts] ── Immutable, cryptographically-signed record of every action and output
    │
    ▼
[Adversarial Validation] ── "Did we actually exploit this?" with proof-of-exploitation receipts
    │
    ▼
[Report Generation] ── A rendered view over the evidence graph; no manual assembly
    │
    ▼
[Retesting] ── Replay the exact evidence chain against the remediated target
```

### 4.3 Value Proposition (for the beachhead user)

| Without Pactyra                                                  | With Pactyra                                                               |
| ---------------------------------------------------------------- | -------------------------------------------------------------------------- |
| Scope lives in a PDF; enforced by operator vigilance             | Scope is machine-readable; enforced at runtime by the control plane        |
| Evidence is screenshots and notes; integrity depends on trust    | Evidence is cryptographic receipts; integrity is mathematically verifiable |
| 30–40% of time spent on integration and documentation labor      | Evidence and reporting are automatic byproducts of operation               |
| Handoffs require Slack narration and knowledge transfer meetings | Attack graph is live, shared, and navigable by any team member             |
| Retesting requires archaeology and operator memory               | Retesting replays the evidence chain; one command                          |
| Client challenges met with screenshots                           | Client challenges met with signed, timestamped receipts                    |
| Junior operators need constant supervision                       | Juniors operate within enforced scope; seniors review the graph            |
| Methodology is tribal knowledge                                  | Methodology is encoded in the scope, graph, and receipts                   |

---

## 5. Differentiator Deep-Dive: Why Scripts + Kali + Notion Cannot Replace This

The combination of Kali Linux, a terminal multiplexer, custom scripts, and a note-taking tool is the de facto standard for small offensive teams. It is also the ceiling. Here is why Pactyra's integration cannot be replicated by gluing these together:

### 5.1 Runtime Scope Enforcement

A bash wrapper can check IPs against a whitelist before running `nmap`. But:

- It cannot enforce scope across heterogeneous tools (Nmap, Burp, custom Python, Nuclei templates) consistently — each tool needs its own wrapper.
- It cannot handle dynamic scope (the operator discovers a new subnet and needs to request scope expansion mid-engagement).
- It cannot produce a verifiable proof that scope was respected — only that the wrapper was running.
- It does not survive `operator runs tool directly` — the wrapper is advisory, not enforced.

Pactyra enforces scope at the orchestration layer: tools run inside isolated workers that receive only the scope-authorized parameters. The operator cannot bypass this because the worker does not have network access outside the authorized target set.

### 5.2 Verifiable Evidence Integrity

Screenshots and terminal logs can be forged or tampered with after the fact. Timestamps can be manipulated. Output can be edited. A determined client or regulator will not accept "here's a PNG of my terminal" as proof that a critical vulnerability was actually exploited on a specific date.

Pactyra produces **immutable evidence receipts** with three independently verifiable layers:

- **Integrity:** Every command, its parameters, its output, and its timing are hashed and chained. Given a receipt, a third party can verify the content hasn't been modified since registration.
- **Provenance:** Each receipt records who ran the command, when, with what tool version, and under which scope authorization. The full lineage is traceable.
- **Reproducibility:** The receipt captures the exact procedure — tool, parameters, target, preconditions — so another authorized operator can repeat the validation and compare results.

The receipt chain is append-only. Prior receipts cannot be modified without breaking the hash chain. Note: integrity proves the artifact wasn't tampered with after capture; it does not prove the operator didn't manipulate the environment before capture, nor does it guarantee the interpretation of the finding is correct. That responsibility remains with the operator.

### 5.3 Attack Graph Construction

An operator can manually diagram an attack path in draw.io or keep a mental model. But:

- The diagram diverges from reality as the engagement progresses.
- It captures only what the operator remembers to document, not what actually happened.
- It cannot be queried: "show me all unexplored ports on hosts where we have credentials."
- It is not shared in real-time across the team.

Pactyra constructs the attack graph **from execution**: every scan result, every successful exploit, every discovered credential updates the graph automatically. The graph is a live, queryable state machine — not a documentation artifact.

### 5.4 Adversarial Validation

A Nessus scan flags 200 "critical" findings. A Nuclei template fires on a header match. Neither proves exploitability. The operator manually validates a subset and documents the result. But the gap between "scanner says vulnerable" and "operator proved it's exploitable" is where false positives live, where client trust erodes, and where report quality degrades.

Pactyra's Validator agent records **proof-of-exploitation receipts**: not just "scanner flagged this," but "here is the exact exploit, its output, and the cryptographic receipt of successful exploitation." Findings without exploitation receipts are labeled as unvalidated. This distinction — between detected and proven — is impossible with a script pipeline.

### 5.5 Retesting from Evidence

The retesting problem is the acid test. Today: operator opens old notes, guesses at methodology, runs tools, compares results manually. With Pactyra: the evidence chain from the original engagement is replayed against the remediated target. Every command, every tool invocation, every parameter is reproduced exactly. New evidence receipts are generated and compared against the originals. The retest report is a diff, not a reconstruction.

### 5.6 The Integration Moat

Individually, each capability (scope enforcement, evidence receipts, attack graphs, validation, retesting) could be approximated with dedicated tools. But the **integration** — scope flows into execution flows into evidence flows into the graph flows into validation flows into the report flows into retesting — is the moat. No script pipeline maintains integrity across all these boundaries simultaneously. Pactyra does.

---

## 6. High-Level Capability Map

### 6.1 Scope Management

- **Scope ingestion**: Parse scope documents (CIDR notation, domain lists, excluded hosts, authorized techniques, time windows, contact information).
- **Scope binding**: Cryptographically bind scope to the engagement. The scope becomes the root of the evidence chain.
- **Runtime enforcement**: Every tool invocation is validated against the active scope before execution. Out-of-scope actions are blocked, not warned.
- **Dynamic scope expansion**: Operator requests scope modification mid-engagement. The PM approves only within pre-authorized boundaries defined in the scope document. A formal expansion envelope prevents ambiguous approvals:

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

The PM does not "freely expand scope"; they activate a capability the customer pre-authorized. Every expansion produces a cryptographic record of the change and its justification.

- **Permission lattice (O0–O6)**: Escalating trust levels. O0 = read-only reconnaissance from external perspective. O6 = full exploitation with credential-based lateral movement. The operator's current level gates what tools and targets are available.

### 6.2 Agent Orchestration

- **Scope Guardian**: Validates every action against the active scope and permission level. The gatekeeper.
- **Recon Planner**: Translates operator intent ("map this subnet's web attack surface") into tool invocations with appropriate parameters. Plans are validated by Scope Guardian before execution.
- **Cartographer**: Consumes tool output and updates the attack graph. Discovers relationships between hosts, services, vulnerabilities, and credentials.
- **Validator**: Takes a finding and produces proof-of-exploitation evidence. Distinguishes between "detected" and "proven."
- **Report Generator**: Renders the evidence graph into a structured report. A view, not a separate artifact.
- **Audit Agent**: Continuously verifies the integrity of the evidence chain. Detects tampering, missing receipts, and scope anomalies.

### 6.3 Attack Graph (Deferred to v0.5+)

> **The attack graph is not part of the v0.1 wedge.** v0.1 focuses on scope enforcement + evidence receipts + professional export. The attack graph is a v0.5 capability. It is described here for architectural coherence only.

- **Live state machine**: Nodes = hosts, services, vulnerabilities, credentials. Edges = network reachability, exploitation, credential reuse, pivot paths.
- **Automatic construction**: Updated by Cartographer as tool output arrives. No manual diagramming.
- **Queryable**: "What hosts have port 8443 open?", "Which services were discovered on host A after we exploited host B?", "Show dead-end pivots."
- **Operator-annotatable**: Operators can add hypotheses, mark nodes as "investigated — dead end," and attach notes to edges.
- **Persistence**: Attack graph survives session disconnections. Shared across the team.

### 6.4 Evidence Receipts

- **Immutable chain**: Every action (tool invocation, scope change, operator note, finding validation) produces a receipt.
- **Receipt contents**: Action type, parameters, output hash, timestamp, operator identity, active scope snapshot, permission level, parent receipt hash.
- **Cryptographic verification**: Any receipt can be verified independently. The chain is append-only; prior receipts cannot be modified without detection.
- **Storage**: Local during engagement. Exportable for client delivery, audit, or retesting.

### 6.5 Adversarial Validation

- **Validation workflow**: Finding detected → Validator agent crafts proof-of-exploitation → exploitation attempt recorded with full receipt → result classified as "proven exploitable," "proven not exploitable," or "inconclusive."
- **Exploitation receipts**: Contain the exact exploit, target, preconditions, postconditions, and output. Verifiable independently.
- **Validation coverage**: Report distinguishes between validated and unvalidated findings explicitly.

### 6.6 Retesting

- **Evidence replay**: The evidence chain from an original engagement is replayed against the remediated target. Same tools, same parameters, same order.
- **Diff report**: Side-by-side comparison of original vs. retest findings. What was fixed, what persists, what's new.
- **Re-validation**: Exploitation receipts are replayed against the remediated target. If the exploit no longer works, the finding is marked as remediated with cryptographic proof of the attempt.

---

## 7. Explicit Scope Boundaries

### 7.1 In Scope (v0.1 → v1.0)

| Capability                                      | v0.1                                | v0.5                                | v1.0                        |
| ----------------------------------------------- | ----------------------------------- | ----------------------------------- | --------------------------- |
| Scope ingestion (CIDR, domain list)             | ✅                                  | ✅                                  | ✅                          |
| Runtime scope enforcement                       | ✅ (basic IP/port gating)           | ✅ (full O0–O6 lattice)             | ✅                          |
| Isolated tool execution (Docker workers)        | ✅                                  | ✅                                  | ✅                          |
| Evidence receipts (append-only, hash-chained)   | ✅                                  | ✅                                  | ✅                          |
| Standalone evidence export (no platform needed) | ✅                                  | ✅                                  | ✅                          |
| Single-operator mode                            | ✅                                  | ✅                                  | ✅                          |
| Terminal-first operator interface (TUI)         | ✅                                  | ✅                                  | ✅                          |
| Report generation (structured output)           | ✅ (basic evidence dump)            | ✅ (executive + technical sections) | ✅ (customizable templates) |
| Multi-agent architecture                        | ✅ (Scope Guardian + Recon Planner) | ✅ (Cartographer, Validator added)  | ✅ (full agent suite)       |
| Basic attack graph (hosts, services, ports)     | —                                   | ✅                                  | ✅                          |
| Team collaboration (shared graph, receipts)     | —                                   | ✅                                  | ✅                          |
| Adversarial validation (proof-of-exploitation)  | —                                   | ✅                                  | ✅                          |
| Retesting from evidence replay                  | —                                   | ✅ (basic replay)                   | ✅ (diff engine)            |
| Dynamic scope expansion workflow                | —                                   | ✅                                  | ✅                          |
| Permission lattice (O0–O6)                      | O0–O3 only                          | O0–O6                               | O0–O6                       |
| Cryptographic receipt verification (standalone) | ✅                                  | ✅                                  | ✅                          |

### 7.2 Explicit Non-Goals (Now and Forever)

These are derived from the strategic decisions and represent permanent boundaries, not just v1 deferrals:

- **Not a Linux distribution.** Pactyra runs on existing Linux installations. It does not ship its own kernel, package manager, or desktop environment.
- **Not a replacement for Burp Suite, BloodHound, Nmap, or Nuclei.** Pactyra orchestrates these tools; it does not reimplement their functionality. The operator brings their tools; Pactyra provides the operations kernel.
- **Not a platform for attacking targets without authorization.** The platform requires a cryptographically-bound scope document to operate. There is no "unscoped mode."
- **Not malware, C2, or persistence infrastructure.** Pactyra does not generate or host payloads designed for unauthorized access.
- **Not an exploit marketplace.** Pactyra does not distribute, sell, or broker exploits.
- **Not a system where an AI agent can elevate its own permissions.** Agents operate within the operator's current permission level. Permission escalation requires explicit operator action.
- **Not autonomous exploitation without supervision.** Agents propose; operators dispose. No finding is published without operator review.
- **Not a defensive SIEM, EDR, or SOAR platform.** Pactyra is for authorized offensive operations, not defensive monitoring.
- **Not a customer portal (v0.1).** The platform is operator-first. Client-facing portals, if any, are post-v1.0 considerations that must not distort the open-source core.
- **Not multi-tenant enterprise from v1.** Single-operator, single-engagement focus. Multi-tenancy is a potential future consideration, not a v1 design constraint.
- **Not a full attack graph in v0.1.** The attack graph is deferred to v0.5. v0.1 focuses on scope enforcement + evidence receipts + professional export.
- **Not dozens of incomplete integrations.** v0.1 supports exactly one tool (Nmap). Tool adapter SDK is v0.5. Quality over quantity.
- **Not a guarantee to find all vulnerabilities.** Pactyra improves operational rigor, traceability, and reproducibility. It does not improve the operator's skill, the tools' detection capabilities, or the completeness of a methodology.
- **Commercial services must not distort the open-source core.** Hosted collaboration, enterprise features, and commercial services may exist in the future, but the self-hostable, AI-optional, provider-agnostic core remains the reference implementation and the public good.

---

## 8. Success Metrics

### 8.1 v0.1 Hard Metrics (Non-Negotiable)

These are the exit criteria for v0.1 "First Light." Every metric must be at the specified threshold before v0.1 ships.

| Metric                                           | Criterion | Rationale                                                         |
| ------------------------------------------------ | --------- | ----------------------------------------------------------------- |
| Actions executed without valid scope             | **0**     | Scope enforcement is the wedge; any bypass is a failed v0.1       |
| Known violations not blocked                     | **0**     | Enforcement, not advisory. Warnings are not acceptable            |
| Artifacts without verifiable provenance          | **0**     | Every artifact must trace to who, when, tool, authorization       |
| Confirmed findings without reproducible evidence | **0**     | Every finding in the report must be backed by replayable evidence |
| Actions bound to `scope_hash`                    | **100%**  | No action exists outside the scope binding                        |
| Usable export without platform access            | **100%**  | The evidence export must be readable and verifiable standalone    |
| Expansions without valid approval                | **0**     | Every scope expansion must have a recorded, valid approval        |

### 8.2 Operator Efficiency (v0.5 target)

| Metric                         | Current Baseline (Alex's firm) | Target (v0.5) | Measurement                                                   |
| ------------------------------ | ------------------------------ | ------------- | ------------------------------------------------------------- |
| **Integration labor ratio**    | 30–40% of engagement time      | ≤15%          | Self-reported or time-in-tool vs. time-in-evidence-production |
| **Report assembly time**       | 4–8 hours per engagement       | ≤1 hour       | Clock time from engagement end to report draft complete       |
| **Retesting preparation time** | 4–6 hours                      | ≤30 minutes   | Clock time from retest request to first replayed action       |
| **Scope violation incidents**  | 1–2 near-misses/year           | 0 incidents   | Runtime enforcement prevents rather than warns                |

### 8.3 Adoption & Community (open-source metric)

| Metric                             | Target (v1.0)                      | Measurement          |
| ---------------------------------- | ---------------------------------- | -------------------- |
| **GitHub stars**                   | 500+                               | GitHub API           |
| **Active contributors**            | 15+ unique contributors/month      | GitHub API           |
| **Public benchmarks reproduced**   | 3+ independent reproductions       | Community reports    |
| **Design partner teams**           | 5+ teams using on real engagements | Self-reported        |
| **Grant/accelerator applications** | 3+ submitted, 1+ accepted          | Application tracking |

---

## 9. Risks & Open Questions

### 9.1 Technical Risks

| Risk                                                                                                                                        | Severity | Mitigation                                                                                                                               |
| ------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| **Scope enforcement bypass via tool bugs**: A worker escaping its sandbox could access out-of-scope targets.                                | Critical | Defense in depth: network-level egress filtering per worker, not just parameter validation. Scope Guardian validates at multiple layers. |
| **Evidence chain corruption or forgery**: If the signing key is compromised, all receipts become suspect.                                   | Critical | Hardware-backed key storage (TPM). Receipt chain verification is standalone — a client can verify without trusting the platform.         |
| **Performance overhead of receipt generation**: Hashing and signing every action could slow operations.                                     | Medium   | Async signing with bounded queues. Batch receipts for high-frequency actions.                                                            |
| **Attack graph state explosion**: Large engagements could produce graphs too large to navigate or store efficiently.                        | Medium   | Pruning strategies. Graph query language with pagination. Differential updates.                                                          |
| **Tool compatibility fragmentation**: Each external tool (Nmap, Nuclei, etc.) needs an adapter for parameter validation and output parsing. | Medium   | Adapter SDK. Focus on 5–8 core tools for v0.1; community adapters later.                                                                 |

### 9.2 Project/Community Risks

| Risk                                                                                                                                               | Severity | Mitigation                                                                                                                                                                                    |
| -------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Operator resistance to control plane**: "I just want to run my tools." The enforcement that protects the operator also constrains them.          | High     | Progressive disclosure of constraints. Fast scope-expansion workflow. The terminal feels like a terminal, not a web form.                                                                     |
| **Open-source community failure**: The project fails to attract contributors and becomes a single-maintainer effort.                               | High     | Public benchmarks, clear contribution guidelines, design partner program, conference talks, grant applications from v0.5.                                                                     |
| **"Why not just scripts?" perception**: Operators with existing custom toolchains may not see the value of a control plane.                        | Medium   | Target the evidence challenge and scope violation pain — these are unsolved even with custom scripts. Public benchmarks demonstrate the gap.                                                  |
| **Competitive replication**: If the thesis is compelling, larger players (commercial or open-source) could replicate the architecture.             | Medium   | The integration moat (scope → execution → evidence → export) is hard to replicate quickly. Community and reputation are compounding advantages. Open-source license protects the core.        |
| **Grant/accelerator dependency**: If the project relies on external funding to sustain development beyond v0.5, funding gaps could stall progress. | Medium   | Keep the core self-hostable and low-cost to maintain. Design partner program provides real-world validation that attracts funding. Multiple funding sources (grants, accelerators, sponsors). |

### 9.3 Execution Risks

| Risk                                                                                                                                                                  | Severity | Mitigation                                                                                                                                    |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| **Multi-language monorepo complexity**: Rust + Go + Python + TypeScript creates integration friction and requires broad expertise.                                    | High     | Start with Rust (trust boundary) and Go (orchestration) only. Python and TypeScript are v0.5+. Single build system (Bazel or Nix).            |
| **Greenfield architecture over-engineering**: Building a permission lattice, evidence chain, attack graph, and multi-agent system before a single tool runs is risky. | High     | v0.1 delivers: scope parsing, basic enforcement, two agents, one tool (Nmap), evidence receipts, and a TUI. Nothing more. Ship in 8–12 weeks. |
| **Design partner recruitment**: Without real operators using it, v0.1 may solve imaginary problems.                                                                   | High     | Recruit 3 design partners before writing code. Their first engagement on Pactyra is the v0.1 acceptance test.                                 |

### 9.4 Resolved Questions (Proposal Question Round — Answered)

1. **Primary pain (pricing):** Evidence challenges — when a client disputes a finding and the operator lacks verifiable proof. ✅
2. **User breadth (v0.1):** Operator-only. PM receives exported report. Client never touches the platform. ✅
3. **Scope expansion:** PM approves within pre-authorized boundaries defined in the scope document's expansion envelope. ✅
4. **Aha moment:** First blocked scope violation — "This just saved me from a legal problem." ✅
5. **Compliance frameworks:** Guidance, not blocker. SOC 2, PCI DSS, ISO 27001 serve as design guidance for evidence format and retention, but are NOT v0.1 blockers. ✅
6. **Failed engagements:** Immutable — no deletion. Evidence chain is immutable by design. Client cannot demand deletion; evidence can be archived/sealed but never destroyed. ✅
7. **Tool ecosystem (v0.1):** CLI-only with structured output (JSON, XML, grep-able). Burp Suite and other GUI tools run alongside, never inside Pactyra workers. ✅
8. **Monetization model:** Open-source core (AGPL-3.0). Commercial services permitted later but must not distort or weaken the open-source reference implementation. ✅

---

## 10. Phased Vision

### v0.1 — "First Light" (8–12 weeks)

**Goal:** Demonstrate the complete golden path against a local cyber range. The thesis is proven when an operator loads a scope, runs discovery, attempts an out-of-scope action, sees it blocked, and exports a verifiable evidence pack.

**The v0.1 Golden Path (8 steps):**

```text
1. Load and compile an authorized scope
2. Generate a scope_hash
3. Execute permitted discovery against a local cyber range
4. Deliberately attempt an out-of-scope action
5. Block it before execution
6. Record actions and decisions
7. Generate verifiable receipts
8. Export a reproducible evidence pack
```

**Technology:**

- Rust: Scope parser (CIDR, domain list), scope enforcement engine, evidence receipt chain (hash-chained, append-only), TUI scaffolding
- Go: Worker orchestration with Docker isolation, Nmap adapter (parameter validation, output parsing)
- Multi-agent: Scope Guardian (enforces scope on every tool invocation), Recon Planner (translates operator intent to Nmap invocations)
- Evidence: Append-only receipt chain with integrity + provenance layers. Exportable in standalone-verifiable format.
- Terminal: TUI with scope view, command input, receipt browser.

**Explicitly NOT in v0.1:** Attack graph, adversarial validation, retesting, permission lattice beyond basic IP/port gating, team collaboration, dynamic scope expansion, web dashboard, AI integration, dozens of tools.

**Deliverable:** A public repository where anyone can clone, build, run the golden path against a documented cyber range, and verify the evidence pack independently — without access to the platform.

### v0.5 — "Community Ready" (6–9 months after v0.1)

**Goal:** Pactyra is usable by external operators on real (authorized) engagements. A community forms around the project.

- Rust: Permission lattice O0–O6, receipt verification tool (standalone CLI), dynamic scope expansion envelope
- Go: Tool adapter SDK (Nuclei, HTTP toolkit, custom script adapter), multi-operator worker isolation
- Python: Cartographer (attack graph construction), Validator (proof-of-exploitation)
- Multi-agent: Full agent suite (Cartographer, Validator, Report Generator, Audit Agent)
- Attack graph: State machine with credential nodes, pivot edges, annotations, query language
- Evidence: Full receipt chain with integrity + provenance + reproducibility. Exportable for client delivery.
- AI integration: Optional provider plugins (OpenAI, Anthropic) for operator assistance — never required.
- Community: Public benchmarks, contribution guidelines, design partner program (5+ teams).
- Terminal: Shared attack graph, team awareness, handoff workflows.

### v1.0 — "Reference Implementation" (12–18 months after v0.1)

**Goal:** Pactyra is the recognized open-source reference implementation for authorized offensive automation.

- TypeScript: Operator console (web dashboard for reporting, review, and engagement management)
- Retesting: Diff engine with side-by-side finding comparison
- Report: Customizable templates, executive summary generation, client-ready PDF export
- Evidence: Third-party verifiable receipt format (standalone, no platform needed)
- Community: Active contributor base. Published benchmarks. Conference talks. Grant funding.
- Tool ecosystem: Community adapter marketplace. 20+ supported tools.
- Commercial services (if any): Hosted collaboration, enterprise features — strictly additive, never replacing or weakening the self-hostable open-source core.

**Repository at v1.0:**

- Stars, forks, contributors as primary health metrics
- Reproducible benchmarks against public cyber ranges
- Active grant and accelerator applications
- Provider-agnostic: works without AI, works with multiple AI providers

---

## 11. SDD-000 Boundaries — What This Change Owns

SDD-000 is the **vision, boundaries, and non-goals** change. It does NOT specify concrete schemas, wire formats, agent contracts, or UI behavior. Those are separate SDD changes.

### 11.1 Owned by SDD-000

| Artifact                 | Description                                                                                                     |
| ------------------------ | --------------------------------------------------------------------------------------------------------------- |
| Public-interest thesis   | Why authorized offensive automation needs an open reference implementation                                      |
| Repository model         | Open-source, self-hostable, AI-optional, provider-agnostic                                                      |
| Actors                   | Who uses Pactyra and in what roles                                                                              |
| Primary user             | Technical offensive security operator (pentester, researcher, red team)                                         |
| Jobs to be done          | What the operator needs to accomplish                                                                           |
| Golden path              | The 8-step v0.1 demonstration: scope → hash → discover → attempt violation → block → record → receipts → export |
| Value proposition        | Three-tier: pain (evidence challenges), trust mechanism (bound evidence), aha moment (blocked violation)        |
| Trust model              | Integrity, provenance, reproducibility layers                                                                   |
| Product invariants       | "Every action must produce an auditable entry bound to an active scope"                                         |
| Success metrics          | v0.1 hard metrics (7 non-negotiable) + v0.5/v1.0 targets                                                        |
| Boundaries               | What's in, what's out, what's forever-out                                                                       |
| Non-goals                | 14 explicit permanent non-goals including commercial services constraint                                        |
| v0.1 scope               | The wedge: scope enforcement + evidence receipts + professional export                                          |
| Future expansion         | v0.5 (Community Ready) and v1.0 (Reference Implementation) horizons                                             |
| Scope expansion envelope | Policy structure (YAML shape, approver roles, limits)                                                           |

### 11.2 Deferred to Future SDD Changes

SDD-000 may include high-level requirements (e.g., "Every action must be bound to a `scope_hash`") but NOT the concrete schemas, wire formats, or implementation details. These are owned by dedicated changes:

| Future SDD | Topic                                                                            | Depends On       |
| ---------- | -------------------------------------------------------------------------------- | ---------------- |
| SDD-002    | Scope model and compiler (CIDR, domain, exclusion, expansion envelope)           | SDD-000          |
| SDD-003    | Permission lattice O0–O6 (levels, capabilities, enforcement rules)               | SDD-000, SDD-002 |
| SDD-005    | Tool and agent contracts (adapter SDK, worker isolation, agent communication)    | SDD-000          |
| SDD-007    | Evidence and receipts (receipt schema, hash chain, verification, export format)  | SDD-000, SDD-002 |
| SDD-010    | Finding lifecycle (detection → validation → confirmation → remediation → retest) | SDD-007          |
| SDD-015    | Attack graph (state machine, query language, annotations)                        | SDD-005, SDD-007 |
| SDD-0XX    | Operator TUI (layout, components, keybindings, navigation)                       | SDD-002, SDD-005 |

### 11.3 What SDD-000 Says About Each (High-Level Only)

- **Scope model:** Must parse CIDR and domain lists. Must enforce at runtime, not warn. Must support expansion envelope with pre-authorized boundaries. Must bind every action to a `scope_hash`.
- **Permission lattice:** O0–O6 levels exist. v0.1 only needs O0–O3 (basic IP/port gating).
- **Tool contracts:** Tools run in isolated Docker workers. Adapters translate parameters and parse output. v0.1: Nmap only.
- **Evidence receipts:** Append-only hash chain. Three layers: integrity, provenance, reproducibility. Exportable in standalone-verifiable format. v0.1 does NOT specify the exact JSON schema.
- **Finding lifecycle:** Operator confirms or rejects. Every finding backed by reproducible evidence. v0.1: basic evidence dump, not structured finding management.
- **Attack graph:** NOT in v0.1. Deferred to v0.5.
- **TUI:** Terminal-first. Scope view, command input, receipt browser. v0.1 does NOT specify layout, keybindings, or navigation model.

---

## 12. Brand

### 12.1 Identity

```text
PACTYRA
Open Offensive Operations Kernel

Bound by scope. Proven by evidence.
```

- **Canonical name:** Pactyra (pronounced: pak-TEE-rah)
- **Parent organization:** Arkelythex
- **Full signature:** Pactyra by Arkelythex
- **Repository:** arkelythex/pactyra
- **CLI binary:** `pactyra`

### 12.2 Construction

**PACT** represents the authorization contract — the Rules of Engagement, the agreement between client and operator, the explicit authority required before any action. **-YRA** provides original identity, international fluency, and continuity with the Arkelythex family.

### 12.3 Narrative (Canonical)

> Pactyra derives its identity from the pact that authorizes every offensive-security operation. It represents a system in which authority precedes execution, scope constrains action, and evidence proves what occurred.

### 12.4 CLI Conventions

```bash
pactyra engagement init
pactyra scope compile scope.yaml
pactyra run discovery
pactyra evidence verify receipt.json
pactyra report export
```

### 12.5 Package Naming

```text
pactyra-core
pactyra-scope
pactyra-policy
pactyra-runner
pactyra-evidence
pactyra-labs
pactyra-sdk
```

### 12.6 What to Avoid

- "Pactyra Offensive Security Platform" — use the descriptor instead
- "Ark-Pactyra", "Pactyra-OS", "Pactyra OS" — implies an OS/distribution (a non-goal)
- Hyphenated or compound variants of the name
- "Ark-Onyx" in any public-facing context (retained only as internal historical codename)

---

## Appendix A: Glossary

| Term                         | Definition                                                                                                                                                                                |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Operations kernel**        | The system that manages authorization, orchestration, and evidence for offensive operations. The operator still decides what to do; the kernel enforces how.                              |
| **O0–O6 permission lattice** | A seven-level trust hierarchy. O0 = external reconnaissance only. O6 = full exploitation with credential reuse and lateral movement. Each level subsumes the permissions of lower levels. |
| **Evidence receipt**         | An immutable, cryptographically-signed record of an action taken during an engagement. Contains the action, its parameters, output hash, timestamp, operator identity, and scope context. |
| **Attack graph**             | A directed graph representing the state of an engagement. Nodes are hosts, services, vulnerabilities, and credentials. Edges represent reachability, exploitation, and pivot paths.       |
| **Adversarial validation**   | The process of proving (not just detecting) that a vulnerability is exploitable. Produces a proof-of-exploitation receipt distinct from a scanner finding.                                |
| **Scope document**           | A machine-readable specification of authorized targets, excluded hosts, permitted techniques, time windows, and contact information. Cryptographically bound to the engagement.           |
| **Design partner**           | A real offensive security team that uses Pactyra during development, provides feedback, and validates that the platform solves real problems.                                             |

---

## Appendix B: Architectural Principles

These principles guide all technical decisions:

1. **The operator is in control.** Agents propose; operators dispose. No autonomous exploitation.
2. **Evidence is a byproduct, not a chore.** Every action produces a receipt automatically. The operator never "creates evidence" — they operate, and receipts emerge.
3. **The terminal is the interface.** A web dashboard is a view, not the control surface. Operations happen in the terminal.
4. **Scope is enforced, not advised.** A scope violation is blocked, not warned. Warnings become noise.
5. **The moat is integration.** No single feature is the differentiator. The integrity chain from scope → execution → evidence → graph → validation → report → retesting is the product.
6. **Start with Rust and Go.** Trust boundary in Rust (correctness, safety). Orchestration in Go (concurrency, networking, Docker SDK). Add Python and TypeScript when the core is proven.
7. **Ship v0.1 in 12 weeks or less.** A working operations kernel running one tool for one operator is infinitely more valuable than a complete architecture document.
