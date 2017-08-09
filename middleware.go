package negroniprometheus

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
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
}

// MiddlewareWithoutPath is like Middleware but will only partition latency and response
// size by status code and method. This is useful for applications with such a wide variety
// of paths that metrics would become extremely large to deliver
type MiddlewareWithoutPath struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewMiddleware returns a new prometheus Middleware handler.
func NewMiddleware(name string, buckets ...float64) negroni.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return &m
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := negroni.NewResponseWriter(rw)
	m.reqs.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Inc()
	m.latency.WithLabelValues(http.StatusText(res.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}

// NewMiddlewareWithoutPath returns a new prometheus MiddlewareWithoutPath handler
func NewMiddlewareWithoutPath(name string, buckets ...float64) negroni.Handler {
	var m MiddlewareWithoutPath
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code and method",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code and method",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(m.latency)
	return &m
}

func (m *MiddlewareWithoutPath) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := negroni.NewResponseWriter(rw)
	m.reqs.WithLabelValues(http.StatusText(res.Status()), r.Method).Inc()
	m.latency.WithLabelValues(http.StatusText(res.Status()), r.Method).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
