package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
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

type LoginRequest struct {
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}
type LoginResponse struct {
	Body struct {
		JWT string `json:"jwt"`
	}
}

func (h *AuthHandler) Login(ctx context.Context, input *LoginRequest) (*LoginResponse, error) {
	resp := new(LoginResponse)

	// Check the credentials
	res, err := h.Repo.VerifyAuth(ctx, repository.VerifyAuthParams{
		Email: input.Body.Email,
		Crypt: input.Body.Password,
	})
	if err != nil {
		return resp, err
	}
	if res != true {
		return resp, huma.Error401Unauthorized("Invalid login credentials: ", err)
	}

	// Create a new JWT token
	jwt, err := internal.GenerateToken(input.Body.Email)
	if err != nil {
		return resp, err
	}

	resp.Body.JWT = jwt

	return resp, err
}
