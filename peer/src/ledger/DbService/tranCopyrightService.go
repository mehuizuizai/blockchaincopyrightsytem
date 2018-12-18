package DbService

import (
	"DbDao"
	"DbUtil"
	"fmt"
	"strconv"
	"time"
)

func UpdateWorkOwner(from_id int, to_id int, WorkID string, txid string) (bool, error) { //txId need Update
	// !!!!NEED TO  SEE IF EXIST?
	err := DbDao.UpdateDb("UPDATE tb_Work set owner =? , txID = ? where WorkID =? and owner = ?", to_id, txid, WorkID, from_id)
	if err == nil {
		return true, err
	}
	return false, err
}
func UpdateWorkOwner_PreExe(from_id int, to_id int, WorkID string) string {
	workInfo := make([]string, 0)
	_, queryWorkInfo := DbDao.UpdateDb_PreExe("UPDATE tb_Work set owner =?  where WorkID =? and owner = ?", to_id, WorkID, from_id)
	for _, v := range queryWorkInfo {
		workInfo = append(workInfo, v)
	}
	resultForConsensus := DbUtil.Str_Sha256(workInfo)
	return resultForConsensus
}
func InsertTransaction(txID string, from_id int, to_id int, timeNow time.Time, hash string) (bool, error) {
	prevTxID, _ := DbDao.QueryForMap("SELECT txID FROM tb_TX where to_id =? ", from_id)
	//this work last Transc
	//	to_id_str := strconv.Itoa(to_id)
	fmt.Println(prevTxID)
	fmt.Println(prevTxID["txID"])
	err := DbDao.InsertToDb("INSERT INTO tb_TX(txID,preTxID,from_id,to_id,timestamp,hash)VALUES(?,?,?,?,?,?)", txID, prevTxID["txID"],
		from_id, to_id, timeNow, hash)
	if err == nil {
		return true, err
	}
	return false, err
}

func createTxhash(from_id int, to_id int, workId string, timeNow time.Time, txid string) string {
	txInfo := make([]string, 0)
	txInfo = append(txInfo, strconv.Itoa(from_id))
	txInfo = append(txInfo, strconv.Itoa(to_id))
	txInfo = append(txInfo, workId)
	txInfo = append(txInfo, timeNow.String())
	txHash := DbUtil.Str_Sha256(txInfo)
	return txHash
}

//func InsertWorkEntry(workName string, ownerid int, adminid int, timeNow time.Time, txID int) {
//	//admin  WorkEntry admin
//	DbDao.InsertToDb("INSERT INTO tb_Work(workname,owner,admin,time,txID) VALUES (?,?,?,?,?)", workName, ownerid, adminid, timeNow, txID)

//}
