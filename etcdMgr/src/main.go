package main

import (
	"config"
	"db"
	"fmt"
	"logging"
	"manager"
	//"github.com/op/go-logging"
)

var logger = logging.MustGetLogger()

func main() {
	forever := make(chan bool)
	config.Initialize()

	logging.Initialize()
	ok := db.Initialize()
	if !ok {
		logger.Error("db init failed!")
		return
	}
	err := manager.Initialize()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Println("start etcdMgr............")
	<-forever
}
