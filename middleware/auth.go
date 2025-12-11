// middleware/auth.go
package middleware

import (
	"net/http"
	"os"

	"blog-api/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("API_KEY")

		if expectedKey != "" && apiKey != expectedKey {
			utils.SendError(w, "API key tidak valid", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
