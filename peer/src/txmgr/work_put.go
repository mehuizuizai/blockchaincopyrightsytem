package txmgr

import (
	"chat"
	pb "chat/proto"
	"consensus"
	"errors"
	"ledger/DbService"
	"math/rand"
	"strconv"
	"time"
)

type workPut struct {
	workID    string
	workName  string
	owner     string
	admin     string
	timestamp time.Time
}

func WorkPutHandler(workName, owner, admin, timestamp string) error {
	//create a session for this tx.
	sessionID := strconv.FormatInt(time.Now().UnixNano(), 16)

	//broadcast tx requst.
	args := &pb.WorkPutTxRequest{
		SessionID: sessionID,
		WorkName:  workName,
		Owner:     owner,
		Admin:     admin,
		Timestamp: timestamp,
	}
	_, err := chat.SendMsg(pb.Request_WORK_PUT, args, "192.168.13.82") // there is no need to get response msg.
	if err != nil {
		logger.Error("Send message error")
		return errors.New("Send message error")
	}

	err = workPutHandler(workName, owner, admin, timestamp, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func workPutHandler(workName, owner, admin, timestamp, sessionID string) error {
	//parse timestamp to time.Time.
	time_format := "2006-01-02 15:04:05"
	timestamp_, err := time.Parse(time_format, timestamp)
	if err != nil {
		logger.Error("parse timestamp error")
		return errors.New("timestamp format is error")
	}
	//create work id.
	workID := createID(timestamp_)

	tx := workPut{
		workID:    workID,
		workName:  workName,
		owner:     owner,
		admin:     admin,
		timestamp: timestamp_,
	}
	txsession.Lock()
	txsession.txSessionMap[sessionID] = tx
	txsession.Unlock()

	//TODO call preexecution interface, and the return value type is []byte
	selfVote := DbService.WorkEntry_PreExe(workID, workName, owner, admin, time.Now(), "")
	//trigger consensus
	isSuccessful, isEqual, _ := consensus.StartConsensus(selfVote, sessionID)
	if !isSuccessful {
		logger.Warning("consensus failed...")
		return errors.New("transaction is not successful")
	}

	if !isEqual {
		logger.Warning("local data have some problems, do synce")
		//TODO synce, and pass mostPeers(goroutine)
		return nil
	}

	//create tx id.
	txID := createID(time.Now())
	//TODO decide whether update db really or synce according consensus result.
	DbService.WorkEntry(workID, workName, owner, admin, time.Now(), txID)
	return nil
}

func createID(timestamp time.Time) string {
	kind := 0
	size := 10
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(timestamp.UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func workPutCallback(args interface{}) (pb.Response_Type, interface{}, error) {
	resMsg, ok := args.(pb.WorkPutTxRequest)
	if !ok {
		logger.Error("assert error...")
		return pb.Response_WORK_PUT, pb.WorkPutTxResponse{}, errors.New("handle workput tx msg error")
	}

	go workPutHandler(resMsg.WorkName, resMsg.Owner, resMsg.Admin, resMsg.Timestamp, resMsg.SessionID)

	return pb.Response_WORK_PUT, pb.WorkPutTxResponse{}, nil
}
