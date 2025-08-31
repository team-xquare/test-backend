package mysql

import (
	"context"
	"database/sql"

	"github.com/team-xquare/deployment-platform/internal/app/github"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type githubRepository struct {
	db *sql.DB
}

func NewGitHubRepository(db *sql.DB) github.Repository {
	return &githubRepository{db: db}
}

func (r *githubRepository) SaveInstallation(ctx context.Context, installation *github.Installation) error {
	query := `
        INSERT INTO github_installations (installation_id, user_id, account_login, account_type, permissions)
        VALUES (?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        account_login = VALUES(account_login),
        account_type = VALUES(account_type),
        permissions = VALUES(permissions),
        updated_at = CURRENT_TIMESTAMP
    `

	result, err := r.db.ExecContext(ctx, query, 
		installation.InstallationID,
		installation.UserID,
		installation.AccountLogin,
		installation.AccountType,
		installation.Permissions,
	)
	if err != nil {
		return errors.Internal("Failed to save GitHub installation")
	}

	if installation.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return errors.Internal("Failed to get installation ID")
		}
		installation.ID = uint(id)
	}

	return nil
}

func (r *githubRepository) FindByInstallationID(ctx context.Context, installationID string) (*github.Installation, error) {
	var installation github.Installation
	query := `
        SELECT id, installation_id, user_id, account_login, account_type, permissions, created_at, updated_at
        FROM github_installations WHERE installation_id = ?
    `

	err := r.db.QueryRowContext(ctx, query, installationID).Scan(
		&installation.ID,
		&installation.InstallationID,
		&installation.UserID,
		&installation.AccountLogin,
		&installation.AccountType,
		&installation.Permissions,
		&installation.CreatedAt,
		&installation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get GitHub installation")
	}

	return &installation, nil
}

func (r *githubRepository) FindByUserID(ctx context.Context, userID uint) ([]*github.Installation, error) {
	query := `
        SELECT id, installation_id, user_id, account_login, account_type, permissions, created_at, updated_at
        FROM github_installations WHERE user_id = ?
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Internal("Failed to get GitHub installations")
	}
	defer rows.Close()

	var installations []*github.Installation
	for rows.Next() {
		var installation github.Installation
		err := rows.Scan(
			&installation.ID,
			&installation.InstallationID,
			&installation.UserID,
			&installation.AccountLogin,
			&installation.AccountType,
			&installation.Permissions,
			&installation.CreatedAt,
			&installation.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Internal("Failed to scan GitHub installation")
		}
		installations = append(installations, &installation)
	}

	return installations, nil
}

func (r *githubRepository) DeleteByInstallationID(ctx context.Context, installationID string) error {
	query := "DELETE FROM github_installations WHERE installation_id = ?"

	result, err := r.db.ExecContext(ctx, query, installationID)
	if err != nil {
		return errors.Internal("Failed to delete GitHub installation")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("GitHub installation not found")
	}

	return nil
}