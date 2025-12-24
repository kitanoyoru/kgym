package service

import userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"

type Options struct {
	ID       *string
	Email    *string
	Username *string
	Role     *userentity.Role
}

type Option func(*Options)

func WithID(id string) Option {
	return func(o *Options) {
		o.ID = &id
	}
}

func WithEmail(email string) Option {
	return func(o *Options) {
		o.Email = &email
	}
}

func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = &username
	}
}

func WithRole(role userentity.Role) Option {
	return func(o *Options) {
		o.Role = &role
	}
}
