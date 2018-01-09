package node

import (
	"log"
	"net"
)

type UDPNeighbor struct {
	host   string
	done   chan bool
	logger *log.Logger
	conn   *net.UDPConn
}

func NewUDPNeighbor(host string, logger *log.Logger) *UDPNeighbor {
	return &UDPNeighbor{
		host:   host,
		done:   make(chan bool, 1),
		logger: logger,
	}
}

func (udp *UDPNeighbor) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", udp.host)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	udp.logger.Printf("listening on udp://%v", addr)
	udp.conn = conn

	go udp.read(conn)

	return nil
}

func (udp *UDPNeighbor) Close() {
	udp.conn.Close()
	<-udp.done
}

func (udp *UDPNeighbor) read(conn *net.UDPConn) {
	var buf [1024 * 10]byte
	for {
		n, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			udp.logger.Printf("error reading UDP packet: %v", err)
			break
		} else {
			udp.handleMessage(buf[:n], addr)
		}
	}
	udp.logger.Printf("udp server closed")
	udp.done <- true
}

func (udp *UDPNeighbor) handleMessage(b []byte, addr *net.UDPAddr) {
	udp.logger.Printf("message from UDP neighbor: %v", addr)

	m, err := ParseUdpBytes(b)
	if err != nil {
		udp.logger.Printf("error parsing message: %v", err)
	} else {
		udp.logger.Println(m.Debug())
	}
}
