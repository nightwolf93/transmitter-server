package main

import (
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/logging"
)

func main() {
	config.SetupConfig()
	logging.SetupLog()
}
