package application

import "context"

type Repository interface {
	Save(ctx context.Context, app *Application) error
	FindByID(ctx context.Context, id uint) (*Application, error)
	FindByProjectID(ctx context.Context, projectID uint) ([]*Application, error)
	Delete(ctx context.Context, id uint) error
}