# Conference and Community Strategy

**Date:** 2026-07-26
**Source:** SDD-000 T-022
**Purpose:** Conference target list, CFP deadlines, talk abstract, and online community engagement plan.

---

## 1. Conference Targets

### Tier 1: Must-Present (High Alignment)

| Conference                 | Location                | Typical Dates   | CFP Deadline       | Audience                                    | Alignment |
| -------------------------- | ----------------------- | --------------- | ------------------ | ------------------------------------------- | --------- |
| **BSides** (various)       | Global (local chapters) | Year-round      | Varies by chapter  | 200–800 practitioners                       | ★★★★★     |
| **DEF CON**                | Las Vegas, USA          | August          | ~May               | 30,000+ hackers, researchers, practitioners | ★★★★☆     |
| **Black Hat USA**          | Las Vegas, USA          | August          | ~April (Briefings) | 20,000+ security professionals              | ★★★★☆     |
| **OffensiveCon**           | Berlin, Germany         | May             | ~January           | 500+ offensive security specialists         | ★★★★★     |
| **Hack in the Box (HITB)** | Amsterdam / Singapore   | April / October | ~December / June   | 2,000+ deep technical                       | ★★★★☆     |

### Tier 2: Should-Consider

| Conference                 | Location                | Typical Dates  | Notes                                    |
| -------------------------- | ----------------------- | -------------- | ---------------------------------------- |
| **RSA Conference**         | San Francisco, USA      | February/March | Large but more vendor/enterprise focused |
| **B-Sides London**         | London, UK              | October        | Strong European practitioner community   |
| **SecTor**                 | Toronto, Canada         | October        | Canadian security community              |
| **Wild West Hackin' Fest** | Deadwood, USA / Virtual | October        | Practitioner-focused, approachable       |
| **ShmooCon**               | Washington DC, USA      | January        | Community-driven, good for thesis talks  |

### Tier 3: Future Opportunities (v1.0+)

| Conference           | Reason for Later                                            |
| -------------------- | ----------------------------------------------------------- |
| **USENIX Security**  | Academic — needs published benchmarks and formal evaluation |
| **ACM CCS**          | Academic — needs rigorous evaluation                        |
| **FIRST Conference** | Incident response focus — v0.5+ when evidence model matures |

---

## 2. CFP Timeline (Next 12 Months)

| Deadline           | Conference          | Action                                      |
| ------------------ | ------------------- | ------------------------------------------- |
| **~January 2027**  | OffensiveCon 2027   | Prepare talk abstract, submit               |
| **~April 2027**    | Black Hat USA 2027  | Prepare Briefings proposal                  |
| **~May 2027**      | DEF CON 35          | Prepare CFP submission                      |
| **~December 2026** | HITB Amsterdam 2027 | Monitor for CFP                             |
| **Rolling**        | BSides chapters     | Identify 3–5 local chapters, submit to each |

---

## 3. Talk Abstract (Reusable)

### Title

**Authorized Offensive Automation with Verifiable Evidence — An Open-Source Reference Implementation**

### Format

- 45-minute talk (adaptable to 25-minute or workshop format)
- Target audience: offensive security practitioners, red team operators, security researchers
- Technical level: Intermediate

### Abstract

> Offensive security assessments are essential infrastructure for the digital economy, yet the operational tooling that powers them is fragmented across tools, scripts, terminals, and notes that share no verifiable model of authorization, execution, and evidence. When a client disputes a finding, the operator's defense is screenshots and memory. When a mistyped IP expands a scan beyond authorized scope, the consequences can be legal as well as reputational.
>
> This talk presents Pactyra, an open-source offensive operations kernel that binds every action to a cryptographically-verified scope, records every decision in an append-only evidence chain, and produces independently-verifiable evidence packs. We'll demonstrate the two-binary architecture (Rust trust core + Go orchestrator), walk through the 8-step golden path from scope loading to evidence export, and discuss why a reference implementation for safe, auditable offensive automation matters now — especially as AI agents enter the offensive security toolchain.
>
> Attendees will leave with: an understanding of the scope enforcement and evidence integrity model, a working demo they can reproduce against a local cyber range, and the architectural patterns for building trust boundaries into offensive tooling.

### Key Takeaways

