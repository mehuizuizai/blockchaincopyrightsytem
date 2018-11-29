package db

import (
	"config"
	"database/sql"
	"logging"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var logger = logging.MustGetLogger()

var datastruct string
var db *sql.DB

func Initialize() bool {
	var err error
	db, err = sql.Open("sqlite3", config.BasePath+"etcd.db")
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	datastruct = `CREATE TABLE membersinfo(
		seq INTEGER PRIMARY KEY AUTOINCREMENT,
		ip VARCHAR(20) ,
		clientPort CHAR(6) NULL,
		peerPort CHAR(6) NULL,
		name VARCHAR(30) ,
		id VARCHAR(20) 
	);
	`

	_, err = db.Exec(datastruct)
	if err != nil {
		if strings.EqualFold(err.Error(), "table membersinfo already exists") {
			logger.Info("table membersinfo already exists")
		} else {
			logger.Error(err.Error())
			return false
		}
	}

	return true
}
