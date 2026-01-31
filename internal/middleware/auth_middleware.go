package middleware

import (
	"net/http"
	"strings"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/danielgtaylor/huma/v2"
)

const AuthClaimsContextKey = "claims"

func AuthMiddleware(api huma.API, config *internal.AppConfig) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		path := ctx.Operation().Path
		if strings.HasPrefix(path, "/auth") {
			next(ctx)
			return
		}

		authHeader := ctx.Header("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid Authorization header format, expecting value of type `Bearer <jwt token>`")
			return
		}

		jwt := parts[1]
		claims, err := internal.ValidateToken(jwt, config)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "JWT token verificatin failed due to", err)
			return
		}

		// Call the next middleware in the chain. This eventually calls the
		// operation handler as well.
		next(huma.WithValue(ctx, AuthClaimsContextKey, claims))
	}
}
