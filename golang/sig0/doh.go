package sig0

import (
	"fmt"
	"net/url"
	"strings"

	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/shynome/doh-client"
)

func SendDOHQuery(server string, m *dns.Msg) (*dns.Msg, error) {
	co := &dns.Conn{Conn: doh.NewConn(nil, nil, server)}

	err := co.WriteMsg(m)
	if err != nil {
		return nil, err
	}

	answer, err := co.ReadMsg()
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func SendUDPQuery(server string, m *dns.Msg) (*dns.Msg, error) {
	co, err := dns.Dial("udp4", server+":53")
	if err != nil {
		return nil, err
	}

	err = co.WriteMsg(m)
	if err != nil {
		return nil, err
	}

	answer, err := co.ReadMsg()
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func FindDOHEndpoint(name string) (*url.URL, error) {
	lookup := dns.Fqdn("_dns." + name)

	svcbQry, err := QueryWithType(lookup, dns.TypeSVCB)
	if err != nil {
		return nil, fmt.Errorf("findDOHEndpoint: failed to construct query for %s: %w", lookup, err)
	}

	answer, err := SendDOHQuery(DefaultDOHResolver, svcbQry)
	if err != nil {
		return nil, fmt.Errorf("findDOHEndpoint: no answer for svcb query for %s: %w", lookup, err)
	}
	if len(answer.Answer) < 1 {
		spew.Dump(answer)
		return nil, fmt.Errorf("findDOHEndpoint: no answer section for %s", lookup)
	}

	// TODO deal with more than one SVCB "_dns." + name in RRSet
	first := answer.Answer[0]
	svcb, ok := first.(*dns.SVCB)
	if !ok {
		return nil, fmt.Errorf("findDOHEndpoint: unable to cast %T answer for %s to SVCB", first, lookup)
	}

	if !dns.IsFqdn(svcb.Target) {
		return nil, fmt.Errorf("findDOHEndpoint: Expected SVCB target to be FQDN but got %q", svcb.Target)
	}
	host := svcb.Target[:len(svcb.Target)-1]

	var alpn *dns.SVCBAlpn
	var dohPath *dns.SVCBDoHPath

	for _, v := range svcb.Value {
		switch tv := v.(type) {
		case *dns.SVCBAlpn:
			alpn = tv
		case *dns.SVCBDoHPath:
			dohPath = tv
		default:
			fmt.Fprintf(os.Stderr, "findDOHEndpoint: unexpected keyValue pair type: %T\n", tv)
		}
	}

	if alpn == nil {
		return nil, fmt.Errorf("findDOHEndpoint: expected alpn in svcb.Value")
	}
	if dohPath == nil {
		return nil, fmt.Errorf("findDOHEndpoint: expected dohPath in svcb.Value")
	}
	if dohPath.Template == "" {
		return nil, fmt.Errorf("findDOHEndpoint: expected dohPath Template value")
	}

	// TODO: write this more genericly but striclty speaky we should check for IN HTTPS to check if we can use https
	var hasHttp = false
	for _, a := range alpn.Alpn {
		if strings.Contains(a, "h1") || strings.Contains(a, "h2") {
			hasHttp = true
		}
	}
	if !hasHttp {
		fmt.Printf("alpn: %#v\n", alpn.Alpn[0])
		return nil, fmt.Errorf("findDOHEndpoint: expected http alpn in svcb.Value")
	}

	const placeholder = "{?dns}"
	if !strings.Contains(dohPath.Template, placeholder) {
		return nil, fmt.Errorf("findDOHEndpoint: dohPath does not contain {?dns}: %s", dohPath.Template)
	}

	var u url.URL
	u.Scheme = "https"
	u.Host = host
	u.Path = strings.Replace(dohPath.Template, placeholder, "", -1)
	return &u, nil
}
