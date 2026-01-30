package internal

import (
	"context"

	"github.com/SomeSuperCoder/HumaExampleProject/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DatabaseConnect(ctx context.Context, config *AppConfig) (*pgxpool.Pool, *repository.Queries) {
	pool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		panic(err)
	}

	repo := repository.New(pool)

	return pool, repo
}
