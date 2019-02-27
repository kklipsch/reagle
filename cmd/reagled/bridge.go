package main

import (
	"context"

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

		//hardware address for smart meter
		hardwareAddress string

		//prom documentation makes it seem like you have to return the same metrcis
		//caching previous to return in case of error
		previousValues metricValues
	}
)

func newPrometheusBridge(ctx context.Context, reg prometheus.Registerer, hardwareAddress string) (rainForestBridge, error) {
	bridge := rainForestBridge{
		contextFactory:  func() context.Context { return ctx },
		hardwareAddress: hardwareAddress,
	}

	err := reg.Register(bridge)
	return bridge, err
}

//we have constant metrics so this can delegate to library code
func (bridge rainForestBridge) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(bridge, ch)
}

//collect makes the call to the api and converts into prometheus metrics
func (bridge rainForestBridge) Collect(ch chan<- prometheus.Metric) {
	/*
		ctx := bridge.contextFactory()

		metricValues, err := getMetricValues(ctx, api, bridge.hardwareAddress)
		if err != nil {
			errorsCount.Inc()
			log.Errorf("error getting metrics:%v", err)
			collectValues(ctx, ch, bridge.previousValues)
		} else {
			collectValues(ctx, ch, bridge.previousValues)
			bridge.previousValues = metricValues
		}
	*/
}

func collectValues(ctx context.Context, ch chan<- prometheus.Metric, values metricValues) {
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
