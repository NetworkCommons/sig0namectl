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
