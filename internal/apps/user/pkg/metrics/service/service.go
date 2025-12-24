package service

import (
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"golang.org/x/sync/errgroup"
)

const (
	UserCreatedMetricName = "kgym.user.service.created"
	UserListMetricName    = "kgym.user.service.list"
	UserDeletedMetricName = "kgym.user.service.deleted"
)

var metrics = []prometheus.MetricConfig{
	{
		Name: UserCreatedMetricName,
		Help: "Number of users created",
		Type: prometheus.Counter,
	},
	{
		Name: UserListMetricName,
		Help: "Number of users listed",
		Type: prometheus.Counter,
	},
	{
		Name: UserDeletedMetricName,
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
