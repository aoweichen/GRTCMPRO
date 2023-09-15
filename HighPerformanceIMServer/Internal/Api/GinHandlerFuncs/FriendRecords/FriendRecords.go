package FriendRecords

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMFriends"
	"HighPerformanceIMServer/DataModels/Models/IMFriendsRecords"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/Api/Services/Clients/Message"
	"HighPerformanceIMServer/Internal/Api/Services/IMMessage"
	"HighPerformanceIMServer/Internal/DAO/Friend"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Internal/DAO/Sessions"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Response"
	"HighPerformanceIMServer/Packages/Utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
}

// GetFriendsRequestRecord 获取好友申请记录
func (FRD *Handler) GetFriendsRequestRecord(ctx *gin.Context) {
	var friendsRecordsList []IMFriendsRecords.ImFriendRecords
	id := ctx.MustGet("id")
	// 使用 model.DB 对象调用 Model 方法指定模型 IMFriendsRecords.ImFriendRecords{}
	// 然后调用 Preload 方法预加载关联表 "Users"，
	if result := MYSQL.DataBase.Model(&IMFriendsRecords.ImFriendRecords{}).Preload("Users").
		// 再使用 Where 方法指定查询条件 "to_id=? or from_id=?" 并传入 id 变量两次，
		Where("to_id=? or from_id=?", id, id).
		// 最后使用 Order 方法指定排序字段 "created_at desc"，
		// 并调用 Find 方法执行查询并将结果存储到 list 变量中。
		// 同时将结果赋值给 result 变量，并通过 result.RowsAffected 获取查询结果的行数。
		Order("created_at desc").Find(&friendsRecordsList); result.RowsAffected == 0 {
		// 如果行数为 0，则执行以下代码块
		// 调用 response.SuccessResponse 方法创建一个成功响应，
		// 并调用 ToJson 方法将响应转换为 JSON 格式并发送给客户端
		Response.SuccessResponse().ToJson(ctx)
		return
		// 结束当前函数的执行
	}
	Response.SuccessResponse(friendsRecordsList).ToJson(ctx)
	return
}

// SendAddFriendsRequest 发起好友申请
// SendAddFriendsRequest函数是FriendRecordHandler结构体的方法，用于发送添加好友请求
func (FRD *Handler) SendAddFriendsRequest(ctx *gin.Context) {
	// 从上下文中获取用户ID
	id := ctx.MustGet("id")
	// 根据请求参数创建一个CreateFriendRequest对象
	params := ApiRequests.CreateFriendRequest{
		ToId:        ctx.PostForm("to_id"),       // 添加好友的目标用户ID
		Information: ctx.PostForm("information"), // 添加好友的附加信息
	}

	// 使用validator验证params参数的结构体字段是否合法
	if err := validator.New().Struct(params); err != nil {
		zap.S().Errorln(err)
		// 如果验证失败，返回参数错误的错误响应并序列化为JSON格式返回给客户端
		Response.ErrorResponse(Enums.ParamError, err.Error()).ToJson(ctx)
		return
	}

	// 查询目标用户是否存在
	var users IMUser.ImUsers
	if result := MYSQL.DataBase.Table("im_users").
		Where("id=?", params.ToId).
		First(&users); result.RowsAffected == 0 {
		zap.S().Errorln("用户不存在")
		// 如果用户不存在，返回用户不存在的错误响应并序列化为JSON格式返回给客户端
		Response.ErrorResponse(Enums.ParamError, "用户不存在").ToJson(ctx)
		return
	}

	// 判断是否已发送过好友请求
	var imFriendsRecord IMFriendsRecords.ImFriendRecords
	if result := MYSQL.DataBase.Table("im_users").
		Where("from_id=? and to_id=? and status=0", id, params.ToId).
		First(&imFriendsRecord); result.RowsAffected > 0 {
		zap.S().Errorln("请勿重复添加...")
		// 如果已发送过好友请求，返回重复添加的错误响应并序列化为JSON格式返回给客户端
		Response.ErrorResponse(Enums.ParamError, "请勿重复添加...").ToJson(ctx)
		return
	}

	// 判断是否已是好友关系
	var imFriends IMFriends.ImFriends
	if result := MYSQL.DataBase.Table("im_friends").
		Where("from_id=? and to_id=?", id, params.ToId).
		First(&imFriends); result.RowsAffected > 0 {
		// 如果已是好友关系，返回已是好友关系的错误响应并序列化为JSON格式返回给客户端
		Response.ErrorResponse(Enums.ParamError, "用户已经是好友关系了...").ToJson(ctx)
		return
	}

	// 创建好友记录并保存到数据库中
	imFriendRecords := IMFriendsRecords.ImFriendRecords{
		FromId:      Utils.InterfaceToInt64(id),
		ToId:        Utils.StringToInt64(params.ToId),
		Status:      IMFriendsRecords.WaitingStatus,
		CreatedAt:   Date.NewDate(),
		Information: params.Information,
	}
	MYSQL.DataBase.Save(&imFriendRecords)

	// 创建消息服务实例
	var messageService IMMessage.MessageService
	// 创建好友消息对象
	message := Message.CreateFriendMessage{
		MsgCode:     Enums.WsCreate,
		ID:          imFriendRecords.Id,
		FromId:      imFriendRecords.FromId,
		Status:      imFriendRecords.Status,
		CreatedAt:   imFriendRecords.CreatedAt,
		ToID:        imFriendRecords.ToId,
		Information: imFriendRecords.Information,
		Users: Message.Users{
			ID:     users.ID,
			Avatar: users.Avatar,
			Name:   users.Name,
		},
	}
	// 调用消息服务的SendFriendActionMessage方法发送好友动作消息
	messageService.SendFriendActionMessage(message)

	// 更新好友记录中的用户信息字段
	imFriendRecords.Users.Name = users.Name
	imFriendRecords.Users.Id = users.ID
	imFriendRecords.Users.Avatar = users.Avatar

	// 返回成功响应，并将好友记录序列化为JSON格式返回给客户端
	Response.SuccessResponse(imFriendRecords).ToJson(ctx)
}

