package restapi

import (
	"net/http"
	"web/framework"
)

type UserRgisterReq struct {
	UserName    string
	Password    []byte
	IDNumber    string
	PhoneNumber string
}

type UserRgisterRes struct {
	Result    bool
	ErrorInfo string
}

var userRgisterAPI = &framework.RESTConfig{
	Path:         API_URL_USER_REGISTER,
	Method:       http.MethodPost,
	BodyTemplate: &UserRgisterReq{},
	Callback:     userRegisterHandleFn,
}

func userRegisterHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	//body := req.Body.(*UserRgisterReq)

	return http.StatusOK, nil, nil
}
