//go:build ignore

package main

import (
	"github.com/miekg/dns"
	"github.com/davecgh/go-spew/spew"
)

import (
	"log"
	"os"
	"fmt"
)


func main() {
	fmt.Println("YO GO!")
	if len(os.Args) < 2 { log.Fatal("Filename required") }

	dk := new(dns.DNSKEY)

	fh, oerr := os.Open(os.Args[1])
	if oerr != nil { log.Fatal(oerr) }
	defer fh.Close()

	privkey, readerr := dk.ReadPrivateKey(fh, os.Args[1])
	spew.Dump(privkey, readerr)
}
