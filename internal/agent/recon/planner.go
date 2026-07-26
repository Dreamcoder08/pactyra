// Package recon implements the Recon Planner agent.
//
// The Recon Planner translates operator intent ("scan web ports on scope targets")
// into concrete Nmap invocations. It proposes invocations to the Scope Guardian
// before execution.
//
// See: openspec/changes/SDD-000/design.md §3.2, §6.2
package recon
