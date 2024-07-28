package sig0

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestKey(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	// TODO: cleanup test keys
	buf := make([]byte, 5)
	_, _ = rand.Read(buf)
	testName := fmt.Sprintf("sig0namectl-test-%x.zenr.io", buf)

	err := RequestKey(testName)
	r.NoError(err)

	t.Log("registered key - checking registration")

	for i := 10; true; i-- {
		accepted, err := QueryAny(testName)
		r.NoError(err)

		answer, err := SendDOHQuery("doh.zenr.io", accepted)
		r.NoError(err)
		t.Log(answer)

		if answer.Rcode != dns.RcodeNameError {
			t.Log("name registered")
			break
		}

		if i == 0 {
			t.Fatal("name registration failed")
		}

		t.Log("waiting...")
		time.Sleep(15 * time.Second)
	}

	// TODO: move checkKey() from wrapper_js.go into sig0
	signer, err := LoadOrGenerateKey(testName)
	r.NoError(err)

	t.Cleanup(func() {
		kn := signer.KeyName()
		t.Log("deleting key:", kn)
		_ = os.Remove(kn + ".private")
		_ = os.Remove(kn + ".key")
	})

	err = signer.StartUpdate("zenr.io")
	r.NoError(err)

	err = signer.UpdateA("foo", testName, "1.2.3.4")
	r.NoError(err)

	updateMsg, err := signer.SignUpdate()
	r.NoError(err)

	answer, err := SendDOHQuery("doh.zenr.io", updateMsg)
	r.NoError(err)
	if !a.True(answer.Rcode == dns.RcodeSuccess, "answer not successful") {
		t.Log(answer)
	}
}
