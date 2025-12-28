package redis

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"time"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
	keyrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/multierr"
)

var _ keyrepo.IRepository = (*Repository)(nil)

type Repository struct {
	rdb *redis.ClusterClient
}

func New(ctx context.Context, rdb *redis.ClusterClient) (*Repository, error) {
	return &Repository{
		rdb: rdb,
	}, nil
}

func (r *Repository) GetCurrentSigningKey(ctx context.Context) (keyentity.Key, error) {
	kid, err := r.rdb.Get(ctx, "jwks:active").Result()
	if err != nil {
		return keyentity.Key{}, err
	}

	data, err := r.rdb.Get(ctx, "jwks:key:"+kid).Bytes()
	if err != nil {
		return keyentity.Key{}, err
	}

	var key keyentity.Key
	if err := json.Unmarshal(data, &key); err != nil {
		return keyentity.Key{}, err
	}

	return key, nil
}

func (r *Repository) GetPublicKeys(ctx context.Context) ([]keyentity.Key, error) {
	kids, err := r.rdb.SMembers(ctx, "jwks:public").Result()
	if err != nil {
		return nil, err
	}

	keys := make([]keyentity.Key, 0, len(kids))

	for _, kid := range kids {
		data, err := r.rdb.Get(ctx, "jwks:key:"+kid).Bytes()
		if err != nil {
			return nil, err
		}

		var key keyentity.Key
		if err := json.Unmarshal(data, &key); err != nil {
			continue
		}

		if !key.Active {
			continue
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// TODO: Create a k8s job to rotate the key every week
func (r *Repository) Rotate(ctx context.Context) (keyentity.Key, error) {
	privateKey, err := r.generateRSAKey()
	if err != nil {
		return keyentity.Key{}, err
	}

	kid := time.Now().UTC().Format("20060102T150405")

	key := keyentity.Key{
		ID:        kid,
		Private:   privateKey,
		Public:    privateKey.Public(),
		Algorithm: "RS256",
		Active:    true,
	}

	data, err := json.Marshal(key)
	if err != nil {
		return keyentity.Key{}, err
	}

	pipe := r.rdb.TxPipeline()

	err = multierr.Combine(
		pipe.Set(ctx, "jwks:key:"+kid, data, 0).Err(),
		pipe.Set(ctx, "jwks:active", kid, 0).Err(),
		pipe.SAdd(ctx, "jwks:public", kid).Err(),
	)
	if err != nil {
		return keyentity.Key{}, err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return keyentity.Key{}, err
	}

	return key, nil
}

func (r *Repository) generateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
