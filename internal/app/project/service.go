package project

import (
	"context"
	"fmt"

	"github.com/team-xquare/deployment-platform/internal/app/github"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"gopkg.in/yaml.v3"
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

	initialConfig := ProjectConfig{
		Applications: []Application{},
		Addons:       []Addon{},
	}

	configYAML, err := yaml.Marshal(initialConfig)
	if err != nil {
		return nil, errors.Internal("Failed to marshal initial config")
	}

	project := &Project{
		Name:       req.Name,
		OwnerID:    userID,
		GitHubRepo: req.GitHubRepo,
		ConfigYAML: string(configYAML),
	}

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:         project.ID,
		Name:       project.Name,
		OwnerID:    project.OwnerID,
		GitHubRepo: project.GitHubRepo,
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
	}, nil
}

func (s *Service) GetProjects(ctx context.Context, userID uint) ([]*ProjectResponse, error) {
	projects, err := s.repo.FindByOwnerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = &ProjectResponse{
			ID:         project.ID,
			Name:       project.Name,
			OwnerID:    project.OwnerID,
			GitHubRepo: project.GitHubRepo,
			CreatedAt:  project.CreatedAt,
			UpdatedAt:  project.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetProject(ctx context.Context, userID, projectID uint) (*ProjectConfigResponse, error) {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.NotFound("Project not found")
	}
	if project.OwnerID != userID {
		return nil, errors.Forbidden("Access denied")
	}

	var config ProjectConfig
	if err := yaml.Unmarshal([]byte(project.ConfigYAML), &config); err != nil {
		return nil, errors.Internal("Failed to parse project config")
	}

	return &ProjectConfigResponse{
		ID:           project.ID,
		Name:         project.Name,
		GitHubRepo:   project.GitHubRepo,
		Applications: config.Applications,
		Addons:       config.Addons,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
	}, nil
}

func (s *Service) DeployApplication(ctx context.Context, userID, projectID uint, req DeployApplicationRequest) error {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.NotFound("Project not found")
	}
	if project.OwnerID != userID {
		return errors.Forbidden("Access denied")
	}

	var config ProjectConfig
	if err := yaml.Unmarshal([]byte(project.ConfigYAML), &config); err != nil {
		return errors.Internal("Failed to parse project config")
	}

	app := Application{
		Name:      req.Name,
		Tier:      req.Tier,
		GitHub:    req.GitHub,
		Build:     req.Build,
		Endpoints: req.Endpoints,
	}

	found := false
	for i, existingApp := range config.Applications {
		if existingApp.Name == req.Name {
			config.Applications[i] = app
			found = true
			break
		}
	}

	if !found {
		config.Applications = append(config.Applications, app)
	}

	configYAML, err := yaml.Marshal(config)
	if err != nil {
		return errors.Internal("Failed to marshal config")
	}

	project.ConfigYAML = string(configYAML)
	if err := s.repo.Update(ctx, project); err != nil {
		return err
	}

	payload := github.ConfigAPIPayload{
		Path:   fmt.Sprintf("projects/%s/applications/%s", project.Name, req.Name),
		Action: "apply",
		Spec: map[string]interface{}{
			"tier":      req.Tier,
			"github":    req.GitHub,
			"build":     req.Build,
			"endpoints": req.Endpoints,
		},
	}

	githubService := github.NewService(s.githubRepo)
	return githubService.TriggerGitHubAction(ctx, "team-xquare", "deployment-platform", payload)
}

func (s *Service) DeployAddon(ctx context.Context, userID, projectID uint, req DeployAddonRequest) error {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.NotFound("Project not found")
	}
	if project.OwnerID != userID {
		return errors.Forbidden("Access denied")
	}

	var config ProjectConfig
	if err := yaml.Unmarshal([]byte(project.ConfigYAML), &config); err != nil {
		return errors.Internal("Failed to parse project config")
	}

	addon := Addon{
		Name:    req.Name,
		Type:    req.Type,
		Tier:    req.Tier,
		Storage: req.Storage,
	}

	found := false
	for i, existingAddon := range config.Addons {
		if existingAddon.Name == req.Name {
			config.Addons[i] = addon
			found = true
			break
		}
	}

	if !found {
		config.Addons = append(config.Addons, addon)
	}

	configYAML, err := yaml.Marshal(config)
	if err != nil {
		return errors.Internal("Failed to marshal config")
	}

	project.ConfigYAML = string(configYAML)
	if err := s.repo.Update(ctx, project); err != nil {
		return err
	}

	payload := github.ConfigAPIPayload{
		Path:   fmt.Sprintf("projects/%s/addons/%s", project.Name, req.Name),
		Action: "apply",
		Spec: map[string]interface{}{
			"type":    req.Type,
			"tier":    req.Tier,
			"storage": req.Storage,
		},
	}

	githubService := github.NewService(s.githubRepo)
	return githubService.TriggerGitHubAction(ctx, "team-xquare", "deployment-platform", payload)
}

func (s *Service) DeleteProject(ctx context.Context, userID, projectID uint) error {
	project, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.NotFound("Project not found")
	}
	if project.OwnerID != userID {
		return errors.Forbidden("Access denied")
	}

	payload := github.ConfigAPIPayload{
		Path:   fmt.Sprintf("projects/%s", project.Name),
		Action: "remove",
		Spec:   nil,
	}

	githubService := github.NewService(s.githubRepo)
	if err := githubService.TriggerGitHubAction(ctx, "team-xquare", "deployment-platform", payload); err != nil {
		return err
	}

	return s.repo.Delete(ctx, projectID)
}