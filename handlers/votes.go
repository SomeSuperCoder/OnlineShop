package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal/middleware"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VotesHandler struct {
	Repo *repository.Queries
	Pool *pgxpool.Pool
}

type VoteRequest struct {
	ReviewID uuid.UUID `path:"id"`
	Body     struct {
		Type repository.VoteType `json:"type" enum:"upvote,downvote"`
	}
}
type VoteResponse struct {
	Body repository.Vote
}

func (h *VotesHandler) Post(ctx context.Context, input *VoteRequest) (*VoteResponse, error) {
	resp := new(VoteResponse)

	res, err := middleware.WithAuthContext(ctx, h.Pool, h.Repo, func(ctx context.Context, q *repository.Queries) (repository.Vote, error) {
		if is, err := q.IsOwnerOrAuthor(ctx, repository.IsOwnerOrAuthorParams{
			ID: input.ReviewID,
		}); err != nil {
			return repository.Vote{}, err
		} else if is {
			return repository.Vote{}, AccessDeniedError
		}
		return q.InsertVote(ctx, repository.InsertVoteParams{
			Review: input.ReviewID,
			Type:   input.Body.Type,
		})
	})

	resp.Body = *res

	return resp, err
}
