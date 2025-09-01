package addon

import (
	"context"

	"github.com/team-xquare/deployment-platform/internal/app/github"
)

type Service struct {
	repo      Repository
	githubSvc *github.Service
}

func NewService(repo Repository, githubSvc *github.Service) *Service {
	return &Service{
		repo:      repo,
		githubSvc: githubSvc,
	}
}

func (s *Service) CreateAddon(ctx context.Context, projectID uint, req CreateAddonRequest) (*AddonResponse, error) {
	addon := &Addon{
		ProjectID: projectID,
		Name:      req.Name,
		Type:      req.Type,
		Tier:      req.Tier,
		Storage:   req.Storage,
	}

	if err := s.repo.Save(ctx, addon); err != nil {
		return nil, err
	}

	// Trigger GitHub Actions workflow for addon deployment
	go s.triggerAddonDeployment(addon, "apply")

	return s.toResponse(addon), nil
}

func (s *Service) GetAddon(ctx context.Context, id uint) (*AddonResponse, error) {
	addon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(addon), nil
}

func (s *Service) GetAddonsByProject(ctx context.Context, projectID uint) ([]*AddonResponse, error) {
	addons, err := s.repo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	responses := make([]*AddonResponse, len(addons))
	for i, addon := range addons {
		responses[i] = s.toResponse(addon)
	}

	return responses, nil
}

func (s *Service) UpdateAddon(ctx context.Context, id uint, req UpdateAddonRequest) (*AddonResponse, error) {
	addon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	addon.Name = req.Name
	addon.Type = req.Type
	addon.Tier = req.Tier
	addon.Storage = req.Storage

	if err := s.repo.Save(ctx, addon); err != nil {
		return nil, err
	}

	return s.toResponse(addon), nil
}

func (s *Service) DeleteAddon(ctx context.Context, id uint) error {
	addon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Trigger GitHub Actions workflow for addon removal
	go s.triggerAddonDeployment(addon, "remove")

	return s.repo.Delete(ctx, id)
}

func (s *Service) toResponse(addon *Addon) *AddonResponse {
	return &AddonResponse{
		ID:        addon.ID,
		ProjectID: addon.ProjectID,
		Name:      addon.Name,
		Type:      addon.Type,
		Tier:      addon.Tier,
		Storage:   addon.Storage,
		CreatedAt: addon.CreatedAt,
		UpdatedAt: addon.UpdatedAt,
	}
}

func (s *Service) triggerAddonDeployment(addon *Addon, action string) {
	if s.githubSvc == nil {
		return
	}

	// Use a default GitHub repo for addons (this would be configurable)
	owner := "team-xquare"  // This should come from config
	repo := "infrastructure-configs"  // This should come from config
	
	projectName := "project-" + string(rune(addon.ProjectID))
	path := "projects/" + projectName + "/addons/" + addon.Name
	
	spec := map[string]interface{}{
		"type":    addon.Type,
		"tier":    addon.Tier,
		"storage": addon.Storage,
	}
	
	payload := github.ConfigAPIPayload{
		Path:   path,
		Action: action,
		Spec:   spec,
	}
	
	ctx := context.Background()
	s.githubSvc.TriggerGitHubAction(ctx, owner, repo, payload)
}