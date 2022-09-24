package main

import (
	"log"

	"github.com/wrigri/ptzxbox/internal/config"
	"github.com/wrigri/ptzxbox/internal/controller"
	"github.com/wrigri/ptzxbox/internal/netcon"

	"github.com/kardianos/service"
)

var sysLogger service.Logger

type program struct{}

func main() {
	if service.Interactive() {
		config.Interactive = true
	} else {
		config.Interactive = false
	}

	svcConfig := &service.Config{
		Name:        "PTZXbox",
		DisplayName: "PTZ Xbox Controller Service",
		Description: "This is a service that sends PTZ camera controls based on Xbox Controller inputs.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	config.Logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	sysLogger, err = s.SystemLogger(nil)
	if err != nil {
		config.Logger.Error(err)
	}
	err = s.Run()
	if err != nil {
		config.Logger.Error(err)
	}
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	if config.Interactive {
		config.Logger.Info("Running Interactive")
	}
	err := config.Config.Update()
	if err != nil {
		sysLogger.Errorf("Error updating config from config file: %v\n", err)
	}

	netcon.Conn.IPAddress = config.Config.Network.IpAddress
	netcon.Conn.Port = config.Config.Network.Port()
	if netcon.Conn.Connection == nil {
		err := netcon.Conn.Connect()
		if err != nil {
			sysLogger.Errorf("Cannot create network connection: %v\n", err)
		}
		defer netcon.Conn.Close()
	}

	controller.ControllerLoop()
}

func (p *program) Stop(s service.Service) error {
	return nil
}
