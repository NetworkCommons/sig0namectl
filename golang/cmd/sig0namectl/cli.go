package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  os.Args[0],
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", EnvVars: []string{"SIG0_HOST"}},
			&cli.StringFlag{Name: "zone", Aliases: []string{"z"}, Usage: "the zone you want to update", EnvVars: []string{"SIG0_ZONE"}},
			&cli.StringFlag{Name: "server", Aliases: []string{"srv"},  EnvVars: []string{"SIG0_SERVER"}},
			&cli.StringFlag{Name: "key-name", Aliases: []string{"kn"}, Usage: "Kso.me.na.me.+aaa+bbbbb", EnvVars: []string{"SIG0_SIG0_KEYFILES"}},
		},
		Commands: []*cli.Command{
			queryCmd, printKeyCmd, updateCmd, requestKeyCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
	}
}
