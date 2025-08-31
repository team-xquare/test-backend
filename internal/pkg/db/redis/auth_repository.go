package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/team-xquare/deployment-platform/internal/app/auth"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"github.com/go-redis/redis/v8"
)

type authRepository struct {
	client *redis.Client
}

func NewAuthRepository(client *redis.Client) auth.Repository {
	return &authRepository{client: client}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	duration := time.Until(expiresAt)

	err := r.client.Set(ctx, key, userID, duration).Err()
	if err != nil {
		return errors.Internal("Failed to save refresh token")
	}

	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (uint, error) {
	key := fmt.Sprintf("refresh_token:%s", token)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, errors.Unauthorized("Invalid refresh token")
	}
	if err != nil {
		return 0, errors.Internal("Failed to get refresh token")
	}

	userID, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, errors.Internal("Failed to parse user ID")
	}

	return uint(userID), nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh_token:%s", token)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.Internal("Failed to delete refresh token")
	}

	return nil
}