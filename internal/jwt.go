package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type Claims struct {
	UUID     uuid.UUID       `json:"uuid"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     repository.Role `json:"role"`
	jwt.RegisteredClaims
}

func generateAccessToken(user repository.User, appConfig *AppConfig) (string, time.Time, error) {
	expirationTime := time.Now().Add(appConfig.AccessTokenExpiry)
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
	tokenString, err := token.SignedString(appConfig.AccessTokenSecret)
	return tokenString, expirationTime, err
}

func RefreshTokenKey(token string) string {
	return fmt.Sprintf("refresh_token:%s", token)
}

func generateRefreshToken(ctx context.Context, userUUID uuid.UUID, rdb *redis.Client, appConfig *AppConfig) (string, time.Time, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(appConfig.RefreshTokenExpiry)

	key := RefreshTokenKey(token)
	err := rdb.Set(ctx, key, userUUID.String(), appConfig.RefreshTokenExpiry).Err()

	return token, expiresAt, err
}

func GenerateTokenPair(ctx context.Context, rdb *redis.Client, user repository.User, appConfig *AppConfig) (TokenPair, error) {
	// Generate the access token
	accessToken, accesssExpiresAt, err := generateAccessToken(user, appConfig)
	if err != nil {
		return TokenPair{}, err
	}

	// Generate the refresh token
	refreshToken, refreshExpiresAt, err := generateRefreshToken(ctx, user.ID, rdb, appConfig)
	if err != nil {
		return TokenPair{}, err
	}

	tokenPair := TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accesssExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}
	return tokenPair, nil
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

func ValidateRefreshToken(ctx context.Context, rdb *redis.Client, token string) (uuid.UUID, error) {
	key := RefreshTokenKey(token)

	userIDStr, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, errors.New("token not found")
		}
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user UUID found in redis")
	}

	return userID, nil
}

func InvalidateRefreshToken(ctx context.Context, rdb *redis.Client, token string) error {
	key := RefreshTokenKey(token)
	return rdb.Del(ctx, key).Err()
}
