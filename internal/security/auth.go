package security

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const ClaimsContextKey contextKey = "orbit_claims"

type AuthMiddleware struct {
	tokenMgr *TokenManager
}

func NewAuthMiddleware(tokenMgr *TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenMgr: tokenMgr}
}

// Wrap guards an http.Handler requiring a valid Bearer token or X-Orbit-Key.
func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow unauthenticated health check, public static UI assets, telemetry, and S3
		p := r.URL.Path
		if p == "/healthz" || p == "/" || p == "/index.html" || p == "/icon.svg" ||
			p == "/manifest.webmanifest" || p == "/sw.js" ||
			strings.HasPrefix(p, "/assets/") || strings.HasPrefix(p, "/s3/") ||
			strings.HasPrefix(p, "/api/v1/telemetry/") {
			next.ServeHTTP(w, r)
			return
		}

		token := extractToken(r)
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing authorization token"})
			return
		}

		claims, err := a.tokenMgr.VerifyToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			return parts[1]
		}
	}
	if xKey := r.Header.Get("X-Orbit-Key"); xKey != "" {
		return xKey
	}
	if queryToken := r.URL.Query().Get("token"); queryToken != "" {
		return queryToken
	}
	return ""
}

func GetClaims(ctx context.Context) *TokenClaims {
	if claims, ok := ctx.Value(ClaimsContextKey).(*TokenClaims); ok {
		return claims
	}
	return nil
}
