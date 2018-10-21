package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func endpoint(cfg Config, errors chan<- error) {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))

	err := http.ListenAndServe(cfg.Address, router)
	errors <- err
}
