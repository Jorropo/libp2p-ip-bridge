package main

import (
	//"log"
	//"net"

	"github.com/miekg/dns"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type handler struct {
	dns *dns.Client
	dht *dht.IpfsDHT
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var msg *dns.Msg
	switch r.Question[0].Qtype {
	// TODO: Support more dns type (with a protobuf encoding ?)
	case dns.TypeA:
		/*m := dns.Msg{}
		m.SetReply(r)
		domain := m.Question[0].Name
		log.Printf("domain %s", domain)
		address, ok := domainsToAddresses[domain]
		if ok {
			m.Authoritative = true
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl},
				A:   net.ParseIP(address),
			})
			msg = &m
			break
		}*/
		fallthrough
	default:
		msg, _, _ = h.dns.Exchange(r, dnsServer)
	}
	w.WriteMsg(msg)
}
