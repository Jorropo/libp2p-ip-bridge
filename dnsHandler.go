package main

import (
	"log"
	"net"

  "github.com/miekg/dns"

	//"github.com/libp2p/go-libp2p"
	//"github.com/libp2p/go-libp2p-core/host"
)

var domainsToAddresses map[string]string = map[string]string{
	"google.com.":       "1.2.3.4",
	"jameshfisher.com.": "104.198.14.52",
}

type handler struct {
	dns *dns.Client
	//p2p host.Host
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var msg *dns.Msg
	switch r.Question[0].Qtype {
	case dns.TypeA:
    m := dns.Msg{}
    m.SetReply(r)
		domain := m.Question[0].Name
		log.Printf("domain %s", domain)
		address, ok := domainsToAddresses[domain]
		if ok {
			m.Authoritative = true
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
			msg = &m
      break
		}
    fallthrough
  default:
		msg, _, _ = h.dns.Exchange(r, dnsServer)
	}
	w.WriteMsg(msg)
}
