package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "smartedge_request_duration_seconds",
			Help:    "Tempo de resposta do backend",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"backend"},
	)
	requestSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smartedge_request_total",
			Help: "Número total de requisições por backend",
		},
		[]string{"backend", "status"},
	)
)

func Init() {
	prometheus.MustRegister(requestDuration, requestSuccess)
}

func ObserveRequest(backend string, duration time.Duration, success bool) {
	requestDuration.WithLabelValues(backend).Observe(duration.Seconds())
	status := "failure"
	if success {
		status = "success"
	}
	requestSuccess.WithLabelValues(backend, status).Inc()
}

func Handler() http.Handler {
	return promhttp.Handler()
}
