# WORKFLOW-SESSION.md — Conduit

**Role:** control — remote agent bridge  
**Version:** v0.1.0  
**Local path:** ~/workspace/projects/engx/services/conduit

---

## Start of session checklist

```bash
cd ~/workspace/projects/engx/services/conduit
git pull
go build ./...
go test ./...
go vet ./...
```

---

## Running locally (dev — identity disabled)

```bash
export CONDUIT_REQUIRE_IDENTITY=false
export CONDUIT_HTTP_ADDR=127.0.0.1:9092
export CONDUIT_STREAM_ADDR=127.0.0.1:9093
go run ./cmd/conduit/
```

---

## go.mod setup (first time)

```bash
cd ~/workspace/projects/engx/services/conduit
go get github.com/Harshmaury/Canon@v1.0.0
go mod tidy
# Then replace config.ServiceTokenHeader etc. with direct Canon imports
```

---

## Never

- Allow anonymous sessions when `CONDUIT_REQUIRE_IDENTITY=true`
- Skip Gate validation — always call POST /gate/validate for remote sessions
- Call write endpoints on any platform service (ADR-020)
- Hardcode `"execute"` scope string — use Canon constant when Canon is imported
