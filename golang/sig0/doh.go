package sig0

import (
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
