package service

import userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"

type ListOptions struct {
	Email    *string
	Username *string
	Role     *userentity.Role
}

type ListOption func(*ListOptions)

func WithEmail(email string) ListOption {
	return func(o *ListOptions) {
		o.Email = &email
	}
}

func WithUsername(username string) ListOption {
	return func(o *ListOptions) {
		o.Username = &username
	}
}

func WithRole(role userentity.Role) ListOption {
	return func(o *ListOptions) {
		o.Role = &role
	}
}
