package main

import (
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"

	"github.com/thebutlah/huggable.us/runner"
)

func main() {

	startCommand := cli.Command{
		Name:  "start",
		Usage: "starts the server",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "phttp, port-http",
				Usage: "sets the port to use for the HTTP listener",
			},
			cli.IntFlag{
				Name:  "phttps, port-https",
				Usage: "sets the port to use for the HTTPS listener",
			},
		},
		Action: func(c *cli.Context) error {
			domains := []string{
				"www.huggable.us",
				"huggable.us",
			}
			options := make([]runner.Option, 0)
			if httpPort := c.Int("phttp"); httpPort != 0 {
				options = append(options, runner.HTTPPort(strconv.Itoa(httpPort)))
			}
			if httpsPort := c.Int("phttps"); httpsPort != 0 {
				options = append(options, runner.HTTPSPort(strconv.Itoa(httpsPort)))
			}
			return runner.Start(domains, options...)
		},
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		startCommand,
	}
	app.Action = startCommand.Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
