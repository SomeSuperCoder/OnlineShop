package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/repository"
)

type AuthHandler struct {
	Repo *repository.Queries
}

type RegisterRequest struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Username string `json:"username" pattern:"^[a-zA-Z0-9._-]+$" patternDescription:"alphanum with period, underscore and dash"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
}
type RegisterResponse struct {
	Body repository.InsertUserRow
}

func (h *AuthHandler) Register(ctx context.Context, input *RegisterRequest) (*RegisterResponse, error) {
	resp := new(RegisterResponse)
	res, err := h.Repo.InsertUser(ctx, repository.InsertUserParams{
		Email:    input.Body.Email,
		Username: input.Body.Username,
		Name:     input.Body.Name,
		Crypt:    input.Body.Password,
	})
	resp.Body = res
	return resp, err
}
