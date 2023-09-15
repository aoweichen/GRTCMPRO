package Base

import "github.com/gin-gonic/gin"

// Person 定义一个名为Person的结构体类型
type Person struct {
	ID string `uri:"id" binding:"required"` // ID字段表示请求中的id参数，binding:"required"表示该字段是必需的
}

// GetPersonId 定义一个名为GetPersonId的函数，参数为gin.Context类型，返回值为error和Person类型
func GetPersonId(cxt *gin.Context) (error, Person) {
	var person Person
	// 使用cxt.ShouldBindUri方法将请求中的参数绑定到person结构体中
	if err := cxt.ShouldBindUri(&person); err != nil {
		return err, person
	}
	return nil, person
}
