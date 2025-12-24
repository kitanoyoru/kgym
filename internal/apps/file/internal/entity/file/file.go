package file

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/file/pkg/validator"
)

type File struct {
	ID        string    `validate:"required,uuid"`
	UserID    string    `validate:"required,uuid"`
	Path      string    `validate:"required,min=1,max=255"`
	Size      int64     `validate:"required,min=1"`
	Extension Extension `validate:"required"`
}

func (f File) Validate(ctx context.Context) error {
	return pkgValidator.Validate.StructCtx(ctx, f)
}
