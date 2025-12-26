package postgres

import (
	sq "github.com/Masterminds/squirrel"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
)

type Filter func(*Filters)

type Filters struct {
	ID        *string
	UserID    *string
	TokenType *tokenmodel.TokenType
}

func (f Filters) SQL() sq.Eq {
	eq := make(sq.Eq)

	if f.ID != nil {
		eq["id"] = *f.ID
	}
	if f.UserID != nil {
		eq["user_id"] = *f.UserID
	}
	if f.TokenType != nil {
		eq["token_type"] = (*f.TokenType).String()
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

func WithTokenType(tokenType tokenmodel.TokenType) Filter {
	return func(f *Filters) {
		f.TokenType = &tokenType
	}
}
