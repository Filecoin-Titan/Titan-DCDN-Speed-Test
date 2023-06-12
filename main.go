package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

const defaultLocatorAddress = "https://locator.titannet.io:5000"

func main() {
	app := &cli.App{
		Name:  "titan-dcdn-speed-test",
		Usage: "titan network toolset",
		Commands: []*cli.Command{
			downloadFileCmd,
			speedTestCmd,
			runCmd,
		},
	}

	if os.Getenv("LOCATOR_API_INFO") == "" {
		os.Setenv("LOCATOR_API_INFO", defaultLocatorAddress)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
