package main

import (
	"encoding/base64"
	"fmt"
	// "os"

	"github.com/miekg/dns"
	"github.com/shynome/doh-client"
)

func main() {
	// name := os.Args[1]
	name := "cryptix.zenr.io"
	server := "zembla.zenr.io"
	q := dns.Question{
		Name:   dns.Fqdn(name),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}

	out, err := m.Pack()
	check(err)

	// fmt.Printf("Q:(TXT): %s\t\t\t%d\n", q.Name, q.Qtype)
	fmt.Printf("Q:(TXT):%v\n", q)
	fmt.Println("Q:(b64):", base64.StdEncoding.EncodeToString(out))
	fmt.Println("Q:(DoH):", "https://"+server+"/dns-query="+base64.StdEncoding.EncodeToString(out))

	inputStr := "lzeFAAABAAIAAAAAB2NyeXB0aXgEemVucgJpbwAAAQABwAwAAQABAAAAPAAEX9ko/cAMAAEAAQAAADwABFSiXZs="
	input, err := base64.StdEncoding.DecodeString(inputStr)
	check(err)
	var resp = new(dns.Msg)

	err = resp.Unpack(input)
	check(err)
	// fmt.Printf("A: +%v\n", resp.Answer)

	co := &dns.Conn{Conn: doh.NewConn(nil, nil, server)}
	if err := co.WriteMsg(m); err != nil {
		panic(err)
	}
	// check(err)

	m, err = co.ReadMsg()
	check(err)
	if len(m.Answer) == 0 {
		panic("answer length must greater than 0")
	}
	// fmt.Printf("A:(TXT):%+v\n", m.Answer)
	for _, str := range m.Answer {
		fmt.Printf("A:(TXT): %s\n",str)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
