package application

import "time"

type CreateApplicationRequest struct {
	Name      string           `json:"name" binding:"required"`
	Tier      string           `json:"tier" binding:"required"`
	GitHub    *GitHubConfig    `json:"github"`
	Build     *BuildConfig     `json:"build"`
	Endpoints []EndpointConfig `json:"endpoints"`
}

type UpdateApplicationRequest struct {
	Name      string           `json:"name" binding:"required"`
	Tier      string           `json:"tier" binding:"required"`
	GitHub    *GitHubConfig    `json:"github"`
	Build     *BuildConfig     `json:"build"`
	Endpoints []EndpointConfig `json:"endpoints"`
}

type GitHubConfig struct {
	Owner          string   `json:"owner"`
	Repo           string   `json:"repo"`
	Branch         string   `json:"branch"`
	InstallationID string   `json:"installationId"`
	Hash           string   `json:"hash"`
	TriggerPaths   []string `json:"triggerPaths,omitempty"`
}

type BuildConfig struct {
	Gradle *GradleBuild `json:"gradle,omitempty"`
	NodeJS *NodeJSBuild `json:"nodejs,omitempty"`
	React  *ReactBuild  `json:"react,omitempty"`
	Vite   *ViteBuild   `json:"vite,omitempty"`
	Vue    *VueBuild    `json:"vue,omitempty"`
	NextJS *NextJSBuild `json:"nextjs,omitempty"`
	Go     *GoBuild     `json:"go,omitempty"`
	Rust   *RustBuild   `json:"rust,omitempty"`
	Maven  *MavenBuild  `json:"maven,omitempty"`
	Django *DjangoBuild `json:"django,omitempty"`
	Flask  *FlaskBuild  `json:"flask,omitempty"`
	Docker *DockerBuild `json:"docker,omitempty"`
}

type ApplicationResponse struct {
	ID        uint             `json:"id"`
	ProjectID uint             `json:"project_id"`
	Name      string           `json:"name"`
	Tier      string           `json:"tier"`
	GitHub    *GitHubConfig    `json:"github,omitempty"`
	Build     *BuildConfig     `json:"build,omitempty"`
	Endpoints []EndpointConfig `json:"endpoints,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}