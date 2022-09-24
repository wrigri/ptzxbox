package netcon

import (
	"fmt"
	"net"
)

var Conn *NetworkConnection = &NetworkConnection{}

type NetworkConnection struct {
	IPAddress  string
	Port       string
	Connection *net.UDPConn
}

func (nc *NetworkConnection) Address() string {
	return fmt.Sprintf("%s:%s", nc.IPAddress, nc.Port)
}

func (nc *NetworkConnection) Connect() error {
	udpAddr, err := net.ResolveUDPAddr("udp4", nc.Address())
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	nc.Connection = conn
	return nil
}

func (nc *NetworkConnection) SendCommand(command []byte) error {
	_, err := nc.Connection.Write((command))
	if err != nil {
		return err
	}
	return nil
}

func (nc *NetworkConnection) Query(command []byte) ([]byte, error) {
	_, err := nc.Connection.Write(command)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 1024)
	n, _, err := nc.Connection.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (nc *NetworkConnection) Close() {
	nc.Connection.Close()
}
