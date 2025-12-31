package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
)

func PbToEntity(pbKey *pb.Key) (*keyentity.Key, error) {
	return &keyentity.Key{
		ID:        pbKey.Kid,
		Public:    pbKey.Public,
		Algorithm: pbKey.Algorithm,
		Active:    pbKey.Active,
	}, nil
}

func EntityToPb(entityKey *keyentity.Key) (*pb.Key, error) {
	privateKey, ok := entityKey.Private.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKey, ok := entityKey.Public.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	publicDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicDER,
	})

	return &pb.Key{
		Kid:       entityKey.ID,
		Private:   string(privatePEM),
		Public:    string(publicPEM),
		Algorithm: entityKey.Algorithm,
		Active:    entityKey.Active,
	}, nil
}
