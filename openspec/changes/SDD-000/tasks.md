# Tasks: Pactyra — SDD-000 Vision & Governance

**Change ID:** SDD-000  
**Phase:** sdd-tasks  
**Date:** 2026-07-26  
**Status:** Reviewed

---

## Review Workload Forecast

| Field                   | Value                                                                                                                                                               |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Estimated changed lines | 8,000–15,000 (SDD-002+ v0.1 implementation)                                                                                                                         |
| 400-line budget risk    | High                                                                                                                                                                |
| Chained PRs recommended | Yes                                                                                                                                                                 |
| Suggested split         | PR 1: Repository + Rust crates (scope, guardian, chain) → PR 2: Go orchestrator + TUI → PR 3: Docker workers + Nmap adapter → PR 4: Integration tests + golden path |
| Delivery strategy       | ask-on-risk                                                                                                                                                         |
| Chain strategy          | pending                                                                                                                                                             |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Forecast Notes

SDD-000 itself produces zero implementation code — it is a governance change. The forecast above covers SDD-002 through SDD-0XX (v0.1 implementation). The Rust trust core alone (4 crates) will exceed 400 changed lines; the Go orchestrator with TUI, Docker workers, and agents will be 5,000+ lines. Integration tests add another 1,500–3,000. Chained PRs are strongly recommended. The chain strategy should be resolved before SDD-002 begins.

---

## Overview

SDD-000 is the **vision, boundaries, and non-goals** change. These tasks are **planning and governance tasks** — they solidify the project foundation and prepare downstream SDDs for implementation. No code is produced in SDD-000. Implementation begins in SDD-002.

### Task Organization

Tasks are grouped into five phases, ordered by dependency:

| Phase | Group                       | Dependency |
| ----- | --------------------------- | ---------- |
| 1     | Artifact Finalization       | None       |
| 2     | Technical Validation        | Phase 1    |
| 3     | Repository & Infrastructure | Phase 2    |
| 4     | Downstream SDD Preparation  | Phase 3    |
| 5     | Community & Funding         | Phase 4    |

---

## Phase 1: Artifact Finalization

### T-001: Resolve Open Questions in Proposal

- [x] **Answered:** compliance frameworks → guidance not blocker; failed engagements → immutable no deletion; tool ecosystem → CLI-only structured output. Documented in proposal §9.4. <!-- sdd-owner: implementation -->
- [x] **Decision documented** in proposal §9.4 with resolution and rationale. All answered now, none deferred. <!-- sdd-owner: implementation -->
- [x] **Not needed** — all questions answered now, none deferred. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- All three open questions have documented resolutions in `openspec/changes/SDD-000/proposal.md`.
- Any deferred question has a target SDD identifier.

**Estimated effort:** 1 session (1–2 hours)

---

### T-002: Finalize License Selection

- [x] **AGPL-3.0 selected.** Research complete — AGPL-3.0 best balances protection from commercial enclosure, grant eligibility, and community adoption. See proposal §2.3 and §9.4. <!-- sdd-owner: implementation -->
- [x] **AGPL-3.0 evaluated.** Strong copyleft protects core, compatible with grants, commercial services can negotiate separate licensing. <!-- sdd-owner: implementation -->
- [x] **License documented** in proposal §2.3 with rationale. <!-- sdd-owner: implementation -->
- [x] `LICENSE` file created with AGPL-3.0. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- A specific license is selected and documented with rationale.
- `LICENSE` file exists in repository root.
- The license choice appears in the proposal's repository model table.

**Estimated effort:** 1 session (1–2 hours)

---

### T-003: Finalize Proposal, Spec, and Design After User Review

- [x] Incorporate user feedback from interactive review of proposal, spec, and design. <!-- sdd-owner: implementation -->
- [x] Ensure all three artifacts are internally consistent (no contradictions between proposal scope, spec requirements, design architecture). <!-- sdd-owner: implementation -->
- [x] Cross-check: every spec requirement traces to a proposal section; every design decision traces to a spec requirement. <!-- sdd-owner: implementation -->
- [x] Update status from "Draft" to "Reviewed" in all three artifacts. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- All three artifacts have status "Reviewed."
- No contradictions across proposal → spec → design.
- Every spec requirement has a clear proposal origin.

