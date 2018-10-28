package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func endpoint(cfg Config, hardwareAddress string, localAPI local.API, errors chan<- error) {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))
	router.Handler("GET", "/local/wifi", instrumentHandler("local_wifi", localWifiHandler(localAPI, errors)))
	router.Handler("GET", "/local/devicelist", instrumentHandler("local_devicelist", localDeviceListHandler(localAPI, errors)))
	router.Handler("GET", "/local/meter", instrumentHandler("local_meter", localMeterHandler(hardwareAddress, localAPI, errors)))

	err := http.ListenAndServe(cfg.Address, router)
	errors <- err
}

func localMeterHandler(address string, api local.API, errors chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		details, err := api.DeviceDetails(r.Context(), address)
		if err != nil {
			writeError(w, err, errors, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, errors, details)
	}
}

func localDeviceListHandler(api local.API, errors chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dl, err := api.DeviceList(r.Context())
		if err != nil {
			writeError(w, err, errors, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, errors, dl)
	}

}

func localWifiHandler(api local.API, errors chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wifi, err := api.WifiStatus(r.Context())
		if err != nil {
			writeError(w, err, errors, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, errors, wifi)
	}
}

func jsonResponse(w http.ResponseWriter, errors chan<- error, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		writeError(w, err, errors, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		writeError(w, err, errors, http.StatusInternalServerError)
		return
	}

	return
}

func writeError(w http.ResponseWriter, err error, errors chan<- error, code int) {
	errors <- err
	http.Error(w, err.Error(), code)
}
