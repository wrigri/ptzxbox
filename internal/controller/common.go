package controller

import (
	"math"

	"github.com/wrigri/ptzxbox/internal/config"

	"github.com/harry1453/go-xinput/xinput"
)

type Joystick struct {
	*xinput.Joystick
	Name string
}

func (js *Joystick) Direction() string {
	angle := math.Atan2(float64(js.NormY()), float64(js.NormX())) * 180 / math.Pi
	direction := ""
	if math.Abs(angle) <= 22.5 {
		direction = "r"
	} else if angle > 22.5 && angle < 67.5 {
		direction = "ur"
	} else if angle < -22.5 && angle > -67.5 {
		direction = "dr"
	} else if angle >= 67.5 && angle < 112.5 {
		direction = "u"
	} else if angle <= -67.5 && angle >= -112.5 {
		direction = "d"
	} else if angle > 112.5 && angle < 157.5 {
		direction = "ul"
	} else if angle < -112.5 && angle > -157.5 {
		direction = "dl"
	} else {
		direction = "l"
	}
	return direction
}

func (js *Joystick) Distance() float32 {
	x := js.NormX()
	y := js.NormY()
	dis := float32(math.Sqrt(float64(x*x + y*y)))
	if dis > 1.0 {
		dis = 1.0
	}
	return dis
}

func (js *Joystick) DeadZones() {
	if math.Abs(float64(js.X)) < float64(config.Config.Controller.DeadZones.Low) {
		js.X = config.Config.Controller.DeadZones.Low
	}
	if math.Abs(float64(js.Y)) < float64(config.Config.Controller.DeadZones.Low) {
		js.Y = config.Config.Controller.DeadZones.Low
	}
	if js.X > config.Config.Controller.DeadZones.High {
		js.X = config.Config.Controller.DeadZones.High
	}
	if js.Y > config.Config.Controller.DeadZones.High {
		js.Y = config.Config.Controller.DeadZones.High
	}
	if js.X < -config.Config.Controller.DeadZones.High {
		js.X = -config.Config.Controller.DeadZones.High
	}
	if js.Y < -config.Config.Controller.DeadZones.High {
		js.Y = -config.Config.Controller.DeadZones.High
	}
}

func (js *Joystick) NormX() float32 {
	return newNormalizeJoystickValues(js.X)
}

func (js *Joystick) NormY() float32 {
	return newNormalizeJoystickValues(js.Y)
}

type DPad struct {
	right bool
	left  bool
	up    bool
	down  bool
}

func (dpad *DPad) Direction() string {
	direction := ""
	if dpad.right {
		if dpad.up {
			direction = "ur"
		} else if dpad.down {
			direction = "dr"
		} else {
			direction = "r"
		}
	} else if dpad.left {
		if dpad.up {
			direction = "ul"
		} else if dpad.down {
			direction = "dl"
		} else {
			direction = "l"
		}
	} else if dpad.up {
		direction = "u"
	} else if dpad.down {
		direction = "d"
	}
	return direction
}

func newNormalizeJoystickValues(value float32) float32 {
	d := config.Config.Controller.DeadZones.High - config.Config.Controller.DeadZones.Low
	var n, norm float32
	if value < 0 {
		n = -value - config.Config.Controller.DeadZones.Low
		norm = -n / d
	} else {
		n = value - config.Config.Controller.DeadZones.Low
		norm = n / d
	}
	return norm
}
