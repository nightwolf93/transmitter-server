package main

import (
	"bufio"
	"os"

	"github.com/nightwolf93/transmitter-server/api"
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/logging"
	"github.com/nightwolf93/transmitter-server/net"
	"github.com/nightwolf93/transmitter-server/net/gateway"
	"github.com/nightwolf93/transmitter-server/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup utilities things
	config.SetupConfig()
	logging.SetupLog()

	log.Info("Starting the Transmitter.Server ...")

	// Initialize storage
	storage.InitDB()
	net.LoadChannelsFromDatabase()

	// Init the gateway
	gateway.InitGateway()

	// Init the api
	api.InitAPI()

	// Waiting for input
	for {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		//TODO: Handle the input command
		_ = cmd
	}
}
