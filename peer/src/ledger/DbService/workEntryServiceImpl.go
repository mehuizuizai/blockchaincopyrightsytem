package DbService

import (
	"fmt"
	"ledger/DbUtil"
	"time"
)

func WorkEntry_PreExe(workID string, workName string, ownerName string, adminName string, timeNow time.Time, txId string) []byte {

	result := InsertWorkEntry_PreExe(workID, workName, ownerName, adminName, timeNow)
	return result
}
func WorkEntry(workID string, workName string, ownerName string, adminName string, timeNow time.Time, txId string) bool {
	txHash := createTxhash(adminName, ownerName, workName, timeNow, txId)
	result, _ := InsertWorkEntry(workID, workName, ownerName, adminName, timeNow, txId, txHash)
	if result == false {
		return result
	}
	//
	postRead := []string{}
	DbUtil.Load(&postRead, "txhash")
	fmt.Println("txhash", postRead)
	postRead = append(postRead, txHash)
	DbUtil.Store(postRead, "txhash")
	//
	return result
}
