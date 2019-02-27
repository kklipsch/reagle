package main

import (
	"context"
	"time"

	"github.com/kklipsch/reagle/local"
	cli "gopkg.in/urfave/cli.v1"
)

type Config struct {
	Address     string        `json:"address"`
	Wait        time.Duration `json:"wait"`
	LocalConfig local.Config
}

func configure(ctx context.Context, cliCtx *cli.Context) (Config, error) {
	cfg := Config{
		Address: cliCtx.String(addressFlag.Name),
		Wait:    cliCtx.Duration(waitFlag.Name),
	}

	localCfg := local.Config{
		Location:         cliCtx.String(locationFlag.Name),
		User:             cliCtx.String(userFlag.Name),
		ModelIDForMeter:  cliCtx.String(modelIDFlag.Name),
		ImprovedFirmware: cliCtx.Bool(improvedFirmwareFlag.Name),
		DebugRequest:     cliCtx.Bool(debugRequestFlag.Name),
		DebugResponse:    cliCtx.Bool(debugResponseFlag.Name),
	}

	cfg.LocalConfig = local.SetPassword(localCfg, cliCtx.String(passwordFlag.Name))

	err := local.ValidateConfig(cfg.LocalConfig)
	return cfg, err
}
