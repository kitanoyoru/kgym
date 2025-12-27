package service

import (
	"context"
	"errors"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/postgres"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockRepo *mocks.MockIRepository
	service  *Service
	ctx      context.Context
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mocks.NewMockIRepository(s.ctrl)
	s.service = New(s.mockRepo)
	s.ctx = context.Background()
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceTestSuite) TestCreateToken() {
	s.Run("should create a token successfully", func() {
		req := CreateTokenRequest{
			UserID:    uuid.New().String(),
			TokenType: tokenentity.TokenTypeRefresh,
			Token:     "refresh-token-value",
		}

		expectedTokenModel := tokenmodel.Token{
			ID:        "",
			UserID:    req.UserID,
			TokenType: tokenmodel.TokenTypeRefresh,
			Token:     req.Token,
		}

		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, token tokenentity.Token) error {
				assert.NotEmpty(s.T(), token.ID)
				assert.Equal(s.T(), expectedTokenModel.UserID, token.UserID)
				assert.Equal(s.T(), tokenentity.TokenTypeRefresh, token.TokenType)
				assert.Equal(s.T(), expectedTokenModel.Token, token.Token)
				return nil
			})

		tokenID, err := s.service.CreateToken(s.ctx, req)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), tokenID)
	})

	s.Run("should not create a token because of invalid token type", func() {
		req := CreateTokenRequest{
			UserID:    uuid.New().String(),
			TokenType: tokenentity.TokenType("invalid"),
			Token:     "refresh-token-value",
		}

		tokenID, err := s.service.CreateToken(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), tokenID)
	})

	s.Run("should not create a token because of empty user_id", func() {
		req := CreateTokenRequest{
			UserID:    "",
			TokenType: tokenentity.TokenTypeRefresh,
			Token:     "refresh-token-value",
		}

		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, token tokenentity.Token) error {
				// Repository might validate or database might reject empty user_id
				return errors.New("invalid user_id")
			})

		tokenID, err := s.service.CreateToken(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), tokenID)
	})

	s.Run("should not create a token because of empty token", func() {
		req := CreateTokenRequest{
			UserID:    uuid.New().String(),
			TokenType: tokenentity.TokenTypeRefresh,
			Token:     "",
		}

		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, token tokenentity.Token) error {
				// Repository might validate or database might reject empty token
				return errors.New("invalid token")
			})

		tokenID, err := s.service.CreateToken(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), tokenID)
	})

	s.Run("should not create a token because of repository error", func() {
		req := CreateTokenRequest{
			UserID:    uuid.New().String(),
			TokenType: tokenentity.TokenTypeRefresh,
			Token:     "refresh-token-value",
		}

		expectedError := errors.New("repository error")
		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			Return(expectedError)

		tokenID, err := s.service.CreateToken(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		assert.Empty(s.T(), tokenID)
	})
}

