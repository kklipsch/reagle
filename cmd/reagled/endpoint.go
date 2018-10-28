package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func endpoint(cfg Config, localAPI local.API, errors chan<- error) {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))
	router.Handler("GET", "/local/wifi", instrumentHandler("local_wifi", localWifiHandler(cfg, localAPI, errors)))

	err := http.ListenAndServe(cfg.Address, router)
	errors <- err
}

type WifiResponse struct {
	Enabled    string `json:"enabled"`
	Type       string `json:"type"`
	SSID       string `json:"ssid"`
	Encryption string `json:"encryption"`
	Channel    string `json:"channel"`
	IPAddress  string `json:"ip_address"`
}

func localWifiHandler(cfg Config, api local.API, errors chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wifi, err := api.WifiStatus(r.Context())
		if err != nil {
			writeError(w, err, errors, http.StatusInternalServerError)
			return
		}

		jsonResponse(w, errors, wifiResponse(wifi))
	}
}

func wifiResponse(wifi local.WifiStatus) WifiResponse {
	return WifiResponse{
		Enabled:    wifi.Enabled,
		Type:       wifi.Type,
		SSID:       wifi.SSID,
		Encryption: wifi.Encryption,
		Channel:    wifi.Channel,
		IPAddress:  wifi.IPAddress,
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