**Estimated effort:** 1 session (1–2 hours)

---

### T-004: Finalize Tasks Artifact (This Document)

- [x] After user review of this task breakdown, incorporate feedback. <!-- sdd-owner: implementation -->
- [x] Ensure all tasks have clear acceptance criteria and estimated effort. <!-- sdd-owner: implementation -->
- [x] Update status to "Reviewed." <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- This document has status "Reviewed."
- All tasks accepted by user.

**Estimated effort:** Tied to interactive review session.

---

## Phase 2: Technical Validation

### T-005: Validate Two-Binary Architecture (Rust + Go)

- [x] Verify the IPC boundary (Unix domain socket) design: confirm that all trust-decision paths flow through Rust and no Go code path can bypass scope validation. <!-- sdd-owner: implementation -->
- [x] Review the design's trust domain diagram (§5.4) against the component interaction contracts (§6.1): confirm every Go → Rust operation is accounted for. <!-- sdd-owner: implementation -->
- [x] Identify any architectural gap where the Go binary could: (a) execute a tool without Rust approval, (b) modify a receipt after creation, (c) access the signing key, (d) suppress receipt generation. <!-- sdd-owner: implementation -->
- [x] Document findings. If gaps found, update design §5 before SDD-002 begins. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Architectural review confirms zero paths for Go to bypass Rust trust decisions.
- Any found gaps are documented and resolved in the design.
- Review notes saved to Engram (`sdd/SDD-000/architecture-review`).

**Estimated effort:** 1 session (2–3 hours)

---

### T-006: Validate v0.1 Golden Path Completeness

- [x] Walk the 8-step golden path from proposal §10 and spec "Core Journey" section. <!-- sdd-owner: implementation -->
- [x] For each step, identify: what component(s) are involved, what data flows between them, what the acceptance test would verify. <!-- sdd-owner: implementation -->
- [x] Confirm that every step is testable with the v0.1 architecture (no v0.5+ dependencies). <!-- sdd-owner: implementation -->
- [x] Document the golden path walkthrough in `docs/golden-path.md` (create if needed). <!-- sdd-owner: implementation -->
- [x] If any step is untestable with v0.1 components, update design to close the gap or defer the step. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- All 8 golden path steps are validated as testable with v0.1 architecture.
- `docs/golden-path.md` exists with step-by-step walkthrough.
- No step requires v0.5+ components.

**Estimated effort:** 1 session (2–3 hours)

---

### T-007: Validate Evidence Model Claims

- [x] Review the three evidence layers (integrity, provenance, reproducibility) against the design's receipt chain architecture. <!-- sdd-owner: implementation -->
- [x] Confirm: integrity layer is achievable with hash chain + Ed25519 signing in v0.1. <!-- sdd-owner: implementation -->
- [x] Confirm: provenance layer is achievable with receipt fields (who, when, tool, authorization) in v0.1. <!-- sdd-owner: implementation -->
- [x] Confirm: reproducibility layer is deferred to v0.5 as stated. <!-- sdd-owner: implementation -->
- [x] Verify that every claim about what each layer does NOT guarantee is architecturally honest (no overpromise). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Integrity and provenance layers are confirmed achievable in v0.1.
- Reproducibility layer is confirmed deferred to v0.5 with clear architectural path.
- No evidence claim in spec is architecturally unsupported.

**Estimated effort:** 1 session (1–2 hours)

---

### T-008: Validate Scope Expansion Envelope Design

