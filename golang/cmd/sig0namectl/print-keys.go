package main

import (
	"fmt"
	"log"
	"os"

	"text/tabwriter"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v2"
)

var printKeyCmd = &cli.Command{
	Name:    "print-key",
	Aliases: []string{"pk"},
	Usage:   "loads the key from the keystore and prints its public part",
	Action:  printKeyAction,
}

func printKeyAction(cCtx *cli.Context) error {
	fname := cCtx.String("key-name")

	signer, err := sig0.LoadKeyFile(fname)
	if err != nil {
		return err
	}

	spew.Dump(signer.Key)

	return nil
}

var listKeyCmd = &cli.Command{
	Name:    "list-keys",
	Aliases: []string{"ls"},
	Usage:   "lists all keys in the store and checks their status",
	Action:  listKeysAction,
}

func listKeysAction(cCtx *cli.Context) error {
	//	fname := cCtx.String("key-name")
	keys, err := sig0.ListKeys(".")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 2, 4, 4, ' ', 0)
	fmt.Fprintf(tw, "Key Name\tExists in DNS\tRequest Queued\n")
	for _, k := range keys {
		// TODO: flags
		stat, err := sig0.CheckKeyStatus(k.Name, "zenr.io", sig0.DefaultDOHResolver)
		if err != nil {
			log.Println("failed to check", k.Name, ":", err)
			continue
		}
		fmt.Fprintf(tw, "%s\t%v\t%v\n", k.Name, stat.KeyRRExists, stat.QueuePTRExists)
	}
	return tw.Flush()
}
