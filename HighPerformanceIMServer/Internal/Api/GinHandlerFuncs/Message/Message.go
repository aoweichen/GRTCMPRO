package Message

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMMessages"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/Api/Services/IMMessage"
	"HighPerformanceIMServer/Internal/DAO/Friend"
	"HighPerformanceIMServer/Internal/DAO/Group"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Internal/DAO/MessageDAO"
	"HighPerformanceIMServer/Packages/Date"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Response"
	"HighPerformanceIMServer/Packages/Utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sort"
)

type Handler struct {
}

type HandlerInterface interface {
}

var (
	messageServices IMMessage.MessageService
	friendDao       Friend.DAO
	messageDao      MessageDAO.Dao
	groupDao        Group.DAO
)

func (MH *Handler) sortByMessage(imMessageList []IMMessages.ImMessages, users IMUser.ImUsers) {
	sort.Slice(imMessageList, func(i, j int) bool {
		imMessageList[i].Users.ID = users.ID
		imMessageList[i].Users.Name = users.Name
		imMessageList[i].Users.Email = users.Email
		imMessageList[i].Users.Avatar = users.Avatar
		return imMessageList[i].Id < imMessageList[j].Id
	})
}

func (MH *Handler) RecallMessage(ctx *gin.Context) {
	Response.SuccessResponse().ToJson(ctx)
}

// GetPrivateChatList GetPrivateChatList函数用于获取私聊消息列表
// 参数:
// - ctx: gin.Context对象，用于处理HTTP请求和响应
func (MH *Handler) GetPrivateChatList(ctx *gin.Context) {
	// 定义变量friendChatsList，用于存储IMMessages.ImMessages对象列表
	var friendChatsList []IMMessages.ImMessages
	// 定义变量total，用于存储消息总数
	var total int64
	// 定义变量users，用于存储IMUser.ImUsers对象
	var users IMUser.ImUsers

	// 获取发送者ID、页码、接收者ID和每页大小
	id, page, toID, pageSize := ctx.MustGet("id"), ctx.Query("page"), ctx.Query("to_id"),
		Utils.StringToInt(ctx.DefaultQuery("pageSize", "50"))

	// 构建查询条件，并按创建时间降序排序
	query := MYSQL.DataBase.Table("im_messages").
		Where("(from_id=? and to_id=?) or (from_id=? and to_id=?)", id, toID, toID, id).
		Order("created_at desc")

	// 统计匹配的消息总数
	query.Count(&total)

	// 查询接收者信息
	MYSQL.DataBase.Table("im_users").Where("id=?", toID).First(&users)

	// 根据页码设置查询条件
	if len(page) > 0 {
		query = query.Where("id<?", page)
	}

	// 查询私聊消息列表
	if result := query.Limit(pageSize).Find(&friendChatsList); result.RowsAffected == 0 {
		// 如果查询结果为空，则返回空列表响应
		Response.SuccessResponse(gin.H{
			"list": struct{}{},
			"mate": gin.H{
				"pageSize": pageSize,
				"page":     page,
				"total":    0,
			},
		}, http.StatusOK).ToJson(ctx)
	}

	// 对私聊消息列表进行排序
	MH.sortByMessage(friendChatsList, users)

	// 返回私聊消息列表响应
	Response.SuccessResponse(gin.H{
		"list": friendChatsList,
		"mate": gin.H{
			"pageSize": pageSize,
			"page":     page,
			"total":    total,
		},
	}, http.StatusOK).ToJson(ctx)
	return
}

