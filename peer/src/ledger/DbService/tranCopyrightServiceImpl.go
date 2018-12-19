package DbService

import (
	"fmt"
	"ledger/DbUtil"

	"time"
)

func TranCopyright_PreExe(from_Name string, to_Name string, workId string, timeNow time.Time, txid string) string {
	result := UpdateWorkOwner_PreExe(from_Name, to_Name, workId)
	return result
}
func TranCopyright(from_Name string, to_Name string, workId string, timeNow time.Time, txid string) (bool, error) {
	// need time , txid , hash
	from_id, to_id, _, _ := UpdateWorkOwner(from_Name, to_Name, workId, txid)
	//-------------
	if from_id == 0 || to_id == 0 {
		return false, nil
	}
	//------------
	txhash := createTxhash(from_Name, to_Name, workId, timeNow, txid)
	result, err := InsertTransaction(txid, from_id, to_id, timeNow, txhash)
	if result == false {
		return result, err
	}
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
