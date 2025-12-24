package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultAccessKey = "minioadmin"
	DefaultSecretKey = "minioadmin"
)

type MinioContainer struct {
	testcontainers.Container

	Endpoint  string
	AccessKey string
	SecretKey string
}

type Option func(*Options)

type Options struct {
	AccessKey string
	SecretKey string
}

func WithAccessKey(accessKey string) Option {
	return func(o *Options) {
		o.AccessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(o *Options) {
		o.SecretKey = secretKey
	}
}

func SetupTestContainer(ctx context.Context, options ...Option) (*MinioContainer, error) {
	opts := Options{
		AccessKey: DefaultAccessKey,
		SecretKey: DefaultSecretKey,
	}
	for _, option := range options {
		option(&opts)
	}

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor: wait.ForHTTP("/minio/health/live").
			WithPort("9000").
			WithStartupTimeout(120 * time.Second),
		Env: map[string]string{
			"MINIO_ACCESS_KEY": opts.AccessKey,
			"MINIO_SECRET_KEY": opts.SecretKey,
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	port, err := container.MappedPort(ctx, "9000")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	return &MinioContainer{
		Container: container,
		Endpoint:  endpoint,
		AccessKey: opts.AccessKey,
		SecretKey: opts.SecretKey,
	}, nil
}
