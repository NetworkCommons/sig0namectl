package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

var loadPrivateKeyCmd = &cli.Command{
	Name:    "load-key",
	Usage:   "load-key <fname.private>",
	Args:    true,
	Aliases: []string{"lk"},
	Action:  loadPrivateKeyAction,
}

func loadPrivateKeyAction(cCtx *cli.Context) error {

	fname := cCtx.Args().First()
	if fname == "" {
		fname = cCtx.String("key-name") + ".private"
	}

	dk := new(dns.DNSKEY)

	fh, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fh.Close()

	privkey, err := dk.ReadPrivateKey(fh, os.Args[1])
	spew.Dump(privkey)
	return err
}
