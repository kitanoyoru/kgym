package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/postgres"
)

var _ IService = (*Service)(nil)

type Service struct {
	repo postgres.IRepository
}

func New(repo postgres.IRepository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) CreateToken(ctx context.Context, req CreateTokenRequest) (string, error) {
	tokenType, err := tokenentity.TokenTypeFromString(req.TokenType.String())
	if err != nil {
		return "", err
	}

	entity := tokenentity.Token{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		TokenType: tokenType,
		Token:     req.Token,
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return "", err
	}

	return entity.ID, nil
}

func (s *Service) ListTokens(ctx context.Context, filters ...Filter) ([]token.Token, error) {
	var serviceFilters Filters
	for _, filter := range filters {
		filter(&serviceFilters)
	}

	// TODO: Refactor somehow
	var postgresFilters []postgres.Filter
	if serviceFilters.ID != nil {
		postgresFilters = append(postgresFilters, postgres.WithID(*serviceFilters.ID))
	}
	if serviceFilters.UserID != nil {
		postgresFilters = append(postgresFilters, postgres.WithUserID(*serviceFilters.UserID))
	}
	if serviceFilters.TokenType != nil {
		tokenType, err := tokenmodel.TokenTypeFromString(serviceFilters.TokenType.String())
		if err != nil {
			return nil, err
		}
		postgresFilters = append(postgresFilters, postgres.WithTokenType(tokenType))
	}

	tokens, err := s.repo.List(ctx, postgresFilters...)
	if err != nil {
		return nil, err
	}

	entities := make([]tokenentity.Token, len(tokens))
	for i, token := range tokens {
		tokenTypeEntity, err := tokenentity.TokenTypeFromString(token.TokenType.String())
		if err != nil {
			return nil, err
		}

		entities[i] = tokenentity.Token{
			ID:        token.ID,
			UserID:    token.UserID,
			TokenType: tokenTypeEntity,
			Token:     token.Token,
		}
	}

	return entities, nil
}

func (s *Service) UpdateToken(ctx context.Context, id, token string) error {
	return s.repo.Update(ctx, id, token)
}

func (s *Service) DeleteToken(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, postgres.WithID(id))
}
