package project

import "context"

type Repository interface {
	Save(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id uint) (*Project, error)
	FindByOwnerID(ctx context.Context, ownerID uint) ([]*Project, error)
	FindByOwnerAndName(ctx context.Context, ownerID uint, name string) (*Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uint) error
}
