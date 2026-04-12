package middleware

import (
	"assignment_2/internal/firebase"
	"net/http"

	"cloud.google.com/go/firestore"
)

// APIKeyAuth returns a middleware that validates the X-API-Key header
// against keys stored in Firestore. It skips the paths listed in publicPaths.
func APIKeyAuth(client *firestore.Client, publicPaths []string) func(http.Handler) http.Handler {
	public := make(map[string]bool, len(publicPaths))
	for _, p := range publicPaths {
		public[p] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow public endpoints through without a key.
			if public[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get("X-API-Key")
			if key == "" {
				http.Error(w, "missing API key", http.StatusUnauthorized)
				return
			}

			ak, err := firebase.GetAPIKey(r.Context(), client, key)
			if err != nil {
				http.Error(w, "failed to validate API key", http.StatusInternalServerError)
				return
			}
			if ak == nil {
				http.Error(w, "invalid or revoked API key", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
