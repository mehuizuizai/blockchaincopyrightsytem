package restapi

import (
	"logging"
	"web/framework"
)

var logger = logging.MustGetLogger()

func InitRestApiService(webService *framework.WebService) {
	webService.RegisterAPI(copyrightTxAPI)
	webService.RegisterAPI(loginAPI)
	webService.RegisterAPI(userQueryAPI)
	webService.RegisterAPI(userRgisterAPI)
	webService.RegisterAPI(workListAPI)
	webService.RegisterAPI(workOriginateAPI)
	webService.RegisterAPI(workPutAPI)
	webService.RegisterAPI(workQueryAPI)
}
