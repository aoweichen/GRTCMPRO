package Group

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMGroupUsers"
	"HighPerformanceIMServer/DataModels/Models/IMGroups"
	"HighPerformanceIMServer/DataModels/Models/IMSessions"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Utils"
	"errors"
)

type DAO struct {
}

type ImGroups struct {
	Id        int64  `gorm:"column:id" json:"id"`                 //群聊id
	UserId    int64  `gorm:"column:user_id" json:"user_id"`       //创建者
	Name      string `gorm:"column:name" json:"name"`             //群聊名称
	CreatedAt string `gorm:"column:created_at" json:"created_at"` //添加时间
	Info      string `gorm:"column:info" json:"info"`             //群聊描述
	Avatar    string `gorm:"column:avatar" json:"avatar"`         //群聊头像
	Password  string `gorm:"column:password" json:"password"`     //密码
	IsPwd     int8   `gorm:"column:is_pwd" json:"is_pwd"`         //是否加密 0 否 1 是
	Hot       int    `gorm:"column:hot" json:"hot"`               //热度
}

// GetGroupListJoinedByUser
// 定义名为GetGroupListJoinedByUser的方法，该方法属于GroupDAO结构体
// 该方法用于获取用户加入的群组列表
// 参数userID是一个接口类型，表示用户ID
// 返回值[]IMGroups.ImGroups表示群组列表
func (*DAO) GetGroupListJoinedByUser(userID interface{}) []IMGroups.ImGroups {
	// 使用Model.MYSQLDB从数据库中查询IMGroupUsers.ImGroupUsers表
	// 并将查询结果赋值给SQLQuery变量
	SQLQuery := MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
		Where("user_id=?", userID).Select("group_id")

	var IMGroupLists []IMGroups.ImGroups

	// 如果查询结果中的行数大于0，则继续查询群组列表
	if result := MYSQL.DataBase.Model(&IMGroups.ImGroups{}).
		Where("id in (?)", SQLQuery).
		Order("hot desc").
		Find(&IMGroupLists); result.RowsAffected > 0 {

		// 查询结果中的行数大于0，说明有群组列表数据
		// 可以在这里进行一些处理
		return IMGroupLists
	}

	// 返回群组列表
	return IMGroupLists
}

// IsGroupsUser 判断用户是否存在于群组中
func (*DAO) IsGroupsUser(userId interface{}, groupId interface{}) bool {
	var count int64 // 定义一个变量count，用于存储查询结果的数量

	// 使用Model函数指定要操作的数据库模型，并使用Where函数指定查询条件，Count函数用于统计满足条件的记录数量，并将结果存储到count变量中
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
		Where("user_id = ? and group_id = ?", userId, groupId).Count(&count)

	// 如果count大于0，则表示用户存在于群组中，返回true；否则，返回false
	if count > 0 {
		return true
	} else {
		return false
	}
}

// DeleteGroupUser DeleteGroupUser函数用于删除群组成员
// 参数:
// - id: 成员ID
// - groupId: 群组ID
func (*DAO) DeleteGroupUser(id interface{}, groupId string) {
	var groupUsers IMGroupUsers.ImGroupUsers

	// 构建删除条件，并执行删除操作
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
		Where("user_id=? and group_id=?", id, groupId).
		Delete(&groupUsers)
}

// GetGroupUsers GetGroupUsers函数用于获取群组成员列表
// 参数:
// - groupId: 群组ID
// 返回值:
// - []IMGroupUsers.ImGroupUsers: IMGroupUsers.ImGroupUsers对象列表，表示群组成员列表
func (*DAO) GetGroupUsers(groupId string) []IMGroupUsers.ImGroupUsers {
	var groupUsers []IMGroupUsers.ImGroupUsers

	// 根据群组ID查询群组成员列表，并预加载关联的Users对象
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
		Where("group_id=?", groupId).Preload("Users").First(&groupUsers)

	return groupUsers
}

// CreateGroup CreateGroup函数用于创建群组
// 参数:
// - params: ApiRequests.CreateGroupRequest对象，包含了创建群组的参数
// 返回值:
// - error: 错误对象，如果创建过程中出现错误，则返回错误对象，否则为nil
// - IMGroups.ImGroups: IMGroups.ImGroups对象，表示创建成功的群组
func (*DAO) CreateGroup(params ApiRequests.CreateGroupRequest) (error, IMGroups.ImGroups) {
	// 创建一个IMGroups.ImGroups对象，设置其属性值
	imGroups := IMGroups.ImGroups{
		UserId:    params.UserId,  // 用户ID
		Name:      params.Name,    // 群名称
		CreatedAt: Date.NewDate(), // 创建时间
		Info:      params.Info,    // 群介绍
		Avatar:    params.Avatar,  // 群头像
	}

	// 在数据库中创建群组，并检查是否出现错误
	if MYSQL.DataBase.Model(&IMGroups.ImGroups{}).Create(&imGroups).Error != nil {
		return errors.New("创建错误"), imGroups
	} else {
		return nil, imGroups
	}
}

