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

var (
	allTypes = []requestType{
		localSpecificVariable,
		localAllVariables,
		localMeterDetails,
		localDeviceList,
		localWifiStatus,
		localBaseMetrics,
	}
)

func typeName(t requestType) string {
	switch t {
	case localSpecificVariable:
		return "specific_variable"
	case localAllVariables:
		return "all_variables"
	case localMeterDetails:
		return "meter_details"
	case localDeviceList:
		return "device_list"
	case localWifiStatus:
		return "wifi_status"
	case localBaseMetrics:
		return "base_metrics"
	default:
		return "unknown"
	}
}

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
			awaitErrors.WithLabelValues(typeName(r.typ)).Inc()
			return nil, err
		}

		return result, nil
	case <-ctx.Done():
		awaitCancelled.WithLabelValues(typeName(r.typ)).Inc()
		return nil, ctx.Err()
	}
}

func sendResult(r Request, result interface{}, err error) {
	name := typeName(r.typ)

	toSend := result
	if err != nil {
		toSend = err
		errors.WithLabelValues(name).Inc()
	} else {
		replies.WithLabelValues(name).Inc()
	}

	select {
	case r.resultsPromise <- toSend:
	default:
		sendErrors.WithLabelValues(name).Inc()
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
		resultsPromise: make(chan interface{}, 1),
	}
}
