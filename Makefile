.PHONY: build test lint run golden-path clean

# Build both binaries: pactyra (Go) and pactyra-core (Rust)
build:
	cargo build --release
	go build -o bin/pactyra ./cmd/pactyra

# Run all tests: Rust + Go + integration
test:
	cargo test
	go test ./...

# Lint: clippy for Rust, golangci-lint for Go
lint:
	cargo clippy -- -D warnings
	golangci-lint run ./...

# Build and launch the TUI with a test scope fixture
run: build
	./bin/pactyra run --scope "tests/fixtures/scope_valid.yaml"

# Full golden path integration test (requires Docker + cyber range)
golden-path:
	@echo "Starting cyber range..."
	docker compose -f cyber-range/docker-compose.yml up -d
	@echo "Waiting for targets to be ready..."
	sleep 10
	@echo "Running golden path integration test..."
	go test -v -run TestGoldenPath ./tests/
	@echo "Stopping cyber range..."
	docker compose -f cyber-range/docker-compose.yml down

# Remove build artifacts
clean:
	cargo clean
	rm -rf bin/
	rm -f go.sum
