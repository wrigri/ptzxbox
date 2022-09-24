package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wrigri/ptzxbox/internal/config"
	"github.com/wrigri/ptzxbox/internal/util"
)

var testCmds map[string]string

func testCommands() {
	testCmds = make(map[string]string)
	// PanTilt query
	testCmds[string([]byte{0x81, 0x09, 0x06, 0x12, 0xFF})] = string([]byte{0x90, 0x50, 0x00, 0x00, 0x0B, 0x02, 0x0F, 0x0F, 0x0B, 0x03, 0xFF})
	// Zoom query
	testCmds[string([]byte{0x81, 0x09, 0x04, 0x47, 0xFF})] = string([]byte{0x90, 0x50, 0x01, 0x02, 0x0E, 0x03, 0xFF})
}

func handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client:", addr)
	fmt.Println("Received from UDP client:", util.GetHexString(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}
	// write message back to client
	message := []byte(testCmds[string(buffer[:n])])
	if len(message) > 0 {
		fmt.Println("Sending Back:", util.GetHexString(message))
		_, err = conn.WriteToUDP(message, addr)
	}
	if err != nil {
		log.Println(err)
	}
}

func main() {
	testCommands()
	err := config.Config.Update()
	if err != nil {
		log.Fatalf("Error updaing config from config file: %v\n", err)
	}

	port := config.Config.Network.Port()
	hostName := config.Config.Network.IpAddress
	service := hostName + ":" + port

	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port", port)
	defer ln.Close()

	for {
		handleUDPConnection(ln)
	}

}
