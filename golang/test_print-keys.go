package main

import (
	"github.com/miekg/dns"
	"github.com/davecgh/go-spew/spew"
)

import (
	"fmt"
	"log"
	"os"
)


func main() {
	if len(os.Args) < 2 { log.Fatal("Filename required") }

	domain, ok := os.LookupEnv("GD_DOMAIN")
	if !ok {
		fmt.Println("No GD_DOMAIN ENV var defined")
	} else {
		fmt.Printf("GD_DOMAIN = %s\n", domain)
	}
	host, ok := os.LookupEnv("GD_HOST")
	if !ok {
		fmt.Println("No GD_HOST ENV var defined")
	} else {
		fmt.Printf("GD_HOST = %s\n", host)
	}
	server, ok := os.LookupEnv("GD_SERVER")
	if !ok {
		fmt.Println("No GD_SERVER ENV var defined")
	} else {
		fmt.Printf("GD_SERVER = %s\n", server)
	}

//	m := new(dns.Msg)
//	rrInsert, err := dns.NewRR("rrsetnotused.free2air.net. 3600 IN A 127.0.0.1")
//	if err != nil {
//		panic(err)
//        }
//        spew.Dump(rrInsert)
//	m.Insert([]dns.RR{rrInsert})
//
//	m.SetUpdate("rrset_not_used.free2air.net. 3600 IN A 127.0.0.1")
//
//	in, err := dns.Exchange(m, server)
//	fmt.Println("++")
//	spew.Dump(in)
//	fmt.Println("++")


  	pubfh, perr := os.Open(os.Args[1]+".key")
  	if perr != nil { log.Fatal(perr) }
  
  	dk, pkerr := dns.ReadRR(pubfh, os.Args[1]+".key")
  	if pkerr != nil { log.Fatal(pkerr) }
  	spew.Dump(dk, pkerr)
  
  	privfh, oerr := os.Open(os.Args[1]+".private")
  	if oerr != nil { log.Fatal(oerr) }
  	defer privfh.Close()
  
  	privkey, readerr := dk.(*dns.KEY).ReadPrivateKey(privfh, os.Args[1]+".private")
  	if readerr == nil {
  		spew.Dump(privkey, readerr)
  		log.Println("OK")
  	} else {
  		spew.Dump(privkey, readerr)
  	}
}
