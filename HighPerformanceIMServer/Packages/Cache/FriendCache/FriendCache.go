package FriendCache

import (
	"HighPerformanceIMServer/DataModels/Models/IMFriends"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"sync"
)

// CacheInterface 定义了CacheInterface接口
// 该接口用于定义缓存相关的方法
type CacheInterface interface {
	// SetFriendCache 设置好友缓存
	SetFriendCache(uid string, friends *[]IMFriends.ImFriends)
	// GetFriendCache 获取好友缓存
	GetFriendCache(uid string) (*[]IMFriends.ImFriends, error)
}

// FriendCacheClients是FriendCacheClient的实例
var (
	CacheClients = CacheClient{
		CacheMap: make(map[string]*[]IMFriends.ImFriends),
	}
	mux sync.Mutex
)

// CacheClient 定义了FriendCacheClient结构体
// 该结构体用于实现CacheInterface接口的方法，并包含一个CacheMap字段
type CacheClient struct {
	CacheMap map[string]*[]IMFriends.ImFriends
}

// SetFriendCache 实现CacheInterface接口的SetFriendCache方法
func (FCC *CacheClient) SetFriendCache(uid string, friends *[]IMFriends.ImFriends) {
	// 加锁，保证并发安全
	mux.Lock()
	// 将好友缓存存入CacheMap
	FCC.CacheMap[uid] = friends
	// 解锁
	mux.Unlock()
}

// GetFriendCache 实现CacheInterface接口的GetFriendCache方法
func (FCC *CacheClient) GetFriendCache(uid string) (*[]IMFriends.ImFriends, error) {
	// 判断指定UID的好友缓存是否存在
	if imFriends, ok := FCC.CacheMap[uid]; ok {
		// 如果存在，直接返回缓存的好友列表
		return imFriends, nil
	}
	// 如果缓存不存在，从数据库中获取好友列表
	var imFriendsList []IMFriends.ImFriends
	MYSQL.DataBase.Table("im_friends").Where("m_id=?", uid).Find(&imFriendsList)
	// 将获取到的好友列表存入缓存
	FCC.SetFriendCache(uid, &imFriendsList)
	// 返回好友列表及错误为nil
	return &imFriendsList, nil
}
