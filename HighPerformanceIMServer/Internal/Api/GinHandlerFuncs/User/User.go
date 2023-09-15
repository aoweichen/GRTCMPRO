package User

import (
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"HighPerformanceIMServer/Internal/DAO/Friend"
	"HighPerformanceIMServer/Internal/DAO/Group"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
}

var (
	friendDao Friend.DAO
	groupDao  Group.DAO
)

// UserInfo
// 定义一个名为UserInfo的方法，该方法属于UserHandler结构体
// 该方法用于处理获取用户信息的请求
// 参数ginCtx是一个指向gin.Context类型的指针，用于处理HTTP请求和响应
func (*Handler) UserInfo(ginCtx *gin.Context) {
	// 创建一个名为person的变量，类型为Person结构体
	var person Person

	// 调用ginCtx的ShouldBindUri方法，将URI参数绑定到person变量
	// 如果绑定失败，则将错误信息打印到日志，并返回参数错误的响应
	if ginCtxShouldBindUriError := ginCtx.ShouldBindUri(&person); ginCtxShouldBindUriError != nil {
		zap.S().Errorln(ginCtxShouldBindUriError)
		Response.FailResponse(Enums.ParamError, ginCtxShouldBindUriError.Error()).WriteTo(ginCtx)
		return
	}

	// 创建一个名为users的变量，类型为UserDetails结构体
	var users Details

	// 从数据库中查询用户信息，并将结果赋值给users变量
	// 如果查询结果中的行数为0，则表示用户不存在，返回用户不存在的响应
	if result := MYSQL.DataBase.Model(&IMUser.ImUsers{}).
		Where("id=?", person.ID).
		First(&users); result.RowsAffected == 0 {
		zap.S().Errorln("用户不存在")
		Response.ErrorResponse(Enums.ParamError, "用户不存在").ToJson(ginCtx)
		return
	}

	// 返回成功的响应，将users变量作为响应数据
	Response.SuccessResponse(users).ToJson(ginCtx)
	return
}

// UserContactList
// 定义名为UserContactList的方法，该方法属于UserHandler结构体
// 该方法用于获取用户的联系人列表
// 参数ginCtx是一个指向gin.Context类型的指针，表示Gin框架的上下文对象
func (*Handler) UserContactList(ginCtx *gin.Context) {
	// 从ginCtx中获取用户ID，并将其赋值给userID变量
	userID := ginCtx.MustGet("id")

	// 调用friendDao的GetFriendLists方法，获取用户的好友列表
	friendDaoGetFriendListsError, friendLists := friendDao.GetFriendLists(userID)

	// 如果获取好友列表时出现错误，则记录错误日志并返回错误响应
	if friendDaoGetFriendListsError != nil {
		zap.S().Errorln("获取用户好友列表失败！", friendDaoGetFriendListsError)
		Response.FailResponse(Enums.ParamError, "获取用户好友列表失败！").ToJson(ginCtx)
		return
	}

	// 调用groupDao的GetGroupListJoinedByUser方法，获取用户加入的群组列表
	groupsLists := groupDao.GetGroupListJoinedByUser(userID)

	// 返回成功响应，包含好友列表
	Response.SuccessResponse(friendLists).ToJson(ginCtx)

	// 返回成功响应，包含好友列表和群组列表
	Response.SuccessResponse(gin.H{
		"friends": friendLists,
		"groups":  groupsLists,
	}).ToJson(ginCtx)
}
