package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	Repo *repository.Queries
}

type CreateReviewRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Comment string `json:"comment"`
		Stars   int32  `json:"stars"`
	}
}
type CreateReviewResponse struct {
	Body *repository.Review
}

func (h *ReviewHandler) Post(ctx context.Context, input *CreateReviewRequest) (*CreateReviewResponse, error) {
	resp := new(CreateReviewResponse)
	res, err := h.Repo.InsertReview(ctx, repository.InsertReviewParams{
		Product: input.ID,
		Comment: &input.Body.Comment,
		Stars:   input.Body.Stars,
	})
	resp.Body = &res
	return resp, err
}

type GetReviewsForProductRequest struct {
	Pagination
	ID uuid.UUID `path:"id"`
}
type GetReviewsForProducttResponse struct {
	Body []repository.Review
}

func (h *ReviewHandler) GetFor(ctx context.Context, input *GetReviewsForProductRequest) (*GetReviewsForProducttResponse, error) {
	resp := new(GetReviewsForProducttResponse)
	res, err := h.Repo.GetReviewsForProduct(ctx, repository.GetReviewsForProductParams{
		Product: input.ID,
		Limit:   input.Limit,
		Offset:  input.Offset,
	})
	resp.Body = res
	return resp, err
}

type DeleteReviewRequest struct {
	ID uuid.UUID `path:"id"`
}
type DeleteReviewResponse struct {
	Body *repository.Review
}

func (h *ReviewHandler) Delete(ctx context.Context, input *DeleteReviewRequest) (*DeleteReviewResponse, error) {
	resp := new(DeleteReviewResponse)
	res, err := h.Repo.DeleteReview(ctx, repository.DeleteReviewParams(*input))
	resp.Body = &res
	return resp, err
}
