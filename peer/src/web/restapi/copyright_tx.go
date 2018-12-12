package restapi

import (
	"net/http"
	"txmgr"
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
	body := req.Body.(*CopyrightTxReq)

	err := txmgr.CopyrightTxHandler(body.WorkId, body.From, body.To)
	resp := &CopyrightTxRes{}
	if err != nil {
		resp.Result = false
		resp.ErrorInfo = "error"
		return http.StatusOK, resp, nil
	}

	resp.Result = true
	resp.ErrorInfo = ""
	return http.StatusOK, resp, nil
}
