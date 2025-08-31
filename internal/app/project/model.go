package project

import "time"

type Project struct {
	ID         uint      `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	OwnerID    uint      `json:"owner_id" db:"owner_id"`
	GitHubRepo string    `json:"github_repo" db:"github_repo"`
	ConfigYAML string    `json:"config_yaml" db:"config_yaml"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Application struct {
	Name      string           `yaml:"name" json:"name"`
	Tier      string           `yaml:"tier" json:"tier"`
	GitHub    *GitHubConfig    `yaml:"github,omitempty" json:"github,omitempty"`
	Build     *BuildConfig     `yaml:"build,omitempty" json:"build,omitempty"`
	Endpoints []EndpointConfig `yaml:"endpoints,omitempty" json:"endpoints,omitempty"`
}

type GitHubConfig struct {
	Owner          string   `yaml:"owner" json:"owner"`
	Repo           string   `yaml:"repo" json:"repo"`
	Branch         string   `yaml:"branch" json:"branch"`
	InstallationID string   `yaml:"installationId" json:"installationId"`
	Hash           string   `yaml:"hash" json:"hash"`
	TriggerPaths   []string `yaml:"triggerPaths,omitempty" json:"triggerPaths,omitempty"`
}

type BuildConfig struct {
	Gradle *GradleBuild `yaml:"gradle,omitempty" json:"gradle,omitempty"`
	NodeJS *NodeJSBuild `yaml:"nodejs,omitempty" json:"nodejs,omitempty"`
	React  *ReactBuild  `yaml:"react,omitempty" json:"react,omitempty"`
	Vite   *ViteBuild   `yaml:"vite,omitempty" json:"vite,omitempty"`
	Vue    *VueBuild    `yaml:"vue,omitempty" json:"vue,omitempty"`
	NextJS *NextJSBuild `yaml:"nextjs,omitempty" json:"nextjs,omitempty"`
	Go     *GoBuild     `yaml:"go,omitempty" json:"go,omitempty"`
	Rust   *RustBuild   `yaml:"rust,omitempty" json:"rust,omitempty"`
	Maven  *MavenBuild  `yaml:"maven,omitempty" json:"maven,omitempty"`
	Django *DjangoBuild `yaml:"django,omitempty" json:"django,omitempty"`
	Flask  *FlaskBuild  `yaml:"flask,omitempty" json:"flask,omitempty"`
	Docker *DockerBuild `yaml:"docker,omitempty" json:"docker,omitempty"`
}

type GradleBuild struct {
	JavaVersion   string `yaml:"javaVersion" json:"javaVersion"`
	JarOutputPath string `yaml:"jarOutputPath" json:"jarOutputPath"`
	BuildCommand  string `yaml:"buildCommand" json:"buildCommand"`
}

type NodeJSBuild struct {
	NodeVersion  string `yaml:"nodeVersion" json:"nodeVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	StartCommand string `yaml:"startCommand" json:"startCommand"`
}

type ReactBuild struct {
	NodeVersion  string `yaml:"nodeVersion" json:"nodeVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	DistPath     string `yaml:"distPath" json:"distPath"`
}

type ViteBuild struct {
	NodeVersion  string `yaml:"nodeVersion" json:"nodeVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	DistPath     string `yaml:"distPath" json:"distPath"`
}

type VueBuild struct {
	NodeVersion  string `yaml:"nodeVersion" json:"nodeVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	DistPath     string `yaml:"distPath" json:"distPath"`
}

type NextJSBuild struct {
	NodeVersion  string `yaml:"nodeVersion" json:"nodeVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	StartCommand string `yaml:"startCommand" json:"startCommand"`
}

type GoBuild struct {
	GoVersion    string `yaml:"goVersion" json:"goVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	BinaryName   string `yaml:"binaryName" json:"binaryName"`
}

type RustBuild struct {
	RustVersion  string `yaml:"rustVersion" json:"rustVersion"`
	BuildCommand string `yaml:"buildCommand" json:"buildCommand"`
	BinaryName   string `yaml:"binaryName" json:"binaryName"`
}

type MavenBuild struct {
	JavaVersion   string `yaml:"javaVersion" json:"javaVersion"`
	BuildCommand  string `yaml:"buildCommand" json:"buildCommand"`
	JarOutputPath string `yaml:"jarOutputPath" json:"jarOutputPath"`
}

type DjangoBuild struct {
	PythonVersion string `yaml:"pythonVersion" json:"pythonVersion"`
	BuildCommand  string `yaml:"buildCommand" json:"buildCommand"`
	StartCommand  string `yaml:"startCommand" json:"startCommand"`
}

type FlaskBuild struct {
	PythonVersion string `yaml:"pythonVersion" json:"pythonVersion"`
	BuildCommand  string `yaml:"buildCommand" json:"buildCommand"`
	StartCommand  string `yaml:"startCommand" json:"startCommand"`
}

type DockerBuild struct {
	DockerfilePath string `yaml:"dockerfilePath" json:"dockerfilePath"`
	ContextPath    string `yaml:"contextPath" json:"contextPath"`
}

type EndpointConfig struct {
	Port   int      `yaml:"port" json:"port"`
	Routes []string `yaml:"routes" json:"routes"`
}

type Addon struct {
	Name    string `yaml:"name" json:"name"`
	Type    string `yaml:"type" json:"type"`
	Tier    string `yaml:"tier" json:"tier"`
	Storage string `yaml:"storage" json:"storage"`
}

type ProjectConfig struct {
	Applications []Application `yaml:"applications" json:"applications"`
	Addons       []Addon       `yaml:"addons" json:"addons"`
}
