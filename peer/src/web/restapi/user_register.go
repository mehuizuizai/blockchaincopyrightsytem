package restapi

import (
	"net/http"
	"txmgr"
	"web/framework"
)

type UserRgisterReq struct {
	UserName    string
	Password    string
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
	body := req.Body.(*UserRgisterReq)

	err := txmgr.UserRegisterHandler(body.UserName, body.Password, body.IDNumber, body.PhoneNumber)

	resp := &UserRgisterRes{}
	if err != nil {
		resp.Result = false
		resp.ErrorInfo = "error"
		return http.StatusOK, resp, nil
	}

	resp.Result = true
	resp.ErrorInfo = ""

	return http.StatusOK, resp, nil
}
