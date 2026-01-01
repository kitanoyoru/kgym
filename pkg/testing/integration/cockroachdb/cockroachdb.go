package cockroachdb

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CockroachDBContainer struct {
	testcontainers.Container

	URI string
}

func SetupTestContainer(ctx context.Context) (*CockroachDBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:v23.1.11",
		ExposedPorts: []string{"26257/tcp", "8080/tcp"},
		Cmd: []string{
			"start-single-node",
			"--insecure",
		},
		WaitingFor: wait.ForHTTP("/health").
			WithPort("8080").
			WithStartupTimeout(120 * time.Second),
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
		_ = container.Terminate(ctx)
		return nil, err
	}

	port, err := container.MappedPort(ctx, "26257")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, err
	}

	uri := fmt.Sprintf("postgresql://root@%s:%s/defaultdb?sslmode=disable", host, port.Port())

	return &CockroachDBContainer{
		Container: container,
		URI:       uri,
	}, nil
}
