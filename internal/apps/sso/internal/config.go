package internal

import "time"

type Config struct {
	GRPC
	Cache
	Database

	ShutdownTimeout time.Duration `env:"KGYM_SSO_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type GRPC struct {
	Endpoint string `env:"KGYM_SSO_GRPC_ENDPOINT" validate:"required"`

	UserEndpoint string `env:"KGYM_SSO_USER_ENDPOINT" validate:"required"`

	MaxSendMsgSize       int           `env:"KGYM_SSO_GRPC_MAX_SEND_MSG_SIZE" envDefault:"1024"`
	MaxRecvMsgSize       int           `env:"KGYM_SSO_GRPC_MAX_RECV_MSG_SIZE" envDefault:"1024"`
	ConnectionTimeout    time.Duration `env:"KGYM_SSO_GRPC_CONNECTION_TIMEOUT" envDefault:"10s"`
	MaxConcurrentStreams uint32        `env:"KGYM_SSO_GRPC_MAX_CONCURRENT_STREAMS" envDefault:"1000"`
}

type Cache struct {
	Address string `env:"KGYM_SSO_CACHE_ADDRESS" validate:"required"`
}

type Database struct {
	ConnectionString string `env:"KGYM_SSO_DATABASE_CONNECTION_STRING" validate:"required"`
}
