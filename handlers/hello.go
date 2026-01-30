package handlers

import (
	"context"
	"fmt"
)

type HelloInput struct {
	Name string `query:"name" example:"Bob" required:"true"`
}

type HelloResponse struct {
	Body struct {
		Message string `json:"message" example:"Hello, World!"`
	}
}

func Hello(ctx context.Context, input *HelloInput) (*HelloResponse, error) {
	resp := new(HelloResponse)
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)

	return resp, nil
}
