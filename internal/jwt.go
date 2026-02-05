package internal

import (
	"context"
	"time"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Claims struct {
	UUID     uuid.UUID       `json:"uuid"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     repository.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(ctx context.Context, repo *repository.Queries, email string, config *AppConfig) (string, error) {
	user, err := repo.UnsafeGetUserByEmail(ctx, repository.UnsafeGetUserByEmailParams{
		Email: email,
	})
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UUID:     user.ID,
		Username: user.Email,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.AccessTokenSecret)
}

func ValidateToken(tokenString string, config *AppConfig) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return config.AccessTokenSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
