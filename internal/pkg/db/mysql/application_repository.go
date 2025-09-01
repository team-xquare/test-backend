package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/team-xquare/deployment-platform/internal/app/application"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type applicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) application.Repository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Save(ctx context.Context, app *application.Application) error {
	triggerPathsJSON, _ := json.Marshal(app.GitHubTriggerPaths)
	buildConfigJSON, _ := json.Marshal(app.BuildConfig)
	endpointsJSON, _ := json.Marshal(app.Endpoints)

	if app.ID == 0 {
		// Insert new application
		query := `
			INSERT INTO applications (
				project_id, name, tier, 
				github_owner, github_repo, github_branch, github_installation_id, github_trigger_paths,
				build_type, build_config, endpoints
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		result, err := r.db.ExecContext(ctx, query,
			app.ProjectID, app.Name, app.Tier,
			app.GitHubOwner, app.GitHubRepo, app.GitHubBranch, app.GitHubInstallationID, string(triggerPathsJSON),
			app.BuildType, string(buildConfigJSON), string(endpointsJSON),
		)
		if err != nil {
			return errors.Internal("Failed to create application")
		}

		id, err := result.LastInsertId()
		if err != nil {
			return errors.Internal("Failed to get application ID")
		}
		app.ID = uint(id)
	} else {
		// Update existing application
		query := `
			UPDATE applications SET
				name = ?, tier = ?,
				github_owner = ?, github_repo = ?, github_branch = ?, github_installation_id = ?, github_trigger_paths = ?,
				build_type = ?, build_config = ?, endpoints = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`
		_, err := r.db.ExecContext(ctx, query,
			app.Name, app.Tier,
			app.GitHubOwner, app.GitHubRepo, app.GitHubBranch, app.GitHubInstallationID, string(triggerPathsJSON),
			app.BuildType, string(buildConfigJSON), string(endpointsJSON),
			app.ID,
		)
		if err != nil {
			return errors.Internal("Failed to update application")
		}
	}

	return nil
}

func (r *applicationRepository) FindByID(ctx context.Context, id uint) (*application.Application, error) {
	query := `
		SELECT id, project_id, name, tier,
			github_owner, github_repo, github_branch, github_installation_id, github_trigger_paths,
			build_type, build_config, endpoints, created_at, updated_at
		FROM applications WHERE id = ?
	`

	var app application.Application
	var triggerPathsJSON, buildConfigJSON, endpointsJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&app.ID, &app.ProjectID, &app.Name, &app.Tier,
		&app.GitHubOwner, &app.GitHubRepo, &app.GitHubBranch, &app.GitHubInstallationID, &triggerPathsJSON,
		&app.BuildType, &buildConfigJSON, &endpointsJSON, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("Application not found")
		}
		return nil, errors.Internal("Failed to get application")
	}

	// Unmarshal JSON fields
	json.Unmarshal([]byte(triggerPathsJSON), &app.GitHubTriggerPaths)
	json.Unmarshal([]byte(buildConfigJSON), &app.BuildConfig)
	json.Unmarshal([]byte(endpointsJSON), &app.Endpoints)

	return &app, nil
}

func (r *applicationRepository) FindByProjectID(ctx context.Context, projectID uint) ([]*application.Application, error) {
	query := `
		SELECT id, project_id, name, tier,
			github_owner, github_repo, github_branch, github_installation_id, github_trigger_paths,
			build_type, build_config, endpoints, created_at, updated_at
		FROM applications WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, errors.Internal("Failed to get applications")
	}
	defer rows.Close()

	var applications []*application.Application
	for rows.Next() {
		var app application.Application
		var triggerPathsJSON, buildConfigJSON, endpointsJSON string

		err := rows.Scan(
			&app.ID, &app.ProjectID, &app.Name, &app.Tier,
			&app.GitHubOwner, &app.GitHubRepo, &app.GitHubBranch, &app.GitHubInstallationID, &triggerPathsJSON,
			&app.BuildType, &buildConfigJSON, &endpointsJSON, &app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Internal("Failed to scan application")
		}

		// Unmarshal JSON fields
		json.Unmarshal([]byte(triggerPathsJSON), &app.GitHubTriggerPaths)
		json.Unmarshal([]byte(buildConfigJSON), &app.BuildConfig)
		json.Unmarshal([]byte(endpointsJSON), &app.Endpoints)

		applications = append(applications, &app)
	}

	return applications, nil
}

func (r *applicationRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM applications WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal("Failed to delete application")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("Application not found")
	}

	return nil
}