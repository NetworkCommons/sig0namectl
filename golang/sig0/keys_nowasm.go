//go:build !wasm
// +build !wasm

package sig0

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func GenerateKeyAndSave(zone string) (*Signer, error) {
	signer, err := GenerateKey(zone)
	if err != nil {
		return nil, err
	}
	_ = signer
	return nil, fmt.Errorf("TODO: not yet implemented persistence")
}

func ListKeys(dir string) ([]string, error) {
	fh, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", dir, err)
	}
	defer fh.Close()

	names, err := fh.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	var keyfiles []string
	for _, name := range names {
		if strings.HasPrefix(name, "K") && strings.HasSuffix(name, ".key") {
			keyName := strings.TrimSuffix(name, ".key")

			_, err := LoadKeyFile(keyName)
			if err != nil {
				log.Printf("DEBUG: trying %s failed: %v", keyName, err)
				continue
			}

			keyfiles = append(keyfiles, keyName)
		}
	}

	return keyfiles, nil
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
