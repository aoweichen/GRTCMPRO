package Group

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMGroupUsers"
	"HighPerformanceIMServer/DataModels/Models/IMGroups"
	"HighPerformanceIMServer/DataModels/Models/IMMessages"
	"HighPerformanceIMServer/DataModels/Models/IMSessions"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Base"
	"HighPerformanceIMServer/Internal/Api/Services/IMMessage"
	"HighPerformanceIMServer/Internal/DAO/Group"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Hash"
	"HighPerformanceIMServer/Packages/Response"
	"HighPerformanceIMServer/Packages/Utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Handler GroupHandler结构体定义了一个群组处理器
type Handler struct {
}

var (
	groupDao       Group.DAO                // 创建GroupDAO对象
	messageService IMMessage.MessageService // 创建ImMessageService对象
)

// GetGroupList GetGroupList方法用于获取群组列表
// 参数：
// - ctx: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) GetGroupList(ctx *gin.Context) {
	groupList := groupDao.GetGroupListJoinedByUser(ctx.MustGet("id")) // 调用GroupDAO的GetGroupListJoinedByUser方法，获取用户加入的群组列表
	Response.SuccessResponse(groupList).ToJson(ctx)                   // 将群组列表返回给客户端
	return
}

// CreateGroup Store Store方法用于创建群组
// 参数：
// - ctx: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) CreateGroup(ctx *gin.Context) {
	id := ctx.MustGet("id") // 获取当前用户的ID
	var selectUser SelectUser

	err := ctx.ShouldBind(&selectUser) // 解析请求体中的JSON数据到selectUser结构体中
	if err != nil {
		return
	}
	selectUser.SelectUser = append(selectUser.SelectUser, Utils.InterfaceToInt64String(id)) // 将当前用户的ID添加到selectUser.SelectUser切片中
	fmt.Println(selectUser.SelectUser)

	// 构建创建群组的请求参数
	params := ApiRequests.CreateGroupRequest{
		UserId:     Utils.InterfaceToInt64(id),                // 当前用户的ID
		Name:       ctx.PostForm("name"),                      // 群组名称
		Info:       ctx.PostForm("info"),                      // 群组信息
		Avatar:     ctx.PostForm("avatar"),                    // 群组头像
		Password:   ctx.PostForm("password"),                  // 群组密码
		IsPwd:      Utils.StringToInt(ctx.PostForm("is_pwd")), // 是否设置密码
		Theme:      ctx.PostForm("theme"),                     // 群组主题
		SelectUser: selectUser.SelectUser,                     // 选定的成员ID列表
	}

	errs := validator.New().Struct(params) // 验证请求参数的有效性
	if errs != nil {
		Response.ErrorResponse(Enums.ParamError, errs.Error()).ToJson(ctx) // 若参数无效，则返回错误响应给客户端
		return
	}

	if params.IsPwd == IMGroups.IS_PWD_YES {
		params.Password, _ = Hash.SaltCryptoHashPassword(params.Password) // 对密码进行加密处理
	}

	err, imGroups := groupDao.CreateGroup(params) // 调用GroupDAO的CreateGroup方法，创建群组
	if err != nil {
		Response.FailResponse(Enums.ApiError, "创建群聊失败！").WriteTo(ctx) // 若创建群组失败，则返回失败响应给客户端
		return
	}

	groupDao.CreateSelectGroupUser(selectUser.SelectUser, int(imGroups.Id), params.Avatar, params.Name) // 创建选定的群组成员及相关会话

	messageService.SendGroupSessionMessage(selectUser.SelectUser, imGroups.Id) // 发送创建群聊消息

	Response.SuccessResponse(imGroups).WriteTo(ctx) // 返回成功响应给客户端
	return
}

