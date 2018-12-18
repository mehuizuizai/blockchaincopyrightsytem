package DbService

import (
	"DbDao"
)

func QueryUsrname(username string) (username1 string, phoneNumber string, idNumber string) {
	QueryResult, _ := DbDao.QueryForMap("select * from tb_User where username=?", username)
	Username := QueryResult["username"]
	PhoneNumber := QueryResult["phoneNumber"]
	IdNumber := QueryResult["idNumber"]
	return Username, PhoneNumber, IdNumber
	//	return QueryResult
}
