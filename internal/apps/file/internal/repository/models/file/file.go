package file

import "time"

const (
	Table = "files"
)

var Columns = []string{
	"id",
	"user_id",
	"path",
	"size",
	"extension",
	"state",
	"created_at",
	"updated_at",
	"deleted_at",
}

type File struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Path      string     `db:"path"`
	Size      int64      `db:"size"`
	Extension Extension  `db:"extension"`
	State     State      `db:"state"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (f File) Values() []any {
	return []any{
		f.ID,
		f.UserID,
		f.Path,
		f.Size,
		f.Extension,
		f.State,
		f.CreatedAt,
		f.UpdatedAt,
		f.DeletedAt,
	}
}
