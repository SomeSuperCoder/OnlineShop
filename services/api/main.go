package main

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/handlers"
	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	appConfig := internal.LoadAppConfig()
	pool, repo := internal.DatabaseConnect(ctx, appConfig)
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

	authHandler := handlers.AuthHandler{Repo: repo, Config: appConfig}
	{
		huma.Post(api, "/auth/register", authHandler.Register)
		huma.Post(api, "/auth/login", authHandler.Login)
	}

	orderHandler := handlers.ProductHandler{Repo: repo}
	{
		huma.Get(api, "/orders", orderHandler.GetAll)
		huma.Get(api, "/orders/{id}", orderHandler.GetByID)
		huma.Post(api, "/orders", orderHandler.Post)
		huma.Delete(api, "/orders/{id}", orderHandler.Delete)
	}

	reviewHandler := handlers.ReviewHandler{Repo: repo}
	{
		huma.Get(api, "/orders/{id}/reviews", reviewHandler.GetFor)
		huma.Post(api, "/orders/{id}/reviews", reviewHandler.Post)
		huma.Delete(api, "/reviews/{id}", reviewHandler.Delete)
	}

	r.Run(fmt.Sprintf(":%s", appConfig.Port))
}
