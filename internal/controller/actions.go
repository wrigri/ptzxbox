package controller

import (
	"github.com/wrigri/ptzxbox/internal/camctrl"
	"github.com/wrigri/ptzxbox/internal/config"
)

func ButtonChange(buttonName string, newValue bool, oldValue bool) {
	if newValue != oldValue {
		if newValue {
			// This controls what happens when a button is pressed.
			switch buttonName {
			case "Start":
				camctrl.AutoFocus(false)
			case "Back":
				config.Config.Update()
				camctrl.Reset()
			case "A":
				camctrl.GoToPresetPosition(config.Config.Presets.A)
			case "B":
				camctrl.GoToPresetPosition(config.Config.Presets.B)
			case "X":
				camctrl.GoToPresetPosition(config.Config.Presets.X)
			case "Y":
				camctrl.GoToPresetPosition(config.Config.Presets.Y)
			case "LeftShoulder":
				if config.ManualFocus {
					camctrl.FocusFar()
				} else {
					camctrl.ZoomOutDefault()
				}
			case "RightShoulder":
				if config.ManualFocus {
					camctrl.FocusNear()
				} else {
					camctrl.ZoomInDefault()
				}
			case "LeftJoystick":
				camctrl.OnePushAF()
			case "RightJoystick":
				camctrl.AutoFocus(true)
			}
		} else {
			// This controls what happens when a button is released.
			switch buttonName {
			case "LeftShoulder":
				camctrl.ZoomStop()
				camctrl.FocusStop()
			case "RightShoulder":
				camctrl.ZoomStop()
				camctrl.FocusStop()
			case "Start":
				camctrl.AutoFocus(true)
			}
		}
	}
}

func DPadChange(newDpad DPad, oldDpad DPad) {
	if newDpad != oldDpad {
		newDirection := newDpad.Direction()
		switch newDirection {
		case "r":
			camctrl.RightDefault()
		case "l":
			camctrl.LeftDefault()
		case "u":
			camctrl.UpDefault()
		case "d":
			camctrl.DownDefault()
		case "ur":
			camctrl.UpRightDefault()
		case "ul":
			camctrl.UpLeftDefault()
		case "dr":
			camctrl.DownRightDefault()
		case "dl":
			camctrl.DownLeftDefault()
		default:
			camctrl.PanTiltStop()
		}
	}
}

func JoystickChange(js string, newJoystick *Joystick, oldJoystick *Joystick) {
	if *newJoystick.Joystick != *oldJoystick.Joystick {
		dis := newJoystick.Distance()

		panScale := float32(config.Config.Camera.MaxPanSpeed - config.Config.Camera.MinPanSpeed + 1)
		tiltScale := float32(config.Config.Camera.MaxTiltSpeed - config.Config.Camera.MinTiltSpeed + 1)

		panSpeed := int((dis*panScale)+0.5) + config.Config.Camera.MinPanSpeed - 1
		tiltSpeed := int((dis*tiltScale)+0.5) + config.Config.Camera.MinTiltSpeed - 1

		if dis > 0.0001 {
			switch newJoystick.Direction() {
			case "r":
				camctrl.Right(panSpeed, tiltSpeed)
			case "l":
				camctrl.Left(panSpeed, tiltSpeed)
			case "u":
				camctrl.Up(panSpeed, tiltSpeed)
			case "d":
				camctrl.Down(panSpeed, tiltSpeed)
			case "ur":
				camctrl.UpRight(panSpeed, tiltSpeed)
			case "ul":
				camctrl.UpLeft(panSpeed, tiltSpeed)
			case "dr":
				camctrl.DownRight(panSpeed, tiltSpeed)
			case "dl":
				camctrl.DownLeft(panSpeed, tiltSpeed)
			default:
				camctrl.PanTiltStop()
			}
		} else {
			camctrl.PanTiltStop()
		}
	}
}

func TriggerChange(trig string, newTrigger float32, oldTrigger float32) {
	if newTrigger != oldTrigger {
		scale := float32(config.Config.Camera.MaxZoomSpeed - config.Config.Camera.MinZoomSpeed + 1)
		speed := int((newTrigger*scale + 0.5))

		if speed != 0 {
			zoomSpeed := speed + config.Config.Camera.MinZoomSpeed - 1
			if trig == "LeftTrigger" {
				camctrl.ZoomOut(zoomSpeed)
			} else {
				camctrl.ZoomIn(zoomSpeed)
			}
		} else {
			camctrl.ZoomStop()
		}
	}
}
