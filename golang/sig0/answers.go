package sig0

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

func ParseBase64Answer(answer string) (*dns.Msg, error) {
	data, err := base64.StdEncoding.DecodeString(answer)
	if err != nil {
		return nil, fmt.Errorf("parse: invalid base64: %w", err)
	}
	var resp = new(dns.Msg)
	err = resp.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("parse: invalid dns data: %w", err)
	}
	return resp, nil
}

func ExpectSOA(answer *dns.Msg) (*dns.SOA, error) {
	if len(answer.Answer) < 1 {
		return nil, fmt.Errorf("expected at least one answer section.")
	}
	firstNS := answer.Answer[0]
	soa, ok := firstNS.(*dns.SOA)
	if !ok {
		return nil, fmt.Errorf("expected SOA but got type of RR: %T: %+v", firstNS, firstNS)
	}
	return soa, nil
}

func ExpectAdditonalSOA(answer *dns.Msg) (*dns.SOA, error) {
	if len(answer.Ns) < 1 {
		return nil, fmt.Errorf("expected at least one authority section.")
	}
	firstNS := answer.Ns[0]
	soa, ok := firstNS.(*dns.SOA)
	if !ok {
		return nil, fmt.Errorf("expected SOA but got type of RR: %T: %+v", firstNS, firstNS)
	}
	return soa, nil
}

func AnySOA(answer *dns.Msg) (*dns.SOA, error) {
	if soa, err := ExpectSOA(answer); err == nil {
		return soa, nil
	}
	if soa, err := ExpectAdditonalSOA(answer); err == nil {
		return soa, nil
	}
	return nil, errors.New("no SOA in either answer or additional")
}
