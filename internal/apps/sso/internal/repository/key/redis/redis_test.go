package redis

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
	rediscontainer "github.com/kitanoyoru/kgym/pkg/testing/integration/redis"
)

type RepositoryTestSuite struct {
	suite.Suite

	rdb        *redis.Client
	container  *rediscontainer.RedisContainer
	repository *Repository
	ctx        context.Context
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	container, err := rediscontainer.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup test container")

	s.container = container

	// Parse URI to extract host and port
	parsedURI, err := url.Parse(container.URI)
	require.NoError(s.T(), err, "failed to parse Redis URI")

	address := strings.TrimPrefix(container.URI, "redis://")
	if parsedURI.Host != "" {
		address = parsedURI.Host
	}

	s.rdb = redis.NewClient(&redis.Options{
		Addr: address,
	})

	repository, err := New(ctx, s.rdb)
	require.NoError(s.T(), err, "failed to create repository")

	s.repository = repository
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.rdb != nil {
		_ = s.rdb.Close()
	}
	if s.container != nil {
		_ = s.container.Terminate(s.T().Context())
	}
}

func (s *RepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	s.cleanupKeys(ctx)
}

func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	s.cleanupKeys(ctx)
}

func (s *RepositoryTestSuite) cleanupKeys(ctx context.Context) {
	_ = s.rdb.Del(ctx, "jwks:active").Err()

	kids, err := s.rdb.SMembers(ctx, "jwks:public").Result()
	if err == nil && len(kids) > 0 {
		for _, kid := range kids {
			_ = s.rdb.Del(ctx, "jwks:key:"+kid).Err()
		}
		_ = s.rdb.Del(ctx, "jwks:public").Err()
	}

	keys, err := s.rdb.Keys(ctx, "jwks:*").Result()
	if err == nil && len(keys) > 0 {
		_ = s.rdb.Del(ctx, keys...).Err()
	}
}

func (s *RepositoryTestSuite) TestRotate() {
	s.Run("should rotate key successfully", func() {
		key, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), key.ID)
		assert.NotNil(s.T(), key.Private)
		assert.NotNil(s.T(), key.Public)
		assert.Equal(s.T(), "RS256", key.Algorithm)
		assert.True(s.T(), key.Active)

		// Verify key is stored in Redis
		activeKid, err := s.rdb.Get(s.ctx, "jwks:active").Result()
		require.NoError(s.T(), err)
		assert.Equal(s.T(), key.ID, activeKid)

		// Verify key data is stored
		data, err := s.rdb.Get(s.ctx, "jwks:key:"+key.ID).Bytes()
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), data)

		// Verify key is in public set
		isMember, err := s.rdb.SIsMember(s.ctx, "jwks:public", key.ID).Result()
		require.NoError(s.T(), err)
		assert.True(s.T(), isMember)
	})

	s.Run("should rotate key multiple times", func() {
		s.cleanupKeys(s.ctx)
		key1, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		time.Sleep(time.Second)
		key2, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		assert.NotEqual(s.T(), key1.ID, key2.ID)

		// Verify latest key is active
		activeKid, err := s.rdb.Get(s.ctx, "jwks:active").Result()
		require.NoError(s.T(), err)
		assert.Equal(s.T(), key2.ID, activeKid)

		// Verify both keys are in public set
		isMember1, err := s.rdb.SIsMember(s.ctx, "jwks:public", key1.ID).Result()
		require.NoError(s.T(), err)
		assert.True(s.T(), isMember1)

		isMember2, err := s.rdb.SIsMember(s.ctx, "jwks:public", key2.ID).Result()
		require.NoError(s.T(), err)
		assert.True(s.T(), isMember2)
	})
}

func (s *RepositoryTestSuite) TestGetCurrentSigningKey() {
	s.Run("should get current signing key successfully", func() {
		// First rotate to create a key
		rotatedKey, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		// Get the current signing key
		key, err := s.repository.GetCurrentSigningKey(s.ctx)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), rotatedKey.ID, key.ID)
		assert.Equal(s.T(), rotatedKey.Algorithm, key.Algorithm)
		assert.True(s.T(), key.Active)
		assert.NotNil(s.T(), key.Private)
		assert.NotNil(s.T(), key.Public)
	})

	s.Run("should return error when no active key exists", func() {
		s.cleanupKeys(s.ctx)
		_, err := s.repository.GetCurrentSigningKey(s.ctx)
		assert.Error(s.T(), err)
	})

	s.Run("should return error when active key ID exists but key data is missing", func() {
		// Set active key ID but don't set the key data
		err := s.rdb.Set(s.ctx, "jwks:active", "missing-key-id", 0).Err()
		require.NoError(s.T(), err)

		_, err = s.repository.GetCurrentSigningKey(s.ctx)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestGetPublicKeys() {
	s.Run("should get public keys successfully", func() {
		s.cleanupKeys(s.ctx)
		// Rotate multiple keys
		key1, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		time.Sleep(time.Second)
		key2, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		// Get public keys
		keys, err := s.repository.GetPublicKeys(s.ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), keys, 2)

		keyIDs := make(map[string]bool)
		for _, k := range keys {
			keyIDs[k.ID] = true
			assert.True(s.T(), k.Active)
			assert.Equal(s.T(), "RS256", k.Algorithm)
			assert.NotNil(s.T(), k.Public)
		}

		assert.True(s.T(), keyIDs[key1.ID])
		assert.True(s.T(), keyIDs[key2.ID])
	})

	s.Run("should return empty list when no public keys exist", func() {
		s.cleanupKeys(s.ctx)
		keys, err := s.repository.GetPublicKeys(s.ctx)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), keys)
	})

	s.Run("should only return active keys", func() {
		s.cleanupKeys(s.ctx)
		// Rotate a key
		key, err := s.repository.Rotate(s.ctx)
		require.NoError(s.T(), err)

		// Manually add an inactive key to the public set
		// Get existing key and modify it to be inactive
		var existingKey keyentity.Key
		data, err := s.rdb.Get(s.ctx, "jwks:key:"+key.ID).Bytes()
		require.NoError(s.T(), err)
		err = json.Unmarshal(data, &existingKey)
		require.NoError(s.T(), err)
		existingKey.Active = false
		existingKey.ID = "inactive-key-id"
		// Marshal the inactive key
		inactiveData, err := json.Marshal(existingKey)
		require.NoError(s.T(), err)
		_ = s.rdb.Set(s.ctx, "jwks:key:inactive-key-id", inactiveData, 0).Err()
		_ = s.rdb.SAdd(s.ctx, "jwks:public", "inactive-key-id").Err()

		// Get public keys - should only return active ones
		keys, err := s.repository.GetPublicKeys(s.ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), keys, 1)
		assert.Equal(s.T(), key.ID, keys[0].ID)
		assert.True(s.T(), keys[0].Active)
	})

	s.Run("should skip keys with invalid JSON", func() {
		s.cleanupKeys(s.ctx)
		// Add invalid JSON to a key
		_ = s.rdb.Set(s.ctx, "jwks:key:invalid-json", "invalid json data", 0).Err()
		_ = s.rdb.SAdd(s.ctx, "jwks:public", "invalid-json").Err()

		// Get public keys - should skip invalid JSON
		keys, err := s.repository.GetPublicKeys(s.ctx)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), keys)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
