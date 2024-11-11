package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	routedHost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

func main_hole_punching() {
	ctx := context.Background()

	memoryPeerstore, err := pstoremem.NewPeerstore()
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		libp2p.EnableHolePunching(),
		libp2p.Peerstore(memoryPeerstore),
		libp2p.EnableAutoRelayWithPeerSource(newPeerstoreAutoRelayPeerSource(memoryPeerstore)),
	)
	if err != nil {
		panic(err)
	}

	// ...
	// routing
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
	// ...

	const HELLO_PROTOCOL = "/golab/p2p-hello"
	stream, err := host.NewStream(ctx, "Artur", HELLO_PROTOCOL)
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func newPeerstoreAutoRelayPeerSource(store peerstore.Peerstore) autorelay.PeerSource {
	return func(ctx context.Context, num int) <-chan peer.AddrInfo {
		fmt.Printf("Autorelay requested %v nodes!\n", num)
		peerChan := make(chan peer.AddrInfo)
		go func() {
			defer close(peerChan)

			peers := store.Peers()
			provideAtMost := min(len(peers), num)

			for _, i := range rand.Perm(provideAtMost) {
				select {
				case peerChan <- store.PeerInfo(peers[i]):
				case <-ctx.Done():
					return
				}
			}

		}()

		return peerChan
	}
}
