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

	var persisted storedKeyData
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

	var data storedKeyData
	err := json.Unmarshal([]byte(keyDataJson), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal key data for %q: %w", keyfile, err)
	}

	return ParseKeyData(data.Key, data.Private)
}

type storedKeyData struct {
	Key, Private string
}

func ListKeys(dir string) ([]string, error) {
	if dir != "." {
		return nil, fmt.Errorf("directories not supported - use '.'")
	}

	n := js.Global().Get("localStorage").Get("length").Int()

	var keys []string
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

		var data storedKeyData
		err := json.Unmarshal([]byte(keyDataJson), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal key data for %q: %w", keyName, err)
		}

		_, err = ParseKeyData(data.Key, data.Private)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key data for %q: %w", keyName, err)
		}

		keys = append(keys, keyName)
	}

	return keys, nil
}
