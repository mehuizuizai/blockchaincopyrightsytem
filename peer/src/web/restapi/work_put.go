package restapi

import (
	"net/http"
	"web/framework"
)

type WorkPutReq struct {
	WorkName string
	Owner    string
	Admin    string
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
	//body := req.Body.(*WorkPutReq)

	return http.StatusOK, nil, nil
}
