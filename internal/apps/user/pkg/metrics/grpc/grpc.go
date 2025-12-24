package grpc

import (
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"golang.org/x/sync/errgroup"
)

const (
	GRPCMethodCreateUser = "kgym.user.api.grpc.create"
	GRPCMethodListUsers  = "kgym.user.api.grpc.list"
	GRPCMethodDeleteUser = "kgym.user.api.grpc.delete"
)

var metrics = []prometheus.MetricConfig{
	{
		Name: GRPCMethodCreateUser,
		Help: "Number of users created",
		Type: prometheus.Counter,
	},
	{
		Name: GRPCMethodListUsers,
		Help: "Number of users listed",
		Type: prometheus.Counter,
	},
	{
		Name: GRPCMethodDeleteUser,
		Help: "Number of users deleted",
		Type: prometheus.Counter,
	},
}

func RegisterMetrics(registry *prometheus.Registry) error {
	var wg errgroup.Group

	for _, metric := range metrics {
		wg.Go(func() error {
			return registry.RegisterMetric(metric)
		})
	}

	return wg.Wait()
}
