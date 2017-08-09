package negroniprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

func Test_Logger(t *testing.T) {
	recorder := httptest.NewRecorder()

	//To create new registry for prometheus so that, despite registering metrics
	// with same name, will be done to different registry a registry shared among tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)

	n := negroni.New()
	m := NewMiddleware("test")
	n.Use(m)
	r := http.NewServeMux()
	r.Handle("/metrics", prometheus.Handler())
	r.HandleFunc(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	n.UseHandler(r)

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	body := recorder.Body.String()
	if !strings.Contains(body, reqsName) {
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", reqsName)
	}
}

func Test_LoggerWithoutPath(t *testing.T) {
	recorder := httptest.NewRecorder()


	//To create new registry for prometheus so that, despite registering metrics
	// with same name, will be done to different registry a registry shared among tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	prometheus.DefaultGatherer = prometheus.DefaultRegisterer.(prometheus.Gatherer)

	n := negroni.New()
	m := NewMiddlewareWithoutPath("test")
	n.Use(m)
	r := http.NewServeMux()
	r.Handle("/metrics", prometheus.Handler())
	r.HandleFunc(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	n.UseHandler(r)

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	body := recorder.Body.String()
	if !strings.Contains(body, reqsName) {
		fmt.Printf("%f", body)
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", reqsName)
	}
}
