package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"lotus-ctrl/sectors"
	"os"
)

var log = logging.Logger("lotus-ctrl")

func main() {
	level := os.Getenv("LOTUS_LOG_LEVEL")
	if level != "" {
		logging.SetLogLevel("*", level)
	} else {
		logging.SetLogLevel("*", "INFO")
	}

	app := &cli.App{
		Name:    "lotus-ctrl",
		Usage:   "lotus-ctrl",
		Version: "0.0.0",
		Flags:   []cli.Flag{},
		Commands: []*cli.Command{
			sectors.SealCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}
}
