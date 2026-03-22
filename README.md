# Conduit

**Remote Agent Bridge**  
`role: control` | `version: v0.1.0` | ADR-042

---

## What it does

Conduit allows remote `engxa` agents to connect to the platform from outside the local machine. Every session requires a Gate-issued identity token — Conduit is a hard identity boundary.

---

## Ports

| Port | Purpose |
|------|---------|
| `9092` | HTTP API — session list, health |
| `9093` | Stream listener — remote engxa connections |

---

## Environment variables

| Variable | Default | Required |
|----------|---------|----------|
| `CONDUIT_HTTP_ADDR` | `0.0.0.0:9092` | No |
| `CONDUIT_STREAM_ADDR` | `0.0.0.0:9093` | No |
| `CONDUIT_SERVICE_TOKEN` | — | Production |
| `GATE_ADDR` | `http://127.0.0.1:8088` | No |
| `NEXUS_ADDR` | `http://127.0.0.1:8080` | No |
| `CONDUIT_REQUIRE_IDENTITY` | `true` | Set `false` for local dev only |

---

## Build

```bash
go build -o conduit ./cmd/conduit/
CONDUIT_REQUIRE_IDENTITY=false ./conduit
```

---

## ADR

ADR-042 — Gate: Platform Identity Authority
