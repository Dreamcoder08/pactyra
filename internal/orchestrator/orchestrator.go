// Package orchestrator manages session state, engagement lifecycle, and agent coordination.
//
// The Orchestrator owns:
// - Active scope reference and current scope_hash
// - Engagement metadata (name, operator, dates)
// - Agent registry and lifecycle
// - IPC client connection to pactyra-core
//
// See: openspec/changes/SDD-000/design.md §3.2

package orchestrator
