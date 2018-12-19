// Initi project Initi.go
package Initi

import (
	//	"database/sql"
	//	"fmt"
	"ledger/DbDao"
	"ledger/DbService"
	"ledger/DbUtil"

	_ "github.com/go-sql-driver/mysql"
)

func Initialize() {
	Create_Db_blockchain()
	Create_Table_Admin()
	Create_Table_User()
	Create_Table_Tx()
	Create_Table_Work()
	Insert_Table_Admin()
	DbDao.NewBlockchain()
	DbService.TimingBlockInit()
	Insert_Table_Admin()
}
func Create_Db_blockchain() {
	db, _ := DbUtil.DbInit()
	sql := `CREATE DATABASE blockchain2`
	smt, _ := db.Prepare(sql)
	smt.Exec()
}
func Create_Table_Admin() {
	db, _ := DbUtil.DbOpen()
	sql := `CREATE TABLE tb_Admin(
	  adminID int(11) NOT NULL,
	  username char(18) NOT NULL,
	  password char(20) NOT NULL,
	  PRIMARY KEY (adminID),
	  UNIQUE KEY username_UNIQUE (username)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	smt, _ := db.Prepare(sql)
	smt.Exec()
}
func Create_Table_User() {
	db, _ := DbUtil.DbOpen()
	sql := `CREATE TABLE tb_User(
  userID int(11) NOT NULL AUTO_INCREMENT,
  username char(18) NOT NULL,
  password char(20) NOT NULL,
  idNumber varchar(18) NOT NULL,
  phoneNumber varchar(11) NOT NULL,
  PRIMARY KEY (userID),
  UNIQUE KEY username_UNIQUE(username)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;`
	smt, _ := db.Prepare(sql)
	smt.Exec()
}
func Create_Table_Tx() {
	db, _ := DbUtil.DbOpen()
	sql := `CREATE TABLE tb_TX (
  txID char(10) NOT NULL,
  preTxID char(10) NOT NULL,
  from_id int(11) NOT NULL,
  to_id int(11) NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  hash char(200) NOT NULL,
  PRIMARY KEY (txID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	smt, _ := db.Prepare(sql)
	smt.Exec()
}
func Create_Table_Work() {
	db, _ := DbUtil.DbOpen()
	sql := `CREATE TABLE tb_Work(
  WorkID char(10) NOT NULL,
  workname varchar(18) NOT NULL,
  owner int(11) NOT NULL,
  admin int(11) NOT NULL,
  time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  txID char(10) NOT NULL,
  PRIMARY KEY (WorkID),
  UNIQUE KEY workname_UNIQUE (workname),
  UNIQUE KEY txID_UNIQUE (txID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	smt, _ := db.Prepare(sql)
	smt.Exec()
}
func Insert_Table_Admin() {
	//	db, _ := DbUtil.DbOpen()
	sql := `INSERT INTO tb_Admin (adminID, username, password) VALUES (?,?,?);`
	DbDao.InsertToDb(sql, 100, "admin", "888888")

}
