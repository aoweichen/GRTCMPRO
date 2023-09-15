package UploadFile

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/Internal/Api/Services/File"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// QINIUYUNHandler 定义一个名为QINIUYUNHandler的结构体类型
type QINIUYUNHandler struct {
}

// QNYResponse 定义一个名为QNYResponse的结构体类型
type QNYResponse struct {
	FileUrl string `json:"file_url"` // FileUrl字段表示文件的URL地址，使用json标签将其映射为file_url
}

// ServiceQiNiuYun 声明一个名为ServiceQiNiuYun的全局变量，表示七牛云服务
var ServiceQiNiuYun File.QiNiuYunService

// UploadFile 定义一个QINIUYUNHandler结构体的UploadFile方法，用于处理文件上传请求
func (QN *QINIUYUNHandler) UploadFile(ctx *gin.Context) {
	// 从请求的表单中获取文件
	file, err := ctx.FormFile("file")
	if err != nil {
		zap.S().Errorln(err)
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ctx)
		return
	}
	// 构建文件保存的路径
	filePath := Configs.ConfigData.Server.FilePath + "/" + file.Filename

	// 将文件保存到指定路径
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		zap.S().Errorln(err)
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ctx)
		return
	}

	// 定义一个名为res的QNYResponse类型变量
	var res QNYResponse

	// 调用ServiceQiNiuYun的UploadFile方法上传文件，并获取文件的URL地址
	fileUrl, _ := ServiceQiNiuYun.UploadFile(filePath, file.Filename)

	// 将文件URL地址赋值给res的FileUrl字段
	res.FileUrl = fileUrl

	// 将成功响应和res的内容作为JSON响应返回给客户端
	Response.SuccessResponse(res).ToJson(ctx)
	return
}
