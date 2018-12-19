package txmgr

import (
	"ledger/DbService"
)

func UserRegisterHandler(username, password, IDNumber, PhoneNumber string) error {
	//TODO call ledger interface.
	ok, err := DbService.InsertRegister()
	if !ok {
		logger.Error(err)
		return err
	}
	return nil
}
