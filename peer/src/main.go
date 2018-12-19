package main

import (
	"chat"
	//"common/etcd"
	"config"
	"consensus"
	"ledger/Initi"
	"logging"
	"txmgr"
	"web"
)

func main() {
	config.Initialize()
	logging.Initialize()
	chat.Initialize()
	consensus.Initialize()
	txmgr.Initialize()
	txmgr.CopyrightTxHandler("123456", "Jane", "Jack")
	Initi.Initialize()
	web.Initialize()
}
