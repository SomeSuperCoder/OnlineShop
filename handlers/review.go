package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

type CreateReviewRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Comment string `json:"comment"`
		Stars   int32  `json:"stars" default:"1"`
	}
}
type CreateReviewResponse struct {
	Body repository.Review
}

func (h *ReviewHandler) Post(ctx context.Context, input *CreateReviewRequest) (*CreateReviewResponse, error) {
	resp := new(CreateReviewResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Review, error) {
		return q.InsertReview(ctx, repository.InsertReviewParams{
			Product: input.ID,
			Comment: &input.Body.Comment,
			Stars:   input.Body.Stars,
		})
	})
	resp.Body = *res
	return resp, err
}

type GetReviewsForProductRequest struct {
	Pagination
	ID uuid.UUID `path:"id"`
}
type GetReviewsForProducttResponse struct {
	Body []repository.GetReviewsForProductRow
}

func (h *ReviewHandler) GetFor(ctx context.Context, input *GetReviewsForProductRequest) (*GetReviewsForProducttResponse, error) {
	resp := new(GetReviewsForProducttResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) ([]repository.GetReviewsForProductRow, error) {
		return q.GetReviewsForProduct(ctx, repository.GetReviewsForProductParams{
			Product: input.ID,
			Limit:   input.Limit,
			Offset:  input.Offset,
		})
	})
	resp.Body = *res
	return resp, err
}

type UpdateReviewRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Comment *string `json:"comment,omitempty"`
		Stars   *int32  `json:"stars,omitempty"`
	}
}
type UpdateReviewResponse struct {
	Body repository.Review
}

func (h *ReviewHandler) Patch(ctx context.Context, input *UpdateReviewRequest) (*UpdateReviewResponse, error) {
	resp := new(UpdateReviewResponse)

	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Review, error) {
		if is, err := q.IsAuthor(ctx, repository.IsAuthorParams{
			ID: input.ID,
		}); err != nil {
			return repository.Review{}, err
		} else if !is {
			return repository.Review{}, AccessDeniedError
		}

		return q.UpdateReview(ctx, repository.UpdateReviewParams{
			ID:      input.ID,
			Comment: input.Body.Comment,
			Stars:   input.Body.Stars,
		})
	})

	resp.Body = *res
	return resp, err
}

type DeleteReviewRequest struct {
	ID uuid.UUID `path:"id"`
}
type DeleteReviewResponse struct {
	Body repository.Review
}

func (h *ReviewHandler) Delete(ctx context.Context, input *DeleteReviewRequest) (*DeleteReviewResponse, error) {
	resp := new(DeleteReviewResponse)
	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Review, error) {
		if is, err := q.IsOwnerOrAuthor(ctx, repository.IsOwnerOrAuthorParams{
			ID: input.ID,
		}); err != nil {
			return repository.Review{}, err
		} else if !is {
			return repository.Review{}, AccessDeniedError
		}

		return q.DeleteReview(ctx, repository.DeleteReviewParams(*input))

	})
	resp.Body = *res
	return resp, err
}
