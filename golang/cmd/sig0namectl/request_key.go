package main

import (
	"fmt"
	"log"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
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
	zr, err := sig0.NewKeyRequest(newName)
	if err != nil {
		return fmt.Errorf("Failed to create request key: %w", err)
	}
	log.Println("Requesting key for", newName, "from", dohServer)

	var answer *dns.Msg
	var i = 0
	for zr.Next() {
		qry := zr.Do(answer)
		if qry == nil {
			break
		}
		spew.Dump(qry)

		answer, err = sig0.SendDOHQuery(dohServer, qry)
		if err != nil {
			return fmt.Errorf("Failed to create request key message: %w", err)
		}

		spew.Dump(answer)

		i++
	}

	err = zr.Err()
	if err != nil {
		return fmt.Errorf("request loop failed: %w", err)
	}

	if answer.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Update failed: %v", answer)
	}

	return nil
}
