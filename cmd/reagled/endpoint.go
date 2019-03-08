package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kklipsch/reagle/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func endpoint(c client.Local) http.Handler {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))
	router.Handler("GET", "/local/wifi", instrumentHandler("local_wifi", clientHandler(c, wifiStatus)))
	router.Handler("GET", "/local/devicelist", instrumentHandler("local_devicelist", clientHandler(c, deviceList)))
	router.Handler("GET", "/local/meter", instrumentHandler("local_meter", clientHandler(c, meterDetails)))
	router.Handler("GET", "/local/variable/:variable", instrumentHandler("variable", clientHandler(c, specificVariable, getVariableFromURL)))
	router.Handler("GET", "/local/variable/", instrumentHandler("variable", clientHandler(c, allVariables)))
	router.Handler("GET", "/local/metrics/", instrumentHandler("variable", clientHandler(c, baseMetrics)))

	return router
}

func wifiStatus(_ interface{}) client.Request   { return client.RequestWifiStatus() }
func deviceList(_ interface{}) client.Request   { return client.RequestDeviceList() }
func meterDetails(_ interface{}) client.Request { return client.RequestMeterDetails() }
func specificVariable(payload interface{}) client.Request {
	return client.RequestSpecificVariable(payload.(string))
}
func allVariables(_ interface{}) client.Request { return client.RequestAllVariables() }
func baseMetrics(_ interface{}) client.Request  { return client.RequestBaseMetrics() }

type payloadFromRequest func(r *http.Request) (interface{}, error)

func clientHandler(c client.Local, req func(interface{}) client.Request, getPayload ...payloadFromRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		var payload interface{}
		if len(getPayload) > 0 {
			payload, err = getPayload[0](r)
			if err != nil {
				writeError(w, fmt.Errorf("unable to get variable: %v", err), http.StatusInternalServerError)
				return
			}
		}

		timeout, clean := context.WithTimeout(r.Context(), time.Second*5)
		defer clean()
		r = r.WithContext(timeout)

		response, err := c.Request(r.Context(), req(payload))
		switch err {
		case nil:
			jsonResponse(w, response)
		case client.ErrRateLimited:
			writeError(w, err, http.StatusServiceUnavailable)
		case context.DeadlineExceeded:
			writeError(w, err, http.StatusServiceUnavailable)
		default:
			writeError(w, err, http.StatusInternalServerError)
		}
	}
}

func getVariableFromURL(r *http.Request) (interface{}, error) {
	ps := httprouter.ParamsFromContext(r.Context())
	if ps == nil {
		return nil, fmt.Errorf("no params in context")
	}

	variable := strings.ToLower(strings.TrimSpace(ps.ByName("variable")))
	if variable == "" {
		return nil, fmt.Errorf("empty variable")
	}

	return variable, nil
}

func jsonResponse(w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	return
}

func writeError(w http.ResponseWriter, err error, code int) {
	instrumentError(err, "endpoint error")
	http.Error(w, err.Error(), code)
}
