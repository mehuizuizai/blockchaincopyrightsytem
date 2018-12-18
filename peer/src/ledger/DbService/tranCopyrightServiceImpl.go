package DbService

import (
	"DbUtil"
	"fmt"
	"time"
)

func TranCopyright_PreExe(from_id int, to_id int, workId string) string {
	result := UpdateWorkOwner_PreExe(from_id, to_id, workId)
	return result
}
func TranCopyright(from_id int, to_id int, workId string, timeNow time.Time, txid string) (bool, error) {
	// need time , txid , hash
	UpdateWorkOwner(from_id, to_id, workId, txid)
	txhash := createTxhash(from_id, to_id, workId, timeNow, txid)
	result, err := InsertTransaction(txid, from_id, to_id, timeNow, txhash)
	//-------------  insert the txhash to
	postRead := []string{}
	DbUtil.Load(&postRead, "txhash")
	fmt.Println("txhash", postRead)
	postRead = append(postRead, txhash)
	DbUtil.Store(postRead, "txhash")
	//------------
	//	Tx_Enter_routinePool(txhash)
	return result, err
}

func TimingBlock() bool {
	postRead := []string{}
	DbUtil.Load(&postRead, "txhash")
	fmt.Println("txhash", postRead)
	fmt.Println("txhash[0]", postRead[0])
	//	if postRead[0] == "" {
	//		return false
	//	}
	Tx_Enter_routinePool(postRead[0])
	return true
}
func TimingBlockInit() {
	for {
		now := time.Now()
		next := now.Add(time.Minute * 10)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		fmt.Println("timeBlock begin", time.Now())
		TimingBlock()
		//		result := TimingBlock()
		//		if result == false {
		//			break
		//		}
	}
	//	fmt.Println("have no block")
}
