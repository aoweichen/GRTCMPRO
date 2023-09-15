package Sessions

import (
	"HighPerformanceIMServer/DataModels/ApiRequests"
	"HighPerformanceIMServer/DataModels/Models/IMSessions"
	"HighPerformanceIMServer/Internal/Api/GinHandlerFuncs/Base"
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"HighPerformanceIMServer/Internal/DAO/Sessions"
	"HighPerformanceIMServer/Packages/Enums"
	"HighPerformanceIMServer/Packages/Response"
	"HighPerformanceIMServer/Packages/Utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
}

type Person struct {
	ID string `uri:"id" binding:"required"`
}

// GetSessionsList
// 定义名为GetSessionsList的方法，该方法属于SessionHandler结构体
// 该方法用于获取会话列表
// 参数ginCtx是一个指向gin.Context类型的指针，表示Gin框架的上下文对象
func (s Handler) GetSessionsList(ginCtx *gin.Context) {
	// 从ginCtx中获取用户ID，并将其赋值给id变量
	id := ginCtx.MustGet("id")

	var sessionsList []IMSessions.ImSessions

	// 查询会话列表
	if result := MYSQL.DataBase.Model(&IMSessions.ImSessions{}).
		Preload("Users"). // 使用Preload方法预加载Users和Groups关联数据，确保在查询结果中包含关联数据。
		Preload("Groups").
		Where("from_id=? and status=0", id).
		Order("top_status desc").
		Find(&sessionsList); result.RowsAffected == 0 {

		// 如果查询结果中的行数为0，则返回成功响应
		zap.S().Info("获取sessions list 成功")
		Response.SuccessResponse().ToJson(ginCtx)
		return
	} else {
		// 如果查询结果中的行数不为0，则返回成功响应，包含会话列表
		Response.SuccessResponse(sessionsList).ToJson(ginCtx)
		return
	}
}

// AddSessions 定义一个名为AddSessions的方法，参数为SessionHandler类型和gin.Context类型
func (s Handler) AddSessions(ginCtx *gin.Context) {
	// 从gin.Context中获取id值，必须存在，否则会panic
	id := ginCtx.MustGet("id")

	// 定义一个名为params的结构体变量，类型为ApiRequests.SessionStore
	// 通过ginCtx.PostForm获取表单中的id和type值，并将其转换为int64和int类型
	params := ApiRequests.SessionStore{
		Id:   Utils.StringToInt64(ginCtx.PostForm("id")),
		Type: Utils.StringToInt(ginCtx.PostForm("type")),
	}

	// 使用validator.New()创建一个新的验证器，并对params进行验证
	// 如果验证出错，则将错误信息打印出来，并返回参数错误的响应给客户端
	validateNewStructError := validator.New().Struct(params)
	if validateNewStructError != nil {
		zap.S().Errorln(validateNewStructError)
		Response.FailResponse(Enums.ParamError, validateNewStructError.Error()).ToJson(ginCtx)
		return
	}

	// 定义一个名为sessions的IMSessions类型变量
	var sessions IMSessions.ImSessions
	// 查询数据库表"im_sessions"，条件为from_id和to_id等于id和params.Id
	// 如果查询结果的行数大于0，则表示sessions已存在，将其返回给客户端
	if result := MYSQL.DataBase.Table("im_sessions").
		Where("from_id=? and to_id=?", id, params.Id).
		First(&sessions); result.RowsAffected > 0 {
		zap.S().Infoln("sessions已存在，返回已存在的sessions")
		Response.SuccessResponse(sessions).ToJson(ginCtx)
		return
	} else {
		// 如果sessions不存在，则创建一个新的sessions，并返回给客户端
		var sessionsDao Sessions.DAO
		sessions := sessionsDao.CreateSession(Utils.InterfaceToInt64(id), params.Id, params.Type)
		zap.S().Infoln("创建新的sessions成功")
		Response.SuccessResponse(sessions).ToJson(ginCtx)
		return
	}
}

// UpdateSessions 定义一个名为UpdateSessions的方法，参数为SessionHandler类型和gin.Context类型
// 置顶更新或者会话备注更新
func (s Handler) UpdateSessions(ginCtx *gin.Context) {
	// 调用BaseHandlerFunc.GetPersonId方法获取person信息，并返回错误和person
	err, person := Base.GetPersonId(ginCtx)
	if err != nil {
		// 如果获取person失败，则记录错误日志并返回参数错误的响应给客户端
		zap.S().Errorln("参数错误")
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ginCtx)
		return
	}

	// 定义一个名为params的结构体变量，类型为ApiRequests.SessionUpdate
	// 通过ginCtx.PostForm获取表单中的top_status和note值，并将其转换为int和string类型
	params := ApiRequests.SessionUpdate{
		TopStatus: Utils.StringToInt(ginCtx.PostForm("top_status")),
		Note:      ginCtx.PostForm("note"),
	}

	// 对params进行验证
	if errs := validator.New().Struct(params); errs != nil {
		// 如果验证出错，则记录错误日志并返回参数错误的响应给客户端
		zap.S().Errorln(errs.Error())
		Response.FailResponse(Enums.ParamError, errs.Error()).ToJson(ginCtx)
		return
	}

	// 更新数据库中id等于person.ID的IMSessions记录，更新内容为params的值
	MYSQL.DataBase.Model(&IMSessions.ImSessions{}).Where("id", person.ID).Updates(&params)

	// 记录更新成功的日志，并返回成功的响应给客户端
	zap.S().Infoln("更新成功")
	Response.SuccessResponse().ToJson(ginCtx)
}

// DeleteSessions 定义一个名为DeleteSessions的方法，参数为SessionHandler类型和gin.Context类型
func (s Handler) DeleteSessions(ginCtx *gin.Context) {
	// 调用BaseHandlerFunc.GetPersonId方法获取person信息，并返回错误和person
	err, person := Base.GetPersonId(ginCtx)
	if err != nil {
		// 如果获取person失败，则记录错误日志并返回参数错误的响应给客户端
		zap.S().Errorln("参数错误")
		Response.FailResponse(Enums.ParamError, err.Error()).ToJson(ginCtx)
		return
	}

	// 在数据库中删除id等于person.ID的IMSessions记录
	MYSQL.DataBase.Delete(&IMSessions.ImSessions{}, person.ID)

	// 记录删除会话成功的日志，并返回成功的响应给客户端
	zap.S().Infoln("删除会话成功")
	Response.SuccessResponse().ToJson(ginCtx)
	return
}
