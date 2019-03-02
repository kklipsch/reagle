package client

import "context"

type (
	requestType int

	//Request is used to send information to the client
	Request struct {
		typ            requestType
		payload        interface{}
		resultsPromise chan interface{}
	}
)

const (
	localSpecificVariable requestType = iota
	localAllVariables
	localMeterDetails
	localDeviceList
	localWifiStatus
	localBaseMetrics
)

//RequestSpecificVariable is a Request to do a device query for the provided variable name on the smart meter
func RequestSpecificVariable(variable string) Request {
	return request(localSpecificVariable, variable)
}

//RequestAllVariables is a Request to do a device query for all available variables on the smart meter
func RequestAllVariables() Request {
	return request(localAllVariables)
}

//RequestMeterDetails is a Request to do a device details on the smart meter
func RequestMeterDetails() Request {
	return request(localMeterDetails)
}

//RequestDeviceList is a Request to do a device list for all devices
func RequestDeviceList() Request {
	return request(localDeviceList)
}

//RequestWifiStatus is a Request to do a wifi status call on the Eagle
func RequestWifiStatus() Request {
	return request(localWifiStatus)
}

//RequestBaseMetrics is a Request to do a device query for a BaseMetrics
func RequestBaseMetrics() Request {
	return request(localBaseMetrics)
}

func awaitResult(ctx context.Context, r Request) (interface{}, error) {
	select {
	case result, ok := <-r.resultsPromise:
		if !ok {
			//should never happen
			panic("results promise closed!")
		}

		err, is := result.(error)
		if is {
			return nil, err
		}

		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func sendResult(ctx context.Context, r Request, result interface{}, err error) {
	toSend := result
	if err != nil {
		toSend = err
	}

	select {
	case r.resultsPromise <- toSend:
	case <-ctx.Done():
	}
}

func request(typ requestType, payload ...interface{}) Request {
	var p interface{}
	if len(payload) > 0 {
		p = payload[0]
	}

	return Request{
		typ:            typ,
		payload:        p,
		resultsPromise: make(chan interface{}),
	}
}
