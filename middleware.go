package negroniprometheus

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
	EnablePathLogging = true
)

const (
	reqsName    = "negroni_requests_total"
	latencyName = "negroni_request_duration_milliseconds"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
	EnablePathLogging bool
}

// NewMiddleware returns a new prometheus Middleware handler.
func NewMiddleware(name string, buckets ...float64) *Middleware {
	var m Middleware
	var labels []string
	var helpEnd string
	if EnablePathLogging {
		labels = []string{"code", "method", "path"}
		helpEnd = "partitioned by status code, method and HTTP path."
	} else {
		labels = []string{"code", "method"}
		helpEnd = "partitioned by status code and method."
	}
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, " + helpEnd,
			ConstLabels: prometheus.Labels{"service": name},
		},
		labels,
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, " + helpEnd,
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		labels,
	)
	prometheus.MustRegister(m.latency)
	return &m
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := negroni.NewResponseWriter(rw)
	if EnablePathLogging {
		m.reqs.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Inc()
		m.latency.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	} else {
		m.reqs.WithLabelValues(http.StatusText(res.Status()), r.Method).Inc()
		m.latency.WithLabelValues(http.StatusText(res.Status()), r.Method).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}

}
