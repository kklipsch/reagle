package main

import (
	"github.com/kklipsch/reagle/local"
)

type apiFactory func() (local.API, error)

func instrumentedAPIFactory(cfg local.Config) apiFactory {
	return func() (local.API, error) {
		localAPI := local.New(cfg)
		localAPI.Client.Transport, err = instrumentClient("local", localAPI.Client.Transport)
		return localAPI, err
	}
}
