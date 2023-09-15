package GlobalVars

import "HighPerformanceIMServerProxyServevr/configs"

var (
	ConsulAddr      string
	AuthServiceName string
	IMServiceName   string
)
var (
	// ServerURLs 全局
	ServerURLs []string
	// ServerStatusMap 全局
	ServerStatusMap map[string]int64
)

func InitGlobalVars() {
	ConsulAddr = configs.ConfigData.Consul.Host + configs.ConfigData.Consul.Port
	AuthServiceName = configs.ConfigData.AuthService.Name
	IMServiceName = configs.ConfigData.IMServers.Name
	ServerStatusMap = make(map[string]int64)
}
