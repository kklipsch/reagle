package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kklipsch/reagle/client"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	instantDemand    = prometheus.NewDesc("instantaneous_demand", "current demand", nil, nil)
	currentDelivered = prometheus.NewDesc("current_summation_delivered", "total provided", nil, nil)
	price            = prometheus.NewDesc("price", "price as provided by the meter", []string{"currency"}, nil)
)

type (

	//implements the prometheus collector interface to generate the metrics on demand instead of on a schedule
	rainForestBridge struct {
		//context necessary for api calls but no great way to inject it in prometheus collector interface
		contextFactory func() context.Context

		//prom documentation makes it seem like you have to return the same metrcis
		//caching previous to return in case of error
		previousValues atomic.Value

		//use the standard client that keeps things threadsafe and throttled
		c client.Local
	}
)

func newPrometheusBridge(ctx context.Context, reg prometheus.Registerer, c client.Local) (*rainForestBridge, error) {
	bridge := &rainForestBridge{
		contextFactory: func() context.Context { return ctx },
		c:              c,
	}

	bridge.previousValues.Store(client.BaseMetrics{})

	err := reg.Register(bridge)
	return bridge, err
}

//we have constant metrics so this can delegate to library code
func (bridge *rainForestBridge) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(bridge, ch)
}

//collect makes the call to the api and converts into prometheus metrics
func (bridge *rainForestBridge) Collect(ch chan<- prometheus.Metric) {
	ctx := bridge.contextFactory()

	timeout, clean := context.WithTimeout(ctx, time.Second*5)
	defer clean()

	response, err := bridge.c.Request(timeout, client.RequestBaseMetrics())
	if err != nil {
		instrumentError(err, "unable to get metrics for prometheus bridge")
		response = bridge.previousValues.Load()
	}

	collectValues(ctx, ch, response)
	bridge.previousValues.Store(response)
}

func collectValues(ctx context.Context, ch chan<- prometheus.Metric, response interface{}) {
	//go ahead and panic cause if this isnt a BaseMetrics it means something disastorous has happened
	values := response.(client.BaseMetrics)

	send(ctx, ch, prometheus.MustNewConstMetric(
		instantDemand,
		prometheus.GaugeValue,
		values.Demand,
	))

	send(ctx, ch, prometheus.MustNewConstMetric(
		currentDelivered,
		prometheus.CounterValue,
		values.Delivered,
	))

	send(ctx, ch, prometheus.MustNewConstMetric(
		price,
		prometheus.GaugeValue,
		values.Price,
		values.Currency,
	))
}

func send(ctx context.Context, ch chan<- prometheus.Metric, metric prometheus.Metric) {
	select {
	case ch <- metric:
	case <-ctx.Done():
	}
}
