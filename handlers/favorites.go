package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/internal/redisclient"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FavoritesHandler struct {
	RedisClient *redis.Client
	AppConfig   *internal.AppConfig
}

type GetFavoritesResponse struct {
	Body struct {
		Favorites []string `json:"favorites"`
	}
}

func (h *FavoritesHandler) Get(ctx context.Context, input *struct{}) (*GetFavoritesResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(GetFavoritesResponse)
	favoritesValue, err := redisclient.GetFavorites(ctx, h.RedisClient, claims.UUID)
	resp.Body.Favorites = favoritesValue
	return resp, nil
}

type AddItemToFavoritesRequest struct {
	Body struct {
		Item uuid.UUID `json:"item" format:"uuid"`
	}
}
type AddItemToFavoritesResponse struct {
	Body redisclient.FavoritesModificationResult
}

func (h *FavoritesHandler) Post(ctx context.Context, input *AddItemToFavoritesRequest) (*AddItemToFavoritesResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(AddItemToFavoritesResponse)
	result, err := redisclient.AddFavorite(ctx, h.RedisClient, input.Body.Item, claims.UUID, h.AppConfig)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}

type RemoveItemFromFavoritesRequest struct {
	ID uuid.UUID `query:"item" format:"uuid"`
}
type RemoveItemFromFavoritesResponse struct {
	Body redisclient.FavoritesModificationResult
}

func (h *FavoritesHandler) Delete(ctx context.Context, input *RemoveItemFromFavoritesRequest) (*RemoveItemFromFavoritesResponse, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(RemoveItemFromFavoritesResponse)
	result, err := redisclient.RemoveFavorite(ctx, h.RedisClient, input.ID, claims.UUID)
	if err != nil {
		return nil, err
	}
	resp.Body = *result
	return resp, nil
}
