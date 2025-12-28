package key

import (
	"context"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
)

type IRepository interface {
	GetCurrentSigningKey(ctx context.Context) (keyentity.Key, error)
	GetPublicKeys(ctx context.Context) ([]keyentity.Key, error)
	Rotate(ctx context.Context) (keyentity.Key, error)
}
