package txmgr

import (
	"chat"
	pb "chat/proto"
	"logging"
	"sync"
)

var logger = logging.MustGetLogger()

type txSession struct {
	sync.RWMutex
	txSessionMap map[string]interface{}
}

var txsession *txSession

func Initialize() {
	//register msg to chat module.
	chat.RegisterMsg(pb.Request_COPYRIGHT_TX, copyrightTxCallback, pb.Response_COPYRIGHT_TX)

	txsession = newTxSession()
}

func newTxSession() *txSession {
	return &txSession{
		txSessionMap: make(map[string]interface{}),
	}
}
