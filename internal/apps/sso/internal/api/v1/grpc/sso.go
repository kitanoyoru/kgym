package grpc

import (
	"context"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	keyserializer "github.com/kitanoyoru/kgym/internal/apps/sso/internal/api/v1/grpc/serializer/key"
	authservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/auth"
	keyservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/key"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GRPCServerPrefix = "kgym.sso.api.grpc"
)

type SSOServer struct {
	pb.UnimplementedSSOServiceServer

	tracer      trace.Tracer
	authService authservice.IService
	keyService  keyservice.IService
}

func NewSSOServer(authService authservice.IService, keyService keyservice.IService) (*SSOServer, error) {
	tracer := otel.Tracer(GRPCServerPrefix)

	return &SSOServer{
		authService: authService,
		keyService:  keyService,
		tracer:      tracer,
	}, nil
}

func (s *SSOServer) GetToken(ctx context.Context, req *pb.GetToken_Request) (*pb.GetToken_Response, error) {
	ctx, span := s.tracer.Start(ctx, "GetToken")
	defer span.End()

	switch req.Grant.(type) {
	case *pb.GetToken_Request_PasswordGrant:
		passwordGrant := req.GetPasswordGrant()

		resp, err := s.authService.PasswordGrant(ctx, authservice.PasswordGrantRequest{
			Email:    passwordGrant.Username,
			Password: passwordGrant.Password,
			ClientID: req.ClientId,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to grant password")
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
			return nil, status.Error(codes.Internal, "failed to grant refresh token")
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
	ctx, span := s.tracer.Start(ctx, "GetJWKS")
	defer span.End()

	keys, err := s.keyService.GetPublicKeys(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get public keys")
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
