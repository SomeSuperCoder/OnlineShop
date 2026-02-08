package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/internal/redisclient"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CartHandler struct {
	RedisClient *redis.Client
	AppConfig   *internal.AppConfig
}

type GetCartResponse struct {
	Body struct {
		Cart []string `json:"cart"`
	}
}

func (h *CartHandler) Get(ctx context.Context, input *struct{}) (*GetCartResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(GetCartResponse)
	cartValue, err := redisclient.GetCart(ctx, h.RedisClient, claims.UUID)
	resp.Body.Cart = cartValue
	return resp, err
}

type AddItemToCartRequest struct {
	Body struct {
		Item uuid.UUID `json:"item" format:"uuid"`
	}
}
type AddItemToCartResponse struct {
	Body redisclient.CartModificationResult
}

func (h *CartHandler) Post(ctx context.Context, input *AddItemToCartRequest) (*AddItemToCartResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(AddItemToCartResponse)
	result, err := redisclient.AddItemToCart(ctx, h.RedisClient, input.Body.Item, claims.UUID, h.AppConfig)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}

type RemoveItemFromCartRequest struct {
	ID uuid.UUID `path:"id" format:"uuid"`
}
type RemoveItemFromCartResponse struct {
	Body redisclient.CartModificationResult
}

func (h *CartHandler) Delete(ctx context.Context, input *RemoveItemFromCartRequest) (*RemoveItemFromCartResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(RemoveItemFromCartResponse)
	result, err := redisclient.RemoveItemFromCart(ctx, h.RedisClient, input.ID, claims.UUID)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}
