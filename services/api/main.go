package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SomeSuperCoder/OnlineShop/handlers"
	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	appConfig := internal.LoadAppConfig()
	pool, repo, _ := internal.DatabaseConnect(ctx, appConfig)
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
	if !appConfig.TestMode {
		log.Println("Adding auth middleware")
		api.UseMiddleware(middleware.AuthMiddleware(api, appConfig))
	}

	MountRoutes(api, repo, appConfig)

	r.Run(fmt.Sprintf(":%s", appConfig.Port))
}

func MountRoutes(api huma.API, repo *repository.Queries, appConfig *internal.AppConfig) {
	authHandler := handlers.AuthHandler{Repo: repo, Config: appConfig}
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
		}, authHandler.Register)
	}

	productHandler := handlers.ProductHandler{Repo: repo}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products",
			Tags:    []string{"Products"},
			Summary: "Get all products",
		}, productHandler.GetAll)
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
			Method:  http.MethodDelete,
			Path:    "/products/{id}",
			Tags:    []string{"Products"},
			Summary: "Delete Product",
		}, productHandler.Delete)
	}

	reviewHandler := handlers.ReviewHandler{Repo: repo}
	{
		huma.Register(api, huma.Operation{
			Method:  http.MethodGet,
			Path:    "/products/{id}/reviews",
			Tags:    []string{"Reviews"},
			Summary: "Get reviews for product",
		}, reviewHandler.GetFor)
		huma.Register(api, huma.Operation{
			Method:  http.MethodPost,
			Path:    "/products/{id}/reviews",
			Tags:    []string{"Reviews"},
			Summary: "Create review for product",
		}, reviewHandler.Post)
		huma.Register(api, huma.Operation{
			Method:  http.MethodDelete,
			Path:    "/reviews/{id}",
			Tags:    []string{"Reviews"},
			Summary: "Delete review",
		}, reviewHandler.Delete)
	}

}
