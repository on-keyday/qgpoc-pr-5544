package util

import (
	"fmt"
	"log"
	"net"

	"github.com/quic-go/quic-go"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

type ObservedPacketConn struct {
	quic.OOBCapablePacketConn
	bt *ipv4.PacketConn
}

// fill to output different interface
func (c *ObservedPacketConn) WriteMsgUDP(b, oob []byte, addr *net.UDPAddr) (n, oobn int, err error) {
	log.Printf("Sending packet to %s OOB Len: %d", addr.String(), len(oob))
	const cmsghdrLen = unix.SizeofCmsghdr
	if addr.IP.To4() != nil {
		if len(oob) >= cmsghdrLen+4 {
			for i := 0; i < 4; i++ {
				oob[cmsghdrLen+i] = 0
			}
		}
	} else if addr.IP.To16() != nil {
		if len(oob) >= cmsghdrLen+16+4 {
			for i := 0; i < 4; i++ {
				oob[cmsghdrLen+16+i] = 0
			}
		}
	}
	return c.OOBCapablePacketConn.WriteMsgUDP(b, oob, addr)
}

func (c *ObservedPacketConn) ReadBatch(ms []ipv4.Message, flags int) (int, error) {
	n, err := c.bt.ReadBatch(ms, flags)
	for i := 0; i < n; i++ {
		src := ms[i].Addr.String()
		log.Printf("Received packet from %s", src)
		_ = src
	}
	return n, err
}

func ListenObservedPacketConn(network, address string) (*ObservedPacketConn, error) {
	pkt, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}
	oobcap, ok := pkt.(quic.OOBCapablePacketConn)
	if !ok {
		return nil, fmt.Errorf("PacketConn %T does not implement OOBCapablePacketConn", pkt)
	}
	return &ObservedPacketConn{OOBCapablePacketConn: oobcap, bt: ipv4.NewPacketConn(pkt)}, nil
}
