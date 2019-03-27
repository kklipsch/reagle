package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kklipsch/reagle/local"
)

func newMediator(api local.API, wait time.Duration) *mediator {
	return &mediator{
		api:     api,
		address: getSmartMeterAddress(NewRateLimit(wait), api), //note this has its own rate limit, so that it doesnt interfere with the other one
		limit:   NewRateLimit(wait),
	}
}

type mediator struct {
	api     local.API
	address smartMeterAddress
	limit   RateLimit
}

func (m *mediator) mediate(ctx context.Context, requests <-chan Request) {
	for {
		select {
		case req, ok := <-requests:
			if !ok {
				panic("request channel closed should not be possible")
			}

			cRequests.WithLabelValues(typeName(req.typ)).Inc()
			result, err := m.request(ctx, req.typ, req.payload)
			sendResult(req, result, err)
		case <-ctx.Done():
			log.Printf("shutting down mediator due to context cancellation")
			return
		}
	}
}

func (m *mediator) request(ctx context.Context, typ requestType, payload interface{}) (interface{}, error) {
	if err := EnforceLimit(ctx, m.limit); err != nil {
		limit.WithLabelValues(typeName(typ)).Inc()

		return nil, err
	}

	return m.query(ctx, typ, payload)
}

func (m *mediator) query(ctx context.Context, typ requestType, payload interface{}) (interface{}, error) {
	address, err := m.getAddress(ctx, typ)
	if err != nil {
		return nil, err
	}

	switch typ {
	case localSpecificVariable:
		variable := payload.(string)
		return m.api.DeviceQuery(ctx, address, variable)
	case localAllVariables:
		details, err := m.api.DeviceDetails(ctx, address)
		if err != nil {
			return nil, err
		}

		variables := local.VariablesFromDetailsResponse(details)
		if len(variables) < 1 {
			return nil, fmt.Errorf("no variables defined")
		}

		return m.api.DeviceQuery(ctx, address, variables...)
	case localMeterDetails:
		return m.api.DeviceDetails(ctx, address)
	case localBaseMetrics:
		return getBaseMetrics(ctx, m.api, address)
	case localDeviceList:
		return m.api.DeviceList(ctx)
	case localWifiStatus:
		return m.api.WifiStatus(ctx)
	}

	panic(fmt.Sprintf("unknown request type: %v", typ))
}

func (m *mediator) getAddress(ctx context.Context, typ requestType) (string, error) {
	switch typ {
	case localWifiStatus, localDeviceList:
		//these query types do not require an address so don't even bothe trying to get it
		return "", nil
	default:
		return m.address(ctx)
	}
}
