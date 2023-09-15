package Response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// JsonResponse TODO 是否移动到DataModels中
// JsonResponse 响应状态码结构体
type JsonResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	HttpCode int         `json:"http_code"`
}

// ToJson 响应json
func (JR *JsonResponse) ToJson(ginCtx *gin.Context) {
	code := 200
	if JR.HttpCode != 200 {
		code = JR.HttpCode
	}
	ginCtx.JSON(code, JR)
}

// FailResponse 失败响应
func FailResponse(code int, message string, data ...interface{}) *JsonResponse {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data
	} else {
		responseData = struct{}{}
	}
	zap.S().Infof("response data: %#v", data)
	returnResponse := &JsonResponse{
		Code:    code,
		Message: message,
		Data:    responseData,
	}
	zap.S().Infof("FailResponse return success,response data: %#v", returnResponse)
	return returnResponse
}

// SuccessResponse 成功响应
func SuccessResponse(data ...interface{}) *JsonResponse {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = struct{}{}
	}
	zap.S().Infof("response data: %#v", data)
	returnResponse := &JsonResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    responseData,
	}
	zap.S().Infof("FailResponse return success,response data: %#v", returnResponse)
	return returnResponse
}

// ErrorResponse 错误响应
func ErrorResponse(status int, message string, data ...interface{}) *JsonResponse {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data
	} else {
		responseData = struct{}{}
	}
	zap.S().Infof("response data: %#v", data)
	returnResponse := &JsonResponse{
		Code:    status,
		Message: message,
		Data:    responseData,
	}
	zap.S().Infof("FailResponse return success,response data: %#v", returnResponse)
	return returnResponse
}

// WriteTo 工具函数 WriteTo 将 json 设为响应体. HTTP 状态码由应用状态码决定
func (JR *JsonResponse) WriteTo(ginCtx *gin.Context) {
	var code int
	if JR.HttpCode == 0 {
		code = http.StatusOK
	} else {
		code = JR.HttpCode
	}
	zap.S().Infof("gin return response data: code => %#v,response => %#v", code, *JR)
	ginCtx.JSON(code, JR)
}

// SetHttpCode 设置JR中httpCode
func (JR *JsonResponse) SetHttpCode(httpCode int) *JsonResponse {
	JR.HttpCode = httpCode
	return JR
}

// responseCode 获取 HTTP 状态码. HTTP 状态码由 应用状态码映射
func (JR *JsonResponse) responseCode() int {
	//	TODO 完善应用状态码对应 http 状态码
	if JR.Code != http.StatusOK {
		return 200
	}
	return 200
}
