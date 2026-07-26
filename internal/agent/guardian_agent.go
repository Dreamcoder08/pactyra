// Package agent provides the Scope Guardian agent (Go-side proxy).
//
// The Scope Guardian Agent is the Go-side proxy for the Rust Scope Guardian.
// It receives validation requests from the Recon Planner or operator,
// forwards to Rust via IPC, and returns allow/deny.
//
// See: openspec/changes/SDD-000/design.md §3.2, §6.2
package agent
