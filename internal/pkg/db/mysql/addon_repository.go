package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/team-xquare/deployment-platform/internal/app/addon"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type addonRepository struct {
	db *sql.DB
}

func NewAddonRepository(db *sql.DB) addon.Repository {
	return &addonRepository{db: db}
}

func (r *addonRepository) Save(ctx context.Context, addon *addon.Addon) error {
	configJSON, _ := json.Marshal(addon.Config)

	if addon.ID == 0 {
		// Insert new addon
		query := `
			INSERT INTO addons (project_id, name, type, tier, storage, config)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		result, err := r.db.ExecContext(ctx, query,
			addon.ProjectID, addon.Name, addon.Type, addon.Tier, addon.Storage, string(configJSON),
		)
		if err != nil {
			return errors.Internal("Failed to create addon")
		}

		id, err := result.LastInsertId()
		if err != nil {
			return errors.Internal("Failed to get addon ID")
		}
		addon.ID = uint(id)
	} else {
		// Update existing addon
		query := `
			UPDATE addons SET
				name = ?, type = ?, tier = ?, storage = ?, config = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`
		_, err := r.db.ExecContext(ctx, query,
			addon.Name, addon.Type, addon.Tier, addon.Storage, string(configJSON), addon.ID,
		)
		if err != nil {
			return errors.Internal("Failed to update addon")
		}
	}

	return nil
}

func (r *addonRepository) FindByID(ctx context.Context, id uint) (*addon.Addon, error) {
	query := `
		SELECT id, project_id, name, type, tier, storage, config, created_at, updated_at
		FROM addons WHERE id = ?
	`

	var addon addon.Addon
	var configJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&addon.ID, &addon.ProjectID, &addon.Name, &addon.Type, &addon.Tier, &addon.Storage, &configJSON,
		&addon.CreatedAt, &addon.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("Addon not found")
		}
		return nil, errors.Internal("Failed to get addon")
	}

	// Unmarshal JSON config
	json.Unmarshal([]byte(configJSON), &addon.Config)

	return &addon, nil
}

func (r *addonRepository) FindByProjectID(ctx context.Context, projectID uint) ([]*addon.Addon, error) {
	query := `
		SELECT id, project_id, name, type, tier, storage, config, created_at, updated_at
		FROM addons WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, errors.Internal("Failed to get addons")
	}
	defer rows.Close()

	var addons []*addon.Addon
	for rows.Next() {
		var addon addon.Addon
		var configJSON string

		err := rows.Scan(
			&addon.ID, &addon.ProjectID, &addon.Name, &addon.Type, &addon.Tier, &addon.Storage, &configJSON,
			&addon.CreatedAt, &addon.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Internal("Failed to scan addon")
		}

		// Unmarshal JSON config
		json.Unmarshal([]byte(configJSON), &addon.Config)

		addons = append(addons, &addon)
	}

	return addons, nil
}

func (r *addonRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM addons WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal("Failed to delete addon")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("Addon not found")
	}

	return nil
}