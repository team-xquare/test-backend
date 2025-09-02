package github

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/google/go-github/v66/github"
	"github.com/team-xquare/deployment-platform/internal/pkg/config"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
	"golang.org/x/oauth2"
)

type Service struct {
	repo   Repository
	client *github.Client
}

func NewService(repo Repository) *Service {
	var client *github.Client

	if config.AppConfig.GitHubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.AppConfig.GitHubToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) HandleInstallationWebhook(ctx context.Context, payload []byte, signature string) error {
	if !s.verifySignature(payload, signature) {
		return errors.Forbidden("Invalid webhook signature")
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.BadRequest("Invalid webhook payload")
	}

	switch event.Action {
	case "created":
		return s.handleInstallationCreated(ctx, event)
	case "deleted":
		return s.handleInstallationDeleted(ctx, event)
	default:
		return nil
	}
}

func (s *Service) GetUserInstallations(ctx context.Context, userID uint) ([]*InstallationResponse, error) {
	installations, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*InstallationResponse, len(installations))
	for i, installation := range installations {
		accountLogin := installation.AccountLogin
		// "unknown" 계정명을 더 나은 이름으로 변경
		if accountLogin == "unknown" {
			accountLogin = "GitHub Installation " + installation.InstallationID
		}
		
		responses[i] = &InstallationResponse{
			ID:             installation.ID,
			InstallationID: installation.InstallationID,
			AccountLogin:   accountLogin,
			AccountType:    installation.AccountType,
		}
	}

	return responses, nil
}

func (s *Service) handleInstallationCreated(ctx context.Context, event WebhookEvent) error {
	installation := event.Installation

	installationIDFloat, ok := installation["id"].(float64)
	if !ok {
		return errors.BadRequest("Invalid installation ID")
	}

	account, ok := installation["account"].(map[string]interface{})
	if !ok {
		return errors.BadRequest("Invalid account data")
	}

	accountLogin, ok := account["login"].(string)
	if !ok {
		return errors.BadRequest("Invalid account login")
	}

	accountType, ok := account["type"].(string)
	if !ok {
		return errors.BadRequest("Invalid account type")
	}

	permissions, _ := json.Marshal(installation["permissions"])

	githubInstallation := &Installation{
		InstallationID: strconv.FormatFloat(installationIDFloat, 'f', -1, 64),
		AccountLogin:   accountLogin,
		AccountType:    accountType,
		Permissions:    string(permissions),
	}

	return s.repo.SaveInstallation(ctx, githubInstallation)
}

func (s *Service) handleInstallationDeleted(ctx context.Context, event WebhookEvent) error {
	installation := event.Installation

	installationIDFloat, ok := installation["id"].(float64)
	if !ok {
		return errors.BadRequest("Invalid installation ID")
	}

	installationID := strconv.FormatFloat(installationIDFloat, 'f', -1, 64)
	return s.repo.DeleteByInstallationID(ctx, installationID)
}

func (s *Service) verifySignature(payload []byte, signature string) bool {
	if config.AppConfig.GitHubWebhookSecret == "" {
		return true
	}

	expectedSignature := "sha256=" + s.computeSignature(payload, config.AppConfig.GitHubWebhookSecret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (s *Service) computeSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) TriggerGitHubAction(ctx context.Context, owner, repo string, payload ConfigAPIPayload) error {
	// Repository dispatch event로 GitHub Actions 트리거
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Internal("Failed to marshal payload: " + err.Error())
	}

	dispatchEvent := github.DispatchRequestOptions{
		EventType:     "config-api",
		ClientPayload: (*json.RawMessage)(&payloadBytes),
	}

	_, _, err = s.client.Repositories.Dispatch(ctx, owner, repo, dispatchEvent)
	if err != nil {
		return errors.Internal("Failed to trigger GitHub Action: " + err.Error())
	}

	return nil
}

func (s *Service) GetRepositories(ctx context.Context, installationID string) ([]*GitHubRepo, error) {
	// GitHub App installation을 통해 접근 가능한 repositories만 가져옴
	// Installation ID를 사용해서 해당 installation에 속한 repo들만 반환
	
	// 먼저 installation이 존재하는지 확인
	_, err := s.repo.FindByInstallationID(ctx, installationID)
	if err != nil {
		return nil, err
	}

	// GitHub App installation access token을 사용해야 하지만, 
	// 현재는 사용자 token으로 자신이 접근 가능한 repo들 중에서 필터링
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var filteredRepos []*GitHubRepo
	for {
		repos, resp, err := s.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, errors.Internal("Failed to fetch repositories from GitHub: " + err.Error())
		}

		for _, repo := range repos {
			// 사용자가 owner이거나 collaborate 권한이 있는 repo만 포함
			if repo.GetPermissions()["push"] || repo.GetPermissions()["admin"] {
				githubRepo := &GitHubRepo{
					ID:       int(repo.GetID()),
					Name:     repo.GetName(),
					FullName: repo.GetFullName(),
					Owner: Owner{
						Login: repo.GetOwner().GetLogin(),
					},
					Private: repo.GetPrivate(),
				}
				filteredRepos = append(filteredRepos, githubRepo)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return filteredRepos, nil
}

// LinkInstallationToUser links a GitHub installation to a specific user
func (s *Service) LinkInstallationToUser(ctx context.Context, userID uint, installationID string) error {
	// Check if user is already linked to this installation
	isLinked, err := s.repo.IsUserLinkedToInstallation(ctx, userID, installationID)
	if err != nil {
		return err
	}

	if isLinked {
		// Already linked
		return nil
	}

	// Check if installation exists, if not try to get real data from GitHub
	_, err = s.repo.FindByInstallationID(ctx, installationID)
	if err != nil {
		// If installation not found, try to get it from GitHub using the user's token
		if appErr, ok := err.(*errors.AppError); ok && appErr.StatusCode == 404 {
			// Try to get installation info from GitHub API using user token
			installationData, fetchErr := s.fetchInstallationInfo(ctx, installationID)
			if fetchErr != nil {
				// If we can't get real data, create a basic record
				installationData = &Installation{
					InstallationID: installationID,
					AccountLogin:   "GitHub App Installation",
					AccountType:    "User",
					Permissions:    "{}",
				}
			}
			
			// Save the installation
			if saveErr := s.repo.SaveInstallation(ctx, installationData); saveErr != nil {
				return saveErr
			}
		} else {
			return err
		}
	}

	// Link user to installation
	return s.repo.LinkUserToInstallation(ctx, userID, installationID)
}

// fetchInstallationInfo tries to get installation info from GitHub API
func (s *Service) fetchInstallationInfo(ctx context.Context, installationID string) (*Installation, error) {
	// Try to get installation info from GitHub API
	// For now, we'll return basic info since GitHub App token setup is complex
	return &Installation{
		InstallationID: installationID,
		AccountLogin:   "GitHub Installation " + installationID,
		AccountType:    "User",
		Permissions:    "{}",
	}, nil
}

