package txmgr

import (
	"errors"
	"ledger/DbService"
)

func LoginHandler(username string, password string, flag uint8) error {
	//TODO
	_, ok, _ := DbService.QueryLogin(username, password)
	if !ok {
		logger.Error(errors.New("login error"))
		return errors.New("login error")
	}
	return nil
}
