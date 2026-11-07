package main

// Optional inbound API key gate for public deployments (e.g. Railway).
//
// The gateway is designed for trusted networks and by default accepts any
// caller. When the API_KEY environment variable is set, every route except
// /healthz requires that key, accepted as "Authorization: Bearer <key>"
// (OpenAI-style clients) or "x-api-key" (Anthropic-style clients). Leaving
// API_KEY unset keeps the historical open behavior for local use.
//
// A tiny constant-key gate instead of CLIProxyAPI's key config keeps this
// deployment concern out of the embedded config file format.

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	sdkginhandlers "github.com/router-for-me/CLIProxyAPI/v7/sdk/api/handlers"
)

// apiKeyEnvVar holds the required inbound key; empty or unset disables auth.
const apiKeyEnvVar = "API_KEY"

// apiKeyMiddleware rejects requests without the configured key. /healthz stays
// open so platform health probes (Railway, Docker HEALTHCHECK) keep working.
func apiKeyMiddleware() gin.HandlerFunc {
	key := strings.TrimSpace(os.Getenv(apiKeyEnvVar))
	return func(c *gin.Context) {
		if key == "" || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}
		provided := inboundAPIKey(c.Request)
		if provided != "" && subtle.ConstantTimeCompare([]byte(provided), []byte(key)) == 1 {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, sdkginhandlers.ErrorResponse{
			Error: sdkginhandlers.ErrorDetail{Message: "invalid or missing API key", Type: "authentication_error"},
		})
	}
}

// inboundAPIKey extracts the caller key from Authorization: Bearer <key> or
// x-api-key, returning "" when neither carries one.
func inboundAPIKey(req *http.Request) string {
	auth := strings.TrimSpace(req.Header.Get("Authorization"))
	if prefix := "Bearer "; len(auth) > len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	return strings.TrimSpace(req.Header.Get("x-api-key"))
}
