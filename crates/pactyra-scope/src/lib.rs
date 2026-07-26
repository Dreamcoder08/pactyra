// pactyra-scope: Scope parsing, compilation, and deterministic hashing.
//
// This crate is the foundation of Pactyra's trust model. It parses scope
// documents (CIDR notation, domain lists, excluded hosts, port ranges),
// compiles them into a canonical machine-enforceable representation, and
// produces a deterministic scope_hash that becomes the root of the evidence chain.
//
// See: openspec/changes/SDD-000/design.md §3.1
// Owned by: SDD-002

pub mod compiler;
pub mod hash;
pub mod parser;
