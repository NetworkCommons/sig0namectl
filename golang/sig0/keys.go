package sig0

import (
	"crypto"
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

var DefaultKeyTTL uint32 = 60

type storedKeyData struct {
	Name    string
	Key     string
	Private string
}

func (data storedKeyData) ParseKey() (*dns.KEY, error) {
	rr, err := dns.NewRR(data.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to read RR from key data: %w", err)
	}

	dnsKey, ok := rr.(*dns.KEY)
	if !ok {
		return nil, fmt.Errorf("expected dns.KEY, instead: %T", rr)
	}
	return dnsKey, nil
}

// wasm<>js bridge helper
func (data storedKeyData) AsMap() map[string]any {
	return map[string]any{
		"Key":  data.Key,
		"Name": data.Name,
		// TODO: could call ParseKey() and add header fields as needed
	}
}

type Signer struct {
	Key     *dns.KEY
	private crypto.PrivateKey

	update *dns.Msg
}

func (s Signer) KeyName() string {
	zone := s.Key.Hdr.Name
	return fmt.Sprintf("K%s+%03d+%d", zone, s.Key.Algorithm, s.Key.KeyTag())
}

func LoadOrGenerateKey(zone string) (*Signer, error) {
	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}

	keys, err := ListKeys(".")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if strings.HasPrefix(key.Name, "K"+zone) {
			signer, err := LoadKeyFile(key.Name)
			if err == nil {
				return signer, nil
			}
		}
	}

	log.Println("generating new key for", zone)
	return GenerateKeyAndSave(zone)

}

// GenerateKey creates a new ED25519 key for the given zone
func GenerateKey(zone string) (*Signer, error) {
	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}

	var k = new(dns.KEY)
	k.Hdr.Name = zone
	k.Hdr.Class = dns.ClassINET
	k.Hdr.Rrtype = dns.TypeKEY
	k.Algorithm = dns.ED25519
	k.Hdr.Ttl = DefaultKeyTTL
	// TODO: find consts for these magic numbers
	k.Flags = 0x200
	k.Protocol = 3

	const keySize = 256
	pk, err := k.Generate(keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return &Signer{
		Key:     k,
		private: pk,
	}, nil
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

	return &Signer{Key: dnsKey, private: privkey}, nil
}