// ApplyJoin ApplyJoin方法用于申请加入群组
// 参数：
// - ctx: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) ApplyJoin(ctx *gin.Context) {
	id := ctx.MustGet("id")              // 获取当前用户的ID
	err, person := Base.GetPersonId(ctx) // 调用GetPersonId方法，获取当前用户信息
	if err != nil {
		Response.FailResponse(Enums.ParamError, "参数错误！").WriteTo(ctx) // 若获取用户信息失败，则返回参数错误响应给客户端
		return
	}
	var group IMGroups.ImGroups
	if result := MYSQL.DataBase.Model(&IMGroups.ImGroups{}).Where("id=?", person.ID).Find(&group); result.RowsAffected == 0 {
		Response.FailResponse(Enums.ParamError, "群聊不存在！").WriteTo(ctx) // 若群聊不存在，则返回群聊不存在响应给客户端
		return
	}

	if groupDao.IsGroupsUser(id, person.ID) {
		Response.FailResponse(Enums.ParamError, "已经是群成员了~").WriteTo(ctx) // 若当前用户已经是群成员，则返回已经是群成员的响应给客户端
		return
	}
	if group.IsPwd == int8(IMGroups.IS_PWD_YES) {
		if !Hash.CheckPassword(ctx.PostForm("password"), group.Password) {
			Response.FailResponse(Enums.ParamError, "入群密码错误~,请联系管理员邀请").WriteTo(ctx) // 若入群密码错误，则返回密码错误的响应给客户端
			return
		}
	}

	groupDao.CreateOneGroupUser(group, int(Utils.InterfaceToInt64(id))) // 将当前用户添加到群组成员中

	name := ctx.MustGet("name") // 获取当前用户的名称

	groupDao.DeleteGroupUser(id, person.ID) // 将当前用户从原有的群组中删除

	// 构建私聊消息的请求参数
	params := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(),            // 消息ID
		MsgCode:     Enums.WsChatMessage,            // 消息类型
		MsgClientId: Date.TimeUnixNano(),            // 消息客户端ID
		FromID:      Utils.InterfaceToInt64(id),     // 发送人ID
		ToID:        Utils.StringToInt64(person.ID), // 接收人ID
		ChannelType: 2,                              // 渠道类型（私聊）
		MsgType:     IMMessages.JOIN_GROUP,          // 消息类型（加入群聊）
		Message:     fmt.Sprintf("%s 加入群聊", name),   // 消息内容
		SendTime:    Date.NewDate(),                 // 发送时间
		Data:        ctx.PostForm("data"),           // 消息数据
	}
	// 发送退群消息
	messageService.SendGroupMessage(params)

	Response.SuccessResponse().WriteTo(ctx) // 返回成功响应给客户端
	return
}

// GetUsers GetUsers方法用于获取群组成员列表
// 参数：
// - ctx: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) GetUsers(ctx *gin.Context) {
	err, person := Base.GetPersonId(ctx) // 调用GetPersonId方法，获取当前用户信息
	if err != nil {
		Response.FailResponse(Enums.ParamError, "参数错误！").WriteTo(ctx) // 若获取用户信息失败，则返回参数错误响应给客户端
		return
	}
	var group ImGroups
	if result := MYSQL.DataBase.Model(&IMGroups.ImGroups{}).Where("id=?", person.ID).Find(&group); result.RowsAffected == 0 {
		Response.FailResponse(Enums.ParamError, "群聊不存在！").WriteTo(ctx) // 若群聊不存在，则返回群聊不存在响应给客户端
		return
	}
	Response.SuccessResponse(&GroupsDate{
		Groups: group,
		Users:  groupDao.GetGroupUsers(person.ID), // 调用GetGroupUsers方法，获取群组成员列表
	}).WriteTo(ctx) // 返回成功响应和群组成员列表给客户端
	return
}

