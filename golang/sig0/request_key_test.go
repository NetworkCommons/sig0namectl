package sig0

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestKey(t *testing.T) {
	r := require.New(t)

	buf := make([]byte, 5)
	rand.Read(buf)
	testSubZone := fmt.Sprintf("sig0namectl-test-%x", buf)

	zoneRequestMsg, dohServer, err := CreateRequestKeyMsg(testSubZone, "zenr.io")
	r.NoError(err)
	r.Equal("doh.zenr.io", dohServer)
	t.Log(zoneRequestMsg)
	// t.FailNow()

	// TODO: cleanup test keys
}
