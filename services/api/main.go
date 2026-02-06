package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/OnlineShop/handlers"
	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	appConfig := internal.LoadAppConfig()
	pool, repo, redisClient := internal.DatabaseConnect(ctx, appConfig)
	defer pool.Close()

	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	humaConfig := huma.DefaultConfig(
		"Online Shop Huma + Gin API",
		"1.0.0",
	)
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8888/api/v1", Description: "Local API version 1"},
	}
	api := humagin.NewWithGroup(r, apiGroup, humaConfig)
	fmt.Printf("appConfig.TestMode: %v\n", appConfig.TestMode)
	api.UseMiddleware(middleware.AuthMiddleware(api, appConfig))

	MountRoutes(api, repo, pool, redisClient, appConfig)

	r.Run(fmt.Sprintf(":%s", appConfig.Port))
}

func MountRoutes(api huma.API, repo *repository.Queries, pool *pgxpool.Pool, redisClient *redis.Client, appConfig *internal.AppConfig) {
	authHandler := handlers.AuthHandler{Repo: repo, AppConfig: appConfig, Redis: redisClient}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Tags:    []string{"Auth"},
			Summary: "Register",
		}, authHandler.Register)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Tags:    []string{"Auth"},
			Summary: "Login",
		}, authHandler.Login)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Tags:    []string{"Auth"},
			Summary: "Refresh",
		}, authHandler.Refresh)
	}

	productHandler := handlers.ProductHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products",
			Tags:    []string{"Products"},
			Summary: "Get products paged",
		}, productHandler.GetPaged)
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products/search",
			Tags:    []string{"Products"},
			Summary: "Search for products",
		}, productHandler.SearchForProducts)
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products/{id}",
			Tags:    []string{"Products"},
			Summary: "Get product by ID",
		}, productHandler.GetByID)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/products",
			Tags:    []string{"Products"},
			Summary: "Create product",
		}, productHandler.Post)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/products/{id}",
			Tags:    []string{"Products"},
			Summary: "Update product",
		}, productHandler.Patch)
		huma.Register(api, huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/products/{id}",
			Tags:    []string{"Products"},
			Summary: "Delete Product",
		}, productHandler.Delete)
	}

	reviewHandler := handlers.ReviewHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products/{id}/reviews",
			Tags:    []string{"Reviews"},
			Summary: "Get reviews for product paged",
		}, reviewHandler.GetFor)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/products/{id}/reviews",
			Tags:    []string{"Reviews"},
			Summary: "Create review for product",
		}, reviewHandler.Post)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/reviews/{id}",
			Tags:    []string{"Reviews"},
			Summary: "Update review",
		}, reviewHandler.Patch)
		huma.Register(api, huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/reviews/{id}",
			Tags:    []string{"Reviews"},
			Summary: "Delete review",
		}, reviewHandler.Delete)
	}

	cartHandler := handlers.CartHandler{RedisClient: redisClient, AppConfig: appConfig}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/cart",
			Tags:    []string{"Cart"},
			Summary: "Get cart for the current user",
		}, cartHandler.Get)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/cart",
			Tags:    []string{"Cart"},
			Summary: "Add item to cart",
		}, cartHandler.Post)
		huma.Register(api, huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/cart/{id}",
			Tags:    []string{"Cart"},
			Summary: "Remove item from cart",
		}, cartHandler.Delete)
	}

	voteHandler := handlers.VotesHandler{Repo: repo, Pool: pool}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/reviews/{id}/votes",
			Tags:    []string{"Votes"},
			Summary: "Upvote or downvote a review",
		}, voteHandler.Post)
	}
	userHandler := handlers.UserHandler{Repo: repo, Pool: pool, Redis: redisClient}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Tags:    []string{"Users"},
			Summary: "Update a user",
		}, userHandler.Patch)
		huma.Register(api, huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Tags:    []string{"Users"},
			Summary: "Delete a user",
		}, userHandler.Delete)
	}
}
