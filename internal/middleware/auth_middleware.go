package middleware

import "github.com/danielgtaylor/huma/v2"

func MyMiddleware(ctx huma.Context, next func(huma.Context)) {
	// Set a custom header on the response.
	ctx.SetHeader("My-Custom-Header", "Hello, world!")

	// Call the next middleware in the chain. This eventually calls the
	// operation handler as well.
	next(ctx)
}
