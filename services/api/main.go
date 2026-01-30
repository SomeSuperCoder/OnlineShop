package main

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/HumaExampleProject/handlers"
	"github.com/SomeSuperCoder/HumaExampleProject/internal"
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
		"Huma + Gin API",
		"1.0.0",
	))

	huma.Get(api, "/hello", handlers.Hello)

	orderHandler := handlers.OrderHandler{Repo: repo}
	huma.Get(api, "/orders", orderHandler.GetAll)
	huma.Post(api, "/orders", orderHandler.Post)

	r.Run(fmt.Sprintf(":%s", config.Port))
}
