// Package worker manages Docker network egress policies for workers.
//
// Each worker gets a dedicated Docker network with egress restricted
// to the authorized target set from the active scope. This is defense-in-depth:
// even if Rust validation has a bug, the network won't pass traffic
// to out-of-scope targets.
//
// See: openspec/changes/SDD-000/design.md §5.3, §6.3
package worker
