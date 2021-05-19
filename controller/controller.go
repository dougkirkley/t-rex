package controller

import "github.com/dougkirkley/trex/config"

// Controller is a controller
type Controller interface {
	Run() error
	Config() *config.TrexConfig
}