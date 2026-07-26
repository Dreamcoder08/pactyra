// Package adapter defines the Adapter interface for tool integration.
//
// The Adapter interface is the contract for all tool adapters:
// - Validate tool-specific parameters before execution
// - Translate structured intent into tool CLI flags
// - Parse tool output (XML, JSON, text) into structured results
//
// v0.1 implements exactly one adapter: Nmap.
// v0.5 introduces the Tool Adapter SDK for community adapters.
//
// See: openspec/changes/SDD-000/design.md §3.2, §13 (Future-Proofing)
package adapter
