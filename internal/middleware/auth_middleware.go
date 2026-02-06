package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

const AuthClaimsContextKey = "claims"
const TestUsername = "test-user"

func GetClaimsFromContext(ctx context.Context) (*internal.Claims, error) {
	val := ctx.Value(AuthClaimsContextKey)
	log.Println(val)
	valParsed, ok := val.(internal.Claims)
	log.Println(ok)
	log.Println(valParsed)

	log.Println("BEGIN DEEP DIVE")
	log.Printf("Type: %T", val)
	log.Printf("Kind: %v", reflect.TypeOf(val).Kind())
	log.Printf("VAlues: %+v", val)
	log.Println("END DEEP DIVE")

	if claims, ok := ctx.Value(AuthClaimsContextKey).(*internal.Claims); ok {
		return claims, nil
	} else {
		return nil, huma.Error401Unauthorized("faield to extract JWT claims form context")
	}
}

func AuthMiddleware(api huma.API, config *internal.AppConfig) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		path := ctx.Operation().Path
		if strings.HasPrefix(path, "/auth") {
			next(ctx)
			return
		}

		if config.TestMode {
			next(huma.WithValue(ctx, AuthClaimsContextKey, internal.Claims{
				UUID:     uuid.MustParse("0a381e5b-7c7e-4fc7-9441-133871d89223"),
				Username: TestUsername,
				Email:    "test@localhost",
				Role:     repository.RoleUser,
			}))

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

		// TODO: remove
		jsonData, _ := json.Marshal(claims)
		log.Println("===== BEGIN auth middleware =====")
		log.Println(string(jsonData))
		log.Println("===== END auth middleware =====")

		next(huma.WithValue(ctx, AuthClaimsContextKey, claims))
	}
}
