package auth

import (
	"context"
	"time"
)

type Repository interface {
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (uint, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
