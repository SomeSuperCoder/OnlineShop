package handlers

import (
	"context"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	Repo      *repository.Queries
	Redis     *redis.Client
	AppConfig *internal.AppConfig
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
	Body internal.TokenPair
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

	// Get the user from the database by email
	user, err := h.Repo.UnsafeGetUserByEmail(ctx, repository.UnsafeGetUserByEmailParams{
		Email: input.Body.Email,
	})

	// Generate a new token pair
	tokenPair, err := internal.GenerateTokenPair(ctx, h.Redis, user, h.AppConfig)
	if err != nil {
		return resp, err
	}

	resp.Body = tokenPair

	return resp, err
}

type RefreshRequest struct {
	Body struct {
		RefreshToken string `json:"refresh_token"`
	}
}
type RefreshResponse struct {
	Body internal.TokenPair
}

func (h *AuthHandler) Refresh(ctx context.Context, input *RefreshRequest) (*RefreshResponse, error) {
	resp := new(RefreshResponse)

	// Verify the refresh token
	userID, err := internal.ValidateRefreshToken(ctx, h.Redis, input.Body.RefreshToken)
	if err != nil {
		return nil, huma.Error401Unauthorized("failed to verify the refresh token", err)
	}

	// Get the user from the database by ID
	user, err := h.Repo.UnsafeGetUserByID(ctx, repository.UnsafeGetUserByIDParams{
		ID: userID,
	})

	// Remove the old token
	err = internal.InvalidateRefreshToken(ctx, h.Redis, input.Body.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Generate a new token pair
	tokenPair, err := internal.GenerateTokenPair(ctx, h.Redis, user, h.AppConfig)
	if err != nil {
		return resp, err
	}
	resp.Body = tokenPair
	return resp, err
}
