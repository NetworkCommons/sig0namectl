package sig0

import (
	"encoding/base64"
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

func ExpectSOA(answer *dns.Msg) (string, error) {
	if len(answer.Answer) < 1 {
		return "", fmt.Errorf("expected at least one authority section.")
	}
	firstNS := answer.Answer[0]
	soa, ok := firstNS.(*dns.SOA)
	if !ok {
		return "", fmt.Errorf("expected SOA but got type of RR: %T: %+v", firstNS, firstNS)
	}
	return soa.Ns, nil
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
