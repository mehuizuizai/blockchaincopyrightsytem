package restapi

import (
	"net/http"
	"web/framework"
)

type WorkOriginateReq struct {
	WorkName string
}

type WorkOriginateRes struct {
	OwnerList []string //user name set
}

var workOriginateAPI = &framework.RESTConfig{
	Path:         API_URL_WORK_ORIGINATE,
	Method:       http.MethodPost,
	BodyTemplate: &WorkOriginateReq{},
	Callback:     workOriginateHandleFn,
}

func workOriginateHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*WorkOriginateReq)

	return http.StatusOK, nil, nil
}
