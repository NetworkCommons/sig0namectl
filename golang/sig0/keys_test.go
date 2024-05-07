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
	a := assert.New(t)
	keyName := createKeyViaBind(t)

	keyContent, err := os.ReadFile(keyName + ".key")
	r.NoError(err)

	privateContent, err := os.ReadFile(keyName + ".private")
	r.NoError(err)

	signer, err := ParseKeyData(string(keyContent), string(privateContent))
	r.NoError(err)
	r.NotNil(signer)

	// t.Logf("pk type: %T", signer.private)

	a.Equal(uint8(0xf), signer.Key.Algorithm, "Algorithm")
	a.Equal(uint16(0x200), signer.Key.Flags, "Flags")
	a.Equal(uint8(0x3), signer.Key.Protocol, "Protocol")
	a.Equal(DefaultKeyTTL, signer.Key.Hdr.Ttl, "TTL")

	// re-create the key string
	k := signer.Key
	out := k.String()
	out = strings.ReplaceAll(out, "\t", " ")
	out += "\n"
	a.Equal(string(keyContent), out)

	// TODO: recreate the private string
	// pk := signer.Key.PrivateKeyString(signer.private)
	// assert.Equal(t, string(privateContent), pk)
}

func TestWritingAndParsingInBind(t *testing.T) {
	r := require.New(t)

	signer, err := GenerateKeyAndSave("go.te.st")
	r.NoError(err)

	kn := signer.KeyName()
	t.Log("KeyName:", kn)
	t.Cleanup(func() {
		os.Remove(kn + ".key")
		os.Remove(kn + ".private")
	})

	nsUpdateVersion, err := exec.Command("nsupdate", "-V").CombinedOutput()
	r.NoError(err)
	// TODO: add regexp for >= 9.18
	r.Contains(string(nsUpdateVersion), "nsupdate 9.18")

	nsUpdateInput := strings.NewReader(`
update add go.te.st 60 A 1.1.1.1
show
`)
	cmd := exec.Command("nsupdate", "-k", kn)
	cmd.Stdin = nsUpdateInput
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	err = cmd.Run()
	r.NoError(err)
}

func TestCompareFlags(t *testing.T) {
	a := assert.New(t)

	bindKey := createAndLoadKey(t)

	ourKey, err := GenerateKey("go.te.st")
	require.NoError(t, err)

	a.Equal(bindKey.Key.Algorithm, ourKey.Key.Algorithm, "Algorithm")
	a.Equal(bindKey.Key.Flags, ourKey.Key.Flags, "Flags")
	a.Equal(bindKey.Key.Protocol, ourKey.Key.Protocol, "Protocol")
	a.Equal(bindKey.Key.Hdr.Ttl, ourKey.Key.Hdr.Ttl, "TTL")
}

func createKeyViaBind(t *testing.T) string {
	var buf bytes.Buffer
	cmd := exec.Command("dnssec-keygen", "-K", "/tmp", "-a", "ED25519", "-n", "HOST", "-T", "KEY", "-L", "60", "go.te.st")
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		t.Log(buf.String())
	}
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
