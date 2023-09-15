package Date

// 导入time包，用于处理时间相关操作
import "time"

// NewDate 函数返回当前时间的字符串格式（年月日 时分秒）
func NewDate() string {
	// 使用time包的Now函数获取当前时间，Unix函数将其转换为Unix时间戳，然后使用Format函数将其格式化为指定的字符串格式
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}

// TimeUnixNano 函数返回当前时间的纳秒级Unix时间戳
func TimeUnixNano() int64 {
	// 使用time包的Now函数获取当前时间，UnixNano函数将其转换为纳秒级Unix时间戳
	return time.Now().UnixNano()
}

// TimeUnix 函数返回当前时间的秒级Unix时间戳
func TimeUnix() int64 {
	// 使用time包的Now函数获取当前时间，Unix函数将其转换为秒级Unix时间戳
	return time.Now().Unix()
}
