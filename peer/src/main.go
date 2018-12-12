package main

import (
	"chat"
	//"common/etcd"
	"config"
	"logging"
	"web"
)

func main() {
	config.Initialize()
	logging.Initialize()
	//etcd.Initialize()
	chat.Initialize()

	web.Initialize()
}