// CreateOneGroupUser CreateOneGroupUser函数用于创建群组成员和会话
// 参数:
// - group: IMGroups.ImGroups对象，表示群组信息
// - id: 成员ID
func (*DAO) CreateOneGroupUser(group IMGroups.ImGroups, id int) {
	// 创建一个IMGroupUsers.ImGroupUsers对象，表示群组成员
	groupUser := IMGroupUsers.ImGroupUsers{
		UserId:    id,             // 成员ID
		CreatedAt: Date.NewDate(), // 创建时间
		GroupId:   int(group.Id),  // 群组ID
		Avatar:    group.Avatar,   // 群头像
		Name:      group.Name,     // 群名称
	}

	// 创建一个IMSessions.ImSessions对象，表示会话
	session := IMSessions.ImSessions{
		FromId:      int64(id),            // 发送者ID
		ToId:        group.Id,             // 接收者ID（群组ID）
		CreatedAt:   Date.NewDate(),       // 创建时间
		ChannelType: IMSessions.GROUPTYPE, // 频道类型（群组类型）
		Name:        group.Name,           // 会话名称（群名称）
		Avatar:      group.Avatar,         // 会话头像（群头像）
	}

	// 在数据库中创建群组成员和会话
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).Create(&groupUser)
	MYSQL.DataBase.Model(&IMSessions.ImSessions{}).Create(&session)

	return
}

// DelSelectGroupUser DelSelectGroupUser函数用于删除选定的群组成员及相关会话
// 参数:
// - userIds: []string类型，表示选定的成员ID列表
// - groupId: int类型，表示群组ID
// - avatar: string类型，表示群组头像
// - name: string类型，表示群组名称
func (*DAO) DelSelectGroupUser(userIds []string, groupId int, avatar string, name string) {
	// 删除选定的群组成员
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).Where("user_id in(?)", userIds).Delete(&IMGroupUsers.ImGroupUsers{})

	// 删除相关会话
	MYSQL.DataBase.Model(&IMSessions.ImSessions{}).Where("group_id=? and from_id in(?)", groupId, userIds).Delete(&IMSessions.ImSessions{})

	return
}

// CreateSelectGroupUser CreateSelectGroupUser函数用于创建选定的群组成员及相关会话
// 参数:
// - userIds: []string类型，表示选定的成员ID列表
// - groupId: int类型，表示群组ID
// - avatar: string类型，表示群组头像
// - name: string类型，表示群组名称
func (*DAO) CreateSelectGroupUser(userIds []string, groupId int, avatar string, name string) {
	count := len(userIds)                                    // 计算成员ID列表的长度
	createdAt := Date.NewDate()                              // 获取当前时间
	var groupUser = make([]IMGroupUsers.ImGroupUsers, count) // 创建一个IMGroupUsers.ImGroupUsers类型的切片
	var sessionsData = make([]IMSessions.ImSessions, count)  // 创建一个IMSessions.ImSessions类型的切片

	// 遍历成员ID列表，设置群组成员和会话的属性值
	for key, id := range userIds {
		groupUser[key].UserId = Utils.StringToInt(id) // 设置群组成员的用户ID
		groupUser[key].CreatedAt = createdAt          // 设置群组成员的创建时间
		groupUser[key].Avatar = avatar                // 设置群组成员的头像
		groupUser[key].GroupId = groupId              // 设置群组成员的群组ID
		groupUser[key].Name = name                    // 设置群组成员的名称

		sessionsData[key].FromId = Utils.StringToInt64(id)   // 设置会话的发送者ID
		sessionsData[key].GroupId = int64(groupId)           // 设置会话的群组ID
		sessionsData[key].CreatedAt = createdAt              // 设置会话的创建时间
		sessionsData[key].ChannelType = IMSessions.GROUPTYPE // 设置会话的频道类型（群组类型）
		sessionsData[key].Name = name                        // 设置会话的名称
		sessionsData[key].Avatar = avatar                    // 设置会话的头像
		sessionsData[key].TopTime = Date.NewDate()           // 设置会话的置顶时间
	}

	// 在数据库中创建群组成员和会话
	MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).Create(&groupUser)
	MYSQL.DataBase.Model(&IMSessions.ImSessions{}).Create(&sessionsData)

	return
}
