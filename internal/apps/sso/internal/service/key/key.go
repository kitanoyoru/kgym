package key

import (
	"context"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
	keyrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key"
)

var _ IService = (*Service)(nil)

type Service struct {
	keyRepository keyrepo.IRepository
}

func NewService(keyRepository keyrepo.IRepository) *Service {
	return &Service{
		keyRepository: keyRepository,
	}
}

func (s *Service) GetPublicKeys(ctx context.Context) ([]keyentity.Key, error) {
	return s.keyRepository.GetPublicKeys(ctx)
}
