// Package export handles evidence pack creation and standalone verification.
//
// The Evidence Exporter (Go side):
// - Requests full chain data from Rust via ExportChain IPC
// - Packages with metadata (engagement name, operator, dates)
// - Produces standalone export directory with receipts, scope, public key, VERIFY.md
//
// The standalone verifier (pactyra verify command):
// - Reads an evidence pack directory
// - Checks hash chain integrity
// - Verifies every signature against the public key
// - Reports results without requiring pactyra-core to be running
//
// See: openspec/changes/SDD-000/design.md §3.2, §4.4, §9.4
package export
