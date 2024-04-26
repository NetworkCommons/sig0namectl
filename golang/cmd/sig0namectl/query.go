package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

var queryCmd = &cli.Command{
	Name:      "query",
	Usage:     "query <name>",
	Aliases:   []string{"q"},
	UsageText: "See flags for usage",
	Action:    queryAction,
}

func queryAction(cCtx *cli.Context) error {
	name := cCtx.Args().First()
	// name := "cryptix.zenr.io"
	server := cCtx.String("server")

	q, err := sig0.QueryA(name)
	if err != nil {
		return err
	}
	fmt.Printf("Q:(TXT):%v\n", name)

	fmt.Println("Q:(B64)", q)

	// send over DoH
	url := fmt.Sprintf("https://%s/dns-query?dns=%s", server, q)
	fmt.Println("Q:(DoH):", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status: %s", resp.Status)
	}

	answerBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dnsAnswer = new(dns.Msg)
	err = dnsAnswer.Unpack(answerBody)
	if err != nil {
		return err
	}

	spew.Dump(dnsAnswer)
	return nil
}
