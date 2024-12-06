package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/urfave/cli/v2"
)

var keysCmd = &cli.Command{
	Name:  "keys",
	Usage: "print sig0 signing key",
	Subcommands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "lists all keys in the keystore",
			Action: listKeysAction,
		},
		{
			Name:   "get",
			Usage:  "search keystore for a key for <my.new.name>",
			Action: getKeyAction,
		},
		{
			Name:   "request",
			Usage:  "request <my.new.name>",
			Action: requestKeyAction,
		},
	},
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
		return fmt.Errorf("failed to create request key: %w", err)
	}

	return nil
}

func getKeyAction(cCtx *cli.Context) error {
	searchDomain := cCtx.Args().First()
	if searchDomain == "" {
		return cli.Exit("searchDomain required", 1)
	}
	keystore := cCtx.String("keystore")
	log.Printf("keystore: %s", keystore)
	log.Printf("searchDomain: %s", searchDomain)
	keys, err := sig0.ListKeysFiltered(keystore, searchDomain)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return cli.Exit("no keys found", 1)
	}
	return printKeys(keys)
}

func listKeysAction(cCtx *cli.Context) error {
	keystore := cCtx.String("keystore")
	log.Printf("keystore: %s", keystore)
	keys, err := sig0.ListKeys(keystore)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return cli.Exit("no keys found", 1)
	}
	return printKeys(keys)
}

func printKeys(keys []sig0.StoredKeyData) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(tw, "Index\tName\tKey Name\tPublic Key\n")
	for i, k := range keys {
		parsed, err := k.ParseKey()
		if err != nil {
			return err
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", i, k.Name, parsed.Hdr.Name, parsed.PublicKey)
	}
	err := tw.Flush()
	if err != nil {
		return err
	}
	return nil
}
