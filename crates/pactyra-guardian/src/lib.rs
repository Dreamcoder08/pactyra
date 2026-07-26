// pactyra-guardian: Runtime scope enforcement.
//
// This crate is the architectural gatekeeper. Given a proposed tool invocation
// (target IP/domain, ports, tool, parameters), it evaluates against the active
// compiled scope and returns allow/deny. Every decision is logged.
//
// The Scope Guardian is the runtime enforcement of Pactyra's core invariant:
// "Every action must be bound to an active scope."
//
// See: openspec/changes/SDD-000/design.md §3.1, §5.1
// Owned by: SDD-002 (scope enforcement logic), SDD-003 (permission lattice O4–O6)

pub mod guardian;
