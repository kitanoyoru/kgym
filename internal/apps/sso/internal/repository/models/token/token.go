package token

import "time"

const (
	Table = "tokens"
)

var Columns = []string{
	"id",
	"user_id",
	"token_type",
	"token",
	"created_at",
	"updated_at",
	"deleted_at",
}

type Token struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	TokenType TokenType  `db:"token_type"`
	Token     string     `db:"token"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (t Token) Values() []any {
	return []any{
		t.ID,
		t.UserID,
		t.TokenType,
		t.Token,
		t.CreatedAt,
		t.UpdatedAt,
		t.DeletedAt,
	}
}
