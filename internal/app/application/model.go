package application

import "time"

type Application struct {
	ID        uint      `json:"id" db:"id"`
	ProjectID uint      `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Tier      string    `json:"tier" db:"tier"`
	
	// GitHub Configuration
	GitHubOwner          string   `json:"github_owner" db:"github_owner"`
	GitHubRepo           string   `json:"github_repo" db:"github_repo"`
	GitHubBranch         string   `json:"github_branch" db:"github_branch"`
	GitHubInstallationID string   `json:"github_installation_id" db:"github_installation_id"`
	GitHubHash           string   `json:"github_hash" db:"github_hash"`
	GitHubTriggerPaths   []string `json:"github_trigger_paths" db:"github_trigger_paths"`
	
	// Build Configuration
	BuildType   string      `json:"build_type" db:"build_type"`
	BuildConfig interface{} `json:"build_config" db:"build_config"`
	
	// Endpoints
	Endpoints []EndpointConfig `json:"endpoints" db:"endpoints"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type EndpointConfig struct {
	Port   int      `json:"port"`
	Routes []string `json:"routes"`
}

// Build configurations for different types
type GradleBuild struct {
	JavaVersion   string `json:"javaVersion"`
	JarOutputPath string `json:"jarOutputPath"`
	BuildCommand  string `json:"buildCommand"`
}

type NodeJSBuild struct {
	NodeVersion  string `json:"nodeVersion"`
	BuildCommand string `json:"buildCommand"`
	StartCommand string `json:"startCommand"`
}

type ReactBuild struct {
	NodeVersion  string `json:"nodeVersion"`
	BuildCommand string `json:"buildCommand"`
	DistPath     string `json:"distPath"`
}

type ViteBuild struct {
	NodeVersion  string `json:"nodeVersion"`
	BuildCommand string `json:"buildCommand"`
	DistPath     string `json:"distPath"`
}

type VueBuild struct {
	NodeVersion  string `json:"nodeVersion"`
	BuildCommand string `json:"buildCommand"`
	DistPath     string `json:"distPath"`
}

type NextJSBuild struct {
	NodeVersion  string `json:"nodeVersion"`
	BuildCommand string `json:"buildCommand"`
	StartCommand string `json:"startCommand"`
}

type GoBuild struct {
	GoVersion    string `json:"goVersion"`
	BuildCommand string `json:"buildCommand"`
	BinaryName   string `json:"binaryName"`
}

type RustBuild struct {
	RustVersion  string `json:"rustVersion"`
	BuildCommand string `json:"buildCommand"`
	BinaryName   string `json:"binaryName"`
}

type MavenBuild struct {
	JavaVersion   string `json:"javaVersion"`
	BuildCommand  string `json:"buildCommand"`
	JarOutputPath string `json:"jarOutputPath"`
}

type DjangoBuild struct {
	PythonVersion string `json:"pythonVersion"`
	BuildCommand  string `json:"buildCommand"`
	StartCommand  string `json:"startCommand"`
}

type FlaskBuild struct {
	PythonVersion string `json:"pythonVersion"`
	BuildCommand  string `json:"buildCommand"`
	StartCommand  string `json:"startCommand"`
}

type DockerBuild struct {
	DockerfilePath string `json:"dockerfilePath"`
	ContextPath    string `json:"contextPath"`
}