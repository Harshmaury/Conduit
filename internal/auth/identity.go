// @conduit-project: conduit
// @conduit-path: internal/auth/identity.go
// Package auth validates Gate identity tokens for remote agent sessions.
// ADR-042: Conduit is a remote boundary — identity validation is mandatory.
// All remote sessions must present a valid Gate token with scope "execute".
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Claim is the validated Gate identity attached to every Conduit session.
type Claim struct {
	Subject   string   `json:"sub"`
	Scopes    []string `json:"scp"`
	ExpiresAt int64    `json:"exp"`
	TokenID   string   `json:"jti"`
}

// HasScope returns true if the claim contains the given scope.
func (c *Claim) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// Validator calls Gate POST /gate/validate for full signature+expiry+revocation check.
type Validator struct {
	gateAddr     string
	serviceToken string
	client       *http.Client
}

// NewValidator creates a Validator.
func NewValidator(gateAddr, serviceToken string) *Validator {
	return &Validator{
		gateAddr:     gateAddr,
		serviceToken: serviceToken,
		client:       &http.Client{Timeout: 3 * time.Second},
	}
}

// Validate calls Gate and returns the validated Claim.
// Returns error if the token is invalid, expired, or revoked.
// Returns (nil, nil) if token is empty — callers decide if anonymous is allowed.
func (v *Validator) Validate(identityToken string) (*Claim, error) {
	if identityToken == "" {
		return nil, nil
	}
	body := fmt.Sprintf(`{"token":%q}`, identityToken)
	req, err := http.NewRequest(http.MethodPost, v.gateAddr+"/gate/validate",
		strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gate validate: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if v.serviceToken != "" {
		req.Header.Set("X-Service-Token", v.serviceToken)
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gate validate: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Valid  bool   `json:"valid"`
		Claim  *Claim `json:"claim,omitempty"`
		Reason string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gate validate decode: %w", err)
	}
	if !result.Valid {
		return nil, fmt.Errorf("identity rejected: %s", result.Reason)
	}
	return result.Claim, nil
}
