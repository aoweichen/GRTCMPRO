package SetAvatar

import (
	"encoding/base64"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func GetAvatarBase64Png(url string) string {
	response, err := http.Get(url)
	if err != nil {
		zap.S().Errorln("请求失败: ", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.S().Panic(err)
			panic(err)
		}
	}(response.Body)
	if response.StatusCode == http.StatusOK {
		imageData, err := io.ReadAll(response.Body)
		if err != nil {
			zap.S().Panic(err)
			panic(err)
		}
		base64PngData := base64.StdEncoding.EncodeToString(imageData)
		dataURI := "data:image/svg+xml;base64," + base64PngData
		return dataURI
	} else {
		zap.S().Panic("请求失败，状态码:", response.StatusCode)
		panic("请求失败，状态码:" + strconv.Itoa(response.StatusCode))
	}
}
