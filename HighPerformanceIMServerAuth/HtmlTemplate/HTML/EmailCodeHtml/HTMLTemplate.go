package EmailCodeHtmlTemlate

import (
	"bytes"
	"go.uber.org/zap"
	"html/template"
	"os"
)

func EmailCodeHTMLTemplate(htmlTemplateFilePath, code string) string {
	// 读取HTML模板文件内容
	templateString, osReadFileError := os.ReadFile(htmlTemplateFilePath)
	if osReadFileError != nil {
		zap.S().Errorf("无法读取{%#v}文件！", htmlTemplateFilePath)
		// 如果读取文件失败，记录错误日志并抛出异常
		panic(htmlTemplateFilePath)
	}

	// 将读取到的模板内容转换为字符串
	templateHtmlString := string(templateString)

	// 创建模板对象，并将模板内容解析为模板
	templateData := template.Must(template.New("EmailCodeTemplate").Parse(templateHtmlString))

	// 创建缓冲区
	var buf bytes.Buffer

	// 准备模板所需的数据
	data := struct{ Code string }{Code: code}

	// 执行模板，并将结果写入缓冲区
	templateDataExecuteError := templateData.Execute(&buf, data)
	if templateDataExecuteError != nil {
		zap.S().Errorf("执行模板时发生错误:", templateDataExecuteError)
		// 如果执行模板失败，记录错误日志并抛出异常
		panic("执行模板时发生错误: " + templateDataExecuteError.Error())
	}

	// 将缓冲区中的内容转换为字符串并返回
	return buf.String()
}
