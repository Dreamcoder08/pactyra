# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Pactyra, please report it to the maintainers directly. **Do not open a public issue.**

**Contact:** <security@arkelythex.dev> — will be active before v0.1 release.

We take all security reports seriously. You can expect:

- Acknowledgment within 72 hours
- Regular updates on the investigation
- Credit in the advisory (unless you request anonymity)
- A coordinated disclosure timeline

---

## Trust Model Summary

Pactyra's security model is built on three trust domains:

### Domain 1: Rust Trust Core (`pactyra-core`)

- **Trusts:** Nothing outside itself.
- **Owns:** Signing key (Ed25519), scope decisions, receipt chain, evidence export.
- **Attack surface:** IPC parser, scope grammar parser, hash chain validator.
- **Key security:** Key generated at first run, stored in `~/.pactyra/keys/` with `0600` permissions. Never exposed via IPC.

### Domain 2: Go Orchestrator (`pactyra`)

- **Trusts:** Rust responses, Docker Engine isolation.
- **Owns:** Operator experience (TUI), agent logic, Docker worker lifecycle, tool adapter output parsing.
- **Cannot:** Execute tools without Rust approval, modify receipts, access signing key, suppress receipt generation.

### Domain 3: Docker Workers

- **Trusts:** Nothing. Ephemeral, single-use, fully isolated.
- **Cannot:** Access host filesystem, signing key, Go/Rust processes, or out-of-scope network targets.

---

## Scope Enforcement

Pactyra enforces scope at multiple layers:

1. **Application layer:** Rust Scope Guardian validates every tool invocation before execution. Out-of-scope actions are blocked, not warned.
2. **Network layer:** Docker workers have network egress restricted to scope-authorized targets via Docker network policies.
3. **Evidence layer:** Every action (allowed or blocked) produces an immutable receipt. Tampering is detectable via hash chain verification.

---

## Cryptographic Guarantees

- **Receipt integrity:** SHA-256 hash chain. Tampering with any receipt breaks the chain.
- **Receipt authenticity:** Ed25519 signatures. Every receipt is signed at creation.
- **Evidence export:** Standalone verification with public key. No platform access required.

**What Pactyra does NOT guarantee:**

- That the operator did not manipulate the target environment before capturing evidence.
- That the operator's methodology was correct.
- That the target will respond identically on retest (environments change).
- That the system clock is accurate (timestamps are system-clock-based).

---

## Supported Versions

| Version | Status         |
| ------- | -------------- |
| v0.1    | In development |

Security patches will be backported to the latest stable release.

---

## License

AGPL-3.0 — GNU Affero General Public License v3.0. Strong copyleft that protects the open-source core. See `openspec/changes/SDD-000/proposal.md` §2.3 for the rationale (SDD-000, T-002).

---

## Acknowledgments

We follow responsible disclosure practices. Security researchers who report valid vulnerabilities will be acknowledged in our security advisories (unless they request anonymity).
