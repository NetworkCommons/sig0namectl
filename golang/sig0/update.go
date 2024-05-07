package sig0

import (
	"crypto"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

func (signer *Signer) UpdateA(host, zone, addr string) (*dns.Msg, error) {
	parsedIP := net.ParseIP(addr)
	if parsedIP.To4() == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %s", addr)
	}
	// TODO make RR generic, for now A record for..?
	myRR := fmt.Sprintf("%s.%s 600 IN A %s", host, zone, addr)

	log.Println("-- Set dns.Msg Structure --")
	m := new(dns.Msg)
	m.SetUpdate(dns.Fqdn(zone))

	log.Println("-- Attach RR to dns.Msg --")
	rrInsert, err := dns.NewRR(myRR)
	if err != nil {
		return nil, err
	}

	m.Insert([]dns.RR{rrInsert})

	// create & fill SIG structure (see sig0_test.go for guidance)
	log.Println("-- Create, fill & attach SIG RR to dns.Msg Structure --")
	now := uint32(time.Now().Unix())
	sig0RR := new(dns.SIG)
	sig0RR.Hdr.Name = "."
	sig0RR.Hdr.Rrtype = dns.TypeSIG
	sig0RR.Hdr.Class = dns.ClassANY
	sig0RR.Algorithm = signer.Key.Algorithm
	sig0RR.Expiration = now + 300
	sig0RR.Inception = now - 300
	sig0RR.KeyTag = signer.Key.KeyTag()
	sig0RR.SignerName = signer.Key.Hdr.Name

	mb, err := sig0RR.Sign(signer.private.(crypto.Signer), m)
	if err != nil {
		algstr := dns.AlgorithmToString[signer.Key.Algorithm]
		return nil, fmt.Errorf("failed to sign %v message: %w", algstr, err)
	}

	if err := m.Unpack(mb); err != nil {
		return nil, fmt.Errorf("failed to unpack message: %w", err)
	}

	return m, nil
}
