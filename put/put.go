package put

import (
	"context"
	"log"
	"time"

	kaDht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-crypto"

	proto "github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-ipns"
)

func Put(priv crypto.PrivKey, pub crypto.PubKey, host host.Host, dht *kaDht.IpfsDHT, toPublish []byte) {
	ipnsRecord, err := ipns.Create(priv, toPublish, 0, time.Now().Add(1*time.Hour))
	if err != nil {
		log.Fatalf("can't create ipns records: %s\n", err)
	}

	err = ipns.EmbedPublicKey(pub, ipnsRecord)
	if err != nil {
		log.Fatalf("can't embed ipns records: %s\n", err)
	}

	log.Println(ipnsRecord.String())

	data, err := proto.Marshal(ipnsRecord)
	if err != nil {
		log.Fatalf("can't Marshal ipns record: %s\n", err)
	}

	err = dht.PutValue(context.Background(), ipns.RecordKey(host.ID()), data)
	if err != nil {
		log.Fatalf("can't put value to dht: %s\n", err)
	}
}
