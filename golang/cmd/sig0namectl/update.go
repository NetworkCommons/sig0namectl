// from https://miek.nl/2014/august/16/go-dns-package/

package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

var updateCmd = &cli.Command{
	Name:      "update",
	Aliases:   []string{"u"},
	UsageText: "See flags for usage",
	Action:    updateAction,
}

func updateAction(cCtx *cli.Context) error {
	ipAddrStr := cCtx.Args().First()
	if ipAddrStr == "" {
		return fmt.Errorf("No IP address defined")
	}

	var (
		err error

		sig0Keyfile string

		zone = cCtx.String("zone")
		host = cCtx.String("host")
	)

	server := cCtx.String("server")
	sig0Keyfile = cCtx.String("key-name")

	if sig0Keyfile == "" {
		return fmt.Errorf("No sig0Keyfile defined")
	}

	log.Println("-- Reading SIG(0) Keyfiles (dnssec-keygen format) --")
	signer, err := sig0.LoadKeyFile(sig0Keyfile)
	if err != nil {
		return err
	}

	err = signer.StartUpdate(zone)
	if err != nil {
		return err
	}

	err = signer.UpdateA(host, zone, ipAddrStr)
	if err != nil {
		return err
	}

	m, err := signer.SignUpdate()
	if err != nil {
		return err
	}
	// spew.Dump(m)

	log.Println("-- Configure DoH client --")
	respMsg, err := sig0.SendDOHQuery(server, m)
	if err != nil {
		return err
	}

	log.Println("-- Response from DNS server --")
	fmt.Println(respMsg)

	return nil
}
