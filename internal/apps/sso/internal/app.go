package internal

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	pbsso "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	pbuser "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	apiv1grpc "github.com/kitanoyoru/kgym/internal/apps/sso/internal/api/v1/grpc"
	keyrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key"
	keyredis "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key/redis"
	tokenrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token"
	tokenpostgres "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/postgres"
	userrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user"
	usergrpc "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/grpc"
	authservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/auth"
	keyservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/key"
	pkgpostgres "github.com/kitanoyoru/kgym/pkg/database/postgres"
	pkgredis "github.com/kitanoyoru/kgym/pkg/database/redis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg Config

	dbPool     *pgxpool.Pool
	rdb        *redis.ClusterClient
	grpcServer *grpc.Server

	keyRepository   keyrepo.IRepository
	tokenRepository tokenrepo.IRepository
	userRepository  userrepo.IRepository

	keyService  keyservice.IService
	authService authservice.IService
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

func (app *App) initRepositories(ctx context.Context) error {
	keyRepository, err := keyredis.New(ctx, app.rdb)
	if err != nil {
		return err
	}
	app.keyRepository = keyRepository

	app.tokenRepository = tokenpostgres.New(app.dbPool)

	client, err := grpc.NewClient(app.cfg.UserEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	app.userRepository = usergrpc.New(pbuser.NewUserServiceClient(client))

	return nil
}

func (app *App) initServices(_ context.Context) error {
	app.keyService = keyservice.NewService(app.keyRepository)
	app.authService = authservice.NewService(app.userRepository, app.tokenRepository, app.keyRepository)

	return nil
}

func (app *App) initGRPCServer(_ context.Context) error {
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(app.cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(app.cfg.MaxSendMsgSize),
		grpc.ConnectionTimeout(app.cfg.ConnectionTimeout),
		grpc.MaxConcurrentStreams(app.cfg.MaxConcurrentStreams),
	)

	ssoServer, err := apiv1grpc.NewSSOServer(app.authService, app.keyService)
	if err != nil {
		return err
	}
	pbsso.RegisterSSOServiceServer(server, ssoServer)

	app.grpcServer = server

	return nil
}
