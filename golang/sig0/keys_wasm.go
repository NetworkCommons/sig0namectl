package sig0

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

func GenerateKeyAndSave(zone string) (*Signer, error) {
	signer, err := GenerateKey(zone)
	if err != nil {
		return nil, err
	}

	var persisted StoredKeyData
	persisted.Key = signer.Key.String()
	persisted.Private = signer.Key.PrivateKeyString(signer.private)

	marshalled, err := json.Marshal(persisted)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key data: %w", err)
	}

	keyName := fmt.Sprintf("K%s+%03d+%d", zone, signer.Key.Algorithm, signer.Key.KeyTag())
	js.Global().Get("localStorage").Call("setItem", keyName, string(marshalled))

	return signer, nil
}

func LoadKeyFile(keyfile string) (*Signer, error) {
	keyDataJson := js.Global().Get("localStorage").Call("getItem", keyfile).String()
	if keyDataJson == "" {
		return nil, fmt.Errorf("key not found")
	}

	var data StoredKeyData
	err := json.Unmarshal([]byte(keyDataJson), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal key data for %q: %w", keyfile, err)
	}

	return ParseKeyData(data.Key, data.Private)
}

// Returns all Keystore public keys and names as array of JSON objects
// where
//
//	Name: <filename prefix compatible with nsupdate>
//	Key:  <public key in DNS RR format compatible with nsupdate>
func ListKeys(dir string) ([]StoredKeyData, error) {
	if dir != "." {
		return nil, fmt.Errorf("directories not supported in wasm - use '.'")
	}

	n := js.Global().Get("localStorage").Get("length").Int()

	var keys []StoredKeyData
	for i := 0; i < n; i++ {
		key := js.Global().Get("localStorage").Call("key", i)
		if key.IsNull() {
			break
		}

		keyName := key.String()
		if !strings.HasPrefix(keyName, "K") {
			continue
		}

		keyDataJson := js.Global().Get("localStorage").Call("getItem", keyName).String()
		if keyDataJson == "" {
			continue
		}

		var data StoredKeyData
		err := json.Unmarshal([]byte(keyDataJson), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal key data for %q: %w", keyName, err)
		}
		data.Name = keyName

		// validate key data
		_, err = ParseKeyData(data.Key, data.Private)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key data for %q: %w", keyName, err)
		}

		keys = append(keys, data)
	}

	return keys, nil
}
