package metrics

import (
	"log"

	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics/service"
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
)

var (
	GlobalRegistry = prometheus.NewRegistry()
)

func init() {
	if err := service.RegisterMetrics(GlobalRegistry); err != nil {
		log.Fatal(err)
	}
}
