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
	users    UserLookup
}

// UserLookup finds a dashboard member by username. Member tokens are checked against it on
// every request, so a disabled or removed member loses access at once, and a role change applies
// without waiting for the token to expire.
type UserLookup interface {
	// Lookup returns the member's role ("admin", "operator", "viewer") and whether the account
	// can sign in. found is false when the member does not exist.
	Lookup(username string) (role string, active bool, found bool)
}

// MemberTokenPrefix marks tokens issued to dashboard members. The username follows it.
const MemberTokenPrefix = "user:"

// adminOnlyPrefixes are areas that only the admin role may change.
var adminOnlyPrefixes = []string{"/api/v1/users", "/api/v1/integrations", "/api/v1/server/settings"}

func NewAuthMiddleware(tokenMgr *TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenMgr: tokenMgr}
}

// WithUsers enables member checks. Without it, tokens behave as before.
func (a *AuthMiddleware) WithUsers(users UserLookup) *AuthMiddleware {
	a.users = users
	return a
}

// roleClaim maps a member role to the role name the rest of the agent checks.
func roleClaim(role string) string {
	switch role {
	case "admin":
		return "admin"
	case "operator":
		return "deployer"
	default:
		return "viewer"
	}
}

// enforceRole returns an error message when the request is outside the member's role.
func enforceRole(claims *TokenClaims, r *http.Request) string {
	p := r.URL.Path
	if p == "/api/v1/users/me/password" {
		return ""
	}
	mutating := r.Method != http.MethodGet && r.Method != http.MethodHead
	if claims.Role == "viewer" && mutating {
		return "this account is read-only"
	}
	if mutating && claims.Role != "admin" {
		for _, prefix := range adminOnlyPrefixes {
			if strings.HasPrefix(p, prefix) {
				return "only the admin role can change this"
			}
		}
	}
	return ""
}

// Wrap guards an http.Handler requiring a valid Bearer token or X-Orbit-Key.
func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow unauthenticated health check, public static UI assets, and S3
		// (S3 enforces its own auth). Telemetry and admin APIs require a token.
		p := r.URL.Path
		if p == "/healthz" || p == "/" || p == "/index.html" || p == "/icon.svg" ||
			p == "/manifest.webmanifest" || p == "/sw.js" ||
			strings.HasPrefix(p, "/assets/") || strings.HasPrefix(p, "/s3/") ||
			p == "/api/v1/auth/status" || p == "/api/v1/auth/setup" || p == "/api/v1/auth/login" {
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

		if a.users != nil && strings.HasPrefix(claims.TokenID, MemberTokenPrefix) {
			name := strings.TrimPrefix(claims.TokenID, MemberTokenPrefix)
			if i := strings.LastIndex(name, ":"); i >= 0 {
				name = name[:i]
			}
			role, active, found := a.users.Lookup(name)
			if !found || !active {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "this account is no longer active"})
				return
			}
			claims = &TokenClaims{TokenID: claims.TokenID, Role: roleClaim(role), IssuedAt: claims.IssuedAt, ExpiresAt: claims.ExpiresAt}
		}
		if msg := enforceRole(claims, r); msg != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
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
	// EventSource cannot send headers, so SSE streams may take the token from the query.
	// Restricted to stream endpoints so the token never lands in ordinary request URLs.
	if isStreamPath(r.URL.Path) {
		if queryToken := r.URL.Query().Get("token"); queryToken != "" {
			return queryToken
		}
	}
	return ""
}

func isStreamPath(p string) bool {
	return strings.HasSuffix(p, "/stream") || strings.HasSuffix(p, "/logs/stream")
}

func GetClaims(ctx context.Context) *TokenClaims {
	if claims, ok := ctx.Value(ClaimsContextKey).(*TokenClaims); ok {
		return claims
	}
	return nil
}
