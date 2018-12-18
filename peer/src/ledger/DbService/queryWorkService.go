package DbService

import (
	"DbDao"
)

func QueryWork(workname1 string) (workId string, workname string, owner string, phoneNumber string, admin string, time string) {
	//	QueryResult, _ := DbDao.QueryForMap("select * from tb_Work where workname=?", workname)
	QueryResult, _ := DbDao.QueryForMap("select * from tb_Work ,tb_User where tb_Work.owner =tb_User.userID and  tb_Work.workname=?", workname1)
	WorkId := QueryResult["WorkID"]
	WorkName := QueryResult["workname"]
	Owner := QueryResult["owner"]
	PhoneNumber := QueryResult["phoneNumber"]
	Admin := QueryResult["admin"]
	Time := QueryResult["time"]
	return WorkId, WorkName, Owner, PhoneNumber, Admin, Time
}
