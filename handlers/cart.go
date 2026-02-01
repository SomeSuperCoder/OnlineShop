package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/internal/redis/cart"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CartHandler struct {
	RedisClient *redis.Client
}

type CartHandlerRequest struct {
	Body struct {
		Item uuid.UUID `json:"item" format:"uuid"`
	}
}
type CartHandlerResponse struct {
	Body cart.AddItemResult
}

func (h *CartHandler) Post(ctx context.Context, input *CartHandlerRequest) (*CartHandlerResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(CartHandlerResponse)
	result, err := cart.AddItem(ctx, h.RedisClient, input.Body.Item, claims.Username)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}
