package restapi

import (
	"net/http"
	"txmgr"
	"web/framework"
)

type WorkPutReq struct {
	WorkName  string
	Owner     string
	Admin     string
	Timestamp string
}

type WorkPutRes struct {
	Result    bool
	ErrorInfo string
}

var workPutAPI = &framework.RESTConfig{
	Path:         API_URL_WORK_PUT,
	Method:       http.MethodPost,
	BodyTemplate: &WorkPutReq{},
	Callback:     workPutHandleFn,
}

func workPutHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	body := req.Body.(*WorkPutReq)

	err := txmgr.WorkPutHandler(body.WorkName, body.Owner, body.Admin, body.Timestamp)
	resp := &WorkPutRes{}
	if err != nil {
		resp.Result = false
		resp.ErrorInfo = "error"
		return http.StatusOK, resp, nil
	}

	resp.Result = true
	resp.ErrorInfo = ""
	return http.StatusOK, resp, nil

	return http.StatusOK, nil, nil
}
