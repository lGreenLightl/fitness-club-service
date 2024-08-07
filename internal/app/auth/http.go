package auth

import (
	"context"
	"net/http"
	"strings"

	httperr "github.com/lGreenLightl/fitness-club-service/internal/app/server/err"

	"firebase.google.com/go/auth"
)

type FirebaseHttpMiddleware struct {
	AuthClient *auth.Client
}

func (f FirebaseHttpMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken := f.tokenFromHeader(r)
		if bearerToken == "" {
			httperr.Unauthorized(r, w, "empty-bearer-token", nil)
		}

		authToken, err := f.AuthClient.VerifyIDToken(ctx, bearerToken)
		if err != nil {
			httperr.Unauthorized(r, w, "unable-to-verify-jwt", err)
			return
		}

		ctx = context.WithValue(ctx, customerContextKey, Customer{
			UUID:  authToken.UID,
			Name:  authToken.Claims["name"].(string),
			Role:  authToken.Claims["role"].(string),
			Email: authToken.Claims["email"].(string),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (f FirebaseHttpMiddleware) tokenFromHeader(r *http.Request) string {
	headerVal := r.Header.Get("Authorization")

	if len(headerVal) > 7 && strings.ToLower(headerVal[:6]) == "bearer" {
		return headerVal[7:]
	}

	return ""
}
