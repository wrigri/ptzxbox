package camctrl

import (
	"encoding/binary"

	"github.com/wrigri/ptzxbox/internal/config"
	"github.com/wrigri/ptzxbox/internal/netcon"
	"github.com/wrigri/ptzxbox/internal/util"
)

var lastCommand = []byte{0, 0, 0, 0, 0, 0}
var sequenceNum uint32 = 100

func cmpByte(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sendCommand(command []byte) {
	if cmpByte(command, lastCommand) {
		return
	}
	fullCommand := make([]byte, 0, 24)

	if config.Config.Network.UseShortCommands {
		fullCommand = command
	} else {
		payloadType := []byte{0x01, 0x00}
		payloadLen := make([]byte, 2)
		binary.BigEndian.PutUint16(payloadLen, uint16(len(command)))
		seqBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(seqBytes, sequenceNum)

		fullCommand = append(fullCommand, payloadType...)
		fullCommand = append(fullCommand, payloadLen...)
		fullCommand = append(fullCommand, seqBytes...)
		fullCommand = append(fullCommand, command...)
	}

	if config.Interactive && config.Config.Debug {
		hexCmd := util.GetHexString(fullCommand)
		config.Logger.Infof("Sending Command: '%s' to '%s'", hexCmd, netcon.Conn.Address())
	}
	netcon.Conn.SendCommand(fullCommand)
	//networkUtils.SendCommand(networkUtils.Connection, fullCommand)
	lastCommand = command
	sequenceNum += 1
}