- [x] Review the expansion envelope policy structure (proposal §6.1, spec "Scope Expansion Envelope Policy"). <!-- sdd-owner: implementation -->
- [x] Confirm the YAML shape is parseable and the policy rules are enforceable in code (not just documentation). <!-- sdd-owner: implementation -->
- [x] Identify edge cases: overlapping expansions, concurrent expansion requests, expansion rollback, expansion after engagement close. <!-- sdd-owner: implementation -->
- [x] Document edge cases for SDD-002 (scope model) to address. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Expansion envelope policy is validated as enforceable.
- Edge cases are documented for SDD-002.
- No architectural blocker for v0.5 implementation.

**Estimated effort:** 1 session (1–2 hours)

---

## Phase 3: Repository & Infrastructure

### T-009: Initialize Monorepo Structure

- [x] Create the directory structure as defined in design §8 (Repository Structure). <!-- sdd-owner: implementation -->
- [x] `crates/pactyra-scope/`, `crates/pactyra-guardian/`, `crates/pactyra-chain/`, `crates/pactyra-core/` with placeholder `Cargo.toml` files. <!-- sdd-owner: implementation -->
- [x] `cmd/pactyra/` with placeholder `main.go`. <!-- sdd-owner: implementation -->
- [x] `internal/` directory tree with placeholder packages: `tui/`, `orchestrator/`, `guardian/`, `agent/`, `worker/`, `adapter/`, `export/`. <!-- sdd-owner: implementation -->
- [x] `tests/` directory with placeholder integration test files. <!-- sdd-owner: implementation -->
- [x] `cyber-range/` directory with placeholder `docker-compose.yml`. <!-- sdd-owner: implementation -->
- [x] `docs/` directory with architecture.md symlink/copy and placeholder files. <!-- sdd-owner: implementation -->
- [x] `openspec/` directory already exists; confirm structure is correct. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Directory structure matches design §8 exactly.
- Placeholder files exist so SDD-002+ can start writing code immediately.
- `Cargo.toml` workspace manifest includes all four crates.
- `go.mod` initialized with module path `github.com/arkelythex/pactyra`.

**Estimated effort:** 1 session (1–2 hours)

---

### T-010: Set Up Build System (Makefile)

- [x] Create `Makefile` with targets: `build`, `test`, `lint`, `run`, `golden-path`, `clean` as specified in design §8.1. <!-- sdd-owner: implementation -->
- [x] `make build`: runs `cargo build --release` for Rust workspace, `go build ./cmd/pactyra` for Go binary. <!-- sdd-owner: implementation -->
- [x] `make test`: runs `cargo test` and `go test ./...`. <!-- sdd-owner: implementation -->
- [x] `make lint`: runs `cargo clippy` and `golangci-lint run`. <!-- sdd-owner: implementation -->
- [x] `make run`: builds both binaries and launches TUI with a test scope fixture. <!-- sdd-owner: implementation -->
- [x] `make golden-path`: runs the full 8-step integration test. <!-- sdd-owner: implementation -->
- [x] `make clean`: removes build artifacts. <!-- sdd-owner: implementation -->
- [x] Verify all targets succeed (even if tests are empty/skipped — build infrastructure must work). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- All six Makefile targets exist and execute without error.
- `make build` produces both binaries (even if functional stubs).
- `make test` runs (even if no real tests yet).
- `make lint` runs (even if no code to lint yet).

**Estimated effort:** 1 session (1–2 hours)

---

### T-011: Set Up CI Pipeline

- [x] Create `.github/workflows/ci.yml` (or equivalent for chosen CI platform). <!-- sdd-owner: implementation -->
- [x] CI pipeline: on push/PR to `main`, run `make lint`, `make test`, `make build`. <!-- sdd-owner: implementation -->
- [x] Ensure CI uses a Rust toolchain (stable) and Go toolchain (1.22+). <!-- sdd-owner: implementation -->
- [x] Add Docker availability to CI for worker tests. <!-- sdd-owner: implementation -->
- [x] Add `make golden-path` as a nightly or manual trigger (not per-PR — needs cyber range). <!-- sdd-owner: implementation -->
- [x] Verify CI green on the repository initialization commit. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- CI pipeline runs on push/PR.
- `lint`, `test`, `build` all pass.
- Golden path test is configured (even if skipped pending implementation).

