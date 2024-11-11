package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p"
	routedHost "github.com/libp2p/go-libp2p/p2p/host/routed"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main_routing() {
	ctx := context.Background()

	host, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// host creation
	// ...
	// Kad-dht comes with a list of bootstrapping peers.
	dhtRouting, err := dht.New(ctx, host, dht.BootstrapPeersFunc(dht.GetDefaultBootstrapPeerAddrInfos))
	if err != nil {
		panic(err)
	}
	defer dhtRouting.Close()

	host = routedHost.Wrap(host, dhtRouting)
	defer host.Close()

	err = dhtRouting.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	const HELLO_PROTOCOL = "/golab/p2p-hello"
	stream, err := host.NewStream(ctx, "Artur", HELLO_PROTOCOL)
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
