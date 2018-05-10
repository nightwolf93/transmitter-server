package storage

import (
	"github.com/nightwolf93/transmitter-server/config"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

// InitDB initialize the database store
func InitDB() {
	c, err := leveldb.OpenFile(config.GetConfig().Storage.DBFilePath, nil)
	if err != nil {
		log.Fatalf("Can't initialize the database: %s", err)
		return
	}
	db = c
	log.Infof("Database initialized (path: %s)", config.GetConfig().Storage.DBFilePath)
}

// GetDB get the database connection
func GetDB() *leveldb.DB {
	return db
}
