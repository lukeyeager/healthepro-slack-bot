package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metrics holds Prometheus counters for external API calls.
type Metrics struct {
	HealtheProRequests *prometheus.CounterVec
	SlackRequests      *prometheus.CounterVec
}

// New registers and returns Prometheus metrics.
func New() *Metrics {
	m := &Metrics{
		HealtheProRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "healthepro_requests_total",
			Help: "Total HTTP requests made to the Health-e Pro API, by status code.",
		}, []string{"status_code"}),
		SlackRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "slack_requests_total",
			Help: "Total HTTP requests made to the Slack webhook, by status code.",
		}, []string{"status_code"}),
	}
	prometheus.MustRegister(m.HealtheProRequests, m.SlackRequests)
	return m
}
