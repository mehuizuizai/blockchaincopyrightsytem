package DbService

import (
	"fmt"
	"ledger/DbDao"
)

//func InsertRegister_PreExe(username string, password string, idNumber string, phoneNumber string) (bool, error) {
//	//admin  WorkEntry admin
//	err := DbDao.InsertToDb_PreExe("INSERT INTO tb_User(username,password,idNumber,phoneNumber) VALUES (?,?,?,?)", username, password, idNumber, phoneNumber)
//	//	err := DbDao.InsertToDb("INSERT INTO tb_User(username,password,idNumber,phoneNumber) VALUES (?,?,?,?)", username, password, idNumber, phoneNumber)
//	fmt.Println("err is :", err)
//	if err == nil {
//		return true, err
//	}
//	return false, err

//}
func InsertRegister(username string, password string, idNumber string, phoneNumber string) (bool, error) {
	//admin  WorkEntry admin
	err := DbDao.InsertToDb("INSERT INTO tb_User(username,password,idNumber,phoneNumber) VALUES (?,?,?,?)", username, password, idNumber, phoneNumber)
	//	err := DbDao.InsertToDb("INSERT INTO tb_User(username,password,idNumber,phoneNumber) VALUES (?,?,?,?)", username, password, idNumber, phoneNumber)
	fmt.Println("err is :", err)
	if err == nil {
		return true, err
	}
	return false, err

}
