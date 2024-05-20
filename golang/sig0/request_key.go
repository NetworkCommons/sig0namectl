package sig0

import (
	"errors"
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

var (
	SignalSubzonePrefix = "_signal"
	DefaultTTL          = 300
	DefaultDOHResolver  = "dns.quad9.net"
)

func CreateRequestKeyMsg(subZone, zone string) (*dns.Msg, string, error) {

	// Determine the zone master using the provided sub zone and base zone
	signalZone := fmt.Sprintf("%s.%s", SignalSubzonePrefix, zone)
	querySOAForSignal, err := QuerySOA(signalZone)
	if err != nil {
		return nil, "", fmt.Errorf("Error: ZONE %s SOA record does not resolve: %w", signalZone, err)
	}

	soaAnswer, err := SendDOHQuery(DefaultDOHResolver, querySOAForSignal)
	if err != nil {
		return nil, "", fmt.Errorf("Error: DOH query failed for %s: %w", DefaultDOHResolver, err)
	}

	zoneSoa, err := ExpectSOA(soaAnswer)
	if err != nil {
		return nil, "", fmt.Errorf("Error: SOA record not found in response for %s: %w", signalZone, err)
	}

	if zoneSoa != "ns1.free2air.org." {
		return nil, "", fmt.Errorf("Unexpected SOA: %s - TODO: Query SVCB to get the zone master's DOH endpoint", zoneSoa)
	}
	var dohUpdateHost = "doh.zenr.io"

	// Check if zone already exists
	newSubZone := fmt.Sprintf("%s.%s", subZone, zone)
	err = checkZoneDoesntExist(dohUpdateHost, newSubZone)
	if err != nil {
		return nil, "", err
	}

	zoneRequest := fmt.Sprintf("%s.%s.%s", subZone, SignalSubzonePrefix, zone)
	err = checkZoneDoesntExist(dohUpdateHost, zoneRequest)
	if err != nil {
		return nil, "", err
	}

	// craft RRs and create signed update
	subZoneSigner, err := LoadOrGenerateKey(newSubZone)
	if err != nil {
		return nil, "", err
	}

	err = subZoneSigner.StartUpdate(zone)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	nsupdateItemPtr := fmt.Sprintf("%s %d IN PTR %s", signalZone, DefaultTTL, zoneRequest)
	err = subZoneSigner.UpdateParsedRR(nsupdateItemPtr)
	if err != nil {
		return nil, "", err
	}

	updateMsg, err := subZoneSigner.UnsignedUpdate(signalZone)
	if err != nil {
		return nil, "", err
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
