package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NetworkCommons/sig0poc1/sig0"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
)

func main() {
	name := os.Args[1]
	// name := "cryptix.zenr.io"
	server := "zembla.zenr.io"

	q := sig0.QueryA(name)
	fmt.Printf("Q:(TXT):%v\n", q)

	fmt.Println("Q:", q)

	// send over DoH
	url := fmt.Sprintf("https://%s/dns-query?dns=%s", server, q)
	fmt.Println("Q:(DoH):", url)
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	answerBody, err := io.ReadAll(resp.Body)
	check(err)

	var dnsAnswer = new(dns.Msg)
	err = dnsAnswer.Unpack(answerBody)
	check(err)

	spew.Dump(dnsAnswer)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
