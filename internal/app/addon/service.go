package addon

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAddon(ctx context.Context, projectID uint, req CreateAddonRequest) (*AddonResponse, error) {
	addon := &Addon{
		ProjectID: projectID,
		Name:      req.Name,
		Type:      req.Type,
		Tier:      req.Tier,
		Storage:   req.Storage,
		Config:    req.Config,
	}

	if err := s.repo.Save(ctx, addon); err != nil {
		return nil, err
	}

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
	addon.Config = req.Config

	if err := s.repo.Save(ctx, addon); err != nil {
		return nil, err
	}

	return s.toResponse(addon), nil
}

func (s *Service) DeleteAddon(ctx context.Context, id uint) error {
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
		Config:    addon.Config,
		CreatedAt: addon.CreatedAt,
		UpdatedAt: addon.UpdatedAt,
	}
}