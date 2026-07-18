// Package middleware provides HTTP middleware for the WMS backend.
package middleware

import (
	"encoding/json"
	"net/http"

	adminAuth "github.com/i56/i56-apps/i56-wms/internal/auth"
)

// APIJSON is the shared JSON response helper.
var APIJSON = func(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// AdminOnly returns a middleware that validates admin sessions via HMAC-signed cookies.
// Redirects to /admin/login on failure.
func AdminOnly(sm *adminAuth.SessionManager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ck, err := r.Cookie("admin_session")
			if err != nil {
				http.Redirect(w, r, "/admin/login", 303)
				return
			}
			session := sm.ValidateSession(ck.Value)
			if session == nil {
				http.Redirect(w, r, "/admin/login", 303)
				return
			}
			next(w, r)
		}
	}
}

// AdminOnlyAPI returns a middleware that validates admin sessions for JSON API endpoints.
// Returns 401 JSON on failure instead of redirect.
func AdminOnlyAPI(sm *adminAuth.SessionManager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ck, err := r.Cookie("admin_session")
			if err != nil {
				APIJSON(w, 401, map[string]string{"error": "unauthorized"})
				return
			}
			session := sm.ValidateSession(ck.Value)
			if session == nil {
				APIJSON(w, 401, map[string]string{"error": "invalid_session"})
				return
			}
			next(w, r)
		}
	}
}

// CORS returns a simple CORS middleware for API endpoints.
func CORS(allowedOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// Skip CORS for same-origin requests — setting empty Access-Control-Allow-Origin
			// breaks crossorigin module scripts (blank white screen).
			if origin == "" {
				next(w, r)
				return
			}
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed || len(allowedOrigins) == 0 {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if r.Method == "OPTIONS" {
				w.WriteHeader(204)
				return
			}
			next(w, r)
		}
	}
}

// Logger returns a simple request logging middleware.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
