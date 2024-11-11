package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/libp2p/go-libp2p"
)

func main_010() {
	ctx := context.Background()

	host, err := libp2p.New()
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
