package cart

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GenerateCartKey(username string) string {
	return fmt.Sprintf("cart:%s", username)
}

type AddItemResult struct {
	Cart []string `json:"cart"`
	Len  int64    `json:"len"`
}

func AddItem(ctx context.Context, rdb *redis.Client, item uuid.UUID, username string) (*AddItemResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateCartKey(username)
	pushCmd := pipeline.LPush(ctx, key, item.String())
	getCmd := pipeline.LRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	newLen, err := pushCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to push a new item to cart: %w", err)
	}

	newCart, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new cart value: %w", err)
	}

	return &AddItemResult{
		Cart: newCart,
		Len:  newLen,
	}, nil
}
