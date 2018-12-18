package DbService

import (
	"DbDao"
	"fmt"
)

func Work_Originagte(workName string) []string {
	work_Originage_Name := make([]string, 0)
	workListReverse := make([]string, 0)
	QueryTxID, _ := DbDao.QueryForMap("select txID,username from tb_Work ,tb_User where tb_Work.owner = tb_User.userID and workname=?  ", workName)
	txID := QueryTxID["txID"]
	username := QueryTxID["username"]
	work_Originage_Name = append(work_Originage_Name, username)
	for txID != "" {
		QueryResult, _ := DbDao.QueryForMap("select preTxID ,username from tb_TX ,tb_User where tb_TX.from_id =tb_User.userID and tb_TX.txID=?", txID)
		preTxId := QueryResult["preTxID"]
		username := QueryResult["username"]
		if username == "" {
			break
		}
		work_Originage_Name = append(work_Originage_Name, username)
		txID = preTxId
	}
	for i := len(work_Originage_Name); i > 0; i-- {
		workListReverse = append(workListReverse, work_Originage_Name[i-1])
		fmt.Println(workListReverse)
	}
	return workListReverse
}
