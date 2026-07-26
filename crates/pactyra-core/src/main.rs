// pactyra-core: Rust trust core binary.
//
// This binary is spawned by the Go `pactyra` binary as a child process.
// It listens on a local Unix domain socket and routes requests from Go
// to the appropriate internal component (Scope Compiler, Scope Guardian,
// Receipt Chain).
//
// The operator never runs this binary directly. It is an implementation
// detail of the two-binary architecture.
//
// See: openspec/changes/SDD-000/design.md §2, §3.1, §5.1

mod handlers;
mod server;

fn main() {
    // TODO(SDD-002): Parse CLI args for socket path, initialize components,
    //                start IPC server loop.
    eprintln!("pactyra-core: not yet implemented (SDD-002+)");
}
