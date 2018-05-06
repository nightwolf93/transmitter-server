package logging

import (
	"io"
	"os"

	"github.com/nightwolf93/transmitter-server/config"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// SetupLog initialize the logger
func SetupLog() {
	logConfig := config.GetConfig().Log

	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   logConfig.Path,
		MaxSize:    logConfig.MaxFileSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(mw)
}
