package MySQL

import (
	"HighPerformanceIMServerAuth/Configs"
	"HighPerformanceIMServerAuth/DataModels/IMUsers"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DataBase *gorm.DB

func InitMySQLDB() {
	//	得到dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		Configs.ConfigData.MySQL.Username,
		Configs.ConfigData.MySQL.Password,
		Configs.ConfigData.MySQL.Host,
		Configs.ConfigData.MySQL.Port,
		Configs.ConfigData.MySQL.Database,
		Configs.ConfigData.MySQL.Charset,
	)
	// gorm连接mysql
	var gormConnectError error
	if DataBase, gormConnectError = gorm.Open(
		mysql.New(mysql.Config{DSN: dsn}),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)}); gormConnectError != nil {
		zap.S().Panic("Mysql 连接异常: ", gormConnectError)
		panic(gormConnectError)
	}
	zap.S().Infoln("mysql 连接成功")
	//	迁移表 IMUser.ImUsers
	AutoMigrateDataStruct(&IMUsers.ImUsers{})
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
