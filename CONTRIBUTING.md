# Contributing to Pactyra

Thank you for your interest in contributing to Pactyra! This document provides guidelines for contributing to the project.

> **Note:** Pactyra is in early development (SDD-000 — governance phase). Contribution workflows will be refined as the project matures through v0.1 and v0.5.

---

## Code of Conduct

Pactyra follows a simple principle: **be professional, be constructive, be respectful.** Our community includes operators, researchers, and contributors from diverse backgrounds. Harassment, personal attacks, and unprofessional behavior are not tolerated.

A formal Code of Conduct will be established before v0.5.

---

## How to Contribute

### Current Phase (SDD-000)

We are in the planning and governance phase. The best ways to contribute right now:

1. **Review the artifacts:** Read the [Proposal](openspec/changes/SDD-000/proposal.md), [Spec](openspec/changes/SDD-000/spec.md), and [Design](openspec/changes/SDD-000/design.md). Provide feedback via issues.
2. **Design partner interest:** If you're an offensive security operator interested in being a design partner, reach out. See [Design Partner Strategy](docs/design-partner-strategy.md).
3. **RFC discussions:** Open an issue for architectural feedback, use-case suggestions, or tool ecosystem ideas.

### Implementation Phase (SDD-002+)

Once implementation begins, we welcome:

- **Code contributions:** Pick a task from the [tasks artifact](openspec/changes/SDD-000/tasks.md) or an open issue.
- **Tool adapters:** Contribute adapters for additional security tools (v0.5+).
- **Documentation:** Improve docs, add examples, translate.
- **Testing:** Write integration tests, expand the cyber range, reproduce benchmarks.

---

## Development Setup

### Prerequisites

- **Rust** (stable toolchain): Install via [rustup](https://rustup.rs/)
- **Go** (1.22+): Install via [golang.org](https://golang.org/dl/)
- **Docker Engine**: Install via [docker.com](https://docker.com/)
- **Nmap** (optional, for adapter testing)

### Building

```bash
git clone https://github.com/arkelythex/pactyra.git
cd pactyra

# Build both binaries
make build

# Run tests
make test

# Run linters
make lint

# Full golden path (requires Docker)
make golden-path
```

### Project Structure

```
ark-onyx-offensive/
├── crates/                    # Rust workspace (trust boundary)
│   ├── pactyra-scope/         # Scope parser, compiler, hash
│   ├── pactyra-guardian/      # Runtime scope enforcement
│   ├── pactyra-chain/         # Receipt chain + signing
│   └── pactyra-core/          # Rust binary (IPC server)
├── cmd/pactyra/               # Go binary entry point
├── internal/                  # Go internal packages
│   ├── tui/                   # Terminal UI
│   ├── orchestrator/          # Session state, agent coordination
│   ├── guardian/              # Rust IPC client
│   ├── agent/                 # Agent implementations
│   ├── worker/                # Docker worker management
│   ├── adapter/               # Tool adapters (Nmap)
│   └── export/                # Evidence pack export
├── tests/                     # Integration tests
├── cyber-range/               # Local test infrastructure
└── docs/                      # Documentation
```

---

## Pull Request Guidelines

1. **Keep PRs focused.** One concern per PR (SDD-000 change, bug fix, feature).
2. **Follow the chained PR strategy.** The project uses stacked-to-main chained PRs for large changes. See [SDD Dependency Map](docs/sdd-dependency-map.md).
3. **Write meaningful commit messages.** Use conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`.
4. **Include tests.** New features should include tests. Bug fixes should include regression tests.
5. **Run `make lint` and `make test`** before pushing.
6. **No AI-generated commit attribution.** Do not include "Co-Authored-By: AI" in commits.

---

## Tool Adapter SDK (v0.5+)

In v0.5, Pactyra will ship a Tool Adapter SDK for community-contributed tool integrations. Adapters implement the `adapter.Adapter` Go interface:

- Validate tool-specific parameters
- Translate intent into tool CLI flags
- Parse tool output into structured results

If you're interested in building an adapter, watch for the SDK release in v0.5.

---

## Design Partner Program

We are recruiting 3+ offensive security teams to use Pactyra on real engagements and provide feedback during development. Design partners get:

- Early access to Pactyra builds
- Direct line to the development team
- Influence over the roadmap
- Name in acknowledgments

See [Design Partner Strategy](docs/design-partner-strategy.md) for details.

---

## Recognition

All contributors will be recognized in our contributors list. Significant contributions (major features, tool adapters, security research) will be acknowledged in release notes.

---

## Questions?

Open an issue or start a discussion. For security-sensitive matters, see [SECURITY.md](SECURITY.md).
