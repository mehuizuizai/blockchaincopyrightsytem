package main

import (
	"config"
	"logging"
	"web"
)

func main() {
	config.Initialize()
	logging.Initialize()
	web.Initialize()
}
