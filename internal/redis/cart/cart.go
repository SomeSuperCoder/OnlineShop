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
	Cart  []string `json:"cart"`
	Len   int      `json:"len"`
	Added int64    `json:"added" description:"the amount of new entries added to the cart"`
}

func AddItem(ctx context.Context, rdb *redis.Client, item uuid.UUID, username string) (*AddItemResult, error) {
	pipeline := rdb.TxPipeline()

	key := GenerateCartKey(username)
	pushCmd := pipeline.ZAdd(ctx, key, redis.Z{
		Score:  0,
		Member: item.String(),
	})
	getCmd := pipeline.ZRange(ctx, key, 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	added, err := pushCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to push a new item to cart: %w", err)
	}

	newCart, err := getCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the new cart value: %w", err)
	}

	return &AddItemResult{
		Cart:  newCart,
		Len:   len(newCart),
		Added: added,
	}, nil
}
