package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metrics holds Prometheus counters for external API calls and fetch health.
type Metrics struct {
	HealtheProRequests *prometheus.CounterVec
	FetchFailures      prometheus.Counter
}

// New registers and returns Prometheus metrics.
func New() *Metrics {
	m := &Metrics{
		HealtheProRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "healthepro_requests_total",
			Help: "Total HTTP requests made to the Health-e Pro API, by status code.",
		}, []string{"status_code"}),
		FetchFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "menu_fetch_failures_total",
			Help: "Total hourly fetch attempts that failed after all retries.",
		}),
	}
	prometheus.MustRegister(m.HealtheProRequests, m.FetchFailures)
	return m
}
