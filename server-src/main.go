package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/on-keyday/qgpoc/util"
	"github.com/quic-go/quic-go"
)

func main() {
	// setup routes and VIP
	util.AddIPRoute("10.200.0.0/24", "10.220.0.254")
	util.AddIPRoute("10.210.0.0/24", "10.220.0.254")
	util.AssignVIPToLoopback(os.Getenv("SERVER_VIP"))
	util.AddIPIPTunnel("10.220.0.254", "10.220.0.2")
	util.DisableRPFilters()

	pkt, err := util.ListenObservedPacketConn("udp4", ":8889")
	if err != nil {
		panic(err)
	}
	transport := &quic.Transport{
		Conn:               pkt,
		ConnectionIDLength: 4,
	}

	// make dummy self-signed cert
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	self := &x509.Certificate{
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	data, err := x509.CreateCertificate(rand.Reader, self, self, &privKey.PublicKey, privKey)
	if err != nil {
		panic(err)
	}
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		panic(err)
	}
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privKey,
	}
	conn, err := transport.Listen(&tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}, &quic.Config{})
	if err != nil {
		panic(err)
	}
	for {
		s, err := conn.Accept(context.Background())
		if err != nil {
			return
		}
		go func() {
			for {
				str, err := s.AcceptStream(context.Background())
				if err != nil {
					return
				}
				go func() {
					buf := make([]byte, 1024)
					for {
						n, err := str.Read(buf)
						if err != nil {
							return
						}
						fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
					}
				}()
			}
		}()
	}
}
