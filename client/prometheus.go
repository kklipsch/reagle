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

	limit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_rate_limits",
		Help: "Count of rate limit results from the client",
	},
		[]string{"type"},
	)
)

func initMetricsForAllTypes() {
	for _, t := range allTypes {
		cRequests.WithLabelValue(typeName(t)).Add(0)
		replies.WithLabelValue(typeName(t)).Add(0)
		errors.WithLabelValue(typeName(t)).Add(0)
		sendErrors.WithLabelValue(typeName(t)).Add(0)
		limit.WithLabelValue(typeName(t)).Add(0)
	}
}
