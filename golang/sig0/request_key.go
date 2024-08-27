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
	//DefaultDOHResolver  = "8.8.8.8"
	// DefaultDOHResolver  = "1.1.1.1"
	// DefaultDOHResolver  = "quad9.zenr.io"
	DefaultDOHResolver = "google.zenr.io"
	// DefaultDOHResolver  = "doh.zenr.io"
)

func RequestKey(newName string) error {
	log.Println("[requestKey] query SOA for", newName)
	query, err := QuerySOA(newName)
	if err != nil {
		return fmt.Errorf("ZONE %s SOA record does not resolve: %w", newName, err)
	}

	var answer *dns.Msg
	answer, err = SendDOHQuery(DefaultDOHResolver, query)
	if err != nil {
		// TODO: add context to err
		return err
	}

	soaForZone, err := AnySOA(answer)
	if err != nil {
		return fmt.Errorf("SOA record not found in response for %s: %w", newName, err)
	}

	zoneOfName := soaForZone.Hdr.Name
	log.Printf("[requestKey] Found zone for new name: %s", zoneOfName)

	newNameFQDN := newName
	if !strings.HasSuffix(newNameFQDN, ".") {
		newNameFQDN += "."
	}

	if !strings.HasSuffix(newNameFQDN, zoneOfName) {
		err = fmt.Errorf("expected new zone to be under it's SOA. Instead got SOA %q for %q", zoneOfName, newNameFQDN)
		return err
	}
	subDomain := strings.TrimSuffix(newNameFQDN, zoneOfName)

	// Determine the zone master using the provided sub zone and base zone
	signalZone := fmt.Sprintf("%s.%s", SignalSubzonePrefix, zoneOfName)
	querySOAForSignal, err := QuerySOA(signalZone)
	if err != nil {
		err = fmt.Errorf("ZONE %s SOA record does not resolve: %w", signalZone, err)
		return err
	}

	answer, err = SendDOHQuery(DefaultDOHResolver, querySOAForSignal)
	if err != nil {
		// TODO: add context to err
		return err
	}

	signalZoneSoa, err := ExpectSOA(answer)
	if err != nil {
		err = fmt.Errorf("SOA record not found in response for %s: %w", signalZone, err)
		return err
	}

	if !strings.HasSuffix(signalZoneSoa.Hdr.Name, soaForZone.Hdr.Name) {
		err = fmt.Errorf("expected signal zone to be under requested zonet got %q and %q", signalZoneSoa.Hdr.Name, soaForZone.Hdr.Name)
		return err
	}

	// get DoH endpoint for signal zone via SVCB
	dohEndpoint, err := FindDOHEndpoint(signalZone)
	if err != nil {
		err = fmt.Errorf("unable to lookup DOH endoint for signal zone: %w", err)
		return err
	}
	log.Printf("[requestKey] Found DOH endoint: %s", dohEndpoint)

	// Check if requested name already exists
	existQuery, err := QueryAny(newName)
	if err != nil {
		return err
	}

	answer, err = SendDOHQuery(dohEndpoint.Host, existQuery)
	if err != nil {
		// TODO: add context to err
		return err
	}

	if answer.Rcode != dns.RcodeNameError {
		return fmt.Errorf("new zone %s already exists: %v", newName, answer)
	}

	log.Printf("[requestKey/debug] %s is not yet taken", newName)

	// check request doesnt exist
	zoneRequest := fmt.Sprintf("%s%s.%s", subDomain, SignalSubzonePrefix, zoneOfName)
	existQuery, err = QueryAny(zoneRequest)
	if err != nil {
		return fmt.Errorf("exists query for zoneRequest %q failed: %w", zoneRequest, err)
	}

	answer, err = SendDOHQuery(dohEndpoint.Host, existQuery)
	if err != nil {
		// TODO: add context to err
		return err
	}

	if answer.Rcode != dns.RcodeNameError {
		return fmt.Errorf("existing zoneRequest for %q already exists: %v", zoneRequest, answer)
	}

	// craft RRs and create signed update
	nameSigner, err := LoadOrGenerateKey(newName)
	if err != nil {
		return err
	}

	log.Printf("[requestKey/debug] creating request with key %s", nameSigner.Key.String())

	err = nameSigner.StartUpdate(zoneOfName)
	if err != nil {
		return fmt.Errorf("unable to start update for zone: %q: %w", zoneOfName, err)
	}

	// Here we split the key details
	// turn it into an RR and split of the first 3 fields
	// so that we can re-use the key for a different zone
	// TODO: there should be a way to get the keyData without full stringification
	keyDetails := strings.TrimSpace(nameSigner.Key.String())
	keyFields := strings.Fields(keyDetails)
	if len(keyFields) < 6 {
		return errors.New("invalid key data")
	}
	keyData := strings.Join(keyFields[3:], " ")

	nsupdateItemSig0Key := fmt.Sprintf("%s %d %s", zoneRequest, DefaultTTL, keyData)
	err = nameSigner.UpdateParsedRR(nsupdateItemSig0Key)
	if err != nil {
		return fmt.Errorf("failed to add KEY RR: %w", err)
	}

	nsupdateItemPtr := fmt.Sprintf("%s %d IN PTR %s", signalZone, DefaultTTL, zoneRequest)
	err = nameSigner.UpdateParsedRR(nsupdateItemPtr)
	if err != nil {
		return fmt.Errorf("failed to add PTR RR: %w", err)
	}

	updateMsg, err := nameSigner.UnsignedUpdate(signalZone)
	if err != nil {
		return fmt.Errorf("unable to create update message: %w", err)
	}

	answer, err = SendDOHQuery(dohEndpoint.Host, updateMsg)
	if err != nil {
		return fmt.Errorf("unable to send update: %w", err)
	}

	if answer == nil {
		return fmt.Errorf("answer is nil")

	}

	if answer.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("update failed: %v", answer)
	}

	return nil
}

