package txmgr

import (
	"chat"
	pb "chat/proto"
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

//var txSessionMap map[string]copyrightTx = make(map[string]copyrightTx)

func CopyrightTxHandler(workId, from, to string) error {
	//put session content into cache.
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 16)

	//broacast tx request.
	args := &pb.CopyrightTxRequest{
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

	err = copyrightTxHandler(workId, from, to, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func copyrightTxHandler(workId, from, to, sessionID string) error {
	tx := copyrightTx{
		WorkID: workId,
		From:   from,
		To:     to,
	}
	txsession.Lock()
	txsession.txSessionMap[sessionID] = tx
	txsession.Unlock()

	//TODO call preexecution interface, and the return value type is []byte.

	//trigger consensus
	//for test
	h := sha256.New()
	h.Write([]byte("hello"))
	selfVote := h.Sum(nil)
	isSuccessful, isEqual, _ := consensus.StartConsensus(selfVote, sessionID)
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
func copyrightTxCallback(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.CopyrightTxRequest)
	if !ok {
		logger.Error("assert error...")
		return pb.Response_COPYRIGHT_TX, pb.CopyrightTxResponse{}, fmt.Errorf("handle copyright tx msg error")
	}

	go copyrightTxHandler(resMsg.WorkID, resMsg.From, resMsg.To, resMsg.SessionID)

	//TODO decide whether update db really or synce according consensus result.
	return pb.Response_COPYRIGHT_TX, pb.CopyrightTxResponse{}, nil
}
