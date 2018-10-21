package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsInFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served.",
	})

	requestsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Count of api requests",
		},
		[]string{"handler", "code", "method"},
	)

	requestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "code", "method"},
	)

	requestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_size_bytes",
			Help:    "A histogram of request sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"handler", "code", "method"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"handler", "code", "method"},
	)
)

func init() {
	prometheus.MustRegister(requestsInFlightGauge, requestsCount, requestsDuration, requestSize, responseSize)
}

func instrumentHandler(handlerName string, h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(requestsInFlightGauge,
		promhttp.InstrumentHandlerDuration(requestsDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerCounter(requestsCount.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				promhttp.InstrumentHandlerRequestSize(requestSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					promhttp.InstrumentHandlerResponseSize(responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
						h)))))
}
