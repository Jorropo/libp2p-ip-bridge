package get

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p-core/peer"

	kaDht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/ipfs/go-ipns"
	ipnsPb "github.com/ipfs/go-ipns/pb"

	proto "github.com/gogo/protobuf/proto"
)

func Get(toRead string, dht *kaDht.IpfsDHT) []byte {
	id, err := peer.Decode(toRead)
	if err != nil {
		log.Fatalf("can't decode id: %s\n")
	}

	val, err := dht.GetValue(context.Background(), ipns.RecordKey(id))
	if err != nil {
		log.Fatalf("can't get value of the dht: %s\n", err)
	}

	ent := &ipnsPb.IpnsEntry{}
	err = proto.Unmarshal(val, ent)
	if err != nil {
		log.Fatalf("can't get decode the result: %s\n", err)
	}

	return ent.GetValue()
}
