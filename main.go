package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	libp "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"context"
	
)

func main() {

	/* pass string address to node
	node, err := libp.New(libp.ListenAddrStrings("/ip4/192.168.64.1/tcp/50000"),)
	*/

	node, err := libp.New(
		libp.ListenAddrStrings("/ip4/192.168.64.1/tcp/0"), // chose random port
		libp.Ping(false), // disable default ping
	)
	if err != nil{
		fmt.Printf("err: ", err)
	}

	peerInfo := peerstore.AddrInfo{
		ID: node.ID(),
		Addrs: node.Addrs(),

	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Printf("P2p Address: ", addrs[0])

	// ping service
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)


	// connect to peer and ping if peer addr is parsed
	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])

		if err != nil{
			panic(err)
		}

		peer, err := peerstore.AddrInfoFromP2pAddr(addr)

		if err != nil{
			panic(err)
		}

		if err := node.Connect(context.Background(), *peer); err != nil{
			panic(err)
		}

		fmt.Printf("Sending ping to peer: ", addr)
		ch := pingService.Ping(context.Background(), peer.ID)

		for i:= 0; i < 5; i++ {
			res := <-ch
			fmt.Println("Ping response", "RTT:", res.RTT)
		}
		
	}else{
		// wait for signal termination
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Ctrl + C to shutdown p2p node")

		fmt.Println("Recieve Signal, shutting down ...")
	}

	if err := node.Close(); err != nil{
		panic(err)
	}
}