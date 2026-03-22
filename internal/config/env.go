// @conduit-project: conduit
// @conduit-path: internal/config/env.go
// All Conduit configuration from environment variables. No config files.
package config

import (
	"os"

	canonid "github.com/Harshmaury/Canon/identity"
)

// Config holds all Conduit runtime configuration.
type Config struct {
	// HTTPAddr is the address for the Conduit HTTP API (agent registration, sessions).
	// Environment: CONDUIT_HTTP_ADDR. Default: 0.0.0.0:9092
	HTTPAddr string

	// StreamAddr is the address for the remote engxa stream listener.
	// Environment: CONDUIT_STREAM_ADDR. Default: 0.0.0.0:9093
	StreamAddr string

	// NexusAddr is the Nexus HTTP API address.
	// Environment: NEXUS_ADDR. Default: http://127.0.0.1:8080
	NexusAddr string

	// GateAddr is the Gate HTTP API address for identity validation (ADR-042).
	// Environment: GATE_ADDR. Default: http://127.0.0.1:8088
	GateAddr string

	// ServiceToken is the platform service mesh token (ADR-008).
	// Environment: CONDUIT_SERVICE_TOKEN. Required in production.
	ServiceToken string

	// RequireIdentity enforces Gate identity on all agent sessions.
	// Environment: CONDUIT_REQUIRE_IDENTITY. Default: true (Conduit is a remote boundary).
	RequireIdentity bool
}

// Load reads all configuration from environment with defaults applied.
func Load() *Config {
	return &Config{
		HTTPAddr:        envOr("CONDUIT_HTTP_ADDR", DefaultHTTPAddr),
		StreamAddr:      envOr("CONDUIT_STREAM_ADDR", DefaultStreamAddr),
		NexusAddr:       envOr("NEXUS_ADDR", canonid.DefaultNexusAddr),
		GateAddr:        envOr("GATE_ADDR", canonid.DefaultGateAddr),
		ServiceToken:    os.Getenv("CONDUIT_SERVICE_TOKEN"),
		RequireIdentity: os.Getenv("CONDUIT_REQUIRE_IDENTITY") != "false", // default true
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
