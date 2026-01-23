package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/razedwell/go-hand/internal/security"
)

type ctxKey string

const AcsKey ctxKey = "access_token"

func Auth(jwt *security.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(h, "Bearer ")

			if jwt.IsBlacklisted(r.Context(), tokenStr) {
				http.Error(w, "token revoked", http.StatusUnauthorized)
				return
			}

			_, err := jwt.Verify(tokenStr)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), AcsKey, tokenStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
