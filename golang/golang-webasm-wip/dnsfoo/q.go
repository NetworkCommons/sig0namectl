package main

import (
	"encoding/base64"
	"fmt"

	"github.com/miekg/dns"
)

func main() {
	q := dns.Question{
		Name:   dns.Fqdn("cryptix.zenr.io"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}

	out, err := m.Pack()
	check(err)

	fmt.Println("Q:", base64.StdEncoding.EncodeToString(out))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
