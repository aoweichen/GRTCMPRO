package Logger

import (
	"HighPerformanceIMServerAuth/Configs"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"time"
)

// InitLogger 定义一个名为InitLogger的函数，该函数接收7个参数：filename（文件名），maxSize（最大文件大小），maxBackup（最大备份数），maxAge（最大文件存活时间），compress（是否压缩），logType（日志类型）和level（日志级别）
func InitLogger() {
	filename := Configs.ConfigData.Logger.LogFilePath
	maxBackup := Configs.ConfigData.Logger.MaxBackup
	maxSize := Configs.ConfigData.Logger.MaxSize
	maxAge := Configs.ConfigData.Logger.MaxAge
	compress := Configs.ConfigData.Logger.Compress
	logType := Configs.ConfigData.Logger.Type
	level := Configs.ConfigData.Logger.Level

	// 调用getWriter函数，传入filename、maxSize、maxBackup、maxAge、compress和logType等参数，返回一个writeSyncer对象，用于写入日志数据
	writeSyncer := getWriter(filename, int(maxSize), int(maxBackup), int(maxAge), compress, logType)
	// 声明一个指向zapcore.Level类型的指针变量logLevel
	logLevel := new(zapcore.Level)
	// 使用UnmarshalText方法将level字符串转换为[]byte类型，并赋值给logLevel指针所指向的对象
	// 如果转换过程中出现错误，则打印错误信息并退出函数
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置有误. 请调整配置")
	}
	// 使用newCore方法创建一个新的zapcore.Core实例，传入getEncoder()函数返回的编码器、writeSyncer和logLevel作为参数
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)
	// 使用ReplaceGlobals方法替换全局的zap.Logger实例，传入zap.New函数创建的新实例以及zap.AddCaller、zap.AddCallerSkip和zap.AddStacktrace方法返回的配置选项作为参数
	zap.ReplaceGlobals(zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel)))
}

// 定义一个名为 getWriter 的函数，参数包括 logFilePath（日志文件路径）、maxSize（最大文件大小）、maxBackup（最大备份数）、maxAge（最大文件存活时间）、compress（是否压缩）、logType（日志类型）
func getWriter(logFilePath string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 如果日志类型为 "daily"，则将当前日期格式化为 "2006-01-02.log" 的形式，并将日志文件名替换为该格式的日期字符串
	if logType == "daily" {
		logName := time.Now().Format("2006-01-02.log")
		logFilePath = strings.ReplaceAll(logFilePath, "logs.log", logName)
	}
	// 创建一个 lumberjack.Logger 实例，用于记录日志到指定文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFilePath, // 设置日志文件路径
		MaxSize:    maxSize,     // 设置最大文件大小
		MaxBackups: maxBackup,   // 设置最大备份数
		MaxAge:     maxAge,      // 设置最大文件存活时间
		Compress:   compress,    // 设置是否压缩日志文件
	}
	// 将日志记录到文件中
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoder 设置日志存储格式
func getEncoder() zapcore.Encoder {
	// 定义一个encoderConfig变量，该变量是zapcore.EncoderConfig类型的，用于配置日志的编码规则
	encoderConfig := zapcore.EncoderConfig{
		// 设置日志的时间戳键名为"time"
		TimeKey: "time",
		// 设置日志的级别键名为"level"
		LevelKey: "level",
		// 设置日志的logger名称键名为"logger"
		NameKey: "logger",
		// 设置日志的调用者信息键名为"caller"，这在日志中会显示为"paginator/paginator.go:148"这样的格式
		CallerKey: "caller",
		// 设置日志的函数名键名为空，表示不记录函数名
		FunctionKey: zapcore.OmitKey,
		// 设置日志的消息内容键名为"message"
		MessageKey: "message",
		// 设置日志的堆栈跟踪信息键名为"stacktrace"
		StacktraceKey: "stacktrace",
		// 设置日志的行结束符为默认的"\n"
		LineEnding: zapcore.DefaultLineEnding,
		// 设置日志的级别编码器为大写，如ERROR、INFO等
		EncodeLevel: zapcore.CapitalLevelEncoder,
		// 设置日志的时间编码器为我们自定义的函数customTimeEncoder
		EncodeTime: customTimeEncoder,
		// 设置日志的执行时间编码器为秒为单位
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// 设置日志的调用者信息编码器为短格式，如types/converter.go:17，长格式为绝对路径
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	// 线上环境使用 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder 自定义友好的时间格式
func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05"))
}
