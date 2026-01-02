package internal

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	pbuser "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	apiv1grpc "github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc"
	userrepository "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user"
	userpostgres "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user/postgres"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service/user"
	pkgpostgres "github.com/kitanoyoru/kgym/pkg/database/postgres"
	pkgredis "github.com/kitanoyoru/kgym/pkg/database/redis"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg Config

	dbPool     *pgxpool.Pool
	rdb        *redis.ClusterClient
	grpcServer *grpc.Server

	userRepository userrepository.IRepository

	userService userservice.IService
}

func New(ctx context.Context, cfg Config) (*App, error) {
	dbPool, err := pkgpostgres.New(ctx, pkgpostgres.Config{
		URI: cfg.ConnectionString,
	})
	if err != nil {
		return nil, err
	}

	rdb, err := pkgredis.New(ctx, pkgredis.Config{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg:    cfg,
		dbPool: dbPool,
		rdb:    rdb,
	}

	err = multierr.Combine(
		app.initRepositories(ctx),
		app.initServices(ctx),
		app.initGRPCServer(ctx),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", app.cfg.Endpoint)
	if err != nil {
		return err
	}

	return app.grpcServer.Serve(listener)
}

func (app *App) Shutdown(ctx context.Context) error {
	app.grpcServer.GracefulStop()
	return nil
}

func (app *App) initRepositories(_ context.Context) error {
	app.userRepository = userpostgres.New(app.dbPool)

	return nil
}

func (app *App) initServices(_ context.Context) error {
	app.userService = userservice.New(app.userRepository)

	return nil
}

func (app *App) initGRPCServer(_ context.Context) error {
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(app.cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(app.cfg.MaxSendMsgSize),
		grpc.ConnectionTimeout(app.cfg.ConnectionTimeout),
		grpc.MaxConcurrentStreams(app.cfg.MaxConcurrentStreams),
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	userServer, err := apiv1grpc.NewUserService(app.userService)
	if err != nil {
		return err
	}
	pbuser.RegisterUserServiceServer(server, userServer)

	reflection.Register(server)

	app.grpcServer = server

	return nil
}
