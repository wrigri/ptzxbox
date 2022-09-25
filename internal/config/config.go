package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"gopkg.in/yaml.v2"
)

const configFileName = "config.yaml"

// Runtime Config
var Interactive bool
var Logger service.Logger
var ManualFocus bool

// Config File
var Config ConfigData

type Network struct {
	IpAddress        string `yaml:"ipAddress"`
	ShortCmdPort     string `yaml:"shortCmdPort"`
	LongCmdPort      string `yaml:"longCmdPort"`
	UseShortCommands bool   `yaml:"useShortCommands"`
}

type DeadZones struct {
	Low  float32 `yaml:"low"`
	High float32 `yaml:"high"`
}

type Controller struct {
	DeadZones DeadZones `yaml:"deadZones"`
}

type Camera struct {
	MaxPanSpeed      int `yaml:"maxPanSpeed"`
	MinPanSpeed      int `yaml:"minPanSpeed"`
	DefaultPanSpeed  int `yaml:"defaultPanSpeed"`
	MinTiltSpeed     int `yaml:"minTiltSpeed"`
	MaxTiltSpeed     int `yaml:"maxTiltSpeed"`
	DefaultTiltSpeed int `yaml:"defaultTiltSpeed"`
	MaxZoomSpeed     int `yaml:"maxZoomSpeed"`
	MinZoomSpeed     int `yaml:"minZoomSpeed"`
	DefaultZoomSpeed int `yaml:"defaultZoomSpeed"`
}

type Preset struct {
	Pan       int `yaml:"pan"`
	Tilt      int `yaml:"tilt"`
	Zoom      int `yaml:"zoom"`
	PanSpeed  int `yaml:"panSpeed"`
	TiltSpeed int `yaml:"tiltSpeed"`
}

type Presets struct {
	PanSpeed  int    `yaml:"panSpeed"`
	TiltSpeed int    `yaml:"tiltSpeed"`
	A         Preset `yaml:"A"`
	B         Preset `yaml:"B"`
	X         Preset `yaml:"X"`
	Y         Preset `yaml:"Y"`
}

type ConfigData struct {
	Debug      bool `yaml:"debug"`
	Network    Network
	Controller Controller
	Camera     Camera
	Presets    Presets
}

func (n *Network) Port() string {
	if n.UseShortCommands {
		return n.ShortCmdPort
	}
	return n.LongCmdPort

}

func (cfg *ConfigData) Update() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(yamlFile), &cfg)
	if err != nil {
		return err
	}

	if Interactive && cfg.Debug {
		Logger.Infof("Updated Config:\n%+v", cfg)
	}
	return nil

}

func getConfigPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir, _ := filepath.Split(execPath)
	return filepath.Join(dir, configFileName), nil
}
