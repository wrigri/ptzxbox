package main

import (
	"fmt"
	"log"

	"github.com/wrigri/ptzxbox/internal/config"
	"github.com/wrigri/ptzxbox/internal/netcon"
	"github.com/wrigri/ptzxbox/internal/util"
)

const outputTemplate string = `Current Position:
    pan:   %s
    tilt:  %s
    zoom:  %s
`

func getZoom() string {
	query := []byte{0x81, 0x09, 0x04, 0x47, 0xFF}
	response, err := netcon.Conn.Query(query)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return util.GetPosHex(response[2:6])
}

func getPanTilt() (string, string) {
	query := []byte{0x81, 0x09, 0x06, 0x12, 0xFF}
	response, err := netcon.Conn.Query(query)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return util.GetPosHex(response[2:6]), util.GetPosHex(response[6:10])
}

func main() {
	err := config.Config.Update()
	if err != nil {
		log.Fatalf("Error updating config from config file: %v\n", err)
	}
	netcon.Conn.IPAddress = config.Config.Network.IpAddress
	netcon.Conn.Port = config.Config.Network.Port()
	if netcon.Conn.Connection == nil {
		err := netcon.Conn.Connect()
		if err != nil {
			log.Fatalf("Cannot create network connection: %v\n", err)
		}
		defer netcon.Conn.Close()
	}

	fmt.Println("Running query")
	zoom := getZoom()
	pan, tilt := getPanTilt()

	fmt.Printf(outputTemplate, pan, tilt, zoom)
}
