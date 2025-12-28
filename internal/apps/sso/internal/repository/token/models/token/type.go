package token

import (
	"github.com/pkg/errors"

	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
)

type TokenType string

const (
	TokenTypeRefresh TokenType = "refresh"
)

func TokenTypeFromString(s string) (TokenType, error) {
	switch s {
	case "refresh":
		return TokenTypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func TokenTypeFromEntity(entity tokenentity.TokenType) (TokenType, error) {
	switch entity {
	case tokenentity.TokenTypeRefresh:
		return TokenTypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func (t TokenType) String() string {
	return string(t)
}
