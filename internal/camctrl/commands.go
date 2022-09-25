package camctrl

import (
	"github.com/wrigri/ptzxbox/internal/config"
	"github.com/wrigri/ptzxbox/internal/util"
)

func ZoomIn(speed int) {
	if speed > 7 {
		speed = 7
	}
	cmdentry := byte(0x20 + speed)
	command := []byte{0x81, 0x01, 0x04, 0x07, cmdentry, 0xFF}
	sendCommand(command)
}

func ZoomOut(speed int) {
	if speed > 7 {
		speed = 7
	}
	cmdentry := byte(0x30 + speed)
	command := []byte{0x81, 0x01, 0x04, 0x07, cmdentry, 0xFF}
	sendCommand(command)
}

func ZoomStop() {
	command := []byte{0x81, 0x01, 0x04, 0x07, 0x00, 0xFF}
	sendCommand(command)
}

func basePanTilt(panSpeed, tiltSpeed int, dirBytes []byte) []byte {
	cmdStart := []byte{0x81, 0x01, 0x06, 0x01, byte(panSpeed), byte(tiltSpeed)}
	command := append(cmdStart, dirBytes...)
	command = append(command, byte(0xFF))
	return command
}

func Right(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x02, 0x03}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func Left(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x01, 0x03}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func Up(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x03, 0x01}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func Down(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x03, 0x02}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func UpRight(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x02, 0x01}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func UpLeft(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x01, 0x01}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func DownRight(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x02, 0x02}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func DownLeft(panSpeed, tiltSpeed int) {
	dirBytes := []byte{0x01, 0x02}
	command := basePanTilt(panSpeed, tiltSpeed, dirBytes)
	sendCommand(command)
}

func PanTiltStop() {
	dirBytes := []byte{0x03, 0x03}
	command := basePanTilt(1, 1, dirBytes)
	sendCommand(command)
}

func FocusNear() {
	command := []byte{0x81, 0x01, 0x04, 0x08, 0x03, 0xFF}
	sendCommand(command)
}

func FocusFar() {
	command := []byte{0x81, 0x01, 0x04, 0x08, 0x02, 0xFF}
	sendCommand(command)
}

func FocusStop() {
	command := []byte{0x81, 0x01, 0x04, 0x08, 0x00, 0xFF}
	sendCommand(command)
}

func AutoFocus(on bool) {
	if on {
		config.ManualFocus = false
		command := []byte{0x81, 0x01, 0x04, 0x38, 0x02, 0xFF}
		sendCommand(command)
	} else {
		config.ManualFocus = true
		command := []byte{0x81, 0x01, 0x04, 0x38, 0x03, 0xFF}
		sendCommand(command)
	}
}

func OnePushAF() {
	command := []byte{0x81, 0x01, 0x04, 0x38, 0x04, 0xFF}
	sendCommand(command)
}

func Home() {
	command := []byte{0x81, 0x01, 0x06, 0x04, 0xFF}
	sendCommand(command)
}

func Reset() {
	command := []byte{0x81, 0x01, 0x06, 0x05, 0xFF}
	sendCommand(command)
}

func GoToPresetPosition(preset config.Preset) {
	panBytes := util.GetBytesFromPosHex(preset.Pan)
	tiltBytes := util.GetBytesFromPosHex(preset.Tilt)
	zoomBytes := util.GetBytesFromPosHex(preset.Zoom)

	var panSpeed int
	if preset.PanSpeed != 0 {
		panSpeed = preset.PanSpeed
	} else if config.Config.Presets.PanSpeed != 0 {
		panSpeed = config.Config.Presets.PanSpeed
	} else {
		panSpeed = config.Config.Camera.DefaultPanSpeed
	}

	var tiltSpeed int
	if preset.TiltSpeed != 0 {
		tiltSpeed = preset.TiltSpeed
	} else if config.Config.Presets.TiltSpeed != 0 {
		tiltSpeed = config.Config.Presets.TiltSpeed
	} else {
		tiltSpeed = config.Config.Camera.DefaultTiltSpeed
	}

	zoomCommand := []byte{0x81, 0x01, 0x04, 0x47}
	zoomCommand = append(zoomCommand, zoomBytes...)
	zoomCommand = append(zoomCommand, 0xFF)
	sendCommand(zoomCommand)

	panTiltCommand := []byte{0x81, 0x01, 0x06, 0x02}
	panTiltCommand = append(panTiltCommand, byte(panSpeed), byte(tiltSpeed))
	panTiltCommand = append(panTiltCommand, panBytes...)
	panTiltCommand = append(panTiltCommand, tiltBytes...)
	panTiltCommand = append(panTiltCommand, 0xFF)
	sendCommand(panTiltCommand)
}

func GoToPosition(pan, tilt, zoom, panSpeed, tiltSpeed int) {
	panBytes := util.GetBytesFromPosHex(pan)
	tiltBytes := util.GetBytesFromPosHex(tilt)
	zoomBytes := util.GetBytesFromPosHex(zoom)

	zoomCommand := []byte{0x81, 0x01, 0x04, 0x47}
	zoomCommand = append(zoomCommand, zoomBytes...)
	zoomCommand = append(zoomCommand, 0xFF)
	sendCommand(zoomCommand)

	panTiltCommand := []byte{0x81, 0x01, 0x06, 0x02}
	panTiltCommand = append(panTiltCommand, byte(panSpeed), byte(tiltSpeed))
	panTiltCommand = append(panTiltCommand, panBytes...)
	panTiltCommand = append(panTiltCommand, tiltBytes...)
	panTiltCommand = append(panTiltCommand, 0xFF)
	sendCommand(panTiltCommand)

}

func ZoomInDefault() {
	ZoomIn(config.Config.Camera.DefaultZoomSpeed)
}
func ZoomOutDefault() {
	ZoomOut(config.Config.Camera.DefaultZoomSpeed)
}
func RightDefault() {
	Right(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func LeftDefault() {
	Left(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func UpDefault() {
	Up(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func DownDefault() {
	Down(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func UpRightDefault() {
	UpRight(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func UpLeftDefault() {
	UpLeft(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func DownRightDefault() {
	DownRight(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
func DownLeftDefault() {
	DownLeft(config.Config.Camera.DefaultPanSpeed, config.Config.Camera.DefaultTiltSpeed)
}
