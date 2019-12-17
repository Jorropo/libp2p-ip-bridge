// This is code is borrowed and slightly modified from github.com/libp2p/go-libp2p-example
package bootstrap

import (
	"context"
	"fmt"
	"log"

	kaDht "github.com/libp2p/go-libp2p-kad-dht"
	kaOpts "github.com/libp2p/go-libp2p-kad-dht/opts"

	"github.com/Jorropo/go-utp-transport"
	ipns "github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-tcp-transport"

	ma "github.com/multiformats/go-multiaddr"
)

func GetNode(libp2pListenPort uint16, out_opts ...libp2p.Option) (host.Host, *kaDht.IpfsDHT) {
	bootstrapPeers := IPFS_PEERS

	ctx := context.Background()

	opts := append([]libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", libp2pListenPort),
			fmt.Sprintf("/ip6/::/tcp/%d", libp2pListenPort),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/utp", libp2pListenPort),
			fmt.Sprintf("/ip6/::/udp/%d/utp", libp2pListenPort),
		),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(utp.NewUtpTransport),
		libp2p.DefaultEnableRelay,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.UserAgent("libp2p-ip-bridge"),
	}, out_opts...)

	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		log.Fatalf("can't do libp2p host %s\n", err)
	}

	// Make the DHT
	dht, err := kaDht.New(ctx, basicHost, kaOpts.Validator(&ipns.Validator{KeyBook: basicHost.Peerstore()}))
	if err != nil {
		log.Fatalf("Error with dht creation: %s\n", err)
	}

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, dht)

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", routedHost.ID().Pretty()))

	// connect to the chosen ipfs nodes
	err = bootstrapConnect(ctx, routedHost, bootstrapPeers)
	if err != nil {
		log.Fatalf("can't bootstrap network %s\n", err)
	}

	// Bootstrap the host
	err = dht.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("can't bootstrap dht %s\n", err)
	}

	// Now we can build a full multiaddress to reach this host by encapsulating
	// both addresses.
	log.Println("I can be reached at:")
	for _, addr := range routedHost.Addrs() {
		log.Println(addr.Encapsulate(hostAddr))
	}

	return routedHost, dht
}
