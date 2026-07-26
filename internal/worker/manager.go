// Package worker manages Docker worker lifecycle for isolated tool execution.
//
// The Worker Manager:
// - Creates ephemeral Docker containers per tool invocation
// - Injects validated tool parameters (pre-approved by Scope Guardian)
// - Enforces network egress via Docker network policies (scope-only targets)
// - Captures stdout/stderr streams
// - Destroys containers after completion or timeout
//
// Workers are Trust Domain 3 — fully untrusted, no host access, no key access.
//
// See: openspec/changes/SDD-000/design.md §3.2, §5.3, §6.3

package worker
