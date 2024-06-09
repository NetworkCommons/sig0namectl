// -build test_network

package sig0

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerySOA(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	const srv = "9.9.9.9"

	zones := []struct {
		zone, soa string
	}{

		{"cryptix.pizza", "ns1.dnsowl.com."},
		{"_signal.zenr.io", "ns1.free2air.org."},
		{"twitter.com", "a.u06.twtrdns.net."},
		//{"github.com", "ns-1622.awsdns-10.co.uk."},
		{"wikipedia.org", "ns0.wikimedia.org."},
	}
	for _, testdata := range zones {
		qry, err := QuerySOA(testdata.zone)
		r.NoError(err, testdata)

		answer, err := SendDOHQuery(srv, qry)
		r.NoError(err, testdata)

		soa, err := ExpectSOA(answer)
		r.NoError(err, testdata)
		a.Equal(testdata.soa, soa.Ns)

		verifyication, err := QueryWithType(testdata.zone, dns.TypeSOA)
		r.NoError(err, testdata)

		verifyAnswer, err := SendUDPQuery(soa.Ns, verifyication)
		r.NoError(err, testdata)
		r.True(verifyAnswer.Authoritative, verifyAnswer)
	}
}
