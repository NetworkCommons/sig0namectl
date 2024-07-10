package sig0

import (
	"crypto"
	"fmt"
	"log"
	"net"
	"strings"
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

func (signer *Signer) UnsignedUpdate(zone string) (*dns.Msg, error) {
	if signer.update == nil {
		return nil, fmt.Errorf("no update in progress")
	}

	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}

	m := signer.update
	m.SetUpdate(zone)
	signer.update = nil
	return m, nil
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

func (signer *Signer) UpdateParsedRR(rr string) error {
	rrInsert, err := dns.NewRR(rr)
	if err != nil {
		return fmt.Errorf("sig0: failed to parse RR: %w", err)
	}

	return signer.UpdateRR(rrInsert)
}

func (signer *Signer) UpdateRR(rr dns.RR) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	signer.update.Insert([]dns.RR{rr})
	return nil
}

func (signer *Signer) RemoveParsedRR(rr string) error {
	rrRemove, err := dns.NewRR(rr)
	if err != nil {
		return fmt.Errorf("sig0: failed to parse RR: %w", err)
	}

	return signer.RemoveRR(rrRemove)
}

func (signer *Signer) RemoveRR(rr dns.RR) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	signer.update.Remove([]dns.RR{rr})
	return nil
}

func (signer *Signer) RemoveParsedRRset(rr string) error {
	rrRemove, err := dns.NewRR(rr)
	if err != nil {
		return fmt.Errorf("sig0: failed to parse RR: %w", err)
	}

	return signer.RemoveRRset(rrRemove)
}

func (signer *Signer) RemoveRRset(rr dns.RR) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	signer.update.RemoveRRset([]dns.RR{rr})
	return nil
}

func (signer *Signer) RemoveParsedName(rr string) error {
	rrRemove, err := dns.NewRR(rr)
	if err != nil {
		return fmt.Errorf("sig0: failed to parse RR: %w", err)
	}

	return signer.RemoveName(rrRemove)
}

func (signer *Signer) RemoveName(rr dns.RR) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	signer.update.RemoveName([]dns.RR{rr})
	return nil
}

// UpdateA is a convenience function to update an A record.
// Need to call StartUpdate first, then UpdateA for each record to update, then SignUpdate.
func (signer *Signer) UpdateA(subZone, zone, addr string) error {
	if signer.update == nil {
		return fmt.Errorf("no update in progress")
	}

	parsedIP := net.ParseIP(addr)
	if parsedIP.To4() == nil {
		return fmt.Errorf("invalid IPv4 address: %s", addr)
	}

	myRR := fmt.Sprintf("%s.%s %d IN A %s", subZone, zone, DefaultTTL, addr)
	return signer.UpdateParsedRR(myRR)
}
