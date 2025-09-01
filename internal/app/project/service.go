package project

import (
	"context"

	"github.com/team-xquare/deployment-platform/internal/app/github"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type Service struct {
	repo       Repository
	githubRepo github.Repository
}

func NewService(repo Repository, githubRepo github.Repository) *Service {
	return &Service{repo: repo, githubRepo: githubRepo}
}

func (s *Service) CreateProject(ctx context.Context, userID uint, req CreateProjectRequest) (*ProjectResponse, error) {
	existing, err := s.repo.FindByOwnerAndName(ctx, userID, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.BadRequest("Project with this name already exists")
	}

	project := &Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
	}

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *Service) GetProject(ctx context.Context, userID, projectID uint) (*ProjectResponse, error) {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != userID {
		return nil, errors.Forbidden("Access denied")
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *Service) GetUserProjects(ctx context.Context, userID uint) ([]*ProjectResponse, error) {
	projects, err := s.repo.FindByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = &ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			OwnerID:     project.OwnerID,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *Service) UpdateProject(ctx context.Context, userID, projectID uint, req UpdateProjectRequest) (*ProjectResponse, error) {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != userID {
		return nil, errors.Forbidden("Access denied")
	}

	project.Name = req.Name
	project.Description = req.Description

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

func (s *Service) DeleteProject(ctx context.Context, userID, projectID uint) error {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return errors.Forbidden("Access denied")
	}

	return s.repo.Delete(ctx, projectID)
}