type KeyStatus struct {
	QueuePTRExists, KeyRRExists bool
}

func CheckKeyStatus(keyName, zone, dohServer string) (*KeyStatus, error) {
	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}
	key, err := LoadKeyFile(keyName)
	if err != nil {
		return nil, fmt.Errorf("load KeyFile failed: %w", err)
	}
	keyFqdn := key.Key.Hdr.Name // key DNS name (FQDN with trailing dot)

	var ks KeyStatus

	// construct query for KEY RRSet at FQDN keyname
	// TODO BUG cannot yet pass RData via QueryKEY() for exact RR
	// as SendDOHQuery errors with dns: bad rdata
	msgKey, err := QueryKEY(keyFqdn)
	if err != nil {
		return nil, err
	}

	answerKeyRR, err := SendDOHQuery(dohServer, msgKey)
	if err != nil {
		return nil, err
	}

	if answerKeyRR.Rcode != dns.RcodeSuccess && answerKeyRR.Rcode != dns.RcodeNameError {
		err = fmt.Errorf("did not get KEY RR success answer\n:%#v", answerKeyRR)
		return nil, err
	}

	var answerKey *dns.KEY
	for i, rr := range answerKeyRR.Answer {
		var ok bool
		answerKey, ok = rr.(*dns.KEY)
		if ok {
			break
		}
		log.Printf("[DEBUG] answer[%d] is not a KEY type: %T", i, rr)
	}

	if answerKey != nil &&
		answerKey.Flags == key.Key.Flags &&
		answerKey.Protocol == key.Key.Protocol &&
		answerKey.Algorithm == key.Key.Algorithm &&
		answerKey.PublicKey == key.Key.PublicKey {
		ks.KeyRRExists = true
	}
	// query for submission queue PTR at _signal.zone and submission queue KEY under ._signal.zone
	signalPtrRRName := SignalSubzonePrefix + "." + zone
	signalKeyRRName := strings.TrimSuffix(keyFqdn, "."+zone) + "." +
		SignalSubzonePrefix + "." + zone

	// construct query for _signal.zone PTR RRset
	msgSigPtr, err := QueryPTR(signalPtrRRName)
	if err != nil {
		return nil, err
	}

	// send & search query results for Queued PTR RRs under _signal
	answerSignalPtr, err := SendDOHQuery(dohServer, msgSigPtr)
	if err != nil {
		return nil, err
	}

	if answerSignalPtr.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("did not get PTR RR success answer\n:%#v", answerSignalPtr)
		return nil, err
	}

	var ptrRR *dns.PTR
	for i, rr := range answerSignalPtr.Answer {
		var ok bool
		ptrRR, ok = rr.(*dns.PTR)
		if ok {
			break
		}
		log.Printf("[DEBUG] answer[%d] is not a PTR type: %T", i, rr)
	}

	if ptrRR != nil && ptrRR.Ptr == signalKeyRRName {
		ks.QueuePTRExists = true
	}

	return &ks, nil
}
