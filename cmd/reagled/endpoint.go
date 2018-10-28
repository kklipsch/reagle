package main

import (
	"encoding/json"
	"net/http"

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

	err := http.ListenAndServe(cfg.Address, router)
	fatals <- err
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
