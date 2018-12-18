package DbService

import (
	"DbDao"
	"DbUtil"
	"time"
)

func InsertWorkEntry_PreExe(workID string, workName string, ownerid int, adminid int, timeNow time.Time) string {
	//admin  WorkEntry admin
	// return  two bool ?
	result := InsertWork_PreExe(workID, workName, ownerid, adminid, timeNow)
	return result
}
func InsertWorkEntry(workID string, workName string, ownerid int, adminid int, timeNow time.Time, txID string, hash string) (bool, error) {
	//admin  WorkEntry admin
	InsertWork(workID, workName, ownerid, adminid, timeNow, txID)
	//	result, err := InsertTransaction(workID, adminid, ownerid, timeNow, txID, hash)
	result, err := InsertTransaction(txID, adminid, ownerid, timeNow, hash)
	return result, err
}

func InsertWork_PreExe(workID string, workName string, ownerid int, adminid int, timeNow time.Time) string {
	workInfo := make([]string, 0)
	queryWorkInfo := DbDao.InsertToDb_PreExe("INSERT INTO tb_Work(WorkID,workname,owner,admin,time) VALUES (?,?,?,?,?)", workID, workName, ownerid, adminid, timeNow)
	for _, v := range queryWorkInfo {
		workInfo = append(workInfo, v)
	}
	resultForConsensus := DbUtil.Str_Sha256(workInfo)
	return resultForConsensus

}
func InsertWork(workID string, workName string, ownerid int, adminid int, timeNow time.Time, txID string) error {
	err := DbDao.InsertToDb("INSERT INTO tb_Work(WorkID,workname,owner,admin,time,txID) VALUES (?,?,?,?,?,?)", workID, workName, ownerid, adminid, timeNow, txID)
	return err
}
