# Design Partner Recruitment Strategy

**Date:** 2026-07-26
**Source:** SDD-000 T-019
**Target:** 3+ design partners committed before SDD-002 implementation

---

## 1. Target Profile

### Ideal Design Partner

A **small specialized offensive security firm** with the following characteristics:

| Attribute          | Ideal                                         | Acceptable                           |
| ------------------ | --------------------------------------------- | ------------------------------------ |
| Team size          | 3–10 operators                                | Solo operator with engagement volume |
| Engagement volume  | 3–4/month                                     | 1–2/month                            |
| Client profile     | Mid-market SaaS/fintech                       | Small business, compliance-focused   |
| Current workflow   | Kali VM + tmux + Notion                       | Any fragmented toolchain             |
| Pain level         | High — frustrated with evidence/documentation | Moderate — sees value but not urgent |
| Technical maturity | Comfortable building from source              | Willing to try with support          |
| Geography          | Any (remote-first)                            | Any                                  |

### Archetype: "Alex Chen's Firm"

Based on the beachhead user profile from the proposal:

- 7-person boutique firm (Aegis Security Partners)
- Senior operators lead, juniors learn on the job
- Clients need SOC 2 / ISO 27001 evidence
- Current pain: 30–40% integration labor, scope violation near-misses, evidence challenges

---

## 2. Value Proposition for Partners

### What Partners Get

| Benefit                 | Details                                                                 |
| ----------------------- | ----------------------------------------------------------------------- |
| **Early access**        | Builds before public release. Direct line to development team.          |
| **Influence**           | Shape the roadmap. Your pain points become our priorities.              |
| **Risk reduction**      | Scope enforcement prevents the near-misses your team experiences today. |
| **Credibility defense** | Evidence receipts give your clients cryptographic proof of methodology. |
| **Recognition**         | Name/company in project acknowledgments, release notes, and website.    |
| **No cost**             | Pactyra is open-source. Partners use it on real engagements for free.   |

### What Partners Commit To

| Commitment                  | Details                                                                |
| --------------------------- | ---------------------------------------------------------------------- |
| **Use on real engagements** | Run Pactyra on at least 2 authorized engagements during v0.1–v0.5.     |
| **Provide feedback**        | Monthly 30-minute call. Share what worked, what broke, what's missing. |
| **Report bugs**             | Use GitHub issues for bugs and feature requests.                       |
| **Permission to reference** | We can say "used by real offensive security teams" (no client names).  |

### What Partners Do NOT Commit To

- No exclusivity. Use whatever tools you want alongside Pactyra.
- No public endorsement required. Your feedback can be private.
- No client data sharing. We never see your clients' names or environments.
- No financial commitment. Pactyra is free and open-source.

---

## 3. Candidate Identification

### Recruitment Channels

| Channel                        | Approach                                                              | Expected Yield               |
| ------------------------------ | --------------------------------------------------------------------- | ---------------------------- |
| **Professional network**       | Direct outreach to known operators and firms                          | 3–5 candidates               |
| **LinkedIn**                   | Search for boutique offensive security firms (<20 employees)          | 5–10 candidates              |
| **Security conferences**       | BSides, DEF CON, OffensiveCon — talk to operators at villages         | 3–5 candidates               |
| **CTF communities**            | Hack The Box, TryHackMe forums — look for operators, not just players | 2–3 candidates               |
| **Reddit**                     | r/netsec, r/redteamsec — post about design partner program            | 5–10 expressions of interest |
| **Offensive Security Discord** | Engage in relevant servers, share the thesis                          | 3–5 candidates               |

### Candidate List (Template)

```
1. [Firm Name] — [Location] — [Team Size] — [Contact Status]
2. ...
```

_Fill 5–10 entries through outreach. Target: 3+ committed._

### Screening Questions

