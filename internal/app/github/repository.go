package github

import "context"

type Repository interface {
	SaveInstallation(ctx context.Context, installation *Installation) error
	FindByInstallationID(ctx context.Context, installationID string) (*Installation, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Installation, error)
	DeleteByInstallationID(ctx context.Context, installationID string) error
}
