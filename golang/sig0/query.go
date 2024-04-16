package sig0

import (
	"encoding/base64"

	"github.com/miekg/dns"
)

// QueryA returns a base64 encoded string of a DNS Question for an A record of the passed domain name
func QueryA(name string) string {
	q := dns.Question{
		Name:   dns.Fqdn(name),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}

	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}

	out, err := m.Pack()
	check(err)

	return base64.URLEncoding.EncodeToString(out)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
