package sig0

import (
	"testing"

	"github.com/miekg/dns"
)

func TestUpdate(t *testing.T) {
	signer := createAndLoadKey(t)

	signedUpdate, err := signer.UpdateA("host", "zone", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}

	// verify signing
	rr, ok := signedUpdate.Extra[0].(*dns.SIG)
	if !ok {
		t.Fatalf("expected SIG RR, instead: %T", signedUpdate.Extra[0])
	}

	mb, err := signedUpdate.Pack()
	if err != nil {
		t.Fatal(err)
	}

	err = rr.Verify(signer.dnsKey, mb)
	if err != nil {
		t.Fatal(err)
	}
}
