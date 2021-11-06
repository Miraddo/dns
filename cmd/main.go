// Package main
package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mobocrat/dns/cmd/cli"
)

const address = "8.8.8.8:53"

// main ...
func main() {

	valid := []string{"A", "AAAA", "NS", "CNAME"}

	check := sort.SearchStrings(valid, cli.Query)

	if check != -1 {
		socket, err := net.Dial("udp", address)

		if err != nil {
			panic(err)
		}
		defer socket.Close()

		req := genPackage(cli.Domain, cli.Query)

		if _, err := socket.Write(req); err != nil {
			panic(err)
		}

		resp := make([]byte, 1024)

		socket.SetReadDeadline(time.Now().Add(3 * time.Second))

		n, err := socket.Read(resp)

		if err != nil {
			panic(err)
		}

		for i := 0; i < n; i += 8 {
			fmt.Printf("%0x", resp[i:i+8])
			fmt.Print(" ")
			fmt.Printf("%s", string(resp[i:i+8]))
			fmt.Print("\n")
		}

		os.Exit(1)
	}

	fmt.Println("Query is not supported")
}

// genPackage ...
func genPackage(domain, query string) []byte {

	qcode := make([]byte, 2)

	switch query {
	case "A":
		qcode = []byte{0, 1}
	case "NS":
		qcode = []byte{0, 2}
	case "CNAME":
		qcode = []byte{0, 5}
	case "AAAA":
		qcode = []byte{0, 28}
	}

	fmt.Println(domain)
	dcode := domainToByte(domain)

	packet := []byte{
		0xdb, 0x42, // ID
		1, 0, // Flags
		0, 1, // QDCOUNT
		0, 0, // ANCOUNT
		0, 0, // NSCOUNT
		0, 0, // ARCOUNT
	}

	packet = append(packet, dcode...)
	packet = append(packet, qcode...)
	packet = append(packet, []byte{0, 1}...)

	return packet
}

// domainToByte ...
func domainToByte(domain string) []byte {

	var dbyte []byte

	d := strings.Split(strings.ToLower(domain), ".")

	for _, x := range d {
		dbyte = append(dbyte, byte(len(x)))

		for _, i := range x {
			dbyte = append(dbyte, byte(i))
		}
	}

	dbyte = append(dbyte, byte(0))

	return dbyte
}
