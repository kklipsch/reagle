package main

import (
	"context"
	"time"

	"github.com/kklipsch/reagle/local"
	cli "gopkg.in/urfave/cli.v1"
)

type Config struct {
	Address        string        `json:"address"`
	MetricSchedule time.Duration `json:"metric_schedule"`
	LocalConfig    local.Config
}

func configure(ctx context.Context, cliCtx *cli.Context) (Config, error) {
	cfg := Config{
		Address:        cliCtx.String(AddressFlag.Name),
		MetricSchedule: cliCtx.Duration(MetricScheduleFlag.Name),
	}

	localCfg := local.Config{
		Location:         cliCtx.String(LocationFlag.Name),
		User:             cliCtx.String(UserFlag.Name),
		ModelIDForMeter:  cliCtx.String(ModelIDFlag.Name),
		ImprovedFirmware: cliCtx.Bool(ImprovedFirmwareFlag.Name),
		DebugRequest:     cliCtx.Bool(DebugRequestFlag.Name),
		DebugResponse:    cliCtx.Bool(DebugResponseFlag.Name),
	}

	cfg.LocalConfig = local.SetPassword(localCfg, cliCtx.String(PasswordFlag.Name))

	err := local.ValidateConfig(cfg.LocalConfig)
	return cfg, err
}
