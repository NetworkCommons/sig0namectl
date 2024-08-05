package sig0

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func QueryKEY(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeKEY)
}

func QueryPTR(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypePTR)
}

// uses ANY query type
func QueryAny(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeANY)
}

func QueryNSEC(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeNSEC)
}

func QueryRRSIG(name string) (*dns.Msg, error) {
	return QueryWithType(name, dns.TypeRRSIG)
}

func QueryWithType(name string, qtype uint16) (*dns.Msg, error) {
	q := dns.Question{
		Name:   dns.Fqdn(name),
		Qtype:  qtype,
		Qclass: dns.ClassINET,
	}

	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true, AuthenticatedData: true},
		Question: []dns.Question{q},
	}
	m.SetEdns0(4096, true)

	if os.Getenv("DEBUG") != "" {
		fmt.Println("DNS Query:")
		spew.Dump(m)
	}

	return m, nil
}

func QueryWithStringType(name, qtype string) (*dns.Msg, error) {
	t, err := QueryTypeFromString(qtype)
	if err != nil {
		return nil, err
	}
	return QueryWithType(name, t)
}

func QueryTypeFromString(value string) (uint16, error) {
	var t uint16
	switch strings.ToLower(value) {
	case "a":
		t = dns.TypeA
	case "aaaa":
		t = dns.TypeAAAA
	case "any":
		t = dns.TypeANY
	case "key":
		t = dns.TypeKEY
	case "ptr":
		t = dns.TypePTR
	case "loc":
		t = dns.TypeLOC
	case "txt":
		t = dns.TypeTXT
	case "svcb":
		t = dns.TypeSVCB
	case "srv":
		t = dns.TypeSRV
	case "soa":
		t = dns.TypeSOA
	case "nsec":
		t = dns.TypeNSEC
	case "rrsig":
		t = dns.TypeRRSIG
	default:
		asNum, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return 0, fmt.Errorf("unhandled dns type: %q: %w", value, err)
		}
		t = uint16(asNum)
	}
	return t, nil
}