func (s *ServiceTestSuite) TestListTokens() {
	s.Run("should return list of tokens successfully with no filters", func() {
		tokenModels := []tokenmodel.Token{
			{
				ID:        uuid.New().String(),
				UserID:    uuid.New().String(),
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-1",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
			{
				ID:        uuid.New().String(),
				UserID:    uuid.New().String(),
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-2",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx).
			Return(tokenModels, nil)

		tokens, err := s.service.ListTokens(s.ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), tokens, 2)

		var ids []string
		for _, t := range tokens {
			ids = append(ids, t.ID)
		}
		assert.Contains(s.T(), ids, tokenModels[0].ID)
		assert.Contains(s.T(), ids, tokenModels[1].ID)
	})

	s.Run("should return list of tokens successfully with id filter", func() {
		tokenID := uuid.New().String()
		tokenModels := []tokenmodel.Token{
			{
				ID:        tokenID,
				UserID:    uuid.New().String(),
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-1",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]tokenmodel.Token, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.ID)
				assert.Equal(s.T(), tokenID, *dbFilters.ID)
				return tokenModels, nil
			})

		tokens, err := s.service.ListTokens(s.ctx, WithID(tokenID))
		require.NoError(s.T(), err)
		assert.Len(s.T(), tokens, 1)
		assert.Equal(s.T(), tokenID, tokens[0].ID)
	})

	s.Run("should return list of tokens successfully with user_id filter", func() {
		userID := uuid.New().String()
		tokenModels := []tokenmodel.Token{
			{
				ID:        uuid.New().String(),
				UserID:    userID,
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-1",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]tokenmodel.Token, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.UserID)
				assert.Equal(s.T(), userID, *dbFilters.UserID)
				return tokenModels, nil
			})

		tokens, err := s.service.ListTokens(s.ctx, WithUserID(userID))
		require.NoError(s.T(), err)
		assert.Len(s.T(), tokens, 1)
		assert.Equal(s.T(), userID, tokens[0].UserID)
	})

	s.Run("should return list of tokens successfully with token type filter", func() {
		tokenType := tokenentity.TokenTypeRefresh
		tokenModels := []tokenmodel.Token{
			{
				ID:        uuid.New().String(),
				UserID:    uuid.New().String(),
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-1",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
			{
				ID:        uuid.New().String(),
				UserID:    uuid.New().String(),
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-2",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]tokenmodel.Token, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.TokenType)
				assert.Equal(s.T(), tokenmodel.TokenTypeRefresh, *dbFilters.TokenType)
				return tokenModels, nil
			})

		tokens, err := s.service.ListTokens(s.ctx, WithTokenType(tokenType))
		require.NoError(s.T(), err)
		assert.Len(s.T(), tokens, 2)
		for _, t := range tokens {
			assert.Equal(s.T(), tokenType, t.TokenType)
		}
	})

	s.Run("should return list of tokens successfully with multiple filters", func() {
		userID := uuid.New().String()
		tokenType := tokenentity.TokenTypeRefresh
		tokenModels := []tokenmodel.Token{
			{
				ID:        uuid.New().String(),
				UserID:    userID,
				TokenType: tokenmodel.TokenTypeRefresh,
				Token:     "refresh-token-1",
				CreatedAt: carbon.Now().StdTime(),
				UpdatedAt: carbon.Now().StdTime(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]tokenmodel.Token, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.UserID)
				assert.NotNil(s.T(), dbFilters.TokenType)
				assert.Equal(s.T(), userID, *dbFilters.UserID)
				assert.Equal(s.T(), tokenmodel.TokenTypeRefresh, *dbFilters.TokenType)
				return tokenModels, nil
			})

		tokens, err := s.service.ListTokens(s.ctx, WithUserID(userID), WithTokenType(tokenType))
		require.NoError(s.T(), err)
		assert.Len(s.T(), tokens, 1)
		assert.Equal(s.T(), userID, tokens[0].UserID)
		assert.Equal(s.T(), tokenType, tokens[0].TokenType)
	})

	s.Run("should return empty list when no tokens found", func() {
		s.mockRepo.EXPECT().
			List(s.ctx).
			Return([]tokenmodel.Token{}, nil)

		tokens, err := s.service.ListTokens(s.ctx)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), tokens)
	})

	s.Run("should return error when repository fails", func() {
		expectedError := errors.New("repository error")
		s.mockRepo.EXPECT().
			List(s.ctx).
			Return(nil, expectedError)

		tokens, err := s.service.ListTokens(s.ctx)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		assert.Nil(s.T(), tokens)
	})

	s.Run("should return error when invalid token type in filter", func() {
		invalidTokenType := tokenentity.TokenType("invalid")
		tokens, err := s.service.ListTokens(s.ctx, WithTokenType(invalidTokenType))
		assert.Error(s.T(), err)
		assert.Nil(s.T(), tokens)
	})
}

func (s *ServiceTestSuite) TestUpdateToken() {
	s.Run("should update a token successfully", func() {
		tokenID := uuid.New().String()
		newTokenValue := "new-refresh-token-value"

		s.mockRepo.EXPECT().
			Update(s.ctx, tokenID, newTokenValue).
			Return(nil)

		err := s.service.UpdateToken(s.ctx, tokenID, newTokenValue)
		require.NoError(s.T(), err)
	})

	s.Run("should return error when repository fails", func() {
		tokenID := uuid.New().String()
		newTokenValue := "new-refresh-token-value"
		expectedError := errors.New("repository error")

		s.mockRepo.EXPECT().
			Update(s.ctx, tokenID, newTokenValue).
			Return(expectedError)

		err := s.service.UpdateToken(s.ctx, tokenID, newTokenValue)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
	})
}

func (s *ServiceTestSuite) TestDeleteToken() {
	s.Run("should delete a token successfully", func() {
		tokenID := uuid.New().String()

		s.mockRepo.EXPECT().
			Delete(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) error {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.ID)
				assert.Equal(s.T(), tokenID, *dbFilters.ID)
				return nil
			})

		err := s.service.DeleteToken(s.ctx, tokenID)
		require.NoError(s.T(), err)
	})

	s.Run("should return error when repository fails", func() {
		tokenID := uuid.New().String()
		expectedError := errors.New("repository error")

		s.mockRepo.EXPECT().
			Delete(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) error {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.ID)
				assert.Equal(s.T(), tokenID, *dbFilters.ID)
				return expectedError
			})

		err := s.service.DeleteToken(s.ctx, tokenID)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
