package application

import (
	"context"
	"encoding/json"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateApplication(ctx context.Context, projectID uint, req CreateApplicationRequest) (*ApplicationResponse, error) {
	// Convert request to application model
	app := &Application{
		ProjectID: projectID,
		Name:      req.Name,
		Tier:      req.Tier,
		Endpoints: req.Endpoints,
	}

	// Set GitHub configuration if provided
	if req.GitHub != nil {
		app.GitHubOwner = req.GitHub.Owner
		app.GitHubRepo = req.GitHub.Repo
		app.GitHubBranch = req.GitHub.Branch
		app.GitHubInstallationID = req.GitHub.InstallationID
		app.GitHubHash = req.GitHub.Hash
		app.GitHubTriggerPaths = req.GitHub.TriggerPaths
	}

	// Set build configuration if provided
	if req.Build != nil {
		app.BuildConfig = req.Build
		app.BuildType = s.determineBuildType(req.Build)
	}

	if err := s.repo.Save(ctx, app); err != nil {
		return nil, err
	}

	return s.toResponse(app), nil
}

func (s *Service) GetApplication(ctx context.Context, id uint) (*ApplicationResponse, error) {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(app), nil
}

func (s *Service) GetApplicationsByProject(ctx context.Context, projectID uint) ([]*ApplicationResponse, error) {
	apps, err := s.repo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	responses := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = s.toResponse(app)
	}

	return responses, nil
}

func (s *Service) UpdateApplication(ctx context.Context, id uint, req UpdateApplicationRequest) (*ApplicationResponse, error) {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	app.Name = req.Name
	app.Tier = req.Tier
	app.Endpoints = req.Endpoints

	// Update GitHub configuration
	if req.GitHub != nil {
		app.GitHubOwner = req.GitHub.Owner
		app.GitHubRepo = req.GitHub.Repo
		app.GitHubBranch = req.GitHub.Branch
		app.GitHubInstallationID = req.GitHub.InstallationID
		app.GitHubHash = req.GitHub.Hash
		app.GitHubTriggerPaths = req.GitHub.TriggerPaths
	} else {
		app.GitHubOwner = ""
		app.GitHubRepo = ""
		app.GitHubBranch = ""
		app.GitHubInstallationID = ""
		app.GitHubHash = ""
		app.GitHubTriggerPaths = nil
	}

	// Update build configuration
	if req.Build != nil {
		app.BuildConfig = req.Build
		app.BuildType = s.determineBuildType(req.Build)
	} else {
		app.BuildType = ""
		app.BuildConfig = nil
	}

	if err := s.repo.Save(ctx, app); err != nil {
		return nil, err
	}

	return s.toResponse(app), nil
}

func (s *Service) DeleteApplication(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) toResponse(app *Application) *ApplicationResponse {
	response := &ApplicationResponse{
		ID:        app.ID,
		ProjectID: app.ProjectID,
		Name:      app.Name,
		Tier:      app.Tier,
		Endpoints: app.Endpoints,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}

	// Add GitHub config if present
	if app.GitHubOwner != "" || app.GitHubRepo != "" {
		response.GitHub = &GitHubConfig{
			Owner:          app.GitHubOwner,
			Repo:           app.GitHubRepo,
			Branch:         app.GitHubBranch,
			InstallationID: app.GitHubInstallationID,
			Hash:           app.GitHubHash,
			TriggerPaths:   app.GitHubTriggerPaths,
		}
	}

	// Add build config if present
	if app.BuildConfig != nil && app.BuildType != "" {
		buildConfig := &BuildConfig{}
		
		// Parse build config based on type
		configData, _ := json.Marshal(app.BuildConfig)
		switch app.BuildType {
		case "gradle":
			var gradle GradleBuild
			json.Unmarshal(configData, &gradle)
			buildConfig.Gradle = &gradle
		case "nodejs":
			var nodejs NodeJSBuild
			json.Unmarshal(configData, &nodejs)
			buildConfig.NodeJS = &nodejs
		case "react":
			var react ReactBuild
			json.Unmarshal(configData, &react)
			buildConfig.React = &react
		case "vite":
			var vite ViteBuild
			json.Unmarshal(configData, &vite)
			buildConfig.Vite = &vite
		case "vue":
			var vue VueBuild
			json.Unmarshal(configData, &vue)
			buildConfig.Vue = &vue
		case "nextjs":
			var nextjs NextJSBuild
			json.Unmarshal(configData, &nextjs)
			buildConfig.NextJS = &nextjs
		case "go":
			var goBuild GoBuild
			json.Unmarshal(configData, &goBuild)
			buildConfig.Go = &goBuild
		case "rust":
			var rust RustBuild
			json.Unmarshal(configData, &rust)
			buildConfig.Rust = &rust
		case "maven":
			var maven MavenBuild
			json.Unmarshal(configData, &maven)
			buildConfig.Maven = &maven
		case "django":
			var django DjangoBuild
			json.Unmarshal(configData, &django)
			buildConfig.Django = &django
		case "flask":
			var flask FlaskBuild
			json.Unmarshal(configData, &flask)
			buildConfig.Flask = &flask
		case "docker":
			var docker DockerBuild
			json.Unmarshal(configData, &docker)
			buildConfig.Docker = &docker
		}

		response.Build = buildConfig
	}

	return response
}

func (s *Service) determineBuildType(build *BuildConfig) string {
	buildTypeMap := map[string]interface{}{
		"gradle": build.Gradle,
		"nodejs": build.NodeJS,
		"react":  build.React,
		"vite":   build.Vite,
		"vue":    build.Vue,
		"nextjs": build.NextJS,
		"go":     build.Go,
		"rust":   build.Rust,
		"maven":  build.Maven,
		"django": build.Django,
		"flask":  build.Flask,
		"docker": build.Docker,
	}

	for buildType, config := range buildTypeMap {
		if config != nil {
			return buildType
		}
	}
	return ""
}