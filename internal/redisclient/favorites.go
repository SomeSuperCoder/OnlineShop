package redisclient

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GenerateFavoritesKey(userUUID uuid.UUID) string {
	return fmt.Sprintf("favorites:%s", userUUID.String())
}

func GetFavorites(ctx context.Context, rdb *redis.Client, userUUID uuid.UUID) ([]string, error) {
	return rdb.ZRange(ctx, GenerateFavoritesKey(userUUID), 0, -1).Result()
}

type FavoritesModificationResult struct {
	Favorites []string `json:"favorites"`
	Len       int      `json:"len"`
	LenChange int64    `json:"len_change" description:"the amount of entries that have been added to or removed from favorites"`
}

func AddFavorite(ctx context.Context, rdb *redis.Client, item uuid.UUID, userUUID uuid.UUID, appConfig *internal.AppConfig) (*FavoritesModificationResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateFavoritesKey(userUUID)
	addCmd := pipeline.ZAdd(ctx, key, redis.Z{
		Score:  0,
		Member: item.String(),
	})
	getCmd := pipeline.ZRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	added, err := addCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to add new item to favorites")
	}

	newFavorites, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new favorites value: %w", err)
	}

	return &FavoritesModificationResult{
		Favorites: newFavorites,
		Len:       len(newFavorites),
		LenChange: added,
	}, nil
}

func RemoveFavorite(ctx context.Context, rdb *redis.Client, item uuid.UUID, userUUID uuid.UUID) (*FavoritesModificationResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateFavoritesKey(userUUID)
	removeCmd := pipeline.ZRem(ctx, key, item.String())
	getCmd := pipeline.ZRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	removed, err := removeCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to remove item from favorites: %w", err)
	}

	newFavorites, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new favorites value: %w", err)
	}

	return &FavoritesModificationResult{
		Favorites: newFavorites,
		Len:       len(newFavorites),
		LenChange: removed,
	}, nil
}

func DeleteFavorites(ctx context.Context, rdb *redis.Client, userUUID uuid.UUID) error {
	return rdb.Del(ctx, GenerateFavoritesKey(userUUID)).Err()
}
