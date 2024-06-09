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
	testName := fmt.Sprintf("sig0namectl-test-%x.zenr.io", buf)

	zoneRequestMsg, dohServer, err := CreateRequestKeyMsg(testName)
	r.NoError(err)
	r.Equal("doh.zenr.io", dohServer)
	t.Log(zoneRequestMsg)
	// t.FailNow()

	// TODO: cleanup test keys
}
