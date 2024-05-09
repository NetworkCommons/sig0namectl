package sig0

import (
	"crypto"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

func (signer *Signer) StartUpdate(zone string) error {
	if signer.update != nil {
		return fmt.Errorf("update already in progress")
	}
	log.Println("-- Set dns.Msg Structure --")
	m := new(dns.Msg)
	m.SetUpdate(dns.Fqdn(zone))

	signer.update = m
	return nil
}

func (signer *Signer) SignUpdate() (*dns.Msg, error) {
	if signer.update == nil {
		return nil, fmt.Errorf("no update in progress")
	}

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

	mb, err := sig0RR.Sign(signer.private.(crypto.Signer), signer.update)
	if err != nil {
		algstr := dns.AlgorithmToString[signer.Key.Algorithm]
		return nil, fmt.Errorf("failed to sign %v message: %w", algstr, err)
	}

	// unpack SIG into update message
	if err := signer.update.Unpack(mb); err != nil {
		return nil, fmt.Errorf("failed to unpack message: %w", err)
	}

	m := signer.update
	signer.update = nil
	return m, nil
}

func (signer *Signer) UpdateRR(rr string) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	log.Println("-- Attach RR to dns.Msg --")
	rrInsert, err := dns.NewRR(rr)
	if err != nil {
		return err
	}

	signer.update.Insert([]dns.RR{rrInsert})
	return nil
}

// UpdateA is a convenience function to update an A record.
// Need to call StartUpdate first, then UpdateA for each record to update, then SignUpdate.
func (signer *Signer) UpdateA(host, zone, addr string) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	parsedIP := net.ParseIP(addr)
	if parsedIP.To4() == nil {
		return fmt.Errorf("invalid IPv4 address: %s", addr)
	}

	myRR := fmt.Sprintf("%s.%s 600 IN A %s", host, zone, addr)
	return signer.UpdateRR(myRR)
}
