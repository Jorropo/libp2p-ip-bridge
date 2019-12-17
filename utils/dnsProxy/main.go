package main

import (
	"log"
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/Jorropo/libp2p-ip-bridge/bootstrap"
)

func main() {
	_, dht := bootstrap.GetNode(libp2pListenPort)

	client := &dns.Client{}
	client.Dialer = &net.Dialer{
		Timeout: 200 * time.Millisecond,
		LocalAddr: &net.UDPAddr{
			IP:   net.ParseIP("[::1]"),
			Port: 0,
			Zone: "",
		},
	}

	srv := &dns.Server{Addr: dnsProxy, Net: "udp"}
	srv.Handler = &handler{dns: client, dht: dht}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err)
	}
}
