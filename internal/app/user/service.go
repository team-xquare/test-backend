package user

import (
	"context"

	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.BadRequest("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	user := &User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	return s.repo.Save(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found")
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		GitHubID:  user.GitHubID,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateUserRequest) error {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	user.Name = req.Name
	user.Password = string(hashedPassword)

	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) FindOrCreateByGitHub(ctx context.Context, githubID, email, name string) (*User, error) {
	user, err := s.repo.FindByGitHubID(ctx, githubID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		user.GitHubID = &githubID
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	user = &User{
		Email:    email,
		Name:     name,
		GitHubID: &githubID,
		Password: "github_oauth",
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
