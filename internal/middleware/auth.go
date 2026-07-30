package middleware

import (
	"context"
	"net/http"
	"strings"

	"coffee-site-server/internal/auth"
	"coffee-site-server/internal/httpx"
)

type contextKey string

const claimsKey contextKey = "claims"

// RequireAuth parses the bearer token and stores the claims in the request
// context. It rejects the request with 401 when the token is missing or invalid.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := claimsFromRequest(secret, r)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "請先登入")
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth parses the bearer token when present but never rejects the request.
func OptionalAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if claims, err := claimsFromRequest(secret, r); err == nil {
				r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin must be chained after RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok || claims.Permission != "admin" {
			httpx.Error(w, http.StatusForbidden, "需要管理員權限")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func claimsFromRequest(secret string, r *http.Request) (*auth.Claims, error) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return nil, errMissingToken
	}
	raw := strings.TrimPrefix(header, prefix)
	return auth.ParseToken(secret, raw)
}

var errMissingToken = &missingTokenError{}

type missingTokenError struct{}

func (e *missingTokenError) Error() string { return "missing bearer token" }

func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*auth.Claims)
	return claims, ok
}
