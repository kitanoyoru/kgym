package postgres

import (
	sq "github.com/Masterminds/squirrel"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
)

type Filter func(*Filters)

type Filters struct {
	ID       *string
	Email    *string
	Username *string
	Role     *usermodel.Role
}

func (f Filters) SQL() sq.Eq {
	eq := make(sq.Eq)

	if f.ID != nil {
		eq["id"] = *f.ID
	}
	if f.Email != nil {
		eq["email"] = *f.Email
	}
	if f.Username != nil {
		eq["username"] = *f.Username
	}
	if f.Role != nil {
		eq["role"] = string(*f.Role)
	}

	return eq
}

func WithID(id string) Filter {
	return func(f *Filters) {
		f.ID = &id
	}
}

func WithEmail(email string) Filter {
	return func(f *Filters) {
		f.Email = &email
	}
}

func WithUsername(username string) Filter {
	return func(f *Filters) {
		f.Username = &username
	}
}

func WithRole(role usermodel.Role) Filter {
	return func(f *Filters) {
		f.Role = &role
	}
}
