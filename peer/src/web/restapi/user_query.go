package restapi

import (
	"net/http"
	"web/framework"
)

type UserQueryReq struct {
	UserName string
}

type UserQueryRes struct {
	UserName    string
	PhoneNumber string
	IDNumber    string
}

var userQueryAPI = &framework.RESTConfig{
	Path:         API_URL_USER_QUERY,
	Method:       http.MethodPost,
	BodyTemplate: &UserQueryReq{},
	Callback:     userQueryHandleFn,
}

func userQueryHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*UserQueryReq)

	return http.StatusOK, nil, nil
}
