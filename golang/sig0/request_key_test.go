package sig0

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestRequestKey(t *testing.T) {
	r := require.New(t)

	// TODO: cleanup test keys
	buf := make([]byte, 5)
	rand.Read(buf)
	testName := fmt.Sprintf("sig0namectl-test-%x.zenr.io", buf)

	zr, err := NewKeyRequest(testName)
	r.NoError(err)

	var answer *dns.Msg
	var i = 0
	for zr.Next() {
		t.Log("Loop", i)
		qry := zr.Do(answer)
		t.Log(qry)
		if qry == nil {
			break
		}

		answer, err = SendDOHQuery("doh.zenr.io", qry)
		r.NoError(err)
		t.Log(answer)
		i++
	}
	r.NoError(zr.Err())

	t.Log("registered key - checking registration")

	for i = 10; true; i-- {
		accepted, err := QueryAny(testName)
		r.NoError(err)

		answer, err = SendDOHQuery("doh.zenr.io", accepted)
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

	signer, err := LoadOrGenerateKey(testName)
	r.NoError(err)

	err = signer.StartUpdate("zenr.io")
	r.NoError(err)

	err = signer.UpdateA("foo", testName, "1.2.3.4")
	r.NoError(err)

	updateMsg, err := signer.SignUpdate()
	r.NoError(err)

	answer, err = SendDOHQuery("doh.zenr.io", updateMsg)
	r.NoError(err)
	t.Log(answer)
	r.True(answer.Rcode == dns.RcodeSuccess)
}
