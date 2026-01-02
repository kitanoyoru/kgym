package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pbFile "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/file/v1"
	pbSSO "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	pbUser "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/gateway/internal/handlers/file"
	"github.com/kitanoyoru/kgym/internal/gateway/internal/middlewares"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func New(ctx context.Context, cfg Config) (*Gateway, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.MaxGRPCMsgSize),
			grpc.MaxCallSendMsgSize(cfg.MaxGRPCMsgSize),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
	}

	healthConn, err := grpc.NewClient(cfg.GRPCEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	healthClient := grpc_health_v1.NewHealthClient(healthConn)

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := map[string]string{
				"x-request-id":  r.Header.Get("X-Request-ID"),
				"x-platform":    r.Header.Get("X-Platform"),
				"x-app-version": r.Header.Get("X-App-Version"),
			}

			authorization := r.Header.Get("Authorization")
			if authorization != "" {
				if token, err := r.Cookie("access_token"); err == nil {
					md["authorization"] = "Bearer " + token.Value
				}
			}
			md["authorization"] = authorization

			return metadata.New(md)
		}),
		runtime.WithHealthzEndpoint(healthClient),
	)

	err = multierr.Combine(
		pbUser.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCEndpoint, opts),
		pbFile.RegisterFileServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCEndpoint, opts),
		pbSSO.RegisterSSOServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCEndpoint, opts),
	)
	if err != nil {
		return nil, err
	}

	fileHandler, err := file.New(ctx, mux, file.Config{
		GRPCEndpoint:    cfg.GRPCEndpoint,
		GRPCDialOptions: opts,
		BodyLimit:       cfg.BodyLimit,
		ChunkSize:       5 * 1024 * 1024, // 5MB
	})
	if err != nil {
		return nil, err
	}

	err = multierr.Combine(
		mux.HandlePath(http.MethodPost, "/api/v1/files/user-avatar", fileHandler.UploadUserAvatar()),
	)
	if err != nil {
		return nil, err
	}

	handler := middlewares.Tracing(
		middlewares.Logging(
			cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedHeaders: []string{"*"},
				AllowedMethods: []string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodDelete,
					http.MethodOptions,
				},
			}).Handler(mux),
		),
	)

	return &Gateway{
		server: &http.Server{
			Addr:         ":" + cfg.HTTPPort,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
			Handler:      handler,
		},
	}, nil
}

type Gateway struct {
	server *http.Server
}

func (g *Gateway) Run(ctx context.Context) error {
	return g.server.ListenAndServe()
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}
