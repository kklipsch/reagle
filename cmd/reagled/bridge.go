package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	instantDemand    = prometheus.NewDesc("instantaneous_demand", "current demand", nil, nil)
	currentDelivered = prometheus.NewDesc("current_summation_delivered", "total provided", nil, nil)
	price            = prometheus.NewDesc("price", "price as provided by the meter", []string{"currency"}, nil)
)

type (
	metricValues struct {
		demand    float64
		delivered float64
		price     float64
		currency  string
	}

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

func getMetricValues(ctx context.Context, localAPI local.API, hardwareAddress string) (metricValues, error) {
	values := metricValues{}

	response, err := localAPI.DeviceQuery(ctx, hardwareAddress, "zigbee:InstantaneousDemand", "zigbee:CurrentSummationDelivered", "zigbee:Price", "zigbee:Currency")
	if err != nil {
		return values, fmt.Errorf("call to api failed: %v", err)
	}

	variables := local.ResultsFromDetailsResponse(response)
	if len(variables) != 1 {
		return values, fmt.Errorf("variables has more components than expected: %v", variables)
	}

	var component map[string]local.Variable
	for _, component = range variables {
		break
	}

	values.demand, err = getValueFloat("zigbee:InstantaneousDemand", component)
	if err != nil {
		return values, err
	}

	values.delivered, err = getValueFloat("zigbee:CurrentSummationReceived", component)
	if err != nil {
		return values, err
	}

	values.price, err = getValueFloat("zigbee:Price", component)
	if err != nil {
		return values, err
	}

	values.currency, err = getValue("zigbee:Currency", component)
	return values, err
}

func getValue(name string, v map[string]local.Variable) (string, error) {
	variable, ok := v[name]
	if !ok {
		return "", fmt.Errorf("%v does not exist: %v", name, v)
	}

	return variable.Value, nil
}

func getValueFloat(name string, v map[string]local.Variable) (float64, error) {
	value, err := getValue(name, v)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}

func collectValues(ctx context.Context, ch chan<- prometheus.Metric, values metricValues) {
	send(ctx, ch, prometheus.MustNewConstMetric(
		instantDemand,
		prometheus.GaugeValue,
		values.demand,
	))

	send(ctx, ch, prometheus.MustNewConstMetric(
		currentDelivered,
		prometheus.GaugeValue,
		values.delivered,
	))

	send(ctx, ch, prometheus.MustNewConstMetric(
		price,
		prometheus.GaugeValue,
		values.price,
		values.currency,
	))
}

func send(ctx context.Context, ch chan<- prometheus.Metric, metric prometheus.Metric) {
	select {
	case ch <- metric:
	case <-ctx.Done():
	}
}
