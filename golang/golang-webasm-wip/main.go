package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {
	name := os.Args[1]
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

	fmt.Println("Q:", base64.StdEncoding.EncodeToString(out))

	inputStr := "9beFAAABAAIAAAAAB2NyeXB0aXgEemVucgJpbwAAAQABwAwAAQABAAAAPAAEX9ko/cAMAAEAAQAAADwABFSiXZs="
	input, err := base64.StdEncoding.DecodeString(inputStr)
	check(err)
	var resp = new(dns.Msg)

	err = resp.Unpack(input)
	check(err)
	fmt.Printf("A: +%v\n", resp.Answer)

	// co := &dns.Conn{Conn: doh.NewConn(nil, nil, "1.1.1.1")}

	// err := co.WriteMsg(m)
	// check(err)

	// m, err = co.ReadMsg()
	// check(err)
	// if len(m.Answer) == 0 {
	// 	panic("answer length must greater than 0")
	// }
	// fmt.Printf("Answer: %+v\n", m.Answer)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
