package middleware

import (
	"net/http"

	"cloud.google.com/go/firestore"
)

// PublicRoute defines a method + path pair that skips key validation.
type PublicRoute struct {
	Method string
	Path   string
}

// APIKeyAuth returns middleware that validates the X-API-Key header
// against keys stored in Firestore. Public paths are skipped.
func APIKeyAuth(client *firestore.Client, public []PublicRoute) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, pr := range public {
				if r.Method == pr.Method && r.URL.Path == pr.Path {
					next.ServeHTTP(w, r)
					return
				}
			}

			key := r.Header.Get("X-API-Key")
			if key == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"missing API key"}`))
				return
			}

			ak, err := getAPIKeyFn(r.Context(), client, key)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"failed to validate API key"}`))
				return
			}
			if ak == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":"invalid or revoked API key"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
