// Package guardian provides the IPC client library for communicating with pactyra-core.
//
// This is the ONLY Go package that talks to the Rust binary. All other Go
// packages that need trust decisions (scope validation, receipt recording,
// chain verification) route through this package.
//
// See: openspec/changes/SDD-000/design.md §3.2, §6.1

package guardian
