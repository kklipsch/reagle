package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	errorsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
		Help: "Count of non-fatal errors",
	},
		[]string{"type"},
	)

	requestsInFlightGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served.",
	})

	requestsCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Count of api requests",
		},
		[]string{"handler", "code", "method"},
	)

	requestsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "code", "method"},
	)

	requestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_size_bytes",
			Help:    "A histogram of request sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"handler", "code", "method"},
	)

	responseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"handler", "code", "method"},
	)
)

func initializeErrorCounts() {
	errorsCount.WithLabelValues("context").Inc()
	errorsCount.WithLabelValues("other").Inc()
}

func instrumentError(err error, desc string) {
	if err == nil {
		return
	}

	if err == context.DeadlineExceeded {
		errorsCount.WithLabelValues("context").Inc()
	} else {
		errorsCount.WithLabelValues("other").Inc()
		applicationLogger.WithFields(log.Fields{"err": err}).Errorln(desc)
	}
}

func instrumentedAPI(cfg local.Config) (local.API, error) {
	var err error
	localAPI := local.New(cfg)
	localAPI.Client.Transport, err = instrumentClient("local", localAPI.Client.Transport)
	return localAPI, err
}

func instrumentHandler(handlerName string, h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(requestsInFlightGauge,
		promhttp.InstrumentHandlerDuration(requestsDuration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerCounter(requestsCount.MustCurryWith(prometheus.Labels{"handler": handlerName}),
				promhttp.InstrumentHandlerRequestSize(requestSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					promhttp.InstrumentHandlerResponseSize(responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
						h)))))
}

func instrumentClient(clientName string, inner http.RoundTripper) (http.RoundTripper, error) {
	if inner == nil {
		inner = http.DefaultTransport
	}

	count, err := requestsCount.CurryWith(prometheus.Labels{"handler": fmt.Sprintf("client_%s", clientName)})
	if err != nil {
		return nil, err
	}

	duration, err := requestsDuration.CurryWith(prometheus.Labels{"handler": fmt.Sprintf("client_%s", clientName)})
	if err != nil {
		return nil, err
	}

	return promhttp.InstrumentRoundTripperInFlight(requestsInFlightGauge,
		promhttp.InstrumentRoundTripperCounter(count,
			promhttp.InstrumentRoundTripperDuration(duration, inner))), nil
}
