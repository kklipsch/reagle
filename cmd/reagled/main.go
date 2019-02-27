package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kklipsch/reagle/local"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	applicationLogger = log.WithFields(log.Fields{"name": "reagled"})

	addressFlag = cli.StringFlag{
		Name:   "address",
		Usage:  "where to serve the endpoints",
		EnvVar: "REAGLED_ADDRESS",
		Value:  ":9000",
	}

	waitFlag = cli.DurationFlag{
		Name:   "wait",
		Usage:  "how much time to ensure between calls to the eagle",
		EnvVar: "REAGLED_WAIT",
		Value:  time.Second,
	}

	locationFlag = cli.StringFlag{
		Name:   "location",
		Usage:  "eagle address",
		EnvVar: local.LocationEnv,
	}

	userFlag = cli.StringFlag{
		Name:   "user",
		Usage:  "eagle user",
		EnvVar: local.UserEnv,
	}

	passwordFlag = cli.StringFlag{
		Name:   "password",
		Usage:  "eagle password",
		EnvVar: local.PasswordEnv,
	}

	modelIDFlag = cli.StringFlag{
		Name:   "model_id",
		Usage:  "what the eagle is reporting for your smart meter model id, can be found by hitting the device_list endpoint. Unlikely to need to be set",
		EnvVar: local.MeterModelIDEnv,
	}

	//oddity of the cli parsing library.  a boolt will be set to true by default, so 'setting' this turns it off, giving a mismatch between the name and
	//the cli ergonomics
	improvedFirmwareFlag = cli.BoolTFlag{
		Name:   "unimproved_firmware",
		Usage:  "if your eagle has the unimproved firmware (it responds with invalid xml for multiplier & divisor queries) this should be set",
		EnvVar: local.ImprovedFirmwareEnv,
	}

	debugRequestFlag = cli.BoolFlag{
		Name:   "debug_request",
		Usage:  "if set requests will be debugged",
		EnvVar: local.DebugRequestEnv,
	}

	debugResponseFlag = cli.BoolFlag{
		Name:   "debug_response",
		Usage:  "if set responses will be debugged",
		EnvVar: local.DebugResponseEnv,
	}

	flags = []cli.Flag{
		addressFlag,
		waitFlag,
		locationFlag,
		userFlag,
		passwordFlag,
		modelIDFlag,
		improvedFirmwareFlag,
		debugRequestFlag,
		debugResponseFlag,
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

const (
	configureErrorCode int = iota
	mediatorErrorCode
	bridgeErrorCode
	shutdownErrorCode
)

func start(cliCtx *cli.Context) error {
	applicationLogger.Infoln("starting up")

	ctx := setSignalCancel(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	config, err := configure(ctx, cliCtx)
	if err != nil {
		err = fmt.Errorf("error configuring: %v", err)
		return cli.NewExitError(err, configureErrorCode)
	}

	applicationLogger.WithFields(log.Fields{"config": config}).Infoln("configured")

	mediator, err := startAPIMediator(ctx, config)
	if err != nil {
		err = fmt.Errorf("error starting api mediator: %v", err)
		return cli.NewExitError(err, mediatorErrorCode)
	}

	_, err = newPrometheusBridge(ctx, prometheus.DefaultRegisterer, mediator)
	if err != nil {
		err = fmt.Errorf("error creating prometheus bridge: %v", err)
		return cli.NewExitError(err, bridgeErrorCode)
	}

	srv := startServer(config, mediator)

	applicationLogger.Infoln("started")

	<-ctx.Done()

	shutdownCtx, clean := context.WithTimeout(context.Background(), time.Second*5)
	defer clean()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		err = fmt.Errorf("error shutting down web server: %v", err)
		return cli.NewExitError(err, shutdownErrorCode)
	}

	applicationLogger.Infoln("done")
	return nil
}

func startServer(config Config, mediator apiMediator) *http.Server {
	srv := &http.Server{Addr: config.Address, Handler: endpoint(mediator)}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			applicationLogger.WithFields(log.Fields{"err": err}).Fatalln("failed at serving")
		}
	}()

	return srv
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
