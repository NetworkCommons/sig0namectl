package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
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

	pubfh, err := os.Open(fname + ".key")
	if err != nil {
		return err
	}

	dk, err := dns.ReadRR(pubfh, fname+".key")
	if err != nil {
		return err
	}
	spew.Dump(dk)

	privfh, err := os.Open(fname + ".private")
	if err != nil {
		return err
	}
	defer privfh.Close()

	privkey, err := dk.(*dns.KEY).ReadPrivateKey(privfh, fname+".private")
	if err != nil {
		return err
	}
	spew.Dump(privkey)
	log.Println("OK")

	return nil
}
