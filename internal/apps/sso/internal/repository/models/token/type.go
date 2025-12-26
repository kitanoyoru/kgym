package token

import "github.com/pkg/errors"

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

func (t TokenType) String() string {
	return string(t)
}
