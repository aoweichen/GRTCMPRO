package Auth

import (
	"HighPerformanceIMServerProxyServevr/DataModels/ApiRequests"
	"HighPerformanceIMServerProxyServevr/HandleFuncs/AUTHGRPC"
	"HighPerformanceIMServerProxyServevr/packages/Enums"
	"HighPerformanceIMServerProxyServevr/packages/Response"
	"HighPerformanceIMServerProxyServevr/packages/Utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SendeEmail(ctx *gin.Context) {
	// 创建一个ApiRequests.SendEmailValidateCodeRequest类型的params变量，并初始化其字段值
	params := &ApiRequests.SendEmailValidateCodeRequest{
		Email:     ctx.PostForm("email"),                         // 从请求中获取email字段的值
		EmailType: Utils.StringToInt(ctx.PostForm("email_type")), // 从请求中获取email_type字段的值，并将其转换为整型
	}
	// 验证 params 结构
	// 使用validator库对params进行结构验证
	validateParamsStructError := validator.New().Struct(params) // 使用validator库对params进行结构验证
	if validateParamsStructError != nil {
		// 打印验证失败的错误信息
		zap.S().Errorln(validateParamsStructError.Error())
		// 返回参数错误的响应
		Response.FailResponse(Enums.ParamError, validateParamsStructError.Error()).WriteTo(ctx)
		return
	}
	// 发送 grpc 请求
	response, err := AUTHGRPC.SendeEmailHandler(params.Email, params.EmailType)
	// if err != nil {
	// 	zap.S().Errorf(err.Error())
	// 	Response.FailResponse(Enums.ApiError, err.Error()).ToJson(ctx)
	// 	return
	// }
	if !response.IsSendEmailCodeSuccess {
		zap.S().Errorf(err.Error())
		Response.FailResponse(Enums.ApiError, err.Error()).ToJson(ctx)
		return
	} else {
		zap.S().Infoln("邮件发送成功，请注意查收！")
		// 返回成功响应
		Response.SuccessResponse().ToJson(ctx)
		return
	}
}

func Login(ctx *gin.Context) {
	// 前端提交的表单
	params := &ApiRequests.LoginForm{
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	zap.S().Infof("login form params: %#v", params)
	// 验证表单结构
	validateError := validator.New().Struct(params)
	if validateError != nil {
		zap.S().Errorf("validate login form params error: %#v", validateError)
		Response.FailResponse(http.StatusInternalServerError, validateError.Error()).WriteTo(ctx)
		return
	}
	zap.S().Infof("login form params struct validate success!")
	// grpc
	response, err := AUTHGRPC.LoginHandler(params.Email, params.Password)
	if err != nil {
		zap.S().Errorf(err.Error())
		Response.FailResponse(Enums.ApiError, err.Error()).ToJson(ctx)
		return
	} else {
		Response.SuccessResponse(&AUTHGRPC.LoginResponse{
			ID:              response.ID,
			UID:             response.UID,
			Name:            response.Name,
			Avatar:          response.Avatar,
			Email:           response.Email,
			Token:           response.Token,
			ExpireTime:      response.ExpireTime,
			TokenTimeToLive: response.TokenTimeToLive,
		}).WriteTo(ctx)
		return
	}

}

func Register(ctx *gin.Context) {
	// 创建一个ApiRequests.RegisteredForm类型的params变量，并初始化其字段值
	params := &ApiRequests.RegisteredForm{
		Email:          ctx.PostForm("email"),                         // 从请求中获取email字段的值
		Name:           ctx.PostForm("name"),                          // 从请求中获取name字段的值
		EmailType:      Utils.StringToInt(ctx.PostForm("email_type")), // 从请求中获取email_type字段的值，并将其转换为整型，默认为1
		Password:       ctx.PostForm("password"),                      // 从请求中获取password字段的值
		PasswordRepeat: ctx.PostForm("password_repeat"),               // 从请求中获取password_repeat字段的值
		Code:           ctx.PostForm("code"),                          // 从请求中获取code字段的值
	}
	// 验证 params 结构
	validateParamsStructError := validator.New().Struct(params) // 使用validator库对params进行结构验证
	if validateParamsStructError != nil {
		zap.S().Errorln(validateParamsStructError.Error())                                      // 打印验证失败的错误信息
		Response.FailResponse(Enums.ParamError, validateParamsStructError.Error()).WriteTo(ctx) // 返回参数错误的响应
		return
	}
	// 进行grpc请求
	_, err := AUTHGRPC.RegisterHandler(params.Email, params.Name, int64(params.EmailType),
		params.Password, params.PasswordRepeat, params.Code)

	if err != nil {
		zap.S().Errorln(err)
		Response.FailResponse(Enums.ParamError, err.Error()).WriteTo(ctx)
		return
	} else {
		zap.S().Infoln("注册成功！")
		Response.SuccessResponse().ToJson(ctx) // 返回成功响应
	}
}
