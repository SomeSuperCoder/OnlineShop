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

type UserUpdateRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name     *string `json:"name,omitempty"`
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
	}
}
type UserUpdateResponse struct {
	Body repository.UpdateUserInfoRow
}

func (h *UserHandler) Patch(ctx context.Context, input *UserUpdateRequest) (*UserUpdateResponse, error) {
	resp := new(UserUpdateResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.UpdateUserInfoRow, error) {
		if is, err := q.Is(ctx, repository.IsParams{
			UserID: input.ID,
		}); err != nil {
			return repository.UpdateUserInfoRow{}, err
		} else if !is {
			return repository.UpdateUserInfoRow{}, AccessDeniedError
		}

		return q.UpdateUserInfo(ctx, repository.UpdateUserInfoParams{
			ID:       input.ID,
			Name:     input.Body.Name,
			Email:    input.Body.Email,
			Username: input.Body.Username,
		})
	})
	resp.Body = *res
	return resp, err
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
