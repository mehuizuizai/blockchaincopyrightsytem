package framework

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"encoding/json"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mapstructure"
)

// Callback is called given the pointer of request info struct
// and expecting a http status code, a response body and an error
// body has the same type of BodyTemplate in RESTConfig
type Callback func(*RESTRequest) (int, interface{}, error) //在echo注册的回调函数的参数为来自页面传来的参数集
type VerifyHandler func(*map[string]string) error          //验签处理函数

// RESTConfig is a configuration pack for one RESTful API
type RESTConfig struct {
	Path          string
	Method        string
	BodyTemplate  interface{}
	Callback      Callback
	VerifyFlag    bool
	verifyHandler VerifyHandler
}

// RESTRequest contains path params, query params and request body of a request
type RESTRequest struct {
	PathParams  map[string]string
	QueryParams map[string][]string
	Body        interface{}

	// in case for some special situation, usually just ignore it
	Ctx echo.Context
}

func packPathParamsAndQueryParamsAndContextIntoRESTRequest(ctx echo.Context) *RESTRequest { //提取来自HTTP请求的数据集的内容
	restRequest := &RESTRequest{}
	restRequest.PathParams = iteratePathParamsAndStoreIntoMap(ctx)
	restRequest.QueryParams = map[string][]string(ctx.QueryParams())
	restRequest.Ctx = ctx
	return restRequest
}

func iteratePathParamsAndStoreIntoMap(ctx echo.Context) map[string]string {
	kv := make(map[string]string)
	ps := ctx.ParamNames() //页面传来数据集的元素名字集合
	for _, paramName := range ps {
		kv[paramName] = ctx.Param(paramName)
	}
	return kv
}

func bind(i interface{}, c echo.Context) (err error) {
	req := c.Request()
	if req.Method == echo.GET {
		//		if err = b.bindData(i, c.QueryParams(), "query"); err != nil {
		//			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		//		}
		return
	}
	ctype := req.Header.Get(echo.HeaderContentType)
	if req.ContentLength == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
	}
	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
		if err = json.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
			} else if se, ok := err.(*json.SyntaxError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
			} else if _, ok := err.(*strconv.NumError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, "Request body format error")
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}
	default:
		return echo.ErrUnsupportedMediaType
	}
	return
}

//需重新改写
func (config *RESTConfig) httpHandlerFn(verifyHandler VerifyHandler) echo.HandlerFunc {
	var handleFn echo.HandlerFunc                          //作为echo的处理函数
	bodyType := reflect.TypeOf(config.BodyTemplate).Elem() //见https://studygolang.com/articles/1251，返回BodyTemplate对应的对象的类型，即为来自页面的数据集——requestBody
	// handler function begins here
	handleFn = func(ctx echo.Context) error {
		//create a new restRequest and fullfil it
		restRequest := packPathParamsAndQueryParamsAndContextIntoRESTRequest(ctx)
		restRequest.Body = reflect.New(bodyType).Interface()
		if ctx.Request().Method == echo.GET {
			if err := ctx.Bind(restRequest.Body); err != nil {
				return err
			}
		} else {
			if config.VerifyFlag { //如果需要验证签名
				//第一步：将body中的Json数据转化为用于验证签名的map格式
				var mapData map[string]string
				if err := bind(&mapData, ctx); err != nil {
					return err
				}
				//第二步：验证签名
				if err := verifyHandler(&mapData); err != nil {
					return err
				}

				//第三步：将map格式的数据转化为具体的结构体对象
				config := &mapstructure.DecoderConfig{
					WeaklyTypedInput: true,
					Base64Decode:     true,
					Result:           restRequest.Body,
				}
				decoder, err := mapstructure.NewDecoder(config)
				if err != nil {
					panic(err)
				}

				err = decoder.Decode(mapData)
				if err != nil {
					panic(err)
				}
			} else { //如果不需要验证签名，则直接将body数据转化为结构体对象
				//使用自定义的bind目的是如果出现strconv.NumError错误，那么当传入的body数据有问题时，可能返回
				//的出错信息内容会很多
				if err := bind(restRequest.Body, ctx); err != nil {
					return err
				}
			}
		}
		statusCode, respBody, err := config.Callback(restRequest) //将来自HTTP请求的context传入回调函数(在相应的接口文件中)进行相应的操作，得到HTTP回应数据集
		if err != nil {
			return ctx.JSON(statusCode, respBody) //&errorMessage{Message: err.Error()}
		}
		if respBody == nil {
			return ctx.NoContent(statusCode)
		}
		return ctx.JSON(statusCode, respBody)
	} // handler function ends here
	return handleFn
}

func (server *WebService) RegisterAPI(config *RESTConfig) error { //对echo的封装，底层为第三方服务：echo
	e := server.server
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		//AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	switch config.Method {
	case http.MethodGet:
		e.GET(config.Path, config.httpHandlerFn(server.verifyHandler))
	case http.MethodPost:
		e.POST(config.Path, config.httpHandlerFn(server.verifyHandler))
	case http.MethodPut:
		e.PUT(config.Path, config.httpHandlerFn(server.verifyHandler))
	case http.MethodPatch:
		e.PATCH(config.Path, config.httpHandlerFn(server.verifyHandler))
	case http.MethodDelete:
		e.DELETE(config.Path, config.httpHandlerFn(server.verifyHandler))
	default:
		return errors.New("the method is not supported")
	}
	return nil
}