**Estimated effort:** 1 session (1–2 hours)

---

### T-012: Create Public README and Project Identity

- [x] Write `README.md` following the spec's "Platform Identity" requirement: canonical descriptor "Open Offensive Operations Kernel," tagline "Bound by scope. Proven by evidence.", parent org "Arkelythex." <!-- sdd-owner: implementation -->
- [x] Include: public-interest thesis summary, repository model table, v0.1 golden path diagram (ASCII), quickstart instructions (`make build && make golden-path`), link to full documentation. <!-- sdd-owner: implementation -->
- [x] Include: explicit statements that AI is optional, provider-agnostic, no mandatory services. <!-- sdd-owner: implementation -->
- [x] Create `SECURITY.md` with security policy, reporting contact, and trust model summary. <!-- sdd-owner: implementation -->
- [x] Create `CONTRIBUTING.md` with contribution guidelines (placeholder for now; adapter SDK docs in v0.5). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- `README.md` uses canonical descriptor and tagline exactly.
- Repository model is transparent (license, self-hostable, AI-optional, provider-agnostic).
- `SECURITY.md` and `CONTRIBUTING.md` exist.

**Estimated effort:** 1 session (1–2 hours)

---

### T-013: Set Up Cyber Range for Local Testing

- [x] Create `cyber-range/docker-compose.yml` with 2–3 vulnerable-by-design target services (e.g., DVWA-like web app, a service with open ports, a service with known CVEs). <!-- sdd-owner: implementation -->
- [x] Ensure the range is self-contained (no internet required to test). <!-- sdd-owner: implementation -->
- [x] Write `cyber-range/README.md` with: how to start, what targets exist, which ports are open, what scope.yaml to use. <!-- sdd-owner: implementation -->
- [x] Create a matching `scope.yaml` fixture in `tests/fixtures/scope_valid.yaml` that authorizes the cyber range targets. <!-- sdd-owner: implementation -->
- [x] Verify that `docker-compose up` starts all target services and they are reachable from the host. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- `docker-compose up` starts target services.
- `scope_valid.yaml` correctly describes the range.
- Range is documented so a new contributor can start testing in <5 minutes.

**Estimated effort:** 1 session (2–4 hours)

---

## Phase 4: Downstream SDD Preparation

### T-014: Validate and Document SDD Dependency Map

- [x] Review the SDD dependency table from proposal §11.2. <!-- sdd-owner: implementation -->
- [x] Add any missing downstream SDDs identified during spec/design phases. <!-- sdd-owner: implementation -->
- [x] For each downstream SDD: confirm the dependency list is correct, estimate sequencing (which can run in parallel, which must be serial). <!-- sdd-owner: implementation -->
- [x] Document the final dependency map in `docs/sdd-dependency-map.md`. <!-- sdd-owner: implementation -->
- [x] Identify the minimal set of SDDs needed to ship v0.1 (likely: SDD-002, SDD-005, SDD-007, SDD-0XX). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Complete SDD dependency map documented.
- v0.1-minimal SDD set identified.
- Parallel vs. serial dependencies clarified.

**Estimated effort:** 1 session (1–2 hours)

---

### T-015: Prepare SDD-002 Brief (Scope Model and Compiler)

- [x] Extract all scope-related requirements from spec and design into a brief for SDD-002. <!-- sdd-owner: implementation -->
- [x] Include: CIDR parsing, domain list parsing, excluded hosts, port ranges, scope_hash generation, canonical representation format, expansion envelope schema (for v0.5). <!-- sdd-owner: implementation -->
- [x] Include: the scope compilation flow from design §4.1. <!-- sdd-owner: implementation -->
- [x] Include: scope enforcement invariants from spec (block, not warn; every action references scope_hash). <!-- sdd-owner: implementation -->
- [x] Identify what SDD-002 owns vs. what SDD-003 (permission lattice) and SDD-005 (tool contracts) own. <!-- sdd-owner: implementation -->
- [x] Save brief to Engram (`sdd/SDD-002/brief`) and as `openspec/changes/SDD-002/brief.md`. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- SDD-002 brief exists with clear boundaries.
- All scope requirements are extracted and assigned.
- Dependencies on SDD-003 and SDD-005 are documented.

