// Package agent defines the Agent interface and provides agent implementations.
//
// Agent contract (architectural invariant):
// - An agent receives a task (structured from operator intent or prior agent output).
// - An agent produces proposals (tool invocations, notes, expansion requests).
// - Every proposal must pass through the Scope Guardian Agent before execution.
// - Agents do not execute tools directly — they ask the Worker Manager through the Scope Guardian.
//
// This contract guarantees that no agent — AI or otherwise — can bypass scope.
//
// See: openspec/changes/SDD-000/design.md §3.2, §6.2

package agent
