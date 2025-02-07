package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  os.Args[0],
		Usage: "secure dynamic DNS tool",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "keystore", Aliases: []string{"ks"}, Value: "keystore", Usage: "path to keystore", EnvVars: []string{"SIG0_KEYSTORE"}},
		},
		Commands: []*cli.Command{
			keysCmd, queryCmd, updateCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
	}
}
