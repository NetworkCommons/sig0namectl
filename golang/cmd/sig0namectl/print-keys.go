package main

import (
	"log"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v2"
)

var printKeyCmd = &cli.Command{
	Name:    "print-key",
	Aliases: []string{"pk"},
	Usage:   "add a task to the list",
	Action:  printKeyAction,
}

func printKeyAction(cCtx *cli.Context) error {
	fname := cCtx.String("key-name")

	signer, err := sig0.LoadKeyFile(fname)
	if err != nil {
		return err
	}

	spew.Dump(signer.Key)
	log.Println("OK")

	return nil
}
