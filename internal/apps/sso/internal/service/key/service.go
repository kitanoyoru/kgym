package key

import (
	"context"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
)

type IService interface {
	GetPublicKeys(ctx context.Context) ([]keyentity.Key, error)
}
