package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/file/v1"
	"github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/gateway/internal/middlewares"
	"github.com/rs/cors"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func New(ctx context.Context, cfg Config) (*Gateway, error) {
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
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.MaxGRPCMsgSize),
			grpc.MaxCallSendMsgSize(cfg.MaxGRPCMsgSize),
		),
	}

	err := multierr.Combine(
		user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCEndpoint, opts),
		file.RegisterFileServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCEndpoint, opts),
	)
	if err != nil {
		return nil, err
	}

	handler := middlewares.Logging(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
	}).Handler(mux))

	return &Gateway{
		server: &http.Server{
			Addr:         ":" + cfg.HTTPPort,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
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
