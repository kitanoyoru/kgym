package token

import (
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	"github.com/pkg/errors"
)

type Type string

const (
	TypeRefresh Type = "refresh"
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case "refresh":
		return TypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func TypeFromEntity(entity tokenentity.Type) (Type, error) {
	switch entity {
	case tokenentity.TypeRefresh:
		return TypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func (t Type) String() string {
	return string(t)
}
