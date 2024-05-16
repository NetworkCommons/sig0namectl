package sig0

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"fmt"

	"github.com/miekg/dns"
)

// QueryA returns a base64 encoded string of a DNS Question for an A record of the passed domain name
func QueryA(name string) (string, error) {
	return QueryWithType(name, dns.TypeA)
}

// uses ANY query type
func QueryAny(name string) (string, error) {
	return QueryWithType(name, dns.TypeANY)
}

func QueryWithType(name string, qtype uint16) (string, error) {
	q := dns.Question{
		Name:   dns.Fqdn(name),
		Qtype:  qtype,
		Qclass: dns.ClassINET,
	}

	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}

	out, err := m.Pack()
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(out), nil
}

func QuerySOA(zone string) (string, error) {
	buf := make([]byte, 5)
	rand.Read(buf)
	zone = strings.Repeat(fmt.Sprintf("%x.", buf), 2) + zone
	return QueryWithType(zone, dns.TypeSOA)
}

func ExpectAdditonalSOA(answer *dns.Msg) (string, error) {
	if len(answer.Ns) < 1 {
		return "", fmt.Errorf("expected at least one authority section.")
	}
	firstNS := answer.Ns[0]
	soa, ok := firstNS.(*dns.SOA)
	if !ok {
		return "", fmt.Errorf("expected SOA but got type of RR: %T: %+v", firstNS, firstNS)
	}

	return soa.Ns, nil
}
