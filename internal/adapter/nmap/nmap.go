// Package nmap implements the Nmap tool adapter.
//
// The Nmap adapter:
// - Validates Nmap-specific parameters (port ranges, scan types, timing templates)
// - Translates Recon Planner intent into Nmap CLI flags
// - Parses Nmap XML output into structured results (hosts, ports, services, OS)
//
// This is the ONLY tool adapter in v0.1.
//
// See: openspec/changes/SDD-000/design.md §3.2
package nmap
