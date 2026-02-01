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
	cartValue, err := cart.GetCart(ctx, h.RedisClient, claims.Username)
	resp.Body.Cart = cartValue
	return resp, err
}

type AddItemToCartRequest struct {
	Body struct {
		Item uuid.UUID `json:"item" format:"uuid"`
	}
}
type AddItemToCartResponse struct {
	Body cart.CartModificationResult
}

func (h *CartHandler) Post(ctx context.Context, input *AddItemToCartRequest) (*AddItemToCartResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(AddItemToCartResponse)
	result, err := cart.AddItem(ctx, h.RedisClient, input.Body.Item, claims.Username)
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
	Body cart.CartModificationResult
}

func (h *CartHandler) Delete(ctx context.Context, input *RemoveItemFromCartRequest) (*RemoveItemFromCartResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(RemoveItemFromCartResponse)
	result, err := cart.RemoveItem(ctx, h.RedisClient, input.ID, claims.Username)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}
