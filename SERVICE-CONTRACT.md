# SERVICE-CONTRACT.md — Conduit

**Role:** control  
**Version:** v0.1.0  
**Ports:** 9092 (HTTP API), 9093 (stream listener)  
**Module:** `github.com/Harshmaury/Conduit`

---

## Stream protocol (port 9093)

Remote engxa sends a JSON handshake:

```json
{"agent_id":"agent-abc","identity_token":"<gate-jwt>"}
```

Conduit responds:

```json
{"ok":true,"session_id":"csn_1a2b3c","owner":"harsh@github"}
```

On error:

```json
{"ok":false,"error":"identity token required"}
```

Session stays open with a 60-second heartbeat read timeout.

## HTTP API (port 9092)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | None | Liveness |
| GET | `/conduit/sessions` | X-Service-Token | List active sessions |

## Invariants

- `CONDUIT_REQUIRE_IDENTITY=true` (default) — anonymous sessions always denied
- Identity token validated against Gate `POST /gate/validate` (signature + expiry + revocation)
- Token must carry `execute` scope — missing scope = connection refused
- Reconnecting agent closes the previous session
- Conduit never calls write endpoints on Nexus, Forge, Atlas, or any observer (ADR-020)
