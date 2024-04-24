package sig0

import (
	"crypto"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type Signer struct {
	dnsKey  *dns.KEY
	private crypto.PrivateKey
}

func LoadKeyFile(keyfile string) (*Signer, error) {
	var (
		pubKeyName    = keyfile + ".key"
		secretKeyName = keyfile + ".private"
	)

	pubfh, err := os.Open(pubKeyName)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", pubKeyName, err)
	}
	defer pubfh.Close()

	rr, err := dns.ReadRR(pubfh, pubKeyName)
	if err != nil {
		return nil, fmt.Errorf("failed to read RR from %q: %w", pubKeyName, err)
	}

	dnsKey, ok := rr.(*dns.KEY)
	if !ok {
		return nil, fmt.Errorf("expected dns.KEY, instead: %T", rr)
	}

	hdr := rr.Header()
	log.Println(keyfile+".key import:", hdr.Name, hdr.Ttl, hdr.Class, hdr.Rrtype, dnsKey.Flags, dnsKey.Protocol, dnsKey.Algorithm, dnsKey.PublicKey)

	privfh, err := os.Open(secretKeyName)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", secretKeyName, err)
	}
	defer privfh.Close()

	privkey, err := dnsKey.ReadPrivateKey(privfh, secretKeyName)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key material from %q: %w", secretKeyName, err)
	}

	return &Signer{dnsKey, privkey}, nil
}

func ParseKeyData(key, private string) (*Signer, error) {
	rr, err := dns.NewRR(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read RR from key data: %w", err)
	}

	dnsKey, ok := rr.(*dns.KEY)
	if !ok {
		return nil, fmt.Errorf("expected dns.KEY, instead: %T", rr)
	}

	hdr := rr.Header()
	log.Println("key import:", hdr.Name, hdr.Ttl, hdr.Class, hdr.Rrtype, dnsKey.Flags, dnsKey.Protocol, dnsKey.Algorithm, dnsKey.PublicKey)

	privkey, err := dnsKey.ReadPrivateKey(strings.NewReader(private), rr.Header().Name+":private")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key material from private key data: %w", err)
	}

	return &Signer{dnsKey, privkey}, nil
}
