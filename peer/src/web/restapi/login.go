package restapi

import (
	"net/http"
	"txmgr"
	"web/framework"
)

type LoginReq struct {
	UserName string
	Password string
	Flag     int
}

type LoginRes struct {
	Result    bool
	ErrorInfo string
}

var loginAPI = &framework.RESTConfig{
	Path:         API_URL_LOGIN,
	Method:       http.MethodPost,
	BodyTemplate: &LoginReq{},
	Callback:     loginHandleFn,
}

func loginHandleFn(req *framework.RESTRequest) (int, interface{}, error) {
	body := req.Body.(*LoginReq)

	err := txmgr.LoginHandler(body.UserName, body.Password, body.Flag)

	resp := &LoginRes{}
	if err != nil {
		resp.Result = false
		resp.ErrorInfo = "error"
		return http.StatusOK, resp, nil
	}

	resp.Result = true
	resp.ErrorInfo = ""

	return http.StatusOK, resp, nil
}
