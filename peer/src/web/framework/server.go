package framework

import (
	"logging"

	"github.com/labstack/echo"
)

var logger = logging.MustGetLogger()

// WebService is the main component which holds the actual HTTP server
//internally and register all entries
type WebService struct {
	listenAddr    string
	conns         int //http conn limit
	server        *echo.Echo
	verifyHandler VerifyHandler
}

// NewWebService creates a new *WebService

func NewWebService(listenAddr string, conns int, verifyHandler VerifyHandler) *WebService {

	service := new(WebService)
	e := echo.New()
	service.server = e
	service.listenAddr = listenAddr
	service.conns = conns
	service.verifyHandler = verifyHandler
	return service
}

// Serve starts the server and blocks
func (server *WebService) Serve() error {
	listener, err := newLimitListener(server.listenAddr, server.conns)
	if err != nil {
		logger.Error("WebService Serve:", err)
		return err
	}
	server.server.Listener = listener
	return server.server.Start(server.listenAddr)
}

func (server *WebService) GetListenAddr() string {
	return server.listenAddr
}
