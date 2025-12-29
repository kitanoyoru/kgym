package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
)

type Key struct {
	ID        string `json:"kid"`
	Private   string `json:"private"`
	Public    string `json:"public"`
	Algorithm string `json:"alg"`
	Active    bool   `json:"active"`
}

func FromEntity(entity keyentity.Key) (Key, error) {
	privateKey, ok := entity.Private.(*rsa.PrivateKey)
	if !ok {
		return Key{}, errors.New("private key is not RSA")
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKey, ok := entity.Public.(*rsa.PublicKey)
	if !ok {
		return Key{}, errors.New("public key is not RSA")
	}

	publicDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return Key{}, err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicDER,
	})

	return Key{
		ID:        entity.ID,
		Private:   string(privatePEM),
		Public:    string(publicPEM),
		Algorithm: entity.Algorithm,
		Active:    entity.Active,
	}, nil
}

func (k Key) ToEntity() (keyentity.Key, error) {
	privateBlock, _ := pem.Decode([]byte(k.Private))
	if privateBlock == nil {
		return keyentity.Key{}, errors.New("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return keyentity.Key{}, err
	}

	publicBlock, _ := pem.Decode([]byte(k.Public))
	if publicBlock == nil {
		return keyentity.Key{}, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return keyentity.Key{}, err
	}

	return keyentity.Key{
		ID:        k.ID,
		Private:   privateKey,
		Public:    publicKey,
		Algorithm: k.Algorithm,
		Active:    k.Active,
	}, nil
}
