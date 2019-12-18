package main

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/Jorropo/libp2p-ip-bridge/get"

	ma "github.com/multiformats/go-multiaddr"
)

type handler struct {
	dns *dns.Client
	dht *dht.IpfsDHT
}

const escapeChar = 0x2d // "-"

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var msg *dns.Msg
	for _, question := range r.Question {
		// TODO: Support more dns type (with a protobuf encoding ?)
		if question.Qtype == dns.TypeA {
			log.Printf("domain %s\n", question.Name)
			var dotCount uint = uint(strings.Count(question.Name, "."))
			if dotCount > 1 {
				// Discard the last dot .dht**.**
				split := strings.Split(question.Name, ".")[:dotCount]
				// Removing the last empty block
				dotCount--
				// Switching over the tld.
				switch split[dotCount] {
				case "dht":
					if dotCount == 1 {
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
				case "maddr":
					if dotCount > 1 {
						var maddrStr string
						var i uint = 0
						goto ReadMaddr
					Escape:
						{
							var maddrStre string
							var j uint = 0
							segmentLen := uint(len(split[i]))
							goto EscapeLoop
						EscapeEscape:
							if j == segmentLen {
								maddrStre += "."
								goto EscapeEnd
							}
							if split[i][j] == escapeChar {
								maddrStre += "-"
							} else {
								maddrStre += "." + string(split[i][j])
							}
							j++
						EscapeLoop:
							if j == segmentLen {
								goto EscapeEnd
							}
							if split[i][j] == escapeChar {
								j++
								goto EscapeEscape
							}
							maddrStre += string(split[i][j])
							j++
							goto EscapeLoop
						EscapeEnd:
							maddrStr += ("/" + maddrStre)
							i++
						}
					ReadMaddr:
						if i == dotCount {
							goto MakeMaddr
						}
						if strings.Contains(split[i], "-") {
							goto Escape
						}
						maddrStr += ("/" + split[i])
						i++
						goto ReadMaddr
					MakeMaddr:
						maddr, err := ma.NewMultiaddr(maddrStr)
						log.Printf("maddr: %s\n", maddr.String())
						if err != nil {
							log.Printf("can't parse maddr: %s,%s\n", err, maddrStr)
							goto Forward
						}
						protocols := maddr.Protocols()
						if len(protocols) == 0 {
							goto Forward
						}
						// We want to allow connection to `/ip4/127.0.0.1/tcp/12345` but we
						// expect the app to take care of tcp with the kernel.
						// That why we only keep the first value.
						m := dns.Msg{}
						m.SetReply(r)
						m.Authoritative = true
						switch protocols[0].Code {
						case 4: // ip4
							ipStr, err := maddr.ValueForProtocol(4)
							if err != nil {
								log.Printf("can't get ip from maddr: %s, %s\n", err, maddr.String())
								goto Forward
							}
							m.Answer = append(m.Answer, &dns.A{
								Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl},
								A:   net.ParseIP(ipStr),
							})
							msg = &m
							goto Send
						}
					}
				}
			}
		}
	}
Forward:
	msg, _, _ = h.dns.Exchange(r, dnsServer)
Send:
	w.WriteMsg(msg)
}
