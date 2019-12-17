# libp2p-ip-bridge
This deamon aim to bridge ip over libp2p allowing non libp2p apps to work over
it.

It does that by binding to some address on `127.0.0.0/8` then it forward each
incoming connection to the libp2p tunnel.

You can add some fixed addr in the config, or use automaticaly through a dns
proxy.

## The dns proxy :
Its listening on `127.0.0.1:53`, you must then set your dns resolver as
`127.0.0.1` and it will just proxy normal dns request to a real dns resolver
(except if the host is finded in `/etc/hosts`).

But if the request is on `Qmmfoo.{dht,ipns}` (***WIP*** may change in the
future) it will resolve in the dht or ipns and return the real ip or create a
new libp2p tunnel on a random unused `127.0.0.1/8` forwarding through the
tunnel. And finaly responding the dns request with the `127.0.0.1/8` address
where the tunnel has been binded.