// Logout Logout方法用于退出群组
// 参数：
// - ctx: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) Logout(ctx *gin.Context) {
	err, person := Base.GetPersonId(ctx) // 调用GetPersonId方法，获取当前用户信息
	if err != nil {
		Response.FailResponse(Enums.ParamError, "参数错误！").WriteTo(ctx) // 若获取用户信息失败，则返回参数错误响应给客户端
		return
	}
	id := ctx.MustGet("id")     // 获取当前用户的ID
	name := ctx.MustGet("name") // 获取当前用户的名称

	groupDao.DeleteGroupUser(id, person.ID) // 将当前用户从群组中删除

	params := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(),            // 消息ID
		MsgCode:     Enums.WsChatMessage,            // 消息类型
		MsgClientId: Date.TimeUnixNano(),            // 消息客户端ID
		FromID:      Utils.InterfaceToInt64(id),     // 发送人ID
		ToID:        Utils.StringToInt64(person.ID), // 接收人ID
		ChannelType: IMSessions.GROUPTYPE,           // 渠道类型（群聊）
		MsgType:     IMMessages.LOGOUT_GROUP,        // 消息类型（退出群聊）
		Message:     fmt.Sprintf("%s 退出群聊", name),   // 消息内容
		SendTime:    Date.NewDate(),                 // 发送时间
		Data:        ctx.PostForm("data"),           // 消息数据
	}
	// 退群消息推送
	messageService.SendGroupMessage(params)

	Response.SuccessResponse().WriteTo(ctx) // 返回成功响应给客户端
	return
}

// AddOrRemoveUser CreateOrRemoveUser CreateOrRemoveUser方法用于创建或删除群组用户
// 参数：
// - cxt: *gin.Context类型，表示HTTP请求的上下文
func (*Handler) AddOrRemoveUser(ctx *gin.Context) {

	var selectUser SelectUser // 定义SelectUser结构体变量，用于存储选择的用户信息

	err := ctx.ShouldBind(&selectUser) // 将请求参数绑定到selectUser变量上
	if err != nil {
		return
	}

	params := ApiRequests.CreateUserToGroupRequest{
		GroupId: Utils.StringToInt64(ctx.PostForm("group_id")), // 将请求参数group_id转换为int64类型，并赋值给params的GroupId字段
		Type:    Utils.StringToInt(ctx.PostForm("type")),       // 将请求参数type转换为int类型，并赋值给params的Type字段
		UserId:  selectUser.SelectUser,                         // 将selectUser的SelectUser字段赋值给params的UserId字段
	}

	userId := ctx.MustGet("id") // 获取当前用户的ID
	name := ctx.MustGet("name") // 获取当前用户的名称
	var group ImGroups
	if result := MYSQL.DataBase.Model(&IMGroups.ImGroups{}).Where("id=?", params.GroupId).Find(&group); result.RowsAffected == 0 {
		Response.FailResponse(Enums.ParamError, "群聊不存在！").WriteTo(ctx) // 若群聊不存在，则返回群聊不存在响应给客户端
		return
	}
	if group.UserId != userId {
		Response.FailResponse(Enums.ParamError, "非群主不可以邀请人入群！").WriteTo(ctx) // 若当前用户不是群主，则返回非群主不可以邀请人入群响应给客户端
		return
	}

	if params.Type == 1 {
		groupDao.CreateSelectGroupUser(selectUser.SelectUser, int(params.GroupId), group.Avatar, group.Name) // 创建群组用户
		// 发送群聊会话消息
		messageService.SendGroupSessionMessage(selectUser.SelectUser, params.GroupId)
	} else {
		groupDao.DelSelectGroupUser(selectUser.SelectUser, int(params.GroupId), group.Avatar, group.Name) // 删除群组用户
	}
	var users []IMUser.ImUsers

	MYSQL.DataBase.Model(&IMUser.ImUsers{}).
		Where("id in(?)", MYSQL.DataBase.Model(&IMGroupUsers.ImGroupUsers{}).
			Where("group_id=?", params.GroupId).Select("user_id")).
		Find(&users)

	groupStr, _ := json.Marshal(group)
	message := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(), // 消息ID
		MsgCode:     Enums.WsChatMessage, // 消息类型
		MsgClientId: Date.TimeUnixNano(), // 消息客户端ID
		FromID:      group.Id,            // 发送人ID
		ChannelType: Enums.GroupMessage,  // 渠道类型（群聊消息）
		MsgType:     Enums.JoinGroup,     // 消息类型（加入群聊）
		Message:     "",                  // 消息内容
		SendTime:    Date.NewDate(),      // 发送时间
		Data:        string(groupStr),    // 消息数据
	}

	messageService.SendCreateUserGroupMessage(users, message, name, params.Type, selectUser.SelectUser) // 发送创建或删除群组用户的消息

	Response.SuccessResponse().WriteTo(ctx) // 返回成功响应给客户端
	return

}
