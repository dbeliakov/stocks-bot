package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "stocks_bot"
)

func NewCounter(name string) prometheus.Counter {
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      name,
	})
	prometheus.MustRegister(c)
	return c
}

func NewHistogram(name string) prometheus.Histogram {
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      name,
		Buckets:   nil,
	})
	prometheus.MustRegister(h)
	return h
}

func RunServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("failed to crease server: %+v", err)
	}
	return nil
}
