# negroni-prometheus [![GoDoc](http://godoc.org/github.com/zbinderen/negroni-prometheus?status.svg)](http://godoc.org/github.com/zbindenren/negroni-prometeus)
[Prometheus](http://prometheus.io) middleware for [negroni](https://github.com/codegangsta/negroni).

## Why
[Logging v. instrumentation](http://peter.bourgon.org/blog/2016/02/07/logging-v-instrumentation.html)

Instead of logging request times, it is considered best practice to provide an instrumentation endpoint (like prometheus) fore those metrics.

## Usage

Use this middleware like the negroni.Logger middleware (after negroni.Recovery before every other middleware).

Take a look at the [example](./example/main.go).
