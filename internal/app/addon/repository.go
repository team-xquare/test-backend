package addon

import "context"

type Repository interface {
	Save(ctx context.Context, addon *Addon) error
	FindByID(ctx context.Context, id uint) (*Addon, error)
	FindByProjectID(ctx context.Context, projectID uint) ([]*Addon, error)
	Delete(ctx context.Context, id uint) error
}