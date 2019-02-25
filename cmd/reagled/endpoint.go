package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func endpoint(cfg Config, hardwareAddress string, localAPI local.API, fatals chan<- error) {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))
	router.Handler("GET", "/local/wifi", instrumentHandler("local_wifi", localWifiHandler(localAPI)))
	router.Handler("GET", "/local/devicelist", instrumentHandler("local_devicelist", localDeviceListHandler(localAPI)))
	router.Handler("GET", "/local/meter", instrumentHandler("local_meter", localMeterHandler(hardwareAddress, localAPI)))
	router.Handler("GET", "/local/variable/:variable", instrumentHandler("variable", localVariableHandler(hardwareAddress, localAPI)))
	router.Handler("GET", "/local/variable/", instrumentHandler("variable", localAllVariablesHandler(hardwareAddress, localAPI)))

	err := http.ListenAndServe(cfg.Address, router)
	fatals <- err
}

func localAllVariablesHandler(address string, api local.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, err := api.DeviceDetails(r.Context(), address)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		variables := local.VariablesFromDetailsResponse(details)
		if len(variables) < 1 {
			writeError(w, fmt.Errorf("no variables defined"), http.StatusInternalServerError)
			return
		}

		results, err := api.DeviceQuery(r.Context(), address, variables...)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, results)
	}
}

func localVariableHandler(address string, api local.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ps := httprouter.ParamsFromContext(r.Context())
		if ps == nil {
			writeError(w, fmt.Errorf("no params in context"), http.StatusInternalServerError)
			return
		}

		variable := strings.ToLower(strings.TrimSpace(ps.ByName("variable")))
		if variable == "" {
			writeError(w, fmt.Errorf("empty variable"), http.StatusInternalServerError)
			return
		}

		details, err := api.DeviceQuery(r.Context(), address, variable)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, details)
	}
}

func localMeterHandler(address string, api local.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, err := api.DeviceDetails(r.Context(), address)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, details)
	}
}

func localDeviceListHandler(api local.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dl, err := api.DeviceList(r.Context())
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, dl)
	}

}

func localWifiHandler(api local.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wifi, err := api.WifiStatus(r.Context())
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, wifi)
	}
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
	errorsCount.Inc()
	applicationLogger.WithFields(log.Fields{"code": code, "error": err}).Errorln("endpoint error")
	http.Error(w, err.Error(), code)
}
