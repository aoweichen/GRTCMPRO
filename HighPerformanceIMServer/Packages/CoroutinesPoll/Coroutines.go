package CoroutinesPoll

import (
	"HighPerformanceIMServer/Configs"
	"github.com/panjf2000/ants/v2"
)

var AntPool *ants.Pool

// ConnectPool 函数用于连接池的创建和初始化
func ConnectPool() *ants.Pool {
	// 使用 ConfigModels 包中的 ConfigData 获取服务器的协程池配置信息，创建一个新的 ants.Pool
	AntPool, _ = ants.NewPool(Configs.ConfigData.Server.CoroutinePoll)

	// 返回创建的连接池
	return AntPool
}
