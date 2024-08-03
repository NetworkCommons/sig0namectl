package main

import (
	"encoding/json"
	"fmt"
	"os"
	// "strconv"
	// "strings"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

var queryCmd = &cli.Command{
	Name:    "query",
	Usage:   "query <name>",
	Aliases: []string{"q"},
	Flags: []cli.Flag{
		&cli.GenericFlag{
			Name:  "type",
			Usage: "type of query to run",
			Value: &dnsRRTypeFlag{},
		},

		&cli.BoolFlag{
			Name:  "json",
			Usage: "output JSON",
			Value: false,
		},
	},
	Action: queryAction,
}

func queryAction(cCtx *cli.Context) error {
	name := cCtx.Args().First()

	tf := cCtx.Generic("type")
	typeFlag, ok := tf.(*dnsRRTypeFlag)
	if !ok {
		return fmt.Errorf("unexpected flag type: %T", tf)
	}
	// default
	if typeFlag.qtype == 0 {
		typeFlag.qtype = dns.TypeA
	}
	fmt.Fprintf(os.Stderr, "[Querying] %d:%v\n", typeFlag.qtype, name)

	q, err := sig0.QueryWithType(name, typeFlag.qtype)
	if err != nil {
		return err
	}

	server := cCtx.String("server")
	answer, err := sig0.SendDOHQuery(server, q)
	if err != nil {
		return err
	}

	if cCtx.Bool("json") {
		out, err := json.Marshal(answer)
		if err != nil {
			return fmt.Errorf("failed to encode answer to %d:%q as JSON: %w", typeFlag.qtype, name, err)
		}
		_, _ = os.Stdout.Write(out)
	} else {
		spew.Dump(answer)
	}
	return nil
}

type dnsRRTypeFlag struct {
	qtype uint16
}

func (f *dnsRRTypeFlag) Set(value string) error {
	var err error
	f.qtype, err = sig0.QueryTypeFromString(value)
	if err != nil {
		return err
	}
	return nil
}

func (f *dnsRRTypeFlag) String() string {
	return fmt.Sprintf("dnsType:%d", f.qtype)
}
