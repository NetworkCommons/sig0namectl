package main

import (
	"fmt"
	"log"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/urfave/cli/v2"
)

var requestKeyCmd = &cli.Command{
	Name:    "requestKey",
	Usage:   "requestKey <my.new.name>",
	Aliases: []string{"rk"},
	Action:  requestKeyAction,
}

func requestKeyAction(cCtx *cli.Context) error {
	newName := cCtx.Args().Get(0)
	if newName == "" {
		return cli.Exit("newName required", 1)
	}
	dohServer := "doh.zenr.io"
	log.Println("Requesting key for", newName, "from", dohServer)

	err := sig0.RequestKey(newName)
	if err != nil {
		return fmt.Errorf("Failed to create request key: %w", err)
	}

	return nil
}
