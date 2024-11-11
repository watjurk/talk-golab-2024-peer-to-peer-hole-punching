package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

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
