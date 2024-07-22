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

	keyName := signer.KeyName()

	keyData := []byte(signer.Key.String())
	err = os.WriteFile(keyName+".key", keyData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write %q: %w", keyName+".key", err)
	}

	privateData := []byte(signer.Key.PrivateKeyString(signer.private))
	err = os.WriteFile(keyName+".private", privateData, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to write %q: %w", keyName+".private", err)
	}

	return signer, nil
}

func ListKeys(dir string) ([]storedKeyData, error) {
	fh, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", dir, err)
	}
	defer fh.Close()

	names, err := fh.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	var keyfiles []storedKeyData
	for _, name := range names {
		if !(strings.HasPrefix(name, "K") && strings.HasSuffix(name, ".key")) {
			continue
		}
		keyName := strings.TrimSuffix(name, ".key")

		sig, err := LoadKeyFile(keyName)
		if err != nil {
			log.Printf("DEBUG: trying %s failed: %v", keyName, err)
			continue
		}

		keyfiles = append(keyfiles, storedKeyData{
			Name: name,
			Key:  sig.Key.PublicKey,
		})
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

	return &Signer{Key: dnsKey, private: privkey}, nil
}
