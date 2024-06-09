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

func CreateRequestKeyMsg(newName string) (*dns.Msg, string, error) {
	querySOAForNewZone, err := QuerySOA(newName)
	if err != nil {
		return nil, "", fmt.Errorf("Error: ZONE %s SOA record does not resolve: %w", newName, err)
	}

	newZoneSOAAnswer, err := SendDOHQuery(DefaultDOHResolver, querySOAForNewZone)
	if err != nil {
		return nil, "", fmt.Errorf("Error: DOH query failed for %s: %w", DefaultDOHResolver, err)
	}

	soaForZone, err := AnySOA(newZoneSOAAnswer)
	if err != nil {
		return nil, "", fmt.Errorf("Error: SOA record not found in response for %s: %w", newName, err)
	}

	zoneOfName := soaForZone.Hdr.Name
	log.Printf("[requestKey] Found zone for new name: %s", zoneOfName)

	newNameFQDN := newName
	if !strings.HasSuffix(newNameFQDN, ".") {
		newNameFQDN += "."
	}

	if !strings.HasSuffix(newNameFQDN, zoneOfName) {
		return nil, "", fmt.Errorf("Error: expected new zone to be under it's SOA. Instead got SOA %q for %q", zoneOfName, newNameFQDN)
	}
	subDomain := strings.TrimSuffix(newNameFQDN, zoneOfName)

	// Determine the zone master using the provided sub zone and base zone
	signalZone := fmt.Sprintf("%s.%s", SignalSubzonePrefix, zoneOfName)
	querySOAForSignal, err := QuerySOA(signalZone)
	if err != nil {
		return nil, "", fmt.Errorf("Error: ZONE %s SOA record does not resolve: %w", signalZone, err)
	}

	soaAnswer, err := SendDOHQuery(DefaultDOHResolver, querySOAForSignal)
	if err != nil {
		return nil, "", fmt.Errorf("Error: DOH query failed for %s: %w", DefaultDOHResolver, err)
	}

	signalZoneSoa, err := ExpectSOA(soaAnswer)
	if err != nil {
		return nil, "", fmt.Errorf("Error: SOA record not found in response for %s: %w", signalZone, err)
	}

	if !strings.HasSuffix(signalZoneSoa.Hdr.Name, soaForZone.Hdr.Name) {
		return nil, "", fmt.Errorf("Expected signal zone to be under requested zonet got %q and %q", signalZoneSoa.Hdr.Name, soaForZone.Hdr.Name)
	}

	if signalZoneSoa.Ns != "ns1.free2air.org." {
		return nil, "", fmt.Errorf("Unexpected SOA: %s - TODO: Query SVCB to get the zone master's DOH endpoint", zoneOfName)
	}
	var dohUpdateHost = "doh.zenr.io"

	// Check if zone already exists
	err = checkZoneDoesntExist(dohUpdateHost, newName)
	if err != nil {
		return nil, "", fmt.Errorf("exists check for new name %q failed: %w", newName, err)
	}

	zoneRequest := fmt.Sprintf("%s%s.%s", subDomain, SignalSubzonePrefix, zoneOfName)
	err = checkZoneDoesntExist(dohUpdateHost, zoneRequest)
	if err != nil {
		return nil, "", fmt.Errorf("exists check for zoneRequest %q failed: %w", zoneRequest, err)
	}

	// craft RRs and create signed update
	subZoneSigner, err := LoadOrGenerateKey(newName)
	if err != nil {
		return nil, "", err
	}

	err = subZoneSigner.StartUpdate(zoneOfName)
	if err != nil {
		return nil, "", fmt.Errorf("unable to start update for zone: %q: %w", zoneOfName, err)
	}

	// Here we split the key details
	// turn it into an RR and split of the first 3 fields
	// so that we can re-use the key for a different zone
	keyDetails := strings.TrimSpace(subZoneSigner.Key.String())
	keyFields := strings.Fields(keyDetails)
	if len(keyFields) < 6 {
		return nil, "", errors.New("Invalid key data")
	}
	keyData := strings.Join(keyFields[3:], " ")

	nsupdateItemSig0Key := fmt.Sprintf("%s %d %s", zoneRequest, DefaultTTL, keyData)
	err = subZoneSigner.UpdateParsedRR(nsupdateItemSig0Key)
	if err != nil {
		return nil, "", fmt.Errorf("failed to add KEY RR: %w", err)
	}

	nsupdateItemPtr := fmt.Sprintf("%s %d IN PTR %s", signalZone, DefaultTTL, zoneRequest)
	err = subZoneSigner.UpdateParsedRR(nsupdateItemPtr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to add PTR RR: %w", err)
	}

	updateMsg, err := subZoneSigner.UnsignedUpdate(signalZone)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create update message")
	}

	return updateMsg, dohUpdateHost, nil
}

func checkZoneDoesntExist(dohServer, zone string) error {
	doesExistQuery, err := QueryAny(zone)
	if err != nil {
		return err
	}

	doesExistAnswer, err := SendDOHQuery(dohServer, doesExistQuery)
	if err != nil {
		return err
	}

	if doesExistAnswer.Rcode != dns.RcodeNameError {
		return fmt.Errorf("new zone %s already exists: %v", zone, doesExistAnswer)
	}

	return nil
}
