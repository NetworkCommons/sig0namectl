// from https://miek.nl/2014/august/16/go-dns-package/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

var updateCmd = &cli.Command{
	Name:      "update",
	Aliases:   []string{"u"},
	UsageText: "See flags for usage",
	Subcommands: []*cli.Command{
		{
			Name:      "a",
			UsageText: "update a <host> <ip>",
			Usage:     "update A record for <host> in <zone>",
			Action:    updateAAction,
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "unset", Value: true, Usage: "unset the A record, before adding one"},
				&cli.BoolFlag{Name: "remove", Usage: "remove the A record"},
			},
		},

		{
			Name:   "rr",
			Usage:  "update rr",
			Action: updateRRAction,
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "remove", Usage: "remove RRs"},
			},
		},
	},
}

func updateAAction(cCtx *cli.Context) error {
	host := cCtx.Args().Get(0)
	if host == "" {
		return cli.Exit("host required", 1)
	}
	if !strings.HasSuffix(host, ".") {
		host = host + "."
	}
	ipAddrStr := cCtx.Args().Get(1)
	if ipAddrStr == "" {
		return cli.Exit("IP address required", 1)
	}
	if ipAddr := net.ParseIP(ipAddrStr); ipAddr == nil {
		return cli.Exit("invalid IP address: "+ipAddrStr, 1)
	}

	keystore := cCtx.String("keystore")

	keys, err := sig0.ListKeysFiltered(keystore, host)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return cli.Exit("no key found for host", 1)
	}

	log.Println("-- Using key:", keys[0].Name)
	// ugh.. what? doubley .key.key
	keyPath := filepath.Join(keystore, keys[0].Name)
	keyPath = strings.TrimSuffix(keyPath, ".key")
	signer, err := sig0.LoadKeyFile(keyPath)
	if err != nil {
		return err
	}

	soa, err := sig0.QuerySOA(host)
	if err != nil {
		return err
	}
	reply, err := sig0.SendDOHQuery(sig0.DefaultDOHResolver, soa)
	if err != nil {
		return err
	}
	soaReply, err := sig0.AnySOA(reply)
	if err != nil {
		return err
	}
	zone := soaReply.Hdr.Name
	log.Println("-- SOA lookup for", host, "found zone:", zone)

	dohServer, err := sig0.FindDOHEndpoint(zone)
	if err != nil {
		return err
	}

	err = signer.StartUpdate(zone)
	if err != nil {
		return err
	}

	parsedIP := net.ParseIP(ipAddrStr)
	if parsedIP.To4() == nil {
		return fmt.Errorf("invalid IPv4 address: %s", ipAddrStr)
	}

	rrStr := fmt.Sprintf("%s.%s %d IN A %s", strings.TrimSuffix(host, "."+zone), zone, sig0.DefaultTTL, ipAddrStr)
	rr, err := dns.NewRR(rrStr)
	if err != nil {
		return err
	}

	if cCtx.Bool("remove") {
		err = signer.RemoveRR(rr)
	} else {
		if cCtx.Bool("unset") {
			// query current ip
			query, err := sig0.QueryA(host)
			if err != nil {
				return err
			}
			reply, err := sig0.SendDOHQuery(dohServer.Host, query)
			if err != nil {
				return err
			}
			log.Println("unsetting A record for", host)
			for _, rr := range reply.Answer {
				rrA, ok := rr.(*dns.A)
				if !ok {
					continue
				}
				log.Println("current A record", rrA.A.String())

				err = signer.RemoveRR(rrA)
				if err != nil {
					return fmt.Errorf("failed to unset A record: %w", err)
				}
			}
		}
		err = signer.UpdateRR(rr)
	}
	if err != nil {
		return err
	}

	m, err := signer.SignUpdate()
	if err != nil {
		return err
	}

	log.Println("-- Configure DoH client --")
	respMsg, err := sig0.SendDOHQuery(dohServer.Host, m)
	if err != nil {
		return err
	}

	log.Println("-- Response from DNS server --")
	fmt.Println(respMsg)

	return nil
}

func updateRRAction(cCtx *cli.Context) error {
	host := cCtx.Args().Get(0)
	if host == "" {
		return cli.Exit("host required", 1)
	}
	if !strings.HasSuffix(host, ".") {
		host = host + "."
	}
	keystore := cCtx.String("keystore")

	keys, err := sig0.ListKeysFiltered(keystore, host)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return cli.Exit("no key found for host", 1)
	}

	log.Println("-- Using key:", keys[0].Name)
	// ugh.. what? doubley .key.key
	keyPath := filepath.Join(keystore, keys[0].Name)
	keyPath = strings.TrimSuffix(keyPath, ".key")
	signer, err := sig0.LoadKeyFile(keyPath)
	if err != nil {
		return err
	}

	soa, err := sig0.QuerySOA(host)
	if err != nil {
		return err
	}
	reply, err := sig0.SendDOHQuery(sig0.DefaultDOHResolver, soa)
	if err != nil {
		return err
	}
	soaReply, err := sig0.AnySOA(reply)
	if err != nil {
		return err
	}
	zone := soaReply.Hdr.Name
	log.Println("-- SOA lookup for", host, "found zone:", zone)

	dohServer, err := sig0.FindDOHEndpoint(zone)
	if err != nil {
		return err
	}

	err = signer.StartUpdate(zone)
	if err != nil {
		return err
	}

	// read RRs from stdin
	count := 0
	log.Println("-- Reading RRs from stdin --")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rrStr := scanner.Text()
		if rrStr == "" {
			log.Println("-- Empty line, stopping --")
			break
		}
		rr, err := dns.NewRR(rrStr)
		if err != nil {
			return err
		}
		if cCtx.Bool("remove") {
			log.Println("-- Removing RR --")
			err = signer.RemoveRR(rr)
		} else {
			log.Println("-- Updating RR --")
			err = signer.UpdateRR(rr)
		}
		log.Printf("-- %+v", rr)
		count++
		if err != nil {
			return err
		}
	}

	if count == 0 {
		return cli.Exit("no RRs to update", 0)
	}
	log.Printf("-- Updated %d RRs --", count)

	m, err := signer.SignUpdate()
	if err != nil {
		return err
	}

	log.Println("-- Configure DoH client --")
	respMsg, err := sig0.SendDOHQuery(dohServer.Host, m)
	if err != nil {
		return err
	}

	log.Println("-- Response from DNS server --")
	fmt.Println(respMsg)

	return nil
}
