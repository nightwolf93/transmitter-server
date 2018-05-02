package main

import (
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.SetupConfig()
	logging.SetupLog()

	log.Infof("Config value : %s", config.GetConfig().Log.Path)
}
