package txmgr

import (
	"chat"
	pb "chat/proto"
	"logging"
)

var logger = logging.MustGetLogger()

func Initialize() {
	//register msg to chat module.
	chat.RegisterMsg(pb.Request_COPYRIGHT_TX, copyrightTxHandler, pb.Response_COPYRIGHT_TX)
}
