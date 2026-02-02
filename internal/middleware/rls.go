package middleware

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
)

func WithAuthContext(ctx context.Context, db *pgx.Conn, repo *repository.Queries, fn func(context.Context, *repository.Queries) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return huma.Error401Unauthorized("Failed to extract JWT claims from context", err)
	}

	qtx := repo.WithTx(tx)
	qtx.SetEmailParam(ctx, repository.SetEmailParamParams{
		SetConfig: claims.Email,
	})

	err = fn(ctx, qtx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
