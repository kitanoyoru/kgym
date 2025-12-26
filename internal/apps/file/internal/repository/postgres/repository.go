package postgres

import (
	"context"

	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
)

type IRepository interface {
	Create(ctx context.Context, file filemodel.File) error
	List(ctx context.Context, filters ...Filter) ([]filemodel.File, error)
	Get(ctx context.Context, id string) (filemodel.File, error)
	Delete(ctx context.Context, filters ...Filter) error
	UpdateState(ctx context.Context, id string, state filemodel.State) error
	UpdateSize(ctx context.Context, id string, size int64) error
}
