package project

import "time"

type CreateProjectRequest struct {
	Name       string `json:"name" binding:"required"`
	GitHubRepo string `json:"github_repo"`
}

type UpdateProjectRequest struct {
	Name       string `json:"name" binding:"required"`
	GitHubRepo string `json:"github_repo"`
}

type ProjectResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	OwnerID    uint      `json:"owner_id"`
	GitHubRepo string    `json:"github_repo"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProjectConfigResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	GitHubRepo   string        `json:"github_repo"`
	Applications []Application `json:"applications"`
	Addons       []Addon       `json:"addons"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type DeployApplicationRequest struct {
	Name      string           `json:"name" binding:"required"`
	Tier      string           `json:"tier" binding:"required"`
	GitHub    *GitHubConfig    `json:"github"`
	Build     *BuildConfig     `json:"build"`
	Endpoints []EndpointConfig `json:"endpoints"`
}

type DeployAddonRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Tier    string `json:"tier" binding:"required"`
	Storage string `json:"storage" binding:"required"`
}
