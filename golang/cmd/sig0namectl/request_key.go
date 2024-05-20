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
	Usage:   "requestKey <subZone> <zone>",
	Aliases: []string{"rk"},
	Action:  requestKeyAction,
}

func requestKeyAction(cCtx *cli.Context) error {
	newSubZone := cCtx.Args().Get(0)
	zone := cCtx.Args().Get(1)
	if newSubZone == "" || zone == "" {
		return cli.Exit("subZone and zone are required", 1)
	}

	reqMsg, dohServer, err := sig0.CreateRequestKeyMsg(newSubZone, zone)
	if err != nil {
		return fmt.Errorf("Failed to create request key message: %w", err)
	}
	log.Println("Requesting key for", newSubZone, "under", zone, "from", dohServer)
	spew.Dump(reqMsg)

	answer, err := sig0.SendDOHQuery(dohServer, reqMsg)
	if err != nil {
		return fmt.Errorf("Failed to send DOH query: %w", err)
	}

	if answer.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Update failed: %v", answer)
	}

	return nil
}
