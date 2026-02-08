package middleware

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func WithAuthContext[T any](ctx context.Context, pool *pgxpool.Pool, repo *repository.Queries, fn func(context.Context, *repository.Queries) (T, error)) (*T, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		logrus.Errorln("Failed to begin a transaction")
		return new(T), err
	}
	defer tx.Rollback(ctx)

	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return new(T), huma.Error401Unauthorized("Failed to extract JWT claims from context", err)
	}

	qtx := repo.WithTx(tx)
	_, err = qtx.SetConfig(ctx, repository.SetConfigParams{
		UserID: claims.UUID.String(),
	})
	if err != nil {
		logrus.Errorln("Failed to set config params")
		return new(T), err
	}

	fResult, err := fn(ctx, qtx)
	if err != nil {
		logrus.Errorf("Inner function failed due to: %s", err.Error())
		return new(T), err
	}

	return &fResult, tx.Commit(ctx)
}
