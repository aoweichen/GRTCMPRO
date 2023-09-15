package Friend

import (
	"HighPerformanceIMServer/DataModels/Models/IMFriends"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
)

type APIUsers struct {
	ID     int64  `gorm:"column:id;primaryKey" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Email  string `gorm:"column:email" json:"email"`
	Avatar string `gorm:"column:avatar" json:"avatar"`
	Status int8   `gorm:"column:status" json:"status"`
	Bio    string `gorm:"column:bio" json:"bio"`
	Sex    int8   `gorm:"column:sex" json:"sex"`
}
type DAO struct {
}

// GetFriendLists
// 定义一个名为GetFriendLists的方法，该方法属于FriendDao结构体
// 该方法用于获取好友列表
// 参数id是一个接口类型，表示用户的ID
// 返回值error表示方法执行过程中可能发生的错误，interface{}表示好友列表数据
func (D *DAO) GetFriendLists(id interface{}) (error, interface{}) {
	var err error
	var friendList []IMFriends.ImFriends

	// 使用Model.DB从数据库中查询好友列表
	// 并将结果赋值给list变量
	result := MYSQL.DataBase.Model(&IMFriends.ImFriends{}).Preload("Users").
		Where("from_id=?", id).
		Order("status desc").
		Order("top_time desc").
		Find(&friendList)

	// 如果查询结果中的行数为0，则表示好友列表为空
	// 返回错误和空的列表数据
	if result.RowsAffected == 0 {
		return err, friendList
	}

	// 返回nil表示没有错误，将list作为好友列表数据返回
	return nil, friendList
}

// AgreeFriendRequest 定义一个名为AgreeFriendRequest的方法，接收一个int64类型的toId和fromId参数
func (D *DAO) AgreeFriendRequest(toID int64, fromID int64) {
	// 创建一个im_friends.ImFriends结构体实例friend，设置其字段值
	friend := IMFriends.ImFriends{
		FromId:    fromID,
		ToId:      toID,
		Note:      "",
		CreatedAt: Date.NewDate(),
		UpdatedAt: Date.NewDate(),
		TopTime:   Date.NewDate(),
		Status:    0,
	}
	MYSQL.DataBase.Save(&friend)
}

// GetNotFriendList 获取非好友数据
// 定义一个名为GetNotFriendList的方法，接收一个id接口类型和一个email字符串参数，返回一个APIUsers类型的切片
func (D *DAO) GetNotFriendList(id interface{}, email string) []APIUsers {
	// 创建一个APIUsers类型的切片users
	var users []APIUsers

	// 构建SQL查询语句，查询im_friends表中from_id等于id的记录的to_id字段
	sqlQuery := MYSQL.DataBase.Table("im_friends").Where("from_id=?", id).Select("to_id")

	// 构建查询语句，查询im_users表中不在sqlQuery结果中的记录，并且user_type为0，同时id不等于给定的id参数
	query := MYSQL.DataBase.Table("im_users").
		Where("id not in(?) and user_type=?", sqlQuery, 0).Where("id!=?", id)

	// 如果email参数不为空，则将查询条件添加email字段的模糊匹配
	if len(email) > 0 {
		query = MYSQL.DataBase.Table("im_users").
			Where("email like ?", email)
	}

	// 执行查询语句，选择id、name、email、avatar、bio、sex、status字段，并限制结果数量为5，将查询结果保存到users切片中
	query.Select("id,name,email,avatar,bio,sex,status").Limit(5).Find(&users)

	// 返回查询结果users
	return users
}

// DelFriends 删除好友
func (D *DAO) DelFriends(toId interface{}, fromId interface{}) error {
	// 定义一个错误变量err
	var err error

	// 定义一个im_friends.ImFriends类型的变量friend
	var friend IMFriends.ImFriends

	// 执行删除操作，删除im_friends表中(to_id=toId and from_id=fromId)或者(to_id=fromId and from_id=toId)的记录，并将结果保存到result中
	result := MYSQL.DataBase.Model(&IMFriends.ImFriends{}).Preload("Users").
		Where("(to_id=? and from_id=?) or (to_id=? and from_id=?)", toId, fromId, toId, fromId).Delete(&friend)

	// 如果删除的记录数为0，则返回错误
	if result.RowsAffected == 0 {
		return err
	}

	// 删除成功，返回nil表示没有错误
	return nil
}

// GetFriends 查询好友详情
func (D *DAO) GetFriends(id interface{}) (error, interface{}) {
	// 定义一个错误变量err
	var err error

	// 定义一个im_friends.ImFriends类型的变量list
	var list IMFriends.ImFriends

	// 执行查询操作，从im_friends表中查询from_id等于id的记录，并预加载关联的Users模型，按照status降序和top_time降序排序，将查询结果保存到list中
	result := MYSQL.DataBase.Model(&IMFriends.ImFriends{}).Preload("Users").
		Where("from_id=?", id).
		Order("status desc").
		Order("top_time desc").
		Find(&list)

	// 如果查询结果为空，则返回错误和list
	if result.RowsAffected == 0 {
		return err, list
	}

	// 返回nil表示没有错误，以及查询结果list
	return nil, list
}

// IsFriends 判断是否是好友关系
func (D *DAO) IsFriends(id interface{}, toId interface{}) bool {
	// 定义一个int64类型的变量count
	var count int64

	// 执行查询操作，统计im_friends表中to_id等于id且from_id等于toId的记录数，将结果保存到count中
	MYSQL.DataBase.Table("im_friends").
		Where("to_id=? and from_id=?", id, toId).
		Count(&count)

	// 如果count等于0，则表示不是好友关系，返回false；否则，表示是好友关系，返回true
	if count == 0 {
		return false
	}
	return true
}
