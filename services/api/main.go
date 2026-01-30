package main

import (
	"github.com/SomeSuperCoder/HumaExampleProject/handlers"
	"github.com/SomeSuperCoder/HumaExampleProject/middleware"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := humagin.New(r, huma.DefaultConfig(
		"Huma + Gin API",
		"1.0.0",
	))
	api.UseMiddleware(middleware.AuthMiddleware)

	huma.Get(api, "/hello", handlers.Hello)

	r.Run(":8888")
}
