// @conduit-project: conduit
// @conduit-path: internal/config/canon.go
// Package config — local constants for Conduit service identity and default ports.
//
// Header constants (ServiceTokenHeader, TraceIDHeader, IdentityTokenHeader) have been
// removed. Import canonid "github.com/Harshmaury/Canon/identity" directly wherever
// those constants are needed. See ADR-016, Rule 2.
//
// DefaultNexusAddr and DefaultGateAddr have been removed. Use canonid.DefaultNexusAddr
// and canonid.DefaultGateAddr from Canon directly (already done in env.go).
package config

// ServiceName is the canonical identifier for this service in logs and events.
const ServiceName = "conduit"

// Default listen addresses for Conduit's two ports.
// Override with CONDUIT_HTTP_ADDR and CONDUIT_STREAM_ADDR environment variables.
const (
	DefaultHTTPAddr   = "0.0.0.0:9092" // HTTP API: session list, health
	DefaultStreamAddr = "0.0.0.0:9093" // Stream listener: remote engxa connections
)
