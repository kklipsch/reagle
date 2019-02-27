package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func endpoint(mediator apiMediator) http.Handler {
	router := httprouter.New()
	router.Handler("GET", "/metrics", instrumentHandler("metrics", promhttp.Handler()))
	router.Handler("GET", "/local/wifi", instrumentHandler("local_wifi", localMediated(mediator, wifiStatus)))
	router.Handler("GET", "/local/devicelist", instrumentHandler("local_devicelist", localMediated(mediator, deviceList)))
	router.Handler("GET", "/local/meter", instrumentHandler("local_meter", localMediated(mediator, meterDetails)))
	router.Handler("GET", "/local/variable/:variable", instrumentHandler("variable", localMediated(mediator, specificVariable, getVariableFromURL)))
	router.Handler("GET", "/local/variable/", instrumentHandler("variable", localMediated(mediator, allVariables)))

	return router
}

type variableFromRequest func(r *http.Request) (string, error)

func localMediated(mediator apiMediator, typ requestType, getVariable ...variableFromRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		variable := ""
		if len(getVariable) > 0 {
			variable, err = getVariable[0](r)
			if err != nil {
				writeError(w, fmt.Errorf("unable to get variable: %v", err), http.StatusInternalServerError)
				return
			}
		}

		timeout, clean := context.WithTimeout(r.Context(), time.Second*5)
		defer clean()
		r = r.WithContext(timeout)

		response, err := mediator.sendReceive(r.Context(), newAPIRequest(typ, variable))
		switch err {
		case nil:
			jsonResponse(w, response)
		case errRateLimited:
			writeError(w, err, http.StatusServiceUnavailable)
		case context.DeadlineExceeded:
			writeError(w, err, http.StatusServiceUnavailable)
		default:
			writeError(w, err, http.StatusInternalServerError)
		}
	}
}

func getVariableFromURL(r *http.Request) (string, error) {
	ps := httprouter.ParamsFromContext(r.Context())
	if ps == nil {
		return "", fmt.Errorf("no params in context")
	}

	variable := strings.ToLower(strings.TrimSpace(ps.ByName("variable")))
	if variable == "" {
		return "", fmt.Errorf("empty variable")
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
	errorsCount.Inc()
	applicationLogger.WithFields(log.Fields{"code": code, "error": err}).Errorln("endpoint error")
	http.Error(w, err.Error(), code)
}
