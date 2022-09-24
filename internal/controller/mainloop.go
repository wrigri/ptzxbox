package controller

import (
	"time"

	"github.com/wrigri/ptzxbox/internal/config"

	"github.com/harry1453/go-xinput/xinput"
)

var oldConnectedControllers []xinput.ControllerIndex

type ControllerState struct {
	*xinput.ControllerState
	*DPad
	LeftJoystick  *Joystick
	RightJoystick *Joystick
}

var oldState *ControllerState

func ControllerLoop() {
	if xinput.LoadError != nil {
		config.Logger.Error(xinput.LoadError)
	}

	for {
		if newConnectedControllers := xinput.GetConnectedControllers(); len(newConnectedControllers) == 0 {
			if len(newConnectedControllers) != len(oldConnectedControllers) {
				config.Logger.Info("No Controllers Connected")
				time.Sleep(time.Second)
				oldConnectedControllers = newConnectedControllers
				continue
			}
		} else {
			if len(newConnectedControllers) != len(oldConnectedControllers) {
				config.Logger.Infof("Connected Controllers: %v", newConnectedControllers)
				info, err := xinput.GetControllerInfo(newConnectedControllers[0])
				if err != nil {
					config.Logger.Warning("Cound not get Controller info")
				}
				config.Logger.Infof("Using Controller: %s (%s)", newConnectedControllers[0], info.Subtype.String())
				oldConnectedControllers = newConnectedControllers
			}
			newStateController, err := xinput.GetControllerState(newConnectedControllers[0])
			if err != nil {
				config.Logger.Warning("Could not get controller state.")
				time.Sleep(time.Second)
				continue
			}

			newDPad := &DPad{
				newStateController.Buttons.DpadRight,
				newStateController.Buttons.DpadLeft,
				newStateController.Buttons.DpadUp,
				newStateController.Buttons.DpadDown,
			}
			newLeftJoystick := new(Joystick)
			newLeftJoystick.Name = "LeftJoystick"
			newLeftJoystick.Joystick = &newStateController.LeftJoystick
			newLeftJoystick.DeadZones()

			newRightJoystick := new(Joystick)
			newRightJoystick.Name = "RightJoystick"
			newRightJoystick.Joystick = &newStateController.RightJoystick
			newRightJoystick.DeadZones()

			newState := &ControllerState{newStateController, newDPad, newLeftJoystick, newRightJoystick}
			if oldState == nil || *newState.ControllerState != *oldState.ControllerState {
				if oldState != nil {
					ButtonChange("A", newState.Buttons.A, oldState.Buttons.A)
					ButtonChange("B", newState.Buttons.B, oldState.Buttons.B)
					ButtonChange("X", newState.Buttons.X, oldState.Buttons.X)
					ButtonChange("Y", newState.Buttons.Y, oldState.Buttons.Y)
					ButtonChange("LeftShoulder", newState.Buttons.LeftShoulder, oldState.Buttons.LeftShoulder)
					ButtonChange("RightShoulder", newState.Buttons.RightShoulder, oldState.Buttons.RightShoulder)
					ButtonChange("LeftJoystick", newState.Buttons.LeftJoystick, oldState.Buttons.LeftJoystick)
					ButtonChange("RightJoystick", newState.Buttons.RightJoystick, oldState.Buttons.RightJoystick)
					ButtonChange("Start", newState.Buttons.Start, oldState.Buttons.Start)
					ButtonChange("Back", newState.Buttons.Back, oldState.Buttons.Back)
					DPadChange(*newState.DPad, *oldState.DPad)
					JoystickChange("LeftJoystick", newState.LeftJoystick, oldState.LeftJoystick)
					JoystickChange("RightJoystick", newState.RightJoystick, oldState.RightJoystick)
					TriggerChange("LeftTrigger", newState.LeftTrigger, oldState.LeftTrigger)
					TriggerChange("RightTrigger", newState.RightTrigger, oldState.RightTrigger)
				}
				oldState = newState
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}
