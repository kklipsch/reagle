package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kklipsch/reagle/local"
	log "github.com/sirupsen/logrus"
)

type (
	//we want to a) protect the eagle from getting slammed and b) ensure appropriate concurrency controls
	//so this struct will mediate interactions with the api
	apiMediator struct {
		requests chan<- apiRequest
	}

	requestType int

	//golang why wont you just give me futures/promises
	apiRequest struct {
		variable string
		typ      requestType

		resultPromise chan interface{}
	}
)

const (
	specificVariable requestType = iota
	allVariables
	meterDetails
	deviceList
	wifiStatus
	baseMetrics
)

var errRateLimited = fmt.Errorf("rate limited")

func newAPIRequest(typ requestType, variable string) apiRequest {
	return apiRequest{variable: variable, typ: typ, resultPromise: make(chan interface{})}
}

func startAPIMediator(ctx context.Context, cfg Config) (apiMediator, error) {
	requests := make(chan apiRequest)
	mediator := apiMediator{requests: requests}

	localAPI, err := instrumentedAPI(cfg.LocalConfig)
	if err != nil {
		return mediator, err
	}

	go handleAPICalls(ctx, cfg.Wait, localAPI, requests)

	return mediator, nil
}

func handleAPICalls(ctx context.Context, wait time.Duration, localAPI local.API, requests <-chan apiRequest) {
	var blocked atomicBool

	hardwareAddress, err := localAPI.GetMeterHardwareAddress(ctx)
	if err != nil {
		applicationLogger.WithFields(log.Fields{"err": err}).Errorln("unable to get hardware address for meter")
	}

	for {
		select {
		case req, ok := <-requests:
			if !ok {
				log.Fatal("request channel closed")
				return
			}

			result, err := handleRequest(ctx, wait, blocked, hardwareAddress, localAPI, req)
			if err == errRateLimited {
				log.Info("request was load shed")
			}

			sendResult(ctx, req.resultPromise, result, err)
		case <-ctx.Done():
			return
		}
	}

}

func handleRequest(ctx context.Context, wait time.Duration, blocked atomicBool, hardwareAddress string, localapi local.API, req apiRequest) (interface{}, error) {
	if blocked.Get() {
		return nil, errRateLimited
	}

	gateRequests(ctx, blocked, wait)

	return queryForRequest(ctx, hardwareAddress, localapi, req)
}

func queryForRequest(ctx context.Context, hardwareAddress string, localapi local.API, req apiRequest) (interface{}, error) {
	err := validateHardwareAddress(req.typ, hardwareAddress)
	if err != nil {
		return nil, err
	}

	switch req.typ {
	case specificVariable:
		return localapi.DeviceQuery(ctx, hardwareAddress, req.variable)
	case allVariables:
		details, err := localapi.DeviceDetails(ctx, hardwareAddress)
		if err != nil {
			return nil, err
		}

		variables := local.VariablesFromDetailsResponse(details)
		if len(variables) < 1 {
			return nil, fmt.Errorf("no variables defined")
		}

		return localapi.DeviceQuery(ctx, hardwareAddress, variables...)
	case meterDetails:
		return localapi.DeviceDetails(ctx, hardwareAddress)
	case baseMetrics:
		return getMetricValues(ctx, localapi, hardwareAddress)
	case deviceList:
		return localapi.DeviceList(ctx)
	case wifiStatus:
		return localapi.WifiStatus(ctx)
	default:
		return nil, fmt.Errorf("unknown request type: %v", req.typ)
	}

}

func validateHardwareAddress(typ requestType, address string) error {
	if address != "" {
		return nil
	}

	switch typ {
	case wifiStatus, deviceList:
		return nil
	default:
		return fmt.Errorf("%v:must have hardware address for that query", typ)
	}

}

func gateRequests(ctx context.Context, blocked atomicBool, wait time.Duration) {
	blocked.Set(true)
	go func() {
		select {
		case <-time.After(wait):
		case <-ctx.Done():
		}

		blocked.Set(false)
	}()
}

func sendResult(ctx context.Context, resultPromise chan<- interface{}, result interface{}, err error) {
	toSend := result
	if err != nil {
		toSend = err
	}

	select {
	case resultPromise <- toSend:
	case <-ctx.Done():
	}
}

func (mediator apiMediator) sendReceive(ctx context.Context, request apiRequest) (interface{}, error) {
	err := mediator.send(ctx, request)
	if err != nil {
		return nil, err
	}

	return mediator.receive(ctx, request)
}

func (mediator apiMediator) send(ctx context.Context, request apiRequest) error {
	select {
	case mediator.requests <- request:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (mediator apiMediator) receive(ctx context.Context, request apiRequest) (interface{}, error) {
	select {
	case result, ok := <-request.resultPromise:
		if !ok {
			return nil, fmt.Errorf("result failed")
		}

		err, ok := result.(error)
		if ok {
			return nil, err
		}

		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func instrumentedAPI(cfg local.Config) (local.API, error) {
	var err error
	localAPI := local.New(cfg)
	localAPI.Client.Transport, err = instrumentClient("local", localAPI.Client.Transport)
	return localAPI, err
}
