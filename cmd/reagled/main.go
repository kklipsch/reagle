package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kklipsch/reagle/local"
	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	applicationLogger = log.WithFields(log.Fields{"name": "reagled"})

	AddressFlag = cli.StringFlag{
		Name:   "address",
		Usage:  "where to serve the endpoints",
		EnvVar: "REAGLED_ADDRESS",
		Value:  ":9000",
	}

	LocationFlag = cli.StringFlag{
		Name:   "location",
		Usage:  "eagle address",
		EnvVar: local.LocationEnv,
	}

	UserFlag = cli.StringFlag{
		Name:   "user",
		Usage:  "eagle user",
		EnvVar: local.UserEnv,
	}

	PasswordFlag = cli.StringFlag{
		Name:   "password",
		Usage:  "eagle password",
		EnvVar: local.PasswordEnv,
	}

	flags = []cli.Flag{
		AddressFlag,
		LocationFlag,
		UserFlag,
		PasswordFlag,
	}
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	app := cli.NewApp()
	app.Name = "reagled"
	app.Usage = "bridge to Rainforest Automation Eagle 200"
	app.Flags = flags
	app.Action = start

	err := app.Run(os.Args)
	if err != nil {
		applicationLogger.WithFields(log.Fields{"error": err}).Fatalf("error during run")
	}
}

func start(cliCtx *cli.Context) error {
	applicationLogger.Infoln("starting up")

	ctx := setSignalCancel(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config, err := configure(ctx, cliCtx)
	if err != nil {
		cancel()
		err = fmt.Errorf("error configuring: %v", err)
		return cli.NewExitError(err, 2)
	}

	applicationLogger.WithFields(log.Fields{"config": config}).Infoln("configured")

	localAPI := local.New(config.LocalConfig)
	localAPI.Client.Transport, err = instrumentClient("local", localAPI.Client.Transport)
	if err != nil {
		err = fmt.Errorf("error instrumenting local client: %v", err)
		return cli.NewExitError(err, 3)
	}

	hardwareAddress, err := getMeterAddress(ctx, localAPI)
	if err != nil {
		applicationLogger.WithFields(log.Fields{"error": err}).Errorln("error getting hardware address")
	}

	errors := make(chan error, 1)
	go endpoint(config, hardwareAddress, localAPI, errors)
	go dataGatherer(ctx, config, errors)

	err = nil
	for errors != nil {
		select {
		case <-ctx.Done():
			//this will cause the main data gatherer to stop and send a nil down the error channel
			//if it doesn't this go routine will send an error down the error channel
			go errorAfter(time.Second*5, errors)
		case e := <-errors:
			if e != nil {
				err = e
				cancel()
			}

			errors = nil
		}
	}

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	applicationLogger.Infoln("done")
	return nil
}

func errorAfter(t time.Duration, errors chan error) {
	timeout, cleanup := context.WithTimeout(context.Background(), t)
	defer cleanup()
	<-timeout.Done()
	errors <- fmt.Errorf("program did not stop within %v of context finish", t)
}

func setSignalCancel(ctx context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sig...)

	go func() {
		<-sigChan
		applicationLogger.WithFields(log.Fields{"signal": sig}).Println("received stop signal")
		cancel()
	}()

	return ctx
}

func getMeterAddress(ctx context.Context, api local.API) (string, error) {
	devices, err := api.DeviceList(ctx)
	if err != nil {
		return "", err
	}

	search := "smart_meter"
	var models []string
	for _, device := range devices {
		if device.ModelID == search {
			return device.HardwareAddress, nil
		}

		models = append(models, device.ModelID)
	}

	return "", fmt.Errorf("no %v found in device list: %v", search, models)
}
