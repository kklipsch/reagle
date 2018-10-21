package main

import (
	"context"

	cli "gopkg.in/urfave/cli.v1"
)

type Config struct {
	Address string
}

func configure(ctx context.Context, cliCtx *cli.Context) (Config, error) {
	cfg := Config{
		Address: cliCtx.String(AddressFlag.Name),
	}

	return cfg, nil
}
