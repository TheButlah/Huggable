package main

import (
	"fmt"
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
				Usage: "the port to use for the HTTP listener",
			},
			cli.IntFlag{
				Name:  "phttps, port-https",
				Usage: "the port to use for the HTTPS listener",
			},
			cli.StringFlag{
				Name: "cmode, cert-mode",
				Usage: "the mode to use for configuring the TLS/HTTPS certs. " +
					"Can be \"auto\" for automatic certificates signed by a CA, " +
					"\"self\" for self-signed certificates generated on the fly, or " +
					"\"provided\" for certs provided by the user as a file.",
				Value: "auto",
			},
		},
		Action: startAction,
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

func startAction(c *cli.Context) error {
	hosts := []string{
		"www.huggable.us",
		"huggable.us",
	}
	options := make([]runner.Option, 0)
	cmode := c.String("cmode")
	switch cmode {
	case "auto":
		options = append(options, runner.AutomaticTLS(hosts))
	case "self":
		options = append(options, runner.SelfSignedTLS(hosts))
	case "provided":
		return fmt.Errorf("cmode=\"%s\" unimplemented", cmode)
	default:
		return fmt.Errorf("cmode=\"%s\" unknown", cmode)
	}

	if httpPort := c.Int("phttp"); httpPort != 0 {
		options = append(options, runner.HTTPPort(strconv.Itoa(httpPort)))
	}
	if httpsPort := c.Int("phttps"); httpsPort != 0 {
		options = append(options, runner.HTTPSPort(strconv.Itoa(httpsPort)))
	}
	return runner.Start(hosts, options...)
}
