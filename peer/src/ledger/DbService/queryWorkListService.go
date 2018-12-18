package DbService

import (
	"DbDao"
	"fmt"
)

//type workList struct {
//	WorkId   string
//	WorkName string
//	PutTime  string
//}

func QueryWorkList(username string) []map[string]string {
	//	QueryResult, _ := DbDao.QueryForMapSlice("select * from tb_Work where owner=? ", owner)
	QueryResult, _ := DbDao.QueryForMapSlice("SELECT  WorkID,workname,time FROM tb_Work , tb_User where tb_Work.owner = tb_User.userID and tb_User.username = ?", username)
	for k, _ := range QueryResult {
		fmt.Println(QueryResult[k]["WorkID"])
		fmt.Println(QueryResult[k]["workname"])
		fmt.Println(QueryResult[k]["time"])
	}
	return QueryResult
}
