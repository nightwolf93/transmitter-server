package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

const (
	configPath = "./transmitter.conf"
)

// Config configuration of the application
type Config struct {
	Log Log `toml:"log"`
}

// Log is the log section of the config
type Log struct {
	Path        string `toml:"path"`
	MaxFileSize int    `toml:"maxFileSize"`
	MaxBackups  int    `toml:"maxBackups"`
	MaxAge      int    `toml:"maxAge"`
}

var config *Config

// SetupConfig load the configuration file
func SetupConfig() {
	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Can't load the configuration file: %s", err)
		os.Exit(-1)
	}
	config = &Config{}
	_, err = toml.Decode(string(dat), &config)
	if err != nil {
		log.Errorf("Can't decode the configuration file: %s", err)
		os.Exit(-1)
	}
}

// GetConfig get the configuration file of the application
func GetConfig() *Config {
	return config
}
