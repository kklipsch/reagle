package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kklipsch/reagle/local"
)

type (
	metricValues struct {
		Demand    float64 `json:"demand"`
		Delivered float64 `json:"delivered"`
		Price     float64 `json:"price"`
		Currency  string  `json:"currency"`
	}
)

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

	values.Demand, err = getValueFloat("zigbee:InstantaneousDemand", component)
	if err != nil {
		return values, err
	}

	values.Delivered, err = getValueFloat("zigbee:CurrentSummationDelivered", component)
	if err != nil {
		return values, err
	}

	values.Price, err = getValueFloat("zigbee:Price", component)
	if err != nil {
		return values, err
	}

	values.Currency, err = getValue("zigbee:Currency", component)
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
