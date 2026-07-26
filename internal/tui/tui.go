// Package tui provides the terminal user interface for Pactyra.
//
// The TUI is the primary operator interface in v0.1. It provides four views:
// - Scope View: active scope summary (IP ranges, domains, ports)
// - Command Input: operator types intents, not raw commands
// - Receipt Browser: scrollable list of evidence receipts
// - Status Line: scope_hash, engagement name, agent status, chain height
//
// Built with Bubble Tea (charmbracelet/bubbletea).
// See: openspec/changes/SDD-000/design.md §6.4, §9.1

package tui
