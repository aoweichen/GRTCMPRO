package LoadBalanceHTTP

import (
	"HighPerformanceIMServerProxyServevr/GlobalVars"
	"HighPerformanceIMServerProxyServevr/packages/ConsulServerFind"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const DefaultResponseTime = 1000 // 默认响应时间，单位为毫秒
// ResponseTimeWriter 自定义响应时间记录器
type ResponseTimeWriter struct {
	gin.ResponseWriter
	startTime time.Time
}

// 重写Write方法，在写入响应数据前记录响应时间
func (w *ResponseTimeWriter) Write(data []byte) (int, error) {
	w.WriteHeaderNow()
	return w.ResponseWriter.Write(data)
}

// 重写 WriteHeader 方法，在写入响应头时同时记录响应时间
func (w *ResponseTimeWriter) WriteHeader(code int) {
	duration := time.Since(w.startTime)
	w.Header().Set("X-Response-Time", duration.String())
	w.ResponseWriter.WriteHeader(code)
}

var (
	// 线程锁
	mutex sync.Mutex
)

func GetTargetIMServerURLFromConsulMiddleware(ctx *gin.Context) {
	zap.S().Infoln("进入 GetTargetIMServerURLFromConsulMiddleware")
	targetURL, err := GetTargetIMServerURLFromConsul()
	if err != nil {
		zap.S().Error("从consul获取目标URL失败:", err)
		if errs := ctx.AbortWithError(http.StatusInternalServerError, err); errs != nil {
			zap.S().Error("AbortWithError:", err)
		}
		return
	}
	ctx.Set("targetURL", targetURL)
	ctx.Next()
	zap.S().Infoln("成功出去了")
}

// GetTargetIMServerURLFromRedis 使用最短响应时间算法获得 IMServer 的 URL 地址
// 使用负载均衡策略得到合适 IMServer 的 URL
func GetTargetIMServerURLFromConsul() (*url.URL, error) {
	var err error
	// 定义目标IMServer的URL
	var targetURL *url.URL
	// 定义最短响应时间
	var minResponseTime int64
	// 定义被选择的URL
	var selectedURL string
	// 得到所有的 ServerURLS
	GlobalVars.ServerURLs, err = ConsulServerFind.ConsulServerFindHandlerHTTP(GlobalVars.ConsulAddr, GlobalVars.IMServiceName)
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	// 加锁，保证并发安全
	mutex.Lock()
	defer mutex.Unlock()
	// 遍历服务器 URL 列表，选择响应时间最小的 URL 作为目标 URL
	for _, serverURL := range GlobalVars.ServerURLs {
		responseTime := GlobalVars.ServerStatusMap[serverURL]
		if targetURL == nil || responseTime < minResponseTime {
			minResponseTime = responseTime
			selectedURL = serverURL
		}
	}

	// 如果找到了目标 URL 则解析为 URL 对象
	if selectedURL != "" {
		targetURL, err := url.Parse(selectedURL)
		if err != nil {
			return nil, err
		}
		return targetURL, nil
	} else {
		err := errors.New("服务器{ " + selectedURL + " }已经退出了!")
		return nil, err
	}
}

func InitializeResponseTimes() {
	// 加锁，保证并发安全
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	GlobalVars.ServerURLs, err = ConsulServerFind.ConsulServerFindHandlerHTTP(GlobalVars.ConsulAddr, GlobalVars.IMServiceName)
	if err != nil {
		zap.S().Errorln(err)
	}
	// 初始化每个服务器的响应时间为默认值
	for _, serverURL := range GlobalVars.ServerURLs {
		GlobalVars.ServerStatusMap[serverURL] = DefaultResponseTime
	}
}

// UpdateResponseTime 更新响应时间
func UpdateResponseTime(serverURL string, duration int64) {
	zap.S().Infoln("serverURL: ", serverURL)
	// 加锁，保证并发安全
	mutex.Lock()
	defer mutex.Unlock()
	// 更新服务器的响应时间
	GlobalVars.ServerStatusMap[serverURL] = duration
	zap.S().Infof("ServerStatusMap: %v", GlobalVars.ServerStatusMap)
}
