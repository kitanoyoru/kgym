package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType int

const (
	Counter MetricType = iota
	Gauge
	Histogram
	Summary
)

type MetricConfig struct {
	Name    string
	Help    string
	Labels  []string
	Buckets []float64 // Only for Histogram/Summary
	Type    MetricType
}

type Metric struct {
	MetricType MetricType

	Counter *prometheus.CounterVec
	Gauge   *prometheus.GaugeVec

	Histogram    prometheus.Histogram
	HistogramVec *prometheus.HistogramVec

	Summary    prometheus.Summary
	SummaryVec *prometheus.SummaryVec
}

type Registry struct {
	metrics map[string]*Metric
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: make(map[string]*Metric),
	}
}

func (r *Registry) RegisterMetric(cfg MetricConfig) error {
	switch cfg.Type {
	case Counter:
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: cfg.Name,
			Help: cfg.Help,
		}, cfg.Labels)
		if err := prometheus.Register(counter); err != nil {
			return err
		}
		r.metrics[cfg.Name] = &Metric{
			MetricType: Counter,
			Counter:    counter,
		}
	case Gauge:
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: cfg.Name,
			Help: cfg.Help,
		}, cfg.Labels)
		if err := prometheus.Register(gauge); err != nil {
			return err
		}
		r.metrics[cfg.Name] = &Metric{
			MetricType: Gauge,
			Gauge:      gauge,
		}
	case Histogram:
		var m *Metric
		if len(cfg.Labels) == 0 {
			hist := prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    cfg.Name,
				Help:    cfg.Help,
				Buckets: cfg.BucketsOrDefault(),
			})
			if err := prometheus.Register(hist); err != nil {
				return err
			}
			m = &Metric{
				MetricType: Histogram,
				Histogram:  hist,
			}
		} else {
			histVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    cfg.Name,
				Help:    cfg.Help,
				Buckets: cfg.BucketsOrDefault(),
			}, cfg.Labels)
			if err := prometheus.Register(histVec); err != nil {
				return err
			}
			m = &Metric{
				MetricType:   Histogram,
				HistogramVec: histVec,
			}
		}
		r.metrics[cfg.Name] = m
	case Summary:
		var m *Metric
		if len(cfg.Labels) == 0 {
			summary := prometheus.NewSummary(prometheus.SummaryOpts{
				Name: cfg.Name,
				Help: cfg.Help,
			})
			if err := prometheus.Register(summary); err != nil {
				return err
			}
			m = &Metric{
				MetricType: Summary,
				Summary:    summary,
			}
		} else {
			summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Name: cfg.Name,
				Help: cfg.Help,
			}, cfg.Labels)
			if err := prometheus.Register(summaryVec); err != nil {
				return err
			}
			m = &Metric{
				MetricType: Summary,
				SummaryVec: summaryVec,
			}
		}
		r.metrics[cfg.Name] = m
	}
	return nil
}

func (r *Registry) GetMetric(name string) *Metric {
	return r.metrics[name]
}

func (cfg MetricConfig) BucketsOrDefault() []float64 {
	if len(cfg.Buckets) > 0 {
		return cfg.Buckets
	}

	return prometheus.DefBuckets
}
