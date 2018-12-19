package DbService

import (
	"fmt"
	"ledger/DbDao"
	"ledger/DbUtil"
	"strconv"
	"time"
)

func UpdateWorkOwner(from_Name, to_Name, WorkID string, txid string) (int, int, bool, error) { //txId need Update
	// !!!!NEED TO  SEE IF EXIST?
	//-------------------------
	queryResult1, _ := DbDao.QueryForMap("SELECT tb_User.userID  from tb_User  where tb_User.username =?", from_Name)
	queryResult2, _ := DbDao.QueryForMap("SELECT tb_User.userID  from tb_User  where tb_User.username =?", to_Name)
	from_id := queryResult1["userID"]
	to_id := queryResult2["userID"]
	from_id_int, _ := strconv.Atoi(from_id)
	to_id_int, _ := strconv.Atoi(to_id)
	//-----------------------------
	if from_id_int == 0 || to_id_int == 0 {
		return 0, 0, false, nil
	}
	//------------------------
	err := DbDao.UpdateDb("UPDATE tb_Work set owner =? , txID = ? where WorkID =? and owner = ?", to_id_int, txid, WorkID, from_id_int)
	if err == nil {
		return from_id_int, to_id_int, true, err
	}
	return 0, 0, false, err
}
func UpdateWorkOwner_PreExe(from_Name, to_Name, WorkID string) []byte {
	workInfo := make([]string, 0)
	//-------------
	queryResult1, _ := DbDao.QueryForMap("SELECT tb_User.userID  from tb_User  where tb_User.username =?", from_Name)
	queryResult2, _ := DbDao.QueryForMap("SELECT tb_User.userID  from tb_User  where tb_User.username =?", to_Name)
	from_id := queryResult1["userID"]
	to_id := queryResult2["userID"]
	from_id_int, _ := strconv.Atoi(from_id)
	to_id_int, _ := strconv.Atoi(to_id)
	//--------------------
	_, queryWorkInfo := DbDao.UpdateDb_PreExe("UPDATE tb_Work set owner =?  where WorkID =? and owner = ?", to_id_int, WorkID, from_id_int)
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

func createTxhash(from_Name string, to_Name string, workId string, timeNow time.Time, txid string) string {
	txInfo := make([]string, 0)
	txInfo = append(txInfo, from_Name)
	txInfo = append(txInfo, to_Name)
	txInfo = append(txInfo, workId)
	txInfo = append(txInfo, timeNow.String())
	txHash := DbUtil.Str_Sha256_String(txInfo)
	return txHash
}

//func InsertWorkEntry(workName string, ownerid int, adminid int, timeNow time.Time, txID int) {
//	//admin  WorkEntry admin
//	DbDao.InsertToDb("INSERT INTO tb_Work(workname,owner,admin,time,txID) VALUES (?,?,?,?,?)", workName, ownerid, adminid, timeNow, txID)

//}
