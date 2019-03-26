package client

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_requests",
		Help: "Count of requests to the client",
	},
		[]string{"type"},
	)

	replies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_replies",
		Help: "Count of replies from the client",
	},
		[]string{"type"},
	)

	errors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_errors",
		Help: "Count of errors from the client",
	},
		[]string{"type"},
	)

	sendErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_send_errors",
		Help: "Count of errors from the client specifically in sending the response",
	},
		[]string{"type"},
	)

	requestCancelled = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_request_errors",
		Help: "Count of cancelled from the client attempting to send the request down the mediation channel",
	},
		[]string{"type"},
	)

	awaitErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_await_errors",
		Help: "Count of errors from the client after sending the request awaiting the response",
	},
		[]string{"type"},
	)

	awaitCancelled = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_await_cancelled",
		Help: "Count of cancels from the client after sending the request awaiting the response",
	},
		[]string{"type"},
	)

	limit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_rate_limits",
		Help: "Count of rate limit results from the client",
	},
		[]string{"type"},
	)
)

func initMetricsForAllTypes() {
	for _, t := range allTypes {
		cRequests.WithLabelValues(typeName(t)).Add(0)
		replies.WithLabelValues(typeName(t)).Add(0)
		errors.WithLabelValues(typeName(t)).Add(0)
		sendErrors.WithLabelValues(typeName(t)).Add(0)
		requestCancelled.WithLabelValues(typeName(t)).Add(0)
		awaitErrors.WithLabelValues(typeName(t)).Add(0)
		awaitCancelled.WithLabelValues(typeName(t)).Add(0)
		limit.WithLabelValues(typeName(t)).Add(0)
	}
}
