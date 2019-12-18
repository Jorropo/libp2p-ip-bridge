package main

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/Jorropo/libp2p-ip-bridge/get"
)

type handler struct {
	dns *dns.Client
	dht *dht.IpfsDHT
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var msg *dns.Msg
	var question dns.Question
	for _, quest := range r.Question {
		// TODO: Support more dns type (with a protobuf encoding ?)
		if quest.Qtype == dns.TypeA {
			question = quest
			goto Parse
		}
	}
	goto Forward
Parse:
	log.Printf("domain %s\n", question.Name)
	if strings.Count(question.Name, ".") == 2 {
		split := strings.Split(question.Name, ".")[:2]
		switch split[1] {
		case "dht":
			addr, err := get.Get(split[0], h.dht)
			if err != nil {
				log.Printf("dht search error: %s\n", err)
				goto Forward
			}
			if len(addr) == 4 {
				m := dns.Msg{}
				m.SetReply(r)
				m.Authoritative = true
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl},
					A:   net.IPv4(addr[0], addr[1], addr[2], addr[3]),
				})
				msg = &m
				goto Send
			}
		}
	}
Forward:
	msg, _, _ = h.dns.Exchange(r, dnsServer)
Send:
	w.WriteMsg(msg)
}
