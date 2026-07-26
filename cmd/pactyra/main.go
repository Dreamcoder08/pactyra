// Pactyra — Open Offensive Operations Kernel
//
// pactyra is the Go orchestrator binary. It provides the operator's terminal
// interface (TUI), manages agent lifecycle, coordinates Docker workers, and
// communicates with the Rust trust core (pactyra-core) via Unix domain socket.
//
// The operator runs this binary. It spawns pactyra-core as a child process.
//
// See: openspec/changes/SDD-000/design.md §2, §3.2, §9.2

package main

import "fmt"

func main() {
	// TODO(SDD-002+): Parse subcommands (scope compile, run, export, verify).
	//                 Launch TUI or headless mode.
	//                 Spawn pactyra-core child process.
	fmt.Println("pactyra: not yet implemented (SDD-002+)")
}
