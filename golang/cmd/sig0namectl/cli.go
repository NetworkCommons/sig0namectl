package main

import (
	"fmt"
	"os"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  os.Args[0],
		Usage: "secure dynamic DNS tool",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "keystore", Aliases: []string{"ks"}, Value: "keystore", Usage: "path to keystore", EnvVars: []string{"SIG0_KEYSTORE"}},
			&cli.StringFlag{Name: "resolver", Aliases: []string{"r"}, Value: sig0.DefaultDOHResolver, Usage: "DoH resolver used for bootstrap lookups", EnvVars: []string{"SIG0_DOH_RESOLVER"}},
		},
		Before: func(cCtx *cli.Context) error {
			sig0.DefaultDOHResolver = cCtx.String("resolver")
			return nil
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
