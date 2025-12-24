package postgres

import (
	sq "github.com/Masterminds/squirrel"
	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
)

type Filter func(*Filters)

type Filters struct {
	ID        *string
	UserID    *string
	Path      *string
	Size      *int64
	Extension *filemodel.Extension
}

func (f Filters) SQL() sq.Eq {
	eq := make(sq.Eq)

	if f.ID != nil {
		eq["id"] = *f.ID
	}
	if f.UserID != nil {
		eq["user_id"] = *f.UserID
	}
	if f.Path != nil {
		eq["path"] = *f.Path
	}
	if f.Size != nil {
		eq["size"] = *f.Size
	}
	if f.Extension != nil {
		eq["extension"] = string(*f.Extension)
	}

	return eq
}

func WithID(id string) Filter {
	return func(f *Filters) {
		f.ID = &id
	}
}

func WithUserID(userID string) Filter {
	return func(f *Filters) {
		f.UserID = &userID
	}
}

func WithPath(path string) Filter {
	return func(f *Filters) {
		f.Path = &path
	}
}

func WithSize(size int64) Filter {
	return func(f *Filters) {
		f.Size = &size
	}
}

func WithExtension(extension filemodel.Extension) Filter {
	return func(f *Filters) {
		f.Extension = &extension
	}
}