**Estimated effort:** 1 session (1–2 hours)

---

### T-016: Prepare SDD-007 Brief (Evidence and Receipts)

- [x] Extract all evidence-related requirements from spec and design into a brief for SDD-007. <!-- sdd-owner: implementation -->
- [x] Include: receipt schema (fields: action type, parameters, output hash, timestamp, operator, scope_hash, parent hash), hash chain structure, Ed25519 signing, append-only guarantee, chain verification, evidence export format. <!-- sdd-owner: implementation -->
- [x] Include: the three-layer trust model (integrity, provenance, reproducibility) and what each layer guarantees and does not guarantee. <!-- sdd-owner: implementation -->
- [x] Include: the evidence export flow from design §4.4 and the standalone verification requirement. <!-- sdd-owner: implementation -->
- [x] Identify what SDD-007 owns vs. what SDD-010 (finding lifecycle) owns. <!-- sdd-owner: implementation -->
- [x] Save brief to Engram (`sdd/SDD-007/brief`) and as `openspec/changes/SDD-007/brief.md`. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- SDD-007 brief exists with clear boundaries.
- All evidence requirements are extracted and assigned.
- Boundary with SDD-010 (finding lifecycle) is clear.

**Estimated effort:** 1 session (1–2 hours)

---

### T-017: Prepare SDD-005 Brief (Tool and Agent Contracts)

- [x] Extract all tool/agent-related requirements from spec and design into a brief for SDD-005. <!-- sdd-owner: implementation -->
- [x] Include: Go `agent.Agent` interface, Go `adapter.Adapter` interface, the "propose → validate → execute → record" contract, Scope Guardian agent role, Recon Planner agent role, Worker Manager contract, Nmap adapter specification. <!-- sdd-owner: implementation -->
- [x] Include: the tool invocation flow from design §4.2. <!-- sdd-owner: implementation -->
- [x] Include: the scope violation flow from design §4.3. <!-- sdd-owner: implementation -->
- [x] Save brief to Engram (`sdd/SDD-005/brief`) and as `openspec/changes/SDD-005/brief.md`. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- SDD-005 brief exists with clear boundaries.
- Agent and adapter interfaces are specified at the architectural level.
- Invocation and violation flows are documented.

**Estimated effort:** 1 session (1–2 hours)

---

### T-018: Prepare SDD-003 Brief (Permission Lattice)

- [x] Extract all permission-related requirements from spec and design into a brief for SDD-003. <!-- sdd-owner: implementation -->
- [x] Include: O0–O6 level definitions, capabilities per level, enforcement rules, escalation procedure. <!-- sdd-owner: implementation -->
- [x] Include: v0.1 scope (O0–O3 only: basic IP/port gating) vs. v0.5 scope (O0–O6 full lattice). <!-- sdd-owner: implementation -->
- [x] Save brief to Engram (`sdd/SDD-003/brief`) and as `openspec/changes/SDD-003/brief.md`. <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- SDD-003 brief exists with clear boundaries.
- O0–O6 levels are defined at the architectural level.
- v0.1 vs. v0.5 scope is clear.

**Estimated effort:** 1 session (1 hour)

---

## Phase 5: Community & Funding

### T-019: Design Partner Recruitment Strategy

