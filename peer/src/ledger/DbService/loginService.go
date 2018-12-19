package DbService

import (
	"ledger/DbDao"
)

func QueryLogin(username string, password string, flag int) (map[string]string, bool, error) {
	//map[string]string {

	if flag == 0 {
		//QueryResult
		QueryResult, err := DbDao.QueryForMap("select * from tb_Admin where username=? and password=?", username, password)
		if QueryResult != nil {
			return QueryResult, true, err
		}
		return QueryResult, false, err
	} else {
		QueryResult, err := DbDao.QueryForMap("select * from tb_User where username=? and password=?", username, password)
		if QueryResult != nil {
			return QueryResult, true, err
		}
		return QueryResult, false, err
	}

}
