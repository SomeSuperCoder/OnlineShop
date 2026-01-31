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
	config := internal.LoadAppConfig()
	pool, repo := internal.DatabaseConnect(ctx, config)
	defer pool.Close()

	r := gin.Default()

	api := humagin.New(r, huma.DefaultConfig(
		"Online Shop Huma + Gin API",
		"1.0.0",
	))

	orderHandler := handlers.ProductHandler{Repo: repo}
	huma.Get(api, "/orders", orderHandler.GetAll)
	huma.Get(api, "/orders/{id}", orderHandler.GetByID)
	huma.Post(api, "/orders", orderHandler.Post)
	huma.Delete(api, "/orders/{id}", orderHandler.Delete)

	reviewHandler := handlers.ReviewHandler{Repo: repo}
	huma.Get(api, "/orders/{id}/reviews", reviewHandler.GetFor)
	huma.Post(api, "/orders/{id}/reviews", reviewHandler.Post)
	huma.Delete(api, "/reviews/{id}", reviewHandler.Delete)

	r.Run(fmt.Sprintf(":%s", config.Port))
}
