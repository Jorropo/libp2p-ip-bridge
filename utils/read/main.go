package main

import (
	"flag"
	"log"
	"net"

	"github.com/libp2p/go-libp2p"

	"github.com/Jorropo/libp2p-ip-bridge/bootstrap"
	"github.com/Jorropo/libp2p-ip-bridge/get"
)

var libp2pListenPort = flag.Uint("l", 0, "The port where libp2p is gonna listen on.")
var toRead = flag.String("v", "", "The value to read of the dht. (REQUIRED)")

func main() {
	flag.Parse()
	_, dht := bootstrap.GetNode(uint16(*libp2pListenPort), libp2p.RandomIdentity)

	addr, err := get.Get(*toRead, dht)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	log.Printf("Got %s\n", net.IP(addr).String())
}
