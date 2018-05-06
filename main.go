package main

import (
	"bufio"
	"os"

	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/logging"
	"github.com/nightwolf93/transmitter-server/net/gateway"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup utilities things
	config.SetupConfig()
	logging.SetupLog()

	log.Info("Starting the Transmitter.Server ...")

	// Init the gateway
	gateway.InitGateway()

	// Waiting for input
	for {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		//TODO: Handle the input command
		_ = cmd
	}
}
