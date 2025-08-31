package auth

import (
	"context"
	"time"

	"github.com/team-xquare/deployment-platform/internal/app/user"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/jwt"
)

type Service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(repo Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

func (s *Service) Login(ctx context.Context, req user.LoginRequest) (*LoginResponse, error) {
	userSvc := user.NewService(s.userRepo)
	authenticatedUser, err := userSvc.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(authenticatedUser.ID, authenticatedUser.Email)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.repo.SaveRefreshToken(ctx, authenticatedUser.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:    authenticatedUser.ID,
			Email: authenticatedUser.Email,
			Name:  authenticatedUser.Name,
		},
	}, nil
}

func (s *Service) Register(ctx context.Context, req user.RegisterRequest) error {
	userSvc := user.NewService(s.userRepo)
	return userSvc.Register(ctx, req)
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error) {
	userID, err := s.repo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("User not found")
	}

	accessToken, newRefreshToken, err := jwt.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	s.repo.DeleteRefreshToken(ctx, req.RefreshToken)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.repo.SaveRefreshToken(ctx, user.ID, newRefreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}
