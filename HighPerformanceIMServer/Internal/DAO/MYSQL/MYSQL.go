package MYSQL

import (
	"HighPerformanceIMServer/Configs"
	"HighPerformanceIMServer/DataModels/Models/IMFriends"
	"HighPerformanceIMServer/DataModels/Models/IMFriendsRecords"
	"HighPerformanceIMServer/DataModels/Models/IMGroupMessages"
	"HighPerformanceIMServer/DataModels/Models/IMGroupUsers"
	"HighPerformanceIMServer/DataModels/Models/IMGroups"
	"HighPerformanceIMServer/DataModels/Models/IMMessages"
	"HighPerformanceIMServer/DataModels/Models/IMOfflineMessage"
	"HighPerformanceIMServer/DataModels/Models/IMSessions"
	"HighPerformanceIMServer/DataModels/Models/IMUser"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DataBase *gorm.DB

// InitMySQLDB 初始化 mysql 数据库
func InitMySQLDB() {
	//	得到dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		Configs.ConfigData.MySQL.Username,
		Configs.ConfigData.MySQL.Password,
		Configs.ConfigData.MySQL.Host,
		Configs.ConfigData.MySQL.Port,
		Configs.ConfigData.MySQL.Database,
		Configs.ConfigData.MySQL.Charset)

	var gormMySQLConnectedError error
	DataBase, gormMySQLConnectedError = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if gormMySQLConnectedError != nil {
		zap.S().Errorf("Mysql 连接异常: ", gormMySQLConnectedError)
		panic(gormMySQLConnectedError)
	}
	zap.S().Info("mysql 连接成功")
	// 迁移表 IMUser.ImUsers
	AutoMigrateDataStruct(&IMUser.ImUsers{})
	// 迁移表 IMSessions.ImSessions{}
	AutoMigrateDataStruct(&IMSessions.ImSessions{})
	// 迁移表 IMGroupUsers.ImGroupUsers{}
	AutoMigrateDataStruct(&IMGroupUsers.ImGroupUsers{})
	// 迁移表 IMGroups.ImGroups{}
	AutoMigrateDataStruct(&IMGroups.ImGroups{})

	// 迁移表 IMOfflineMessage.ImOfflineMessages{}
	AutoMigrateDataStruct(&IMOfflineMessage.ImOfflineMessages{})
	// 迁移表 IMOfflineMessage.ImGroupOfflineMessages{}
	AutoMigrateDataStruct(&IMOfflineMessage.ImGroupOfflineMessages{})
	// 迁移表 IMMessages.ImMessages{}
	AutoMigrateDataStruct(&IMMessages.ImMessages{})
	// 迁移表 IMGroupMessages.ImGroupMessages{}
	AutoMigrateDataStruct(&IMGroupMessages.ImGroupMessages{})
	// 迁移表 IMFriendsRecords.ImFriendRecords{}
	AutoMigrateDataStruct(&IMFriendsRecords.ImFriendRecords{})
	// 迁移表 IMFriends.ImFriends{}
	AutoMigrateDataStruct(&IMFriends.ImFriends{})
}

// AutoMigrateDataStruct
// 定义名为MYSQLAutoMigrateDataStruct的函数
// 该函数用于自动迁移数据结构到MySQL数据库
// 参数dataStruct是一个接口类型，表示要迁移的数据结构
func AutoMigrateDataStruct(dataStruct interface{}) {
	// 调用MYSQLDB的AutoMigrate方法，将数据结构迁移到数据库中
	err := DataBase.AutoMigrate(dataStruct)

	// 如果迁移过程中出现错误，则记录错误日志并抛出panic
	if err != nil {
		zap.S().Errorf("迁移表 %#v 出错: %#v", dataStruct, err)
		panic(err)
	}

	// 记录迁移成功的日志
	zap.S().Infof("迁移表 %#v 成功", dataStruct)
}
