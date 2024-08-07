package auth

import (
	"context"
	"net/http"

	httperr "github.com/lGreenLightl/fitness-club-service/internal/app/server/err"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

func HTTPMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mapClaims jwt.MapClaims

		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(_ *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			},
			request.WithClaims(&mapClaims),
		)
		if err != nil {
			httperr.BadRequest(r, w, "unable-to-get-jwt", err)
			return
		}
		if !token.Valid {
			httperr.BadRequest(r, w, "invalid-jwt", nil)
			return
		}

		ctx := context.WithValue(r.Context(), customerContextKey, Customer{
			UUID:  mapClaims["customer_uuid"].(string),
			Name:  mapClaims["name"].(string),
			Role:  mapClaims["role"].(string),
			Email: mapClaims["email"].(string),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
