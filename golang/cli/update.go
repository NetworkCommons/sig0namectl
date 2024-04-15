// from https://miek.nl/2014/august/16/go-dns-package/

package main

import (
	"crypto"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/shynome/doh-client"
	"github.com/urfave/cli/v2"
)

var updateCmd = &cli.Command{
	Name:      "update",
	Aliases:   []string{"u"},
	UsageText: "See flags for usage",
	Action:    updateAction,
}

func updateAction(cCtx *cli.Context) error {
	var sig0Keyfile string

	zone := cCtx.String("zone")
	host := cCtx.String("host")
	server := cCtx.String("server")
	sig0Keyfile = cCtx.String("key-name")

	// TODO make RR generic, for now A record for..?
	myRR := fmt.Sprintf("%s.%s 600 IN A 2.3.4.2", host, zone)

	log.Println("-- Set dns.Msg Structure --")
	m := new(dns.Msg)
	m.SetUpdate(dns.Fqdn(zone))

	log.Println("-- Attach RR to dns.Msg --")
	rrInsert, err := dns.NewRR(myRR)
	if err != nil {
		return err
	}

	// log.Println(spew.Sdump(rrInsert))

	m.Insert([]dns.RR{rrInsert})

	if sig0Keyfile == "" {
		return fmt.Errorf("No sig0Keyfile defined")
	}
	log.Println("-- Reading SIG(0) Keyfiles (dnssec-keygen format) --")
	// log.Printf("GD_SIG0_KEYFILES = %s", sig0Keyfile)
	pubfh, perr := os.Open(sig0Keyfile + ".key")
	if perr != nil {
		return perr
	}
	defer pubfh.Close()

	dk, pkerr := dns.ReadRR(pubfh, sig0Keyfile+".key")
	if pkerr != nil {
		return pkerr
	}

	// TODO: extract alg number more eloquently! :/
	// test := fmt.Sprintln(dk)
	// TODO: how best to get public key to insert into keyRR more elegantly! :/
	// keyFields := strings.Fields(test)
	keyFields := strings.Fields(dk.String())
	keyName := keyFields[0]
	keyTTL := keyFields[1]
	keyClass := keyFields[2]
	keyType := keyFields[3]
	keyFlags := keyFields[4]
	keyVersion := keyFields[5]
	keyAlgorithm := keyFields[6]
	keyPublicKey := keyFields[7]

	keyAlgNum, err := strconv.ParseUint(keyAlgorithm, 10, 8)
	if err != nil {
		return err
	}

	log.Println(sig0Keyfile+".key import:", keyName, keyTTL, keyClass, keyType, keyFlags, keyVersion, keyAlgorithm, keyPublicKey)

	privfh, err := os.Open(sig0Keyfile + ".private")
	if err != nil {
		return err
	}
	defer privfh.Close()

	privkey, err := dk.(*dns.KEY).ReadPrivateKey(privfh, sig0Keyfile+".private")
	if err != nil {
		return err
	}
	log.Println("Private Key OK")

	// // fill KEY structure for keyfiles key see dns_test.go
	// keyRR := &dns.DNSKEY{Flags: 257, Protocol: 3, Algorithm: dns.ED25519}
	// keyRR.Hdr = dns.RR_Header{Name: "vortex.zenr.io.", Rrtype: dns.TypeDNSKEY, Class: dns.ClassINET, Ttl: 3600}
	// // vortex.zenr.io. IN KEY 512 3 15 2MK3KZkUgYQVumU9bhy1KzIZ2FhFQZ8yLP2nFMJRCEQ=

	// create & fill KEY structure (see sig0_test.go for guidance)
	log.Println("-- TODO In progress ... Create and fill KEY structure from dnssec-keygen keyfiles --")
	keyRR := new(dns.KEY)
	keyRR.Hdr.Name = "cryptix.zenr.io." // TODO set to RRset 1st space separated field of dnssec-keygen .key file eg vortex.zenr.io.
	keyRR.Hdr.Rrtype = dns.TypeKEY
	keyRR.Hdr.Class = dns.ClassINET
	keyRR.Hdr.Ttl = 600
	keyRR.Flags = 512 // Take from RR Header
	keyRR.Protocol = 3
	keyRR.Algorithm = uint8(keyAlgNum)
	keyRR.PublicKey = keyPublicKey

	// spew.Dump(keyRR)

	// create & fill SIG structure (see sig0_test.go for guidance)
	log.Println("-- TODO Create, fill & attach SIG RR to dns.Msg Structure --")
	now := uint32(time.Now().Unix())
	sig0RR := new(dns.SIG)
	sig0RR.Hdr.Name = "."
	sig0RR.Hdr.Rrtype = dns.TypeSIG
	sig0RR.Hdr.Class = dns.ClassANY
	sig0RR.Algorithm = uint8(keyAlgNum)
	sig0RR.Expiration = now + 300
	sig0RR.Inception = now - 300
	sig0RR.KeyTag = keyRR.KeyTag()
	sig0RR.SignerName = keyRR.Hdr.Name
	mb, err := sig0RR.Sign(privkey.(crypto.Signer), m)

	algstr := dns.AlgorithmToString[keyRR.Algorithm]
	if err != nil {
		return fmt.Errorf("failed to sign %v message: %v", algstr, err)
	}

	if err := m.Unpack(mb); err != nil {
		return fmt.Errorf("failed to unpack message: %v", err)
	}

	spew.Dump(m)

	// verify signing
	var sigrrwire *dns.SIG
	switch rr := m.Extra[0].(type) {
	case *dns.SIG:
		sigrrwire = rr
	default:
		return fmt.Errorf("expected SIG RR, instead: %v", rr)
	}

	for _, rr := range []*dns.SIG{sig0RR, sigrrwire} {
		id := "sig0RR"
		if rr == sigrrwire {
			id = "sigrrwire"
		}
		if err := rr.Verify(keyRR, mb); err != nil {
			return fmt.Errorf("failed to verify %q signed SIG(%s): %v", algstr, id, err)
		}
	}

	// spew.Dump(sig0RR)

	// log.Println(spew.Sdump(m))

	log.Println("-- Configure client DNS method --")
	co := &dns.Conn{Conn: doh.NewConn(nil, nil, server)}
	if err := co.WriteMsg(m); err != nil {
		panic(err)
	}
	respMsg, err := co.ReadMsg()
	if err != nil {
		panic(err)
	}
	fmt.Println(respMsg)
	return nil
}
