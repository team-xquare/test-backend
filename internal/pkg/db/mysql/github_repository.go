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
        INSERT INTO github_installations (installation_id, account_login, account_type, permissions)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        account_login = VALUES(account_login),
        account_type = VALUES(account_type),
        permissions = VALUES(permissions),
        updated_at = CURRENT_TIMESTAMP
    `

	result, err := r.db.ExecContext(ctx, query,
		installation.InstallationID,
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
	query := `
        SELECT id, installation_id, account_login, account_type, permissions, created_at, updated_at
        FROM github_installations WHERE installation_id = ?
    `

	var installation github.Installation
	err := r.db.QueryRowContext(ctx, query, installationID).Scan(
		&installation.ID,
		&installation.InstallationID,
		&installation.AccountLogin,
		&installation.AccountType,
		&installation.Permissions,
		&installation.CreatedAt,
		&installation.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("GitHub installation not found")
		}
		return nil, errors.Internal("Failed to get GitHub installation")
	}

	return &installation, nil
}

func (r *githubRepository) FindByUserID(ctx context.Context, userID uint) ([]*github.Installation, error) {
	query := `
        SELECT gi.id, gi.installation_id, gi.account_login, gi.account_type, gi.permissions, gi.created_at, gi.updated_at
        FROM github_installations gi
        INNER JOIN user_github_installations ugi ON gi.installation_id = ugi.installation_id
        WHERE ugi.user_id = ?
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
	// Start transaction to delete from both tables
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Internal("Failed to start transaction")
	}
	defer tx.Rollback()

	// Delete user links first
	_, err = tx.ExecContext(ctx, "DELETE FROM user_github_installations WHERE installation_id = ?", installationID)
	if err != nil {
		return errors.Internal("Failed to delete GitHub installation user links")
	}

	// Delete installation
	result, err := tx.ExecContext(ctx, "DELETE FROM github_installations WHERE installation_id = ?", installationID)
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

	if err = tx.Commit(); err != nil {
		return errors.Internal("Failed to commit transaction")
	}

	return nil
}

func (r *githubRepository) LinkUserToInstallation(ctx context.Context, userID uint, installationID string) error {
	query := `
        INSERT INTO user_github_installations (user_id, installation_id)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE created_at = created_at
    `

	_, err := r.db.ExecContext(ctx, query, userID, installationID)
	if err != nil {
		return errors.Internal("Failed to link user to GitHub installation")
	}

	return nil
}

func (r *githubRepository) IsUserLinkedToInstallation(ctx context.Context, userID uint, installationID string) (bool, error) {
	query := "SELECT 1 FROM user_github_installations WHERE user_id = ? AND installation_id = ?"

	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, installationID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Internal("Failed to check user installation link")
	}

	return true, nil
}
