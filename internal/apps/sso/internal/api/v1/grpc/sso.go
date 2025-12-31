package grpc

import (
	"context"
	"fmt"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	keyserializer "github.com/kitanoyoru/kgym/internal/apps/sso/internal/api/v1/grpc/serializer/key"
	authservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/auth"
	keyservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/key"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GRPCServiceMetricsPrefix = "kgym.sso.api.grpc"
)

type SSOServer struct {
	pb.UnimplementedSSOServiceServer

	authService authservice.IService
	keyService  keyservice.IService
}

func NewSSOServer(authService authservice.IService, keyService keyservice.IService) (*SSOServer, error) {
	methods := []string{
		"GetToken",
		"GetJWKS",
	}

	for _, method := range methods {
		if err := metrics.GlobalRegistry.RegisterMetric(prometheus.MetricConfig{
			Name: fmt.Sprintf("%s.%s", GRPCServiceMetricsPrefix, method),
			Type: prometheus.Counter,
		}); err != nil {
			return nil, err
		}
	}

	return &SSOServer{
		authService: authService,
		keyService:  keyService,
	}, nil
}

func (s *SSOServer) GetToken(ctx context.Context, req *pb.GetToken_Request) (*pb.GetToken_Response, error) {
	switch req.Grant.(type) {
	case *pb.GetToken_Request_PasswordGrant:
		passwordGrant := req.GetPasswordGrant()

		resp, err := s.authService.PasswordGrant(ctx, authservice.PasswordGrantRequest{
			Email:    passwordGrant.Username,
			Password: passwordGrant.Password,
			ClientID: req.ClientId,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to grant password: %v", err)
		}

		return &pb.GetToken_Response{
			Token: &pb.Token{
				AccessToken:  resp.AccessToken,
				RefreshToken: resp.RefreshToken,
				TokenType:    pb.TokenType_TOKEN_TYPE_BEARER,
			},
		}, nil
	case *pb.GetToken_Request_RefreshTokenGrant:
		refreshTokenGrant := req.GetRefreshTokenGrant()

		resp, err := s.authService.RefreshTokenGrant(ctx, authservice.RefreshTokenGrantRequest{
			RefreshToken: refreshTokenGrant.RefreshToken,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to grant refresh token: %v", err)
		}

		return &pb.GetToken_Response{
			Token: &pb.Token{
				AccessToken:  resp.AccessToken,
				RefreshToken: resp.RefreshToken,
				TokenType:    pb.TokenType_TOKEN_TYPE_BEARER,
			},
		}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid grant type")
	}
}

func (s *SSOServer) GetJWKS(ctx context.Context, req *pb.GetJWKS_Request) (*pb.GetJWKS_Response, error) {
	keys, err := s.keyService.GetPublicKeys(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get public keys: %v", err)
	}

	pbKeys := make([]*pb.Key, 0, len(keys))
	for _, key := range keys {
		pbKey, err := keyserializer.EntityToPb(&key)
		if err != nil {
			continue
		}

		pbKeys = append(pbKeys, pbKey)
	}

	return &pb.GetJWKS_Response{
		Keys: pbKeys,
	}, nil
}
