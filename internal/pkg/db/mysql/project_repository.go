package mysql

import (
	"context"
	"database/sql"

	"github.com/team-xquare/deployment-platform/internal/app/project"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) project.Repository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Save(ctx context.Context, proj *project.Project) error {
	query := `
        INSERT INTO projects (name, owner_id, github_repo, config_yaml)
        VALUES (?, ?, ?, ?)
    `

	result, err := r.db.ExecContext(ctx, query, proj.Name, proj.OwnerID, proj.GitHubRepo, proj.ConfigYAML)
	if err != nil {
		return errors.Internal("Failed to create project")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get project ID")
	}

	proj.ID = uint(id)
	return nil
}

func (r *projectRepository) FindByID(ctx context.Context, id uint) (*project.Project, error) {
	var p project.Project
	query := `
        SELECT id, name, owner_id, github_repo, config_yaml, created_at, updated_at
        FROM projects WHERE id = ?
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.OwnerID,
		&p.GitHubRepo,
		&p.ConfigYAML,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get project")
	}

	return &p, nil
}

func (r *projectRepository) FindByOwnerID(ctx context.Context, ownerID uint) ([]*project.Project, error) {
	query := `
        SELECT id, name, owner_id, github_repo, config_yaml, created_at, updated_at
        FROM projects WHERE owner_id = ?
    `

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, errors.Internal("Failed to get projects")
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		var p project.Project
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.OwnerID,
			&p.GitHubRepo,
			&p.ConfigYAML,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Internal("Failed to scan project")
		}
		projects = append(projects, &p)
	}

	return projects, nil
}

func (r *projectRepository) FindByOwnerAndName(ctx context.Context, ownerID uint, name string) (*project.Project, error) {
	var p project.Project
	query := `
        SELECT id, name, owner_id, github_repo, config_yaml, created_at, updated_at
        FROM projects WHERE owner_id = ? AND name = ?
    `

	err := r.db.QueryRowContext(ctx, query, ownerID, name).Scan(
		&p.ID,
		&p.Name,
		&p.OwnerID,
		&p.GitHubRepo,
		&p.ConfigYAML,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get project")
	}

	return &p, nil
}

func (r *projectRepository) Update(ctx context.Context, proj *project.Project) error {
	query := `
        UPDATE projects 
        SET name = ?, github_repo = ?, config_yaml = ?
        WHERE id = ?
    `

	result, err := r.db.ExecContext(ctx, query, proj.Name, proj.GitHubRepo, proj.ConfigYAML, proj.ID)
	if err != nil {
		return errors.Internal("Failed to update project")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("Project not found")
	}

	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM projects WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal("Failed to delete project")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("Project not found")
	}

	return nil
}