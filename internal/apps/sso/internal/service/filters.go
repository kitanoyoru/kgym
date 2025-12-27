package service

import "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"

type Filter func(*Filters)

type Filters struct {
	ID        *string
	UserID    *string
	TokenType *token.TokenType
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

func WithTokenType(tokenType token.TokenType) Filter {
	return func(f *Filters) {
		f.TokenType = &tokenType
	}
}
