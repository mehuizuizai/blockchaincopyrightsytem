package DbService

import (
	"DbUtil"
	"fmt"
	"time"
)

func WorkEntry_PreExe(workID string, workName string, ownerid int, adminid int, timeNow time.Time) string {

	result := InsertWorkEntry_PreExe(workID, workName, ownerid, adminid, timeNow)
	return result
}
func WorkEntry(workID string, workName string, ownerid int, adminid int, timeNow time.Time, txId string) bool {
	txHash := createTxhash(adminid, ownerid, workID, timeNow, txId)
	result, _ := InsertWorkEntry(workID, workName, ownerid, adminid, timeNow, txId, txHash)
	//
	postRead := []string{}
	DbUtil.Load(&postRead, "txhash")
	fmt.Println("txhash", postRead)
	postRead = append(postRead, txHash)
	DbUtil.Store(postRead, "txhash")
	//
	return result
}
