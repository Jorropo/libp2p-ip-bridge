package main

import (
	"flag"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"

	"github.com/Jorropo/libp2p-ip-bridge/bootstrap"
	"github.com/Jorropo/libp2p-ip-bridge/put"
)

var libp2pListenPort = flag.Uint("l", 0, "The port where libp2p is gonna listen on.")
var toPublish = flag.String("v", "", "The value to publish into the dht. (REQUIRED)")

func main() {
	flag.Parse()

	// Generate a private key to sign the IPNS record with. Most of the time,
	// however, you'll want to retrieve an already-existing key from IPFS using the
	// go-ipfs/core/coreapi CoreAPI.KeyAPI() interface.
	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		log.Fatalf("can't create key: %s\n", err)
	}

	host, dht := bootstrap.GetNode(uint16(*libp2pListenPort), libp2p.Identity(priv))

	put.Put(priv, pub, host, dht, []byte(*toPublish))

	log.Printf("Serving %s\n", host.ID().Pretty())
	select {}
}
