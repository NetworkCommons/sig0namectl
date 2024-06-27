package sig0

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	r := require.New(t)
	signer := createAndLoadKey(t)

	err := signer.StartUpdate("zone")
	r.NoError(err)

	err = signer.UpdateA("host", "zone", "1.2.3.4")
	r.NoError(err)

	signedUpdate, err := signer.SignUpdate()
	r.NoError(err)

	// verify signing
	rr, ok := signedUpdate.Extra[0].(*dns.SIG)
	r.True(ok, "expected SIG RR, instead: %T", signedUpdate.Extra[0])

	mb, err := signedUpdate.Pack()
	r.NoError(err)

	err = rr.Verify(signer.Key, mb)
	r.NoError(err)
}
