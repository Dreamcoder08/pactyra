// pactyra-chain: Append-only evidence receipt chain.
//
// This crate implements the evidence model's integrity and provenance layers:
// - Receipt creation with Ed25519 signing
// - Append-only hash chain (each receipt hashes its parent)
// - Chain integrity verification (detects tampering)
// - Chain traversal and export
//
// The signing key never leaves this crate. It is the most security-critical
// code in the entire Pactyra system.
//
// See: openspec/changes/SDD-000/design.md §3.1, §5.1
// Owned by: SDD-007

pub mod chain;
pub mod export;
pub mod receipt;
pub mod signer;
