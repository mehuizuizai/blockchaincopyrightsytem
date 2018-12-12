package restapi

import (
	"net/http"
	"web/framework"
)

type WorkQueryReq struct {
	WorkName string
}

type WorkQueryRes struct {
	WorkId     string
	WorkName   string
	Owner      string
	OwnerPhone string
	Admin      string
	PutTime    string
}

var workQueryAPI = &framework.RESTConfig{
	Path:         API_URL_WORK_QUERY,
	Method:       http.MethodPost,
	BodyTemplate: &WorkQueryReq{},
	Callback:     workQueryHandleFn,
}

func workQueryHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*WorkQueryReq)

	return http.StatusOK, nil, nil
}
