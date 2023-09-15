package ApiRequests

import (
	"HighPerformanceIMServer/Internal/DAO/MYSQL"
	"fmt"
)

// IsTableFliedExits
// 判断指定数据表中是否存在某个字段的值为特定值的记录。通过传入字段名、值和数据表名作为参数，函数会在指定的数据表中进行查询，并返回一个布尔值，表示是否存在满足条件的记录。
func IsTableFliedExits(filed string, value string, table string) bool {
	// 声明一个变量count，用于存储查询结果的数量
	var count int64

	// 在指定的数据表中查询满足条件的记录数量，并将结果存储到count变量中
	MYSQL.DataBase.Table(table).Where(fmt.Sprintf("%s=?", filed), value).Count(&count)

	// 如果count大于0，则表示存在满足条件的记录，返回true
	if count > 0 {
		return true
	}

	// 如果count等于0，则表示不存在满足条件的记录，返回false
	return false
}
