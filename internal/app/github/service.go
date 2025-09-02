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
		
		// "unknown"이거나 기본값인 경우 실제 GitHub owner name으로 업데이트 시도
		if accountLogin == "unknown" || accountLogin == "GitHub App Installation" || 
		   accountLogin == "GitHub Installation "+installation.InstallationID {
			// 실제 GitHub owner name 추정
			if realLogin, err := s.guessAccountLoginFromRepos(ctx); err == nil && realLogin != "" {
				accountLogin = realLogin
				// DB에도 업데이트 (비동기로)
				go s.updateInstallationLogin(ctx, installation.InstallationID, realLogin)
			} else {
				accountLogin = "installation-" + installation.InstallationID
			}
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
				// If we can't get real data, try to guess from user's repos
				if realLogin, guessErr := s.guessAccountLoginFromRepos(ctx); guessErr == nil && realLogin != "" {
					installationData = &Installation{
						InstallationID: installationID,
						AccountLogin:   realLogin,
						AccountType:    "User",
						Permissions:    "{}",
					}
				} else {
					installationData = &Installation{
						InstallationID: installationID,
						AccountLogin:   "installation-" + installationID,
						AccountType:    "User",
						Permissions:    "{}",
					}
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
	// GitHub API로부터 실제 계정 정보를 가져오려 시도
	// 현재는 사용자의 첫 번째 repo owner name을 사용해서 추정
	accountLogin, err := s.guessAccountLoginFromRepos(ctx)
	if err != nil || accountLogin == "" {
		accountLogin = "installation-" + installationID
	}
	
	return &Installation{
		InstallationID: installationID,
		AccountLogin:   accountLogin,
		AccountType:    "User",
		Permissions:    "{}",
	}, nil
}

// guessAccountLoginFromRepos tries to guess the account login from user's repositories
func (s *Service) guessAccountLoginFromRepos(ctx context.Context) (string, error) {
	// 사용자의 repository들을 가져와서 가장 많이 나오는 owner name 추정
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	
	repos, _, err := s.client.Repositories.List(ctx, "", opts)
	if err != nil {
		return "", err
	}
	
	// owner 이름들을 count해서 가장 많이 나오는 것 선택
	ownerCount := make(map[string]int)
	for _, repo := range repos {
		if repo.GetOwner() != nil {
			ownerName := repo.GetOwner().GetLogin()
			ownerCount[ownerName]++
		}
	}
	
	// 가장 많이 나오는 owner name 반환
	maxCount := 0
	mostFrequentOwner := ""
	for owner, count := range ownerCount {
		if count > maxCount {
			maxCount = count
			mostFrequentOwner = owner
		}
	}
	
	return mostFrequentOwner, nil
}

// updateInstallationLogin updates installation account login in database
func (s *Service) updateInstallationLogin(ctx context.Context, installationID, accountLogin string) {
	// installation을 찾아서 account login 업데이트
	installation, err := s.repo.FindByInstallationID(ctx, installationID)
	if err != nil {
		return // 에러 무시 (비동기 업데이트)
	}
	
	installation.AccountLogin = accountLogin
	s.repo.SaveInstallation(ctx, installation) // 에러 무시 (비동기)
}

