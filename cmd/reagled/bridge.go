package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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

		//use the standard mediator that keeps things threadsafe and throttled
		mediator apiMediator
	}
)

func newPrometheusBridge(ctx context.Context, reg prometheus.Registerer, mediator apiMediator) (*rainForestBridge, error) {
	bridge := &rainForestBridge{
		contextFactory: func() context.Context { return ctx },
		mediator:       mediator,
	}

	bridge.previousValues.Store(metricValues{})

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

	response, err := bridge.mediator.sendReceive(timeout, newAPIRequest(baseMetrics, ""))
	if err != nil {
		errorsCount.Inc()
		applicationLogger.WithFields(log.Fields{"err": err}).Errorln("unable to get metrics for prometheus bridge")
		response = bridge.previousValues.Load()
	}

	collectValues(ctx, ch, response)
	bridge.previousValues.Store(response)
}

func collectValues(ctx context.Context, ch chan<- prometheus.Metric, response interface{}) {
	//go ahead and panic cause if this isnt a metricvalues it means something disastorous has happened
	values := response.(metricValues)

	send(ctx, ch, prometheus.MustNewConstMetric(
		instantDemand,
		prometheus.GaugeValue,
		values.Demand,
	))

	send(ctx, ch, prometheus.MustNewConstMetric(
		currentDelivered,
		prometheus.GaugeValue,
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