- [x] Draft a design partner recruitment brief: target profile (small offensive security firms, 3–10 operators, similar to Alex Chen's Aegis Security Partners), value proposition for partners, what partners commit to (use on real engagements, provide feedback), what Pactyra provides (early access, direct line to development, name in acknowledgments). <!-- sdd-owner: implementation -->
- [x] Identify 5–10 candidate firms or individual operators through professional network, LinkedIn, security conferences, CTF communities. <!-- sdd-owner: implementation -->
- [x] Prepare outreach materials: one-page overview, demo concept (the golden path), FAQ. <!-- sdd-owner: implementation -->
- [x] Set target: 3+ design partners committed before SDD-002 implementation begins. <!-- sdd-owner: implementation -->
- [x] Save strategy to Engram (`sdd/SDD-000/design-partner-strategy`). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Recruitment brief documented.
- Candidate list of 5–10 firms/operators identified.
- Outreach materials prepared.
- Target of 3+ partners defined with timeline.

**Estimated effort:** 2 sessions (3–5 hours)

---

### T-020: Grant and Funding Landscape Research

- [x] Research grant programs compatible with open-source security tooling: NLnet, Open Technology Fund, Sovereign Tech Fund, Google OSS Pilot, GitHub Accelerator, AWS Open Source Credits. <!-- sdd-owner: implementation -->
- [x] Research security-specific accelerators and incubators: DARPA, Cyber Catalyst, NATO DIANA, YC (security track). <!-- sdd-owner: implementation -->
- [x] For each program, document: eligibility criteria, funding amount, application window, requirements (open source, public good, security focus). <!-- sdd-owner: implementation -->
- [x] Prioritize 3–5 programs that best align with Pactyra's public-interest thesis. <!-- sdd-owner: implementation -->
- [x] Save research to Engram (`sdd/SDD-000/grant-landscape`). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Grant landscape documented with 10+ programs.
- Top 3–5 programs prioritized with rationale.
- Application windows noted for the next 12 months.

**Estimated effort:** 2 sessions (3–5 hours)

---

### T-021: Prepare Initial Grant Application Package

- [x] Draft a reusable grant narrative: problem statement, solution overview, public-interest thesis, technical approach, team background, milestones, budget estimate. <!-- sdd-owner: implementation -->
- [x] Include: why this is public infrastructure, not a commercial product; why open-source; why now. <!-- sdd-owner: implementation -->
- [x] Include: v0.1 proof-of-concept as first milestone, v0.5 community-ready as second, v1.0 reference implementation as third. <!-- sdd-owner: implementation -->
- [x] Prepare a technical appendix: architecture diagram, trust boundary model, golden path walkthrough. <!-- sdd-owner: implementation -->
- [x] Save narrative package to Engram (`sdd/SDD-000/grant-narrative`). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Reusable grant narrative exists covering all key sections.
- Technical appendix is ready to attach.
- Narrative can be adapted to specific grant applications with <1 hour of tailoring.

**Estimated effort:** 2 sessions (4–6 hours)

---

### T-022: Conference and Community Presence Strategy

- [x] Identify 3–5 security conferences suitable for presenting Pactyra's thesis (BSides, DEF CON, Black Hat, OffensiveCon, Hack in the Box). <!-- sdd-owner: implementation -->
- [x] Identify CFP (Call for Papers) deadlines for the next 12 months. <!-- sdd-owner: implementation -->
- [x] Draft a talk abstract: "Authorized Offensive Automation with Verifiable Evidence — An Open-Source Reference Implementation." <!-- sdd-owner: implementation -->
- [x] Identify online communities for initial visibility: r/netsec, r/redteamsec, Hack The Box forums, Offensive Security Discord servers. <!-- sdd-owner: implementation -->
- [x] Save strategy to Engram (`sdd/SDD-000/community-strategy`). <!-- sdd-owner: implementation -->

**Acceptance criteria:**

- Conference target list with CFP deadlines.
- Talk abstract drafted.
- Online community engagement plan documented.

**Estimated effort:** 1 session (2–3 hours)

---

## Dependency Graph

```
Phase 1 (Artifacts)
  T-001 ──┐
  T-002 ──┤
  T-003 ──┤──> Phase 2 (Validation)
  T-004 ──┘       T-005 ──┐
                  T-006 ──┤
                  T-007 ──┤──> Phase 3 (Infrastructure)
                  T-008 ──┘       T-009 ──┐
                                  T-010 ──┤
                                  T-011 ──┤──> Phase 4 (Downstream SDDs)
                                  T-012 ──┤       T-014 ──┐
                                  T-013 ──┘       T-015 ──┤
                                                  T-016 ──┤──> Phase 5 (Community)
                                                  T-017 ──┤       T-019 ──┐
                                                  T-018 ──┘       T-020 ──┤
                                                                  T-021 ──┤
                                                                  T-022 ──┘
```

### Parallelizable Tasks

Within each phase, these tasks can run in parallel:

- **Phase 1:** T-001, T-002, T-003 (independent concerns)
- **Phase 2:** T-005, T-006, T-007, T-008 (different review areas)
- **Phase 3:** T-012, T-013 (independent of T-009, T-010, T-011)
- **Phase 4:** T-015, T-016, T-017, T-018 (different downstream SDDs)
- **Phase 5:** T-019, T-020, T-022 (different research areas)

---

## Task Summary

| ID    | Task                                | Phase              | Est. Effort    | Depends On   |
| ----- | ----------------------------------- | ------------------ | -------------- | ------------ |
| T-001 | Resolve open proposal questions     | 1: Artifacts       | 1–2h           | None         |
| T-002 | Finalize license selection          | 1: Artifacts       | 1–2h           | None         |
| T-003 | Finalize proposal, spec, design     | 1: Artifacts       | 1–2h           | T-001, T-002 |
| T-004 | Finalize tasks artifact             | 1: Artifacts       | Review session | T-003        |
| T-005 | Validate two-binary architecture    | 2: Validation      | 2–3h           | Phase 1      |
| T-006 | Validate golden path completeness   | 2: Validation      | 2–3h           | Phase 1      |
| T-007 | Validate evidence model claims      | 2: Validation      | 1–2h           | Phase 1      |
| T-008 | Validate scope expansion envelope   | 2: Validation      | 1–2h           | Phase 1      |
| T-009 | Initialize monorepo structure       | 3: Infrastructure  | 1–2h           | Phase 2      |
| T-010 | Set up build system (Makefile)      | 3: Infrastructure  | 1–2h           | T-009        |
| T-011 | Set up CI pipeline                  | 3: Infrastructure  | 1–2h           | T-010        |
| T-012 | Create public README and identity   | 3: Infrastructure  | 1–2h           | Phase 2      |
| T-013 | Set up cyber range                  | 3: Infrastructure  | 2–4h           | Phase 2      |
| T-014 | Validate SDD dependency map         | 4: Downstream SDDs | 1–2h           | Phase 3      |
| T-015 | Prepare SDD-002 brief (scope)       | 4: Downstream SDDs | 1–2h           | T-014        |
| T-016 | Prepare SDD-007 brief (evidence)    | 4: Downstream SDDs | 1–2h           | T-014        |
| T-017 | Prepare SDD-005 brief (agents)      | 4: Downstream SDDs | 1–2h           | T-014        |
| T-018 | Prepare SDD-003 brief (permissions) | 4: Downstream SDDs | 1h             | T-014        |
| T-019 | Design partner recruitment strategy | 5: Community       | 3–5h           | Phase 4      |
| T-020 | Grant landscape research            | 5: Community       | 3–5h           | Phase 4      |
| T-021 | Prepare grant application package   | 5: Community       | 4–6h           | T-020        |
| T-022 | Conference and community strategy   | 5: Community       | 2–3h           | Phase 4      |

**Total estimated effort:** 30–50 hours across 22 tasks.

---

## Post-Tasks: Bounded Review (Parent-Owned)

- [ ] Start or reuse bounded review after all Phase 1–5 tasks complete. <!-- sdd-owner: parent -->
- [ ] Validate SDD-000 artifacts are internally consistent and ready for archive. <!-- sdd-owner: parent -->
- [ ] Gate: confirm all downstream SDD briefs are complete before `sdd-archive`. <!-- sdd-owner: parent -->
