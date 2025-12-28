package token

import (
	"time"

	"github.com/dromara/carbon/v2"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
)

const (
	Table = "tokens"
)

var Columns = []string{
	"id",
	"subject",
	"client_id",
	"token_type",
	"token_hash",
	"expires_at",
	"revoked",
	"created_at",
	"updated_at",
	"deleted_at",
}

func TokenFromEntity(entity tokenentity.Token) Token {
	tokenType, err := TokenTypeFromEntity(entity.TokenType)
	if err != nil {
		return Token{}
	}

	now := carbon.Now().StdTime()

	return Token{
		ID:        entity.ID,
		Subject:   entity.Subject,
		ClientID:  entity.ClientID,
		TokenType: tokenType,
		TokenHash: entity.TokenHash,
		ExpiresAt: entity.ExpiresAt,
		Revoked:   entity.Revoked,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type Token struct {
	ID        string     `db:"id"`
	Subject   string     `db:"subject"`
	ClientID  string     `db:"client_id"`
	TokenType TokenType  `db:"token_type"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	Revoked   bool       `db:"revoked"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (t Token) Values() []any {
	return []any{
		t.ID,
		t.Subject,
		t.ClientID,
		t.TokenType,
		t.TokenHash,
		t.ExpiresAt,
		t.Revoked,
		t.CreatedAt,
		t.UpdatedAt,
		t.DeletedAt,
	}
}
