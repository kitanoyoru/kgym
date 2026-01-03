package internal

import (
	"context"
	"net"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	pbuser "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	apiv1grpc "github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc"
	userrepository "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user"
	userpostgres "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user/postgres"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service/user"
	pkgpostgres "github.com/kitanoyoru/kgym/pkg/database/postgres"
	pkgredis "github.com/kitanoyoru/kgym/pkg/database/redis"
	pkglogging "github.com/kitanoyoru/kgym/pkg/logging"
	pkgmetrics "github.com/kitanoyoru/kgym/pkg/metrics"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	Namespace   = "kgym"
	ServiceName = "user"
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
	srvMetrics := grpcprometheus.NewServerMetrics(
		grpcprometheus.WithServerCounterOptions(
			grpcprometheus.WithNamespace(Namespace),
			grpcprometheus.WithSubsystem(ServiceName),
		),
		grpcprometheus.WithServerHandlingTimeHistogram(
			grpcprometheus.WithHistogramNamespace(Namespace),
			grpcprometheus.WithHistogramSubsystem(ServiceName),
		),
		grpcprometheus.WithContextLabels(pkgmetrics.AllMetadataFields...),
	)

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(app.cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(app.cfg.MaxSendMsgSize),
		grpc.ConnectionTimeout(app.cfg.ConnectionTimeout),
		grpc.MaxConcurrentStreams(app.cfg.MaxConcurrentStreams),
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(
				grpcprometheus.WithLabelsFromContext(pkgmetrics.ExtractLabelsFromMetadata),
			),
			logging.UnaryServerInterceptor(
				pkglogging.NewInterceptorLogger(Namespace, ServiceName),
			),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(
				grpcprometheus.WithLabelsFromContext(pkgmetrics.ExtractLabelsFromMetadata),
			),
			logging.StreamServerInterceptor(
				pkglogging.NewInterceptorLogger(Namespace, ServiceName),
			),
		),
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	srvMetrics.InitializeMetrics(server)

	userServer, err := apiv1grpc.NewUserService(app.userService)
	if err != nil {
		return err
	}
	pbuser.RegisterUserServiceServer(server, userServer)
	pbhealth.RegisterHealthServer(server, apiv1grpc.NewHealthzService())

	reflection.Register(server)

	app.grpcServer = server

	return nil
}
