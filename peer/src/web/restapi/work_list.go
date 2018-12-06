package restapi

import (
	"net/http"
	"web/framework"
)

type WorkListReq struct {
	UserName string
}

type Work struct {
	WorkId   string
	WorkName string
	PutTime  string
}

type WorkListRes struct {
	Works []Work
}

var workListAPI = &framework.RESTConfig{
	Path:         API_URL_WORK_LIST,
	Method:       http.MethodPost,
	BodyTemplate: &WorkListReq{},
	Callback:     workListHandleFn,
}

func workListHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*WorkListReq)

	return http.StatusOK, nil, nil
}