// AgreeOrRejectFriendRequest 处理接受或拒绝好友请求的函数
func (FRD *Handler) AgreeOrRejectFriendRequest(ctx *gin.Context) {
	// 定义变量
	var users IMUser.ImUsers                     // 用户信息
	var messageService IMMessage.MessageService  // 消息服务
	var friends IMFriends.ImFriends              // 好友关系
	var messageCode int                          // 消息代码
	var records IMFriendsRecords.ImFriendRecords // 好友记录

	// 从上下文中获取用户ID
	id := ctx.MustGet("id")
	// 获取请求参数
	params := ApiRequests.UpdateFriendRequest{
		Status: Utils.StringToInt(ctx.PostForm("status")), // 好友请求的状态，接受或拒绝
		ID:     ctx.PostForm("id"),                        // 好友请求的ID
	}

	// 使用验证器验证参数的合法性
	err := validator.New().Struct(params)
	if err != nil {
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ctx)
		zap.S().Errorln(err)
		return
	}

	// 从数据库中查询好友记录
	if result := MYSQL.DataBase.Table("im_friend_records").
		Where("id=? and status=0", params.ID).
		First(&records); result.RowsAffected == 0 {
		zap.S().Errorln("数据不存在")
		Response.ErrorResponse(http.StatusInternalServerError, "数据不存在").ToJson(ctx)
	}

	// 检查好友关系是否已存在
	if result := MYSQL.DataBase.Table("im_friends").Where("from_id=? and to_id=?", records.FromId, records.ToId); result.RowsAffected > 0 {
		zap.S().Errorln("用户已经是好友关系了...")
		Response.ErrorResponse(Enums.ParamError, "用户已经是好友关系了...").ToJson(ctx)
	}

	// 更新好友记录的状态
	MYSQL.DataBase.Table("im_users").Where("id=?", id).Find(&users)
	records.Status = params.Status
	MYSQL.DataBase.Updates(&records)

	// 根据请求状态进行相应的操作
	if params.Status == 1 {
		var friendDao Friend.DAO
		var sessionDao Sessions.DAO

		messageCode = Enums.WsFriendOk
		friendDao.AgreeFriendRequest(records.FromId, records.ToId)
		friendDao.AgreeFriendRequest(records.ToId, records.FromId)

		sessionDao.CreateSession(records.FromId, records.ToId, 1)
		sessionDao.CreateSession(records.ToId, records.FromId, 1)

	} else {
		messageCode = Enums.WsFriendError
	}

	// 构建消息对象
	message := &Message.CreateFriendMessage{
		MsgCode:     messageCode,
		ID:          records.Id,
		FromId:      records.ToId,
		Status:      records.Status,
		CreatedAt:   records.CreatedAt,
		ToID:        records.FromId,
		Information: records.Information,
		Users: Message.Users{
			Name:   users.Name,
			ID:     users.ID,
			Avatar: users.Avatar,
		},
	}

	// 发送好友操作消息
	messageService.SendFriendActionMessage(*message)
	friends.Status = params.Status
	Response.SuccessResponse(friends).WriteTo(ctx)
	return
}

// UserQuery 处理用户查询的函数
func (FRD *Handler) UserQuery(ctx *gin.Context) {
	// 定义变量
	var friendDao Friend.DAO                                          // 声明一个 FriendDAO 类型的变量 friendDao
	id := ctx.MustGet("id")                                           // 从上下文中获取用户ID，并赋值给变量 id
	params := ApiRequests.QueryUserRequest{Email: ctx.Query("email")} // 创建一个 QueryUserRequest 结构体实例，并将 email 参数赋值给结构体的 Email 字段

	// 使用验证器验证参数的合法性
	if err := validator.New().Struct(params); err != nil { // 使用验证器验证 params 结构体的合法性，如果验证不通过，将错误信息赋值给变量 err
		zap.S().Errorln(err.Error())                                      // 使用 zap.S().Errorln() 打印日志错误信息
		Response.ErrorResponse(Enums.ParamError, err.Error()).ToJson(ctx) // 返回参数错误的响应
		return
	}

	// 获取非好友用户列表
	users := friendDao.GetNotFriendList(id, params.Email) // 调用 friendDao 的 GetNotFriendList 方法，传入 id 和 params.Email 参数，获取非好友用户列表

	// 返回成功的响应
	Response.SuccessResponse(users).ToJson(ctx) // 返回成功的响应，并将用户列表写入上下文
	return
}