1. "Tell me about your current engagement workflow — what tools, how do you document?"
2. "Have you ever had a scope violation near-miss? How did you catch it?"
3. "When a client challenges a finding, what evidence do you provide?"
4. "How do you handle retesting six months after an engagement?"
5. "What's your biggest frustration with your current toolchain?"

If they answer all five with clear pain points matching Pactyra's thesis → strong candidate.

---

## 4. Outreach Materials

### One-Page Overview

**Title:** Pactyra: Open Offensive Operations Kernel
**Subtitle:** Bound by scope. Proven by evidence.

**The problem:** Offensive results are disputable because evidence is fragmented, incomplete, or does not verifiably demonstrate what occurred.

**The solution:** Pactyra binds every action and piece of evidence to a verifiable scope, identity, tool, version, authorization, and timeline.

**What it does:**

- Enforces scope at runtime — blocks out-of-scope actions before network traffic is sent
- Records every action in an immutable, signed evidence chain
- Produces exportable, independently-verifiable evidence packs
- One-command retesting from evidence replay (v0.5)

**What we need:** 3+ offensive security teams to use Pactyra on real engagements and provide feedback.

**What you get:** Early access, roadmap influence, risk reduction, and recognition.

**Contact:** [TBD]

### FAQ

**Q: Is this a commercial product?**
A: No. Pactyra is open-source. It will always be self-hostable and free. Commercial services may exist in the future but will not weaken the open-source core.

**Q: Do I need AI to use it?**
A: No. AI is an optional accelerator, never a requirement. Pactyra works fully without any AI integration.

**Q: What tools does it support?**
A: v0.1 supports Nmap. v0.5 adds a Tool Adapter SDK for community contributions (Nuclei, Burp, custom scripts).

**Q: Can I use it alongside my current tools?**
A: Yes. Pactyra orchestrates tools; it doesn't replace them. You can run Pactyra for scoped operations and your existing toolchain for everything else.

**Q: What's the time commitment?**
A: Use it on 2+ engagements. Monthly 30-minute feedback call. Report bugs as you find them.

**Q: Is my client data safe?**
A: Yes. Pactyra runs entirely on your machine. No data leaves your environment. We never see your clients, targets, or findings.

---

## 5. Timeline

| Milestone                | Target Date                   | Dependency                    |
| ------------------------ | ----------------------------- | ----------------------------- |
| Candidate list (5–10)    | SDD-000 completion            | Outreach begins               |
| First contact round      | Week 1 post-SDD-000           | Candidate list ready          |
| Follow-up / demo concept | Week 2–3                      | One-page overview + FAQ ready |
| Commitments secured (3+) | Before SDD-002 implementation | 3+ signed on                  |
| First partner engagement | During SDD-002–SDD-0XX        | v0.1 golden path working      |

---

## 6. Demo Concept (The Golden Path)

For partner outreach, demonstrate the 8-step golden path:

1. **Load scope:** Show a real-looking scope YAML. "This is your engagement scope."
2. **Compile:** `pactyra scope compile` → `scope_hash: a1b2c3...` "Every action will reference this."
3. **Discover:** "scan web ports on scope targets" → Nmap executes in isolated Docker worker.
4. **Violate:** "scan 10.0.0.0/8" → **BLOCKED.** "No traffic sent. Receipt recorded."
5. **Review:** Scroll through receipts in TUI. "Every action has a signed receipt."
6. **Export:** `pactyra export` → evidence pack with receipts, scope, public key, VERIFY.md.
7. **Verify:** `pactyra verify ./evidence/` → "Chain valid. No tampering."
8. **Narrative:** "A client challenges a finding. You send them the evidence pack. They verify it independently. No platform access needed."

The demo takes <5 minutes and demonstrates the complete value chain.

---

## 7. Success Criteria

- ✅ 3+ design partners committed
- ✅ At least 1 partner runs a real engagement on Pactyra v0.1
- ✅ Feedback loop established (monthly calls, GitHub issues)
- ✅ At least 1 partner testimonial for v0.5 public launch
