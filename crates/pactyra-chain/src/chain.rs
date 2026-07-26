// Append-only hash chain. Each receipt links to its parent via SHA-256 hash.
// Supports: append, verify (full chain walk), get_by_id, iterate.
// See: openspec/changes/SDD-000/design.md §7.6 (Hash Chain rationale)
