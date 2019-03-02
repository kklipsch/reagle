package client

import (
	"context"

	"github.com/kklipsch/reagle/local"
)

type smartMeterAddress func(context.Context) (string, error)

func getSmartMeterAddress(limit RateLimit, api local.API) smartMeterAddress {
	var address string
	return func(ctx context.Context) (string, error) {
		if address != "" {
			return address, nil
		}

		if err := EnforceLimit(ctx, limit); err != nil {
			return "", nil
		}

		var err error
		address, err = api.GetMeterHardwareAddress(ctx)
		return address, err
	}
}
