package Router

import (
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Friend"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/FriendRecords"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Group"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/GroupMessage"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Message"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Sessions"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/User"
	"HighPerformanceIMServer/Internal/Middleware"

	"github.com/gin-gonic/gin"
)

var (
	users         User.Handler
	sessions      Sessions.Handler
	friends       Friend.Handler
	friendRecord  FriendRecords.Handler
	messages      Message.Handler
	groupMessages GroupMessage.Handler
	groups        Group.Handler
)

func RegisterApiRoutes(routerEngine *gin.Engine) {
	// 设置允许跨域资源访问
	routerEngine.Use(Middleware.CORS())
	routerEngine.Use(Middleware.SetData())

	// 设置主路由组
	apiGroup := routerEngine.Group("/api/v1")
	// 用户接口
	userGroup := apiGroup.Group("/user")
	{
		// 获得 id 所对应的用户信息
		userGroup.GET("/:id", users.UserInfo)
		// 获得用户对应的通讯录列表，包含好友和群聊
		userGroup.Any("/contact/list", users.UserContactList)
	}
	// 会话接口
	sessionsGroup := apiGroup.Group("/sessions")
	{
		//  获得某用户所有会话列表
		sessionsGroup.GET("/list", sessions.GetSessionsList)
		// 新增会话
		sessionsGroup.POST("/add", sessions.AddSessions)
		// 更新会话
		sessionsGroup.PUT("/update/:id", sessions.UpdateSessions)
		// 删除会话
		sessionsGroup.DELETE("/delete/:id", sessions.DeleteSessions)
	}
	// 好友接口
	friendsGroup := apiGroup.Group("/friends")
	{
		// 获取相关信息
		friendsGroup.Any("/list", friends.GetFriendsList)
		friendsGroup.GET("/information/:id", friends.ShowFriendInformation)
		friendsGroup.DELETE("/delete/:id", friends.DeleteFriend)
		friendsGroup.GET("/status/:id", friends.GetUserStatus)
		// 相关请求
		friendsGroup.POST("/add/request", friendRecord.SendAddFriendsRequest)
		friendsGroup.GET("/add/request/list", friendRecord.GetFriendsRequestRecord)
		friendsGroup.PUT("/add/agree", friendRecord.AgreeOrRejectFriendRequest)
		friendsGroup.GET("/query", friendRecord.UserQuery)
	}
	// 消息接口
	messageGroup := apiGroup.Group("/message")
	{
		messageGroup.GET("/messages/private/list", messages.GetPrivateChatList)
		messageGroup.GET("/messages/group/list", groupMessages.GetGroupList)
		messageGroup.POST("/messages/private", messages.SendMessage)
		messageGroup.POST("/messages/group", messages.SendMessage)
		messageGroup.POST("/messages/video", messages.SendVideoMessage)
		messageGroup.POST("/messages/recall", messages.RecallMessage)
	}
	//	群聊接口
	groupGroup := apiGroup.Group("/groups")
	{
		groupGroup.POST("/createGroup", groups.CreateGroup)
		groupGroup.POST("/applyJoin/:id", groups.ApplyJoin)
		groupGroup.POST("/AddOrRemoveUser", groups.AddOrRemoveUser)
		groupGroup.GET("/lists", groups.GetGroupList)
		groupGroup.GET("/users/:id", groups.GetUsers)
		groupGroup.DELETE("/logout/:id", groups.Logout)
	}
	// 文件上传接口
	uploadGroup := apiGroup.Group("/upload")
	{
		uploadGroup.POST("/file")
	}
}
