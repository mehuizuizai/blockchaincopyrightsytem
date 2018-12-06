package restapi

import (
	"net/http"
	"web/framework"
)

type LoginReq struct {
	UserName string
	Password []byte
	Flag     uint8
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
	//body := req.Body.(*UserRgisterReq)

	return http.StatusOK, nil, nil
}
