package sig0

import (
	"os"

	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
)

// QueryA returns a base64 encoded string of a DNS Question for an A record of the passed domain name
func QuerySOA(zone string) (*dns.Msg, error) {
	return QueryWithType(zone, dns.TypeSOA)
}

func QueryA(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeA)
}

// uses ANY query type
func QueryAny(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeANY)
}

func QueryWithType(name string, qtype uint16) (*dns.Msg, error) {
	q := dns.Question{
		Name:   dns.Fqdn(name),
		Qtype:  qtype,
		Qclass: dns.ClassINET,
	}

	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}

	if os.Getenv("DEBUG") != "" {
		fmt.Println("DNS Query:")
		spew.Dump(m)
	}

	return m, nil
}
