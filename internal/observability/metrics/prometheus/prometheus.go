package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type prometheusHTTPObserver struct {
	totalRequest prometheus.Counter
	reqDuration  *prometheus.HistogramVec
}

func NewPrometheusHTTPObserver() *prometheusHTTPObserver {

	totalRequest := promauto.NewCounter(prometheus.CounterOpts{
		Name: "service_courier_requests_total",
		Help: "total amount of requests",
	})

	reqDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_courier_request_duration",
			Help:    "duration of HTTP-request in sec",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 2},
		},
		[]string{"method", "path", "status"},
	)

	return &prometheusHTTPObserver{totalRequest: totalRequest, reqDuration: reqDuration}

}

func (p *prometheusHTTPObserver) IncTotalRequests() {
	p.totalRequest.Inc()
}

func (p *prometheusHTTPObserver) NewRequest(method, path, status string, durationSec float64) {
	p.reqDuration.WithLabelValues(method, path, status).Observe(durationSec)
}
