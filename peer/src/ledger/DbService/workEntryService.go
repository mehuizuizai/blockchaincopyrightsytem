package DbService

import (
	"ledger/DbDao"
	"ledger/DbUtil"
	"strconv"
	"time"
)

func InsertWorkEntry_PreExe(workID string, workName string, ownerName string, adminName string, timeNow time.Time) []byte {
	//admin  WorkEntry admin
	// return  two bool ?
	result := InsertWork_PreExe(workID, workName, ownerName, adminName, timeNow)
	return result
}
func InsertWorkEntry(workID string, workName string, ownerName string, adminName string, timeNow time.Time, txID string, hash string) (bool, error) {
	//admin  WorkEntry admin
	ownerid, adminid, _ := InsertWork(workID, workName, ownerName, adminName, timeNow, txID)
	//-----------------
	if ownerid == 0 || adminid == 0 {
		return false, nil
	}
	//------------
	//	result, err := InsertTransaction(workID, adminid, ownerid, timeNow, txID, hash)
	result, err := InsertTransaction(txID, adminid, ownerid, timeNow, hash)
	return result, err
}

func InsertWork_PreExe(workID string, workName string, ownerName string, adminName string, timeNow time.Time) []byte {
	workInfo := make([]string, 0)
	//--- ounername  ,admin name
	queryResult, _ := DbDao.QueryForMap("SELECT tb_User.userID , tb_Admin.adminID  from tb_User ,tb_Admin where tb_User.username =? and tb_Admin.username =?", ownerName, adminName)
	ownerid := queryResult["userID"]
	adminid := queryResult["adminID"]
	ownerid_int, _ := strconv.Atoi(ownerid)
	adminid_int, _ := strconv.Atoi(adminid)
	//
	queryWorkInfo := DbDao.InsertToDb_PreExe("INSERT INTO tb_Work(WorkID,workname,owner,admin,time) VALUES (?,?,?,?,?)", workID, workName, ownerid_int, adminid_int, timeNow)
	for _, v := range queryWorkInfo {
		workInfo = append(workInfo, v)
	}
	resultForConsensus := DbUtil.Str_Sha256(workInfo)
	return resultForConsensus

}
func InsertWork(workID string, workName string, ownerName string, adminName string, timeNow time.Time, txID string) (int, int, error) {
	queryResult, _ := DbDao.QueryForMap("SELECT tb_User.userID , tb_Admin.adminID  from tb_User ,tb_Admin where tb_User.username =? and tb_Admin.username =?", ownerName, adminName)
	ownerid := queryResult["userID"]
	ownerid_int, _ := strconv.Atoi(ownerid)
	adminid := queryResult["adminID"]
	adminid_int, _ := strconv.Atoi(adminid)
	// ---
	if adminid_int == 0 || ownerid_int == 0 {
		return ownerid_int, adminid_int, nil
	}
	//--
	err := DbDao.InsertToDb("INSERT INTO tb_Work(WorkID,workname,owner,admin,time,txID) VALUES (?,?,?,?,?,?)", workID, workName, ownerid_int, adminid_int, timeNow, txID)
	return ownerid_int, adminid_int, err
}
