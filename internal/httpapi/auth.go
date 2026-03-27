package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// authMiddleware protects endpoints with ADMIN_TOKEN.
//
// Supported:
// - Authorization: Bearer <token>
// - X-Admin-Token: <token>
func authMiddleware(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provided := readToken(r)
			if !secureEqual(provided, adminToken) {
				log.Warn().Str("remote", r.RemoteAddr).Msg("unauthorized")
				w.Header().Set("WWW-Authenticate", "Bearer")
				writeError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

var errUnauthorized = httpError("unauthorized")

type httpError string

func (e httpError) Error() string { return string(e) }

func readToken(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("X-Admin-Token")); v != "" {
		return v
	}

	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "
	if strings.HasPrefix(auth, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(auth, prefix))
	}

	return ""
}

func secureEqual(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
