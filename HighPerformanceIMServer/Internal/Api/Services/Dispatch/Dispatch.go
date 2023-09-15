package Dispatch

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/Internal/DAO/Redis"
	"context"
	"go.uber.org/zap"
	"time"
)

type Service struct {
}

type ServiceInterface interface {
	// SetDispatchNode 设置当前节点信息
	SetDispatchNode(uid string, node string)
	// GetDispatchNode 获取当前节点信息
	GetDispatchNode(uid string, node string)

	MessageDispatch(uid string, node string)

	IsDispatchNode(uid string, node string)
	// DeleteDispatchNode DeleteDispatch 删除当前节点
	DeleteDispatchNode(uid string)
}

// IsDispatchNode IsDispatch 判断当前节点是否分派服务节点
func (DS *Service) IsDispatchNode(uid string) (bool, string) {
	// 调用redis.RedisDB的Exists方法，检查uid是否存在于Redis数据库中，结果赋值给numDispatchNode变量
	numDispatchNode, redisREDISDBExistsResultError := Redis.DataBase.Exists(context.Background(), uid).Result()
	if redisREDISDBExistsResultError != nil {
		zap.S().Errorln("系统错误:", redisREDISDBExistsResultError)
		panic(redisREDISDBExistsResultError)
	}
	// 如果n大于0，表示uid存在于Redis数据库中
	if numDispatchNode > 0 {
		// 调用Service的GetDispatchNode方法，传入uid作为参数，获取uid对应的分派节点，结果赋值给uNode变量
		uNode := DS.GetDispatchNode(uid)

		// 如果当前服务器节点与uNode相同，表示当前服务器节点为分派节点，返回false和空字符串
		if Configs.ConfigData.Server.Node == uNode {
			return false, ""
		}
		// 否则，表示当前服务器节点不是分派节点，返回true和uNode
		return true, uNode
	} else {
		// 如果n等于0，表示uid不存在于Redis数据库中，返回false和空字符串
		return false, ""
	}
}

// GetDispatchNode 是 Service 结构体的方法，用于获取调度节点
// 参数 uid 是用户 ID，表示要获取调度节点的用户
// 返回值是一个字符串，表示调度节点的信息
func (DS *Service) GetDispatchNode(uid string) string {
	// 使用 Redis.REDISDB.Get 方法从 Redis 数据库中获取指定用户 ID 的值，并调用 Val 方法获取其字符串表示
	return Redis.DataBase.Get(context.Background(), uid).Val()
}

// SetDispatchNode 设置分派节点
func (DS *Service) SetDispatchNode(uid string) {
	// 调用Redis.REDISDB的Set方法，将uid作为key，ConfigModels.ConfigData.Server.Node作为value，设置过期时间为24小时
	Redis.DataBase.Set(context.Background(), uid, Configs.ConfigData.Server.Node, time.Hour*24)
}

// DeleteDispatchNode 是 Service 结构体的方法，用于删除调度节点
// 参数 uid 是用户 ID，表示要删除的调度节点的用户
func (DS *Service) DeleteDispatchNode(uid string) {
	// 使用 Redis.REDISDB.Del 方法从 Redis 数据库中删除指定用户 ID 的键值对
	Redis.DataBase.Del(context.Background(), uid)
}

// MessageDispatch TODO 消息分发
func (DS *Service) MessageDispatch() {

}
