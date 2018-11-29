package main

import (
	"common/etcd"
	"config"
	"logging"
	"web"
)

func main() {
	config.Initialize()
	logging.Initialize()
	etcd.Initialize()

	web.Initialize()
}
