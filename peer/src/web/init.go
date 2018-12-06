package web

import (
	"common/utils"
	"config"
	"logging"
	"web/framework"
	"web/restapi"
)

var logger = logging.MustGetLogger()

func verifyHandler(mapData *map[string]string) error {
	return nil
}

func Initialize() {

	ip, _ := utils.GetlocalIP()

	port := config.GetWebServPort()
	addr := ip + ":" + port
	conns := config.GetWebServConns()
	webService := framework.NewWebService(addr, conns, verifyHandler)

	restapi.InitRestApiService(webService)

	webService.Serve()

}
