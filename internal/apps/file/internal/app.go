package internal

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	pbFile "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/file/v1"
	apiv1grpc "github.com/kitanoyoru/kgym/internal/apps/file/internal/api/v1/grpc"
	fileminio "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/minio"
	filepostgres "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/postgres"
	fileservice "github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
	pkgminio "github.com/kitanoyoru/kgym/pkg/database/minio"
	pkgpostgres "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg Config

	dbPool      *pgxpool.Pool
	minioClient *minio.Client
	grpcServer  *grpc.Server

	fileMinioRepository    fileminio.IRepository
	filePostgresRepository filepostgres.IRepository

	fileService fileservice.IService
}

func New(ctx context.Context, cfg Config) (*App, error) {
	dbPool, err := pkgpostgres.New(ctx, pkgpostgres.Config{
		URI: cfg.ConnectionString,
	})
	if err != nil {
		return nil, err
	}

	minioClient, err := pkgminio.New(ctx, pkgminio.Config{
		Endpoint:  cfg.Static.Endpoint,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		Secure:    cfg.Secure,
	})
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg:         cfg,
		dbPool:      dbPool,
		minioClient: minioClient,
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
	fileMinioRepository, err := fileminio.New(
		app.minioClient,
		fileminio.WithBuckets(
			lo.Values(app.cfg.Buckets)...,
		),
	)
	if err != nil {
		return err
	}
	app.fileMinioRepository = fileMinioRepository

	app.filePostgresRepository = filepostgres.New(app.dbPool)

	return nil
}

func (app *App) initServices(_ context.Context) error {
	app.fileService = fileservice.New(fileservice.Config{
		Buckets: app.cfg.Buckets,
	}, app.fileMinioRepository, app.filePostgresRepository)

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

	fileServer, err := apiv1grpc.NewFileService(app.fileService)
	if err != nil {
		return err
	}
	pbFile.RegisterFileServiceServer(server, fileServer)

	reflection.Register(server)

	app.grpcServer = server

	return nil
}
