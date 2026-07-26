# Pactyra Cyber Range

Local test infrastructure for v0.1 golden path validation. Provides vulnerable-by-design target services for controlled testing.

---

## Quickstart

```bash
# Start all targets
docker compose -f cyber-range/docker-compose.yml up -d

# Wait for services to be ready
sleep 10

# Verify targets are reachable
ping -c 2 10.0.1.10
curl -s -o /dev/null -w "%{http_code}" http://10.0.1.10

# Run the golden path test
make golden-path

# Stop all targets
docker compose -f cyber-range/docker-compose.yml down
```

---

## Target Inventory

| Container              | Hostname               | IP          | Ports Open     | Service            |
| ---------------------- | ---------------------- | ----------- | -------------- | ------------------ |
| `pactyra-target-web`   | `test-target.local`    | `10.0.1.10` | 80, 8080, 8443 | DVWA (web app)     |
| `pactyra-target-nginx` | `vulnerable-app.local` | `10.0.1.20` | 80, 443, 8080  | Nginx              |
| `pactyra-target-httpd` | `extra-service.local`  | `10.0.1.30` | 80, 443, 8443  | Apache httpd       |
| `pactyra-gateway`      | `gateway.local`        | `10.0.1.1`  | (none)         | Gateway (excluded) |

---

## Scope Coverage

The matching scope fixture is at `tests/fixtures/scope_valid.yaml`:

```yaml
targets:
  cidr:
    - '10.0.1.0/24'
  domains:
    - 'test-target.local'
    - 'vulnerable-app.local'

excluded:
  hosts:
    - '10.0.1.1' # Gateway excluded from scanning

ports:
  include:
    - '1-1000'
    - '8080'
    - '8443'
```

This scope:

- Authorizes the full `10.0.1.0/24` subnet (254 usable hosts)
- Explicitly excludes the gateway (`10.0.1.1`)
- Authorizes ports 1–1000 plus 8080 and 8443
- Authorizes two domains for DNS-based targeting

**Out-of-scope targets** (for violation testing):

- `10.0.0.0/8` — broader subnet (exceeds authorized CIDR)
- `10.0.1.1` — gateway (explicitly excluded)
- Any port outside 1–1000, 8080, 8443

---

## Network Topology

```
10.0.1.0/24 (pactyra-cyber-range bridge)
├── 10.0.1.1   Gateway (excluded from scope)
├── 10.0.1.10  test-target.local    (DVWA, ports 80/8080/8443)
├── 10.0.1.20  vulnerable-app.local (Nginx, ports 80/443/8080)
└── 10.0.1.30  extra-service.local  (Apache, ports 80/443/8443)
```

All targets are reachable from the host. The Docker network is isolated — targets cannot reach the internet or the host's other networks.

---

## Golden Path Test Scenario

1. **Load scope:** `pactyra scope compile tests/fixtures/scope_valid.yaml`
2. **Start engagement:** `pactyra run --scope <scope_hash>`
3. **Discover:** "scan web ports on scope targets" → should scan `10.0.1.10,.20,.30` ports 80,443,8080,8443
4. **Violate:** "scan 10.0.0.0/8" → should be BLOCKED (exceeds authorized CIDR)
5. **Violate:** "scan 10.0.1.1" → should be BLOCKED (excluded host)
6. **Record:** All actions (allowed and blocked) appear as receipts
7. **Export:** `pactyra export --output ./evidence/`
8. **Verify:** `pactyra verify ./evidence/` → chain valid, all receipts signed

---

## Cleanup

```bash
# Stop and remove all containers and networks
docker compose -f cyber-range/docker-compose.yml down -v

# Remove images if needed
docker rmi vulnerables/web-dvwa nginx:alpine httpd:alpine alpine:latest
```

---

## Extending the Range

To add new targets for testing additional tools (Nuclei, custom scripts in v0.5):

1. Add a new service to `docker-compose.yml`
2. Assign an IP in the `10.0.1.0/24` range
3. Update `tests/fixtures/scope_valid.yaml` if the target should be in scope
4. Document the new target in this README
