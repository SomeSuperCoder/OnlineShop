package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

type UserDeleteRequest struct {
	ID uuid.UUID `path:"id"`
}
type UserDeleteResponse struct {
	Body repository.DeleteUserRow
}

func (h *UserHandler) Delete(ctx context.Context, input *UserDeleteRequest) (*UserDeleteResponse, error) {
	resp := new(UserDeleteResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.DeleteUserRow, error) {
		if is, err := q.Is(ctx, repository.IsParams{
			UserID: input.ID,
		}); err != nil {
			return repository.DeleteUserRow{}, err
		} else if !is {
			return repository.DeleteUserRow{}, AccessDeniedError
		}

		return q.DeleteUser(ctx, repository.DeleteUserParams{
			ID: input.ID,
		})
	})
	resp.Body = *res
	return resp, err
}
