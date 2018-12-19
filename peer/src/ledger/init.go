package ledger

import (
	"DbDao"
	"DbUtil"
	//	"DbService"
)

func Initialize() {
	DbUtil.Create_Table_Admin()
	DbUtil.Create_Table_User()
	DbUtil.Create_Table_Tx()
	DbUtil.Create_Table_Work()
	DbUtil.DbDao.NewBlockchain()
	//	DbService.TimingBlockInit()
}
