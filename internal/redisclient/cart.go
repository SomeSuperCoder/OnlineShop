package redisclient

import (
	"context"
	"fmt"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GenerateCartKey(userUUID uuid.UUID) string {
	return fmt.Sprintf("cart:%s", userUUID.String())
}

func GetCart(ctx context.Context, rdb *redis.Client, userUUID uuid.UUID) ([]string, error) {
	return rdb.ZRange(ctx, GenerateCartKey(userUUID), 0, -1).Result()
}

type CartModificationResult struct {
	Cart      []string `json:"cart"`
	Len       int      `json:"len"`
	LenChange int64    `json:"len_change" description:"the amount of new entries added to or removed from the cart"`
}

func AddItemToCart(ctx context.Context, rdb *redis.Client, item uuid.UUID, userUUID uuid.UUID, appConfig *internal.AppConfig) (*CartModificationResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateCartKey(userUUID)
	addCmd := pipeline.ZAdd(ctx, key, redis.Z{
		Score:  0,
		Member: item.String(),
	})
	expireCmd := pipeline.Expire(ctx, key, appConfig.UserCartExpiry)
	getCmd := pipeline.ZRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	added, err := addCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to add a new item to cart: %w", err)
	}

	if expireCmd.Err() != nil {
		return nil, fmt.Errorf("faild to set cart expiry: %w", err)
	}

	newCart, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new cart value: %w", err)
	}

	return &CartModificationResult{
		Cart:      newCart,
		Len:       len(newCart),
		LenChange: added,
	}, nil
}

func RemoveItemFromCart(ctx context.Context, rdb *redis.Client, item uuid.UUID, userUUID uuid.UUID) (*CartModificationResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateCartKey(userUUID)
	removeCmd := pipeline.ZRem(ctx, key, item.String())
	getCmd := pipeline.ZRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	removed, err := removeCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}

	newCart, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new cart value: %w", err)
	}

	return &CartModificationResult{
		Cart:      newCart,
		Len:       len(newCart),
		LenChange: removed,
	}, nil
}

func DeleteCart(ctx context.Context, rdb *redis.Client, userUUID uuid.UUID) error {
	return rdb.Del(ctx, GenerateCartKey(userUUID)).Err()
}