1. Scope enforcement must be blocking, not advisory. Warning-based enforcement fails in practice.
2. Evidence is a byproduct of operations, not a separate documentation chore.
3. A two-binary model (Rust for trust decisions, Go for orchestration) creates a clean, auditable trust boundary.
4. The integration moat — scope → execution → evidence → export — is the differentiator, not any single feature.

### Demo Elements

- Live golden path against a cyber range: load scope, scan, violate, block, export, verify
- Evidence pack verification with zero platform access
- Architecture walkthrough with trust domain diagram

---

## 4. Online Community Engagement

### Initial Presence (Pre-v0.1)

| Platform                               | Strategy                                                                                            | Timeline                      |
| -------------------------------------- | --------------------------------------------------------------------------------------------------- | ----------------------------- |
| **r/netsec** (Reddit)                  | Post the public-interest thesis. Ask for feedback on the architecture.                              | After SDD-000 completion      |
| **r/redteamsec** (Reddit)              | Share the golden path concept. Ask: "How do you handle scope enforcement today?"                    | After v0.1 demo working       |
| **Hack The Box forums**                | Share cyber range setup. "Here's a local range for testing scope enforcement."                      | After cyber range is runnable |
| **Offensive Security Discord servers** | Join relevant servers. Don't pitch — participate. Answer questions about scope/evidence challenges. | Ongoing                       |

### Growth Phase (v0.1–v0.5)

| Platform               | Strategy                                                                                             |
| ---------------------- | ---------------------------------------------------------------------------------------------------- |
| **GitHub Discussions** | Enable on the repo. Use for RFCs, feature requests, community support.                               |
| **Twitter/X**          | Share development updates, architectural decisions, golden path demos.                               |
| **YouTube**            | Record golden path demo. "Pactyra v0.1 in 5 minutes."                                                |
| **Blog / Dev.to**      | Technical deep-dives: "Why Rust for the trust boundary," "Hash chains vs. blockchains for evidence." |
| **Podcasts**           | Pitch to security podcasts (Risky Business, Darknet Diaries, Security Weekly).                       |

### Mature Community (v0.5+)

| Platform                        | Strategy                                                                                  |
| ------------------------------- | ----------------------------------------------------------------------------------------- |
| **Conference workshops**        | "Build Your Own Tool Adapter" — hands-on SDK workshop.                                    |
| **CTF integration**             | Design a Pactyra-themed CTF challenge: "Find the scope violation in this evidence chain." |
| **Contributor mentorship**      | Pair experienced contributors with newcomers for first PRs.                               |
| **Design partner testimonials** | Publish case studies from design partners (anonymized).                                   |

---

## 5. Community Values

All community engagement should reflect Pactyra's core identity:

1. **Practitioner-first.** The community is for operators, by operators. We build what operators need.
2. **Open by default.** Discussions, RFCs, and decisions happen in public. No backchannel governance.
3. **AI-optional.** AI tools are welcome as accelerators, never as gatekeepers. Human operators are the decision-makers.
4. **Professional discourse.** Passionate but respectful. We debate architecture, not people.
5. **Evidence-driven.** Claims about security, performance, or methodology should be backed by reproducible evidence.

---

## 6. Content Calendar (First 6 Months)

| Month   | Content                                                                    | Channel                     |
| ------- | -------------------------------------------------------------------------- | --------------------------- |
| Month 1 | "Introducing Pactyra: Open Offensive Operations Kernel" — thesis post      | r/netsec, blog              |
| Month 2 | "The Two-Binary Architecture: Why Rust Owns Trust Decisions"               | Blog, r/redteamsec          |
| Month 3 | Golden path demo video (v0.1 alpha)                                        | YouTube, Twitter/X          |
| Month 4 | "How Pactyra's Scope Enforcement Works" — technical walkthrough            | Blog, Hack The Box forums   |
| Month 5 | "Verifying Evidence Packs Without the Platform" — standalone verifier demo | Blog, GitHub Discussions    |
| Month 6 | Design partner case study (anonymized)                                     | Blog, conference submission |

---

## 7. Success Metrics (v1.0)

| Metric                         | Target                              |
| ------------------------------ | ----------------------------------- |
| GitHub stars                   | 500+                                |
| Active contributors/month      | 15+                                 |
| Conference talks delivered     | 3+ (BSides, DEF CON, OffensiveCon)  |
| Design partner teams           | 5+ on real engagements              |
| Public benchmark reproductions | 3+ independent parties              |
| Online community members       | 200+ across Discord, GitHub, Reddit |
