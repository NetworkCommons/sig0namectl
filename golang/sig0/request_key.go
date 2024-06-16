package sig0

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

var (
	SignalSubzonePrefix = "_signal"
	DefaultTTL          = 300
	DefaultDOHResolver  = "dns.quad9.net"
)

type KeyRequest struct {
	newName    string
	soaForZone *dns.SOA
	signalZone string
	subDomain  string

	err   error
	steps []dnsProcess
}

type dnsProcess func(*dns.Msg) *dns.Msg

func NewKeyRequest(newName string) (*KeyRequest, error) {
	kr := &KeyRequest{
		newName: newName,
	}
	kr.steps = []dnsProcess{
		kr.querySOAForNewZone,
		kr.findSOAForNewZone,
		kr.compareSignalSOA,
		kr.checkNewZoneDoesntExist,
		kr.createRequest,
	}

	return kr, nil
}

func (kr *KeyRequest) Next() bool {
	if kr.err != nil {
		return false
	}
	return len(kr.steps) > 0
}

func (kr *KeyRequest) Err() error {
	return kr.err
}

func (kr *KeyRequest) Do(answer *dns.Msg) *dns.Msg {
	has := len(kr.steps) > 0
	if !has {
		if kr.err == nil {
			kr.err = fmt.Errorf("steps out of bounds")
		}
		return nil
	}
	var next dnsProcess
	next, kr.steps = kr.steps[0], kr.steps[1:]
	return next(answer)
}

func (kr *KeyRequest) querySOAForNewZone(_ *dns.Msg) *dns.Msg {
	log.Println("[requestKey] query SOA for", kr.newName)
	query, err := QuerySOA(kr.newName)
	if err != nil {
		kr.err = fmt.Errorf("Error: ZONE %s SOA record does not resolve: %w", kr.newName, err)
		return nil
	}
	return query
}

func (kr *KeyRequest) findSOAForNewZone(answer *dns.Msg) *dns.Msg {
	var err error
	kr.soaForZone, err = AnySOA(answer)
	if err != nil {
		kr.err = fmt.Errorf("Error: SOA record not found in response for %s: %w", kr.newName, err)
		return nil
	}

	zoneOfName := kr.soaForZone.Hdr.Name
	log.Printf("[requestKey] Found zone for new name: %s", zoneOfName)

	newNameFQDN := kr.newName
	if !strings.HasSuffix(newNameFQDN, ".") {
		newNameFQDN += "."
	}

	if !strings.HasSuffix(newNameFQDN, zoneOfName) {
		kr.err = fmt.Errorf("Error: expected new zone to be under it's SOA. Instead got SOA %q for %q", zoneOfName, newNameFQDN)
		return nil
	}
	kr.subDomain = strings.TrimSuffix(newNameFQDN, zoneOfName)

	// Determine the zone master using the provided sub zone and base zone
	kr.signalZone = fmt.Sprintf("%s.%s", SignalSubzonePrefix, zoneOfName)
	querySOAForSignal, err := QuerySOA(kr.signalZone)
	if err != nil {
		kr.err = fmt.Errorf("Error: ZONE %s SOA record does not resolve: %w", kr.signalZone, err)
		return nil
	}

	return querySOAForSignal
}

func (kr *KeyRequest) compareSignalSOA(answer *dns.Msg) *dns.Msg {
	signalZoneSoa, err := ExpectSOA(answer)
	if err != nil {
		kr.err = fmt.Errorf("Error: SOA record not found in response for %s: %w", kr.signalZone, err)
		return nil
	}

	if !strings.HasSuffix(signalZoneSoa.Hdr.Name, kr.soaForZone.Hdr.Name) {
		kr.err = fmt.Errorf("Expected signal zone to be under requested zonet got %q and %q", signalZoneSoa.Hdr.Name, kr.soaForZone.Hdr.Name)
		return nil
	}

	if signalZoneSoa.Ns != "ns1.free2air.org." {
		zoneOfName := kr.soaForZone.Hdr.Name
		kr.err = fmt.Errorf("Unexpected SOA: %s - TODO: Query SVCB to get the zone master's DOH endpoint", zoneOfName)
		return nil
	}

	// Check if requested name already exists
	existQuery, err := QueryAny(kr.newName)
	if err != nil {
		kr.err = err
		return nil
	}
	return existQuery
}

func (kr *KeyRequest) checkNewZoneDoesntExist(answer *dns.Msg) *dns.Msg {
	if answer.Rcode != dns.RcodeNameError {
		kr.err = fmt.Errorf("new zone %s already exists: %v", kr.newName, answer)
		return nil
	}

	log.Printf("[requestKey/debug] %s is not yet taken", kr.newName)

	// check request doesnt exist
	var err error
	zoneOfName := kr.soaForZone.Hdr.Name
	zoneRequest := fmt.Sprintf("%s%s.%s", kr.subDomain, SignalSubzonePrefix, zoneOfName)
	existQuery, err := QueryAny(zoneRequest)
	if err != nil {
		kr.err = fmt.Errorf("exists query for zoneRequest %q failed: %w", zoneRequest, err)
		return nil
	}
	return existQuery
}

func (kr *KeyRequest) createRequest(answer *dns.Msg) *dns.Msg {
	zoneOfName := kr.soaForZone.Hdr.Name
	zoneRequest := fmt.Sprintf("%s%s.%s", kr.subDomain, SignalSubzonePrefix, zoneOfName)
	if answer.Rcode != dns.RcodeNameError {
		kr.err = fmt.Errorf("existing zoneRequest for %q already exists: %v", zoneRequest, answer)
		return nil
	}

	// craft RRs and create signed update
	nameSigner, err := LoadOrGenerateKey(kr.newName)
	if err != nil {
		kr.err = err
		return nil
	}

	log.Printf("[requestKey/debug] creating request with key %s", nameSigner.Key.String())

	err = nameSigner.StartUpdate(zoneOfName)
	if err != nil {
		kr.err = fmt.Errorf("unable to start update for zone: %q: %w", zoneOfName, err)
		return nil
	}

	// Here we split the key details
	// turn it into an RR and split of the first 3 fields
	// so that we can re-use the key for a different zone
	keyDetails := strings.TrimSpace(nameSigner.Key.String())
	keyFields := strings.Fields(keyDetails)
	if len(keyFields) < 6 {
		kr.err = errors.New("Invalid key data")
		return nil
	}
	keyData := strings.Join(keyFields[3:], " ")

	nsupdateItemSig0Key := fmt.Sprintf("%s %d %s", zoneRequest, DefaultTTL, keyData)
	err = nameSigner.UpdateParsedRR(nsupdateItemSig0Key)
	if err != nil {
		kr.err = fmt.Errorf("failed to add KEY RR: %w", err)
		return nil
	}

	nsupdateItemPtr := fmt.Sprintf("%s %d IN PTR %s", kr.signalZone, DefaultTTL, zoneRequest)
	err = nameSigner.UpdateParsedRR(nsupdateItemPtr)
	if err != nil {
		kr.err = fmt.Errorf("failed to add PTR RR: %w", err)
		return nil
	}

	updateMsg, err := nameSigner.UnsignedUpdate(kr.signalZone)
	if err != nil {
		kr.err = fmt.Errorf("unable to create update message: %w", err)
		return nil
	}

	return updateMsg
}