// SendVideoMessage SendVideoMessage函数用于发送视频消息
// 参数:
// - ctx: gin.Context对象，用于处理HTTP请求和响应
func (MH *Handler) SendVideoMessage(ctx *gin.Context) {
	// 定义变量users，用于存储IMUser.ImUsers对象
	var users IMUser.ImUsers

	// 获取发送者ID和接收者ID
	id, toID := ctx.MustGet("id"), ctx.PostForm("to_id")

	// 判断是否为好友关系
	if !friendDao.IsFriends(id, toID) {
		Response.FailResponse(Enums.WsNotFriend, "非好友关系，不能聊天......").ToJson(ctx)
		return
	}

	// 查询发送者信息
	MYSQL.DataBase.Table("im_users").Where("id=?", id).First(&users)

	// 构建视频消息参数
	params := ApiRequests.VideoMessageRequest{
		MsgCode:  Enums.VideoChantMessage,    // 消息代码
		FromID:   Utils.InterfaceToInt64(id), // 发送者ID
		ToID:     Utils.StringToInt64(toID),  // 接收者ID
		Message:  "视频请求......",               // 消息内容
		SendTime: Date.NewDate(),             // 发送时间
		Users: ApiRequests.Users{
			Email:  users.Email,  // 发送者邮箱
			Name:   users.Name,   // 发送者姓名
			Avatar: users.Avatar, // 发送者头像
		},
	}

	// 发送视频消息
	if ok := messageServices.SendVideoMessage(params); !ok {
		Response.FailResponse(http.StatusInternalServerError, "用户不在线").ToJson(ctx)
		return
	} else {
		Response.SuccessResponse(params).ToJson(ctx)
		return
	}
}

// SendMessage SendMessage函数用于发送私聊消息
// 参数:
// - ctx: gin.Context对象，用于处理HTTP请求和响应
func (MH *Handler) SendMessage(ctx *gin.Context) {

	// 获取用户ID
	id := ctx.MustGet("id")

	// 构建私聊消息参数
	params := ApiRequests.PrivateMessageRequest{
		MsgId:       Date.TimeUnixNano(),                                // 消息ID
		MsgClientId: Utils.StringToInt64(ctx.PostForm("msg_client_id")), // 消息客户端ID
		MsgCode:     Enums.WsChatMessage,                                // 消息代码
		FromID:      Utils.InterfaceToInt64(id),                         // 发送者ID
		ToID:        Utils.StringToInt64(ctx.PostForm("to_id")),         // 接收者ID
		MsgType:     Utils.StringToInt(ctx.PostForm("msg_type")),        // 消息类型
		ChannelType: Utils.StringToInt(ctx.PostForm("channel_type")),    // 频道类型，默认为1
		Message:     ctx.PostForm("message"),                            // 消息内容
		SendTime:    Date.NewDate(),                                     // 发送时间
		Data:        ctx.PostForm("data"),                               // 数据
	}

	// 参数校验
	if err := validator.New().Struct(params); err != nil {
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ctx)
		return
	}

	// 根据频道类型进行不同的处理
	switch params.ChannelType {
	case 1:
		var users IMUser.ImUsers
		messageDao.CreateMessage(params)
		MYSQL.DataBase.Model(&IMUser.ImUsers{}).Where("id=?", params.ToID).Find(&users)

		// 如果接收者是机器人用户
		if users.UserType == IMUser.BotType {
			if messageServices.IsUserOnline(Utils.Int64ToString(users.ID)) {
				messageServices.SendPrivateMessage(params)
			} else {
				messageServices.SendChatMessage(params)
			}
			Response.SuccessResponse(params).ToJson(ctx)
			return
		} else {
			// 如果接收者是普通用户
			if !friendDao.IsFriends(id, params.ToID) {
				Response.FailResponse(Enums.WsNotFriend, "非好友关系，无法聊天......")
				return
			} else {
				if ok, message := messageServices.SendPrivateMessage(params); !ok {
					Response.FailResponse(http.StatusInternalServerError, message).ToJson(ctx)
					return
				}
			}
		}
	case 2:
		// 群聊消息
		if !groupDao.IsGroupsUser(id, params.ToID) {
			Response.FailResponse(Enums.WsNotFriend, "你不是此群成员......").ToJson(ctx)
			return
		}

		if ok := messageServices.SendGroupMessage(params); !ok {
			Response.FailResponse(http.StatusInternalServerError, "群聊消息投递异常！").ToJson(ctx)
			return
		}
	}
	Response.SuccessResponse(params).ToJson(ctx)
	return
}
