// @conduit-project: conduit
// @conduit-path: internal/config/canon.go
// Canon constants mirrored for Conduit until go get Canon@v1.0.0 runs.
// Replace with direct Canon imports after: go get github.com/Harshmaury/Canon@v1.0.0
package config

// Header constants — ADR-016, ADR-042.
const (
	ServiceTokenHeader  = "X-Service-Token"
	TraceIDHeader       = "X-Trace-ID"
	IdentityTokenHeader = "X-Identity-Token"
)

// Service constants.
const (
	ServiceName     = "conduit"
	DefaultNexusAddr = "http://127.0.0.1:8080"
	DefaultGateAddr  = "http://127.0.0.1:8088"
)

// Default ports — Conduit is the remote agent bridge.
const (
	DefaultHTTPAddr   = "0.0.0.0:9092" // Conduit HTTP API (agent registration, session management)
	DefaultStreamAddr = "0.0.0.0:9093" // Conduit stream listener (remote engxa connections)
)
