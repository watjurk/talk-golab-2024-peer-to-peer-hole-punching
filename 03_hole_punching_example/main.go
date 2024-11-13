package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	routedHost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/protocol/holepunch"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

const PRINT_MESSAGE_PROTOCOL_ID protocol.ID = "/watjurk/print_message"

func main() {
	connectToID := flag.String("connect-to", "", "The ID of the hosts that this node should connect to")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ps, err := pstoremem.NewPeerstore()
	if err != nil {
		panic(err)
	}

	basicHost, err := libp2p.New(
		libp2p.Peerstore(ps),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithPeerSource(newPeerstoreAutoRelayPeerSource(ps)),
	)
	if err != nil {
		panic(err)
	}

	// Kad-dht comes with a list of bootstrapping peers.
	dhtRouting, err := dht.New(ctx, basicHost, dht.BootstrapPeersFunc(dht.GetDefaultBootstrapPeerAddrInfos))
	if err != nil {
		panic(err)
	}
	defer dhtRouting.Close()

	host := routedHost.Wrap(basicHost, dhtRouting)
	defer host.Close()

	fmt.Printf("This hosts ID is %s\n", host.ID().String())

	// Bootstrap the routing
	err = dhtRouting.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	// Wait for the bootstrapping to finish
	for dhtRouting.RoutingTable().Size() == 0 {
		fmt.Println("Waiting for bootstrap.")
		time.Sleep(time.Second)
	}

	dhtRouting.ForceRefresh()

	// Wait for the hole puncher to be ready
	for !slices.Contains(host.Mux().Protocols(), holepunch.Protocol) {
		fmt.Println("Waiting for hole puncher.")
		time.Sleep(time.Second)
	}

	fmt.Println("Successful bootstrapped the node!")
	host.SetStreamHandler(PRINT_MESSAGE_PROTOCOL_ID, func(s network.Stream) {
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
			fmt.Println(scanner.Text())

		}
	})

	go func() {
		if *connectToID != "" {
			connectToIDDecoded, err := peer.Decode(*connectToID)
			if err != nil {
				panic(err)
			}

			for {
				fmt.Println("Trying to connect to the other node.")

				connectCtx, cancel := context.WithTimeout(ctx, time.Second*15)
				err = host.Connect(connectCtx, peer.AddrInfo{ID: connectToIDDecoded})
				cancel()

				if err != nil {
					fmt.Println(err)
					time.Sleep(time.Second)
					continue
				}

				break
			}

			fmt.Println("Connection successful!")

			var outboundStream network.Stream
			for {
				fmt.Println("Trying to open a stream!")
				outboundStream, err = host.NewStream(ctx, connectToIDDecoded, PRINT_MESSAGE_PROTOCOL_ID)
				if err != nil {
					fmt.Println(err)
					continue
				}

				break
			}

			fmt.Println("Stream open!")
			fmt.Println("Writing data!")

			for i := 0; ; i++ {
				_, err = outboundStream.Write([]byte(fmt.Sprintf("Hello from Florence! %v\n", i)))
				if err != nil {
					fmt.Println(err)
				}

				time.Sleep(time.Second)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
