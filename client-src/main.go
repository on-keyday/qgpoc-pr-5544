package main

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"time"

	"github.com/on-keyday/qgpoc/util"
	"github.com/quic-go/quic-go"
)

func main() {
	util.AddIPRoute(os.Getenv("DEST_NET"), os.Getenv("ROUTE_VIA"))

	pkt, err := util.ListenObservedPacketConn("udp4", "")
	if err != nil {
		panic(err)
	}
	transport := &quic.Transport{
		Conn:               pkt,
		ConnectionIDLength: 4,
	}

	for {
		conn, err := transport.Dial(context.Background(), &net.UDPAddr{
			IP:   net.ParseIP(os.Getenv("SERVER_VIP")),
			Port: 8889,
		}, &tls.Config{
			InsecureSkipVerify: true,
		}, &quic.Config{})
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			panic(err)
		}
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			_, err := stream.Write([]byte("hello"))
			if err != nil {
				break
			}
		}
	}
}
