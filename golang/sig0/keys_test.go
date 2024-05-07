package sig0

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadKey(t *testing.T) {
	r := require.New(t)
	keyName := createKeyViaBind(t)

	signer, err := LoadKeyFile(keyName)
	r.NoError(err)
	r.NotNil(signer)
}

func TestParseKeyFile(t *testing.T) {
	r := require.New(t)
	keyName := createKeyViaBind(t)

	keyContent, err := os.ReadFile(keyName + ".key")
	r.NoError(err)

	privateContent, err := os.ReadFile(keyName + ".private")
	r.NoError(err)

	signer, err := ParseKeyData(string(keyContent), string(privateContent))
	r.NoError(err)
	r.NotNil(signer)
}

func TestCompareFlags(t *testing.T) {
	a := assert.New(t)

	bindKey := createAndLoadKey(t)

	ourKey, err := GenerateKey("go.te.st")
	require.NoError(t, err)

	a.Equal(bindKey.Key.Algorithm, ourKey.Key.Algorithm, "Algorithm")
	a.Equal(bindKey.Key.Flags, ourKey.Key.Flags, "Flags")
	a.Equal(bindKey.Key.Protocol, ourKey.Key.Protocol, "Protocol")
}

func createKeyViaBind(t *testing.T) string {
	var buf bytes.Buffer
	cmd := exec.Command("dnssec-keygen", "-K", "/tmp", "-a", "ED25519", "-n", "HOST", "-T", "KEY", "go.te.st")
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	err := cmd.Run()
	require.NoError(t, err)

	keyName := filepath.Join("/tmp", strings.TrimSpace(buf.String()))

	t.Log("created key file:", keyName)

	t.Cleanup(func() {
		os.Remove(keyName + ".key")
		os.Remove(keyName + ".private")
	})

	return keyName
}

func createAndLoadKey(t *testing.T) *Signer {
	keyName := createKeyViaBind(t)

	signer, err := LoadKeyFile(keyName)
	require.NoError(t, err)
	require.NotNil(t, signer)

	return signer
}
