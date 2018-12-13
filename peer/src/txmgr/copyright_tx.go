package txmgr

import (
	"chat"
	pb "chat/proto"
	"common/utils"
	"consensus"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type copyrightTx struct {
	WorkID string
	From   string
	To     string
}

var sessionMap map[string]copyrightTx = make(map[string]copyrightTx)

func CopyrightTxHandler(workId, from, to string) error {
	//put session content into cache.
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 16)
	tx := copyrightTx{
		WorkID: workId,
		From:   from,
		To:     to,
	}
	sessionMap[sessionID] = tx

	//broacast tx request.
	args := pb.CopyrightTxRequest{
		SessionID: sessionID,
		WorkID:    workId,
		From:      from,
		To:        to,
	}
	_, err := chat.SendMsg(pb.Request_COPYRIGHT_TX, args, "192.168.13.82") // there is no need to get response msg.
	if err != nil {
		logger.Error("Send message error")
		return fmt.Errorf("Send message error")
	}

	//TODO call preexecution interface, and the return value type is []byte.
	//for test
	h := sha256.New()
	h.Write([]byte("hello"))
	selfVote := h.Sum(nil)

	//trigger consensus
	selfIP, _ := utils.GetlocalIP()
	isSuccessful, isEqual, _ := consensus.StartConsensus(selfVote, selfIP, sessionID)
	if !isSuccessful {
		logger.Warning("consensus failed...")
		return fmt.Errorf("transaction is not successful")
	}

	if !isEqual {
		logger.Warning("local data have some problems, do synce")
		//TODO synce, and pass mostPeers(goroutine)
		return nil
	}

	//TODO decide whether update db really or synce according consensus result.
	return nil
}

//callback fuction, it handles copyright tx request from other peer.
func copyrightTxHandler(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.CopyrightTxRequest)
	if !ok {
		logger.Error("assert error...")
		return pb.Response_COPYRIGHT_TX, nil, fmt.Errorf("handle copyright tx msg error")
	}

	//put tx request into session map.
	sessionMap[resMsg.SessionID] = copyrightTx{
		WorkID: resMsg.WorkID,
		From:   resMsg.From,
		To:     resMsg.To,
	}

	//TODO call preexecution interface, and the return value type is []byte.

	//TODO trigger consensus

	//TODO decide whether update db really or synce according consensus result.
	return pb.Response_COPYRIGHT_TX, nil, nil
}
