# PACTYRA

## Open Offensive Operations Kernel

**Bound by scope. Proven by evidence.**

Pactyra is an open-source, terminal-first offensive operations kernel for authorized, policy-constrained, and evidence-backed security assessments. It is an **architecture and specification preview** developed under [Arkelythex](https://github.com/arkelythex) that defines a technical thesis: authorized offensive automation can be safe, open, auditable, and reproducible.

---

## Why Pactyra?

Offensive security operations are fragmented across tools, scripts, terminals, and notes that share no verifiable model of authorization, execution, and evidence. Pactyra provides that model. It **binds every action to a cryptographically-verified scope**, records every decision in an **append-only evidence chain**, and produces **exportable, independently-verifiable evidence packs**.

**The problem Pactyra solves:** When a client disputes a finding — "we fixed that," "that's out of scope," "prove you actually exploited this" — the operator's defense today is screenshots, terminal scrollback, and memory. None of it is cryptographically verifiable. Pactyra changes that.

---

## Repository Model

| Property                | Commitment                                                                                                                  |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **License**             | AGPL-3.0 — GNU Affero General Public License v3.0. See [LICENSE](LICENSE). Strong copyleft — protects the open-source core. |
| **Distribution**        | Self-hostable. Clone, build, run. No mandatory services.                                                                    |
| **Interface**           | Terminal-first. Operations happen in the terminal.                                                                          |
| **AI dependency**       | None. AI is an optional accelerator, never a requirement.                                                                   |
| **Provider lock-in**    | None. Provider-agnostic by design.                                                                                          |
| **Commercial services** | Permitted in the future. Must not distort or weaken the open-source core.                                                   |

---

## v0.1 Golden Path (First Light — planned)

Pactyra v0.1 is designed to demonstrate the complete golden path against a local cyber range. The executable path is under development; the repository currently documents the contracts, trust boundaries, and acceptance criteria:

```text
1. Load and compile an authorized scope
2. Generate a cryptographic scope_hash
3. Execute permitted discovery against a local cyber range
4. Deliberately attempt an out-of-scope action
5. Block it before execution (no network traffic sent)
6. Record every action and decision as an evidence receipt
7. Generate verifiable receipts (integrity + provenance)
8. Export a reproducible evidence pack (verifiable without the platform)
```

---

## Architecture (Two-Binary Model)

```
┌──────────────────────────────────────────┐
│              pactyra (Go)                 │
│  TUI ◄──► Orchestrator ◄──► Agents       │
│                │                          │
│                ▼                          │
│           Worker Manager                  │
│                │                          │
└────────────────┼──────────────────────────┘
                 │  Unix Socket (local only)
┌────────────────┼──────────────────────────┐
│           pactyra-core (Rust)             │
│  Scope Compiler  │  Scope Guardian        │
│                  │                        │
│  Receipt Chain (append-only, Ed25519)     │
│  Key Storage: ~/.pactyra/keys/           │
└───────────────────────────────────────────┘
                 │  Docker API
┌────────────────┼──────────────────────────┐
│         Docker Workers (ephemeral)         │
│  Nmap container — scope-only egress       │
└───────────────────────────────────────────┘
```

**Rust owns every trust decision.** The Go binary orchestrates but must ask Rust before any tool hits the network. Docker workers are fully isolated with network egress restricted to scope-authorized targets.

---

## Quickstart

> **Status:** Specification and architecture preview. No runnable golden path is published yet; implementation is planned in the SDD roadmap.

```bash
# Prerequisites: Rust, Go 1.22+, Docker

git clone https://github.com/Dreamcoder08/pactyra.git
cd pactyra
make build
make golden-path
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

---

## Documentation

- [Design Document](docs/architecture.md) — System architecture, trust boundaries, data flow
- [Golden Path Walkthrough](docs/golden-path.md) — Step-by-step v0.1 demonstration
- [Evidence Model Review](docs/evidence-model-review.md) — Integrity, provenance, reproducibility
- [Architecture Review](docs/architecture-review.md) — Two-binary trust boundary validation
- [Scope Expansion Review](docs/scope-expansion-review.md) — Expansion envelope edge cases
- [SDD Dependency Map](docs/sdd-dependency-map.md) — Downstream SDD relationships
- [Cyber Range](cyber-range/README.md) — Local test infrastructure

---

## Roadmap

| Version | Name            | Timeline                | Key Capabilities                                                   |
| ------- | --------------- | ----------------------- | ------------------------------------------------------------------ |
| v0.1    | First Light     | 8–12 weeks              | Scope enforcement, evidence receipts, Nmap, TUI, standalone export |
| v0.5    | Community Ready | 6–9 months after v0.1   | Attack graph, adversarial validation, retesting, tool adapter SDK  |
| v1.0    | Reference Impl. | 12–18 months after v0.1 | Web dashboard, diff engine, 20+ tools, community marketplace       |

---

## Community

- **Design Partner Program:** We're recruiting 3+ offensive security teams to use Pactyra on real engagements and provide feedback. See [Design Partner Strategy](docs/design-partner-strategy.md).
- **Funding:** Pactyra is a public-interest project. See [Grant Landscape](docs/grant-landscape.md) and [Grant Narrative](docs/grant-narrative.md).
- **Conferences:** See [Community Strategy](docs/community-strategy.md) for upcoming talks and CFP deadlines.

---

## Security

See [SECURITY.md](SECURITY.md) for our security policy, trust model summary, and reporting contact.

---

## License

AGPL-3.0 — GNU Affero General Public License v3.0. See [LICENSE](LICENSE) for the full text. Strong copyleft that protects the open-source core from commercial enclosure, compatible with grant and accelerator programs.
