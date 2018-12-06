package restapi

import (
	"net/http"
	"web/framework"
)

type CopyrightTxReq struct {
	WorkId string
	From   string
	To     string
}

type CopyrightTxRes struct {
	Result    bool
	ErrorInfo string
}

var copyrightTxAPI = &framework.RESTConfig{
	Path:         API_URL_COPYRIGHT_TX,
	Method:       http.MethodPost,
	BodyTemplate: &CopyrightTxReq{},
	Callback:     copyrightTxHandleFn,
}

func copyrightTxHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*CopyrightTxReq)

	return http.StatusOK, nil, nil
}
