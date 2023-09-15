package Utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// GetNowFormatTodayTime 函数返回当前日期的格式化字符串（年-月-日）
func GetNowFormatTodayTime() string {
	// 获取当前时间
	now := time.Now()
	// 使用fmt.Sprintf函数将年、月、日格式化为字符串，并使用"-"作为分隔符
	dateStr := fmt.Sprintf("%02d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
	return dateStr
}

// CreateEmailCode 函数生成一个六位数的验证码字符串
func CreateEmailCode() string {
	// 使用rand.NewSource函数创建一个随机数生成器，并使用当前时间的纳秒级Unix时间戳作为种子
	// 使用rand.New函数创建一个新的随机数生成器
	// 使用Int31n方法生成一个0到9999之间的随机整数，并使用fmt.Sprintf函数将其格式化为四位数的字符串
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

// GetDayTime 函数返回指定天数后的时间戳（秒级）
func GetDayTime(days int) int64 {
	// 使用time.Now函数获取当前时间，并使用Format方法将其格式化为"年-月-日 时:分:秒"的字符串
	nowTimeStr := time.Now().Format("2006-01-02 15:04:05")
	// 使用time.ParseInLocation函数将字符串解析为时间，并指定时区为本地时区
	timeS, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	// 使用AddDate方法将时间增加指定的天数，并使用Unix方法将其转换为时间戳（秒级）
	timeStamp := timeS.AddDate(0, 0, days).Unix()
	return timeStamp
}

// Int64ToString 函数将int64类型的整数转换为字符串
func Int64ToString(int64_ int64) string {
	// 使用strconv.Itoa 函数将int类型的整数转换为字符串
	return strconv.Itoa(int(int64_))
}

// Float64ToString 函数将float64类型的浮点数转换为字符串
func Float64ToString(float64_ float64) string {
	// 使用strconv.Itoa函数将int类型的整数转换为字符串
	return strconv.Itoa(int(float64_))
}

// StringToInt 函数接受一个字符串参数 str，并尝试将其转换为整数类型
func StringToInt(str string) int {
	// 使用 strconv.Atoi 函数将字符串 str 转换为整数 num
	num, strconvAtoiError := strconv.Atoi(str)

	// 如果转换过程中出现错误（即 strconvAtoiError 不为 nil），则输出错误信息并抛出异常
	if strconvAtoiError != nil {
		zap.S().Errorln("字符串转为Int数字错误：", strconvAtoiError)
		panic(strconvAtoiError.Error())
	}

	// 如果没有出现错误，返回转换后的整数 num
	return num
}

// StringToInt64 函数将字符串转换为int64类型的整数
func StringToInt64(str string) int64 {
	// 使用strconv.Atoi函数将字符串转换为int类型的整数
	num, strconvAtoiError := strconv.Atoi(str)
	// 如果转换过程中出现错误（即 strconvAtoiError 不为 nil），则输出错误信息并抛出异常
	if strconvAtoiError != nil {
		zap.S().Errorln("字符串转为Int64数字错误：", strconvAtoiError)
		panic(strconvAtoiError.Error())
	}
	return int64(num)
}

// FirstElement 函数返回字符串切片的第一个元素
func FirstElement(args []string) string {
	// 判断字符串切片的长度是否大于0
	if len(args) > 0 {
		// 返回第一个元素
		return args[0]
	}
	// 若长度为0，则返回空字符串
	return ""
}

// Explode 函数将字符串按照指定的分隔符拆分为字符串切片
func Explode(delimiter, text string) []string {
	// 判断分隔符的长度是否大于待拆分的字符串的长度
	if len(delimiter) > len(text) {
		// 使用strings.Split函数将分隔符和待拆分的字符串拆分为字符串切片
		return strings.Split(delimiter, text)
	} else {
		// 使用strings.Split函数将待拆分的字符串按照分隔符拆分为字符串切片
		return strings.Split(text, delimiter)
	}
}

// InterfaceToInt64 函数将接口类型转换为int64类型的整数
func InterfaceToInt64(inter interface{}) int64 {
	// 将接口类型断言为int64类型，并返回结果
	return inter.(int64)
}

// InterfaceToInt64String 函数将接口类型转换为int64类型的整数，并将其转换为字符串
func InterfaceToInt64String(inter interface{}) string {
	// 将接口类型断言为int64类型，并将其转换为字符串
	int64Val := inter.(int64)
	return Int64ToString(int64Val)
}

// InterfaceToString 函数将接口类型转换为字符串
func InterfaceToString(inter interface{}) string {
	// 将接口类型断言为字符串类型，并返回结果
	return inter.(string)
}

// ErrorHandler 函数用于处理错误，如果err不为nil，则将错误信息记录到日志中
func ErrorHandler(err error) {
	// 判断err是否为nil
	if err != nil {
		// 将错误信息记录到日志中
		zap.S().Error(err.Error())
		return
	}
}
