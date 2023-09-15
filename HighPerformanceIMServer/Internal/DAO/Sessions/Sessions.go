package Sessions

import (
	"HighPerformanceIMServer/DataModels/Models/IMSessions"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
)

type DAO struct {
}

// CreateSession 函数用于创建会话，并返回会话对象
func (SD *DAO) CreateSession(fromID int64, toID int64, channelType int) (sessions *IMSessions.ImSessions) {
	// 创建一个ImUsers对象
	var users IMUser.ImUsers

	// 从数据库中查询toID对应的用户信息，并将结果赋值给users对象
	MYSQL.DataBase.Table("im_users").Where("id=?", toID).First(&users)

	// 创建一个ImSessions对象，设置其属性值
	session := &IMSessions.ImSessions{
		FromId:      fromID,
		ToId:        toID,
		CreatedAt:   Date.NewDate(),
		TopStatus:   IMSessions.TopStatus,
		TopTime:     Date.NewDate(),
		Note:        users.Name,
		ChannelType: channelType,
		Name:        users.Name,
		Avatar:      users.Avatar,
		Status:      IMSessions.SessionStatusOk,
	}

	// 将session对象保存到数据库中
	MYSQL.DataBase.Save(session)

	// 返回创建的会话对象
	return session
}
