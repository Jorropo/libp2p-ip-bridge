package main

import (
	"log"
	"net"
  "time"

	"github.com/miekg/dns"
)

func main() {
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
	srv.Handler = &handler{dns: client}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
