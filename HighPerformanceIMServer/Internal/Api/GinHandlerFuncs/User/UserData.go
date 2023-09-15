package User

// Person 定义一个名为Person的结构体
type Person struct {
	ID string `uri:"id" binding:"required"` // ID字段用于接收URI中的"id"参数，并且此参数是必需的
}

// Details 定义一个名为UserDetails的结构体
type Details struct {
	ID            int64  `gorm:"column:id;primaryKey" json:"id"`                // ID字段是一个64位整数，作为主键，并且在数据库中的列名为"id"，在JSON序列化时字段名为"id"
	Name          string `gorm:"column:name" json:"name"`                       // Name字段是一个字符串，在数据库中的列名为"name"，在JSON序列化时字段名为"name"
	Email         string `gorm:"column:email" json:"email"`                     // Email字段是一个字符串，在数据库中的列名为"email"，在JSON序列化时字段名为"email"
	Avatar        string `gorm:"column:avatar" json:"avatar"`                   // Avatar字段是一个字符串，在数据库中的列名为"avatar"，在JSON序列化时字段名为"avatar"
	Status        int8   `gorm:"column:status" json:"status"`                   // Status字段是一个8位整数，在数据库中的列名为"status"，在JSON序列化时字段名为"status"
	Bio           string `gorm:"column:bio" json:"bio"`                         // Bio字段是一个字符串，在数据库中的列名为"bio"，在JSON序列化时字段名为"bio"
	Sex           int8   `gorm:"column:sex" json:"sex"`                         // Sex字段是一个8位整数，在数据库中的列名为"sex"，在JSON序列化时字段名为"sex"
	Age           int    `gorm:"column:age" json:"age"`                         // Age字段是一个整数，在数据库中的列名为"age"，在JSON序列化时字段名为"age"
	LastLoginTime string `gorm:"column:last_login_time" json:"last_login_time"` // LastLoginTime字段是一个字符串，在数据库中的列名为"last_login_time"，在JSON序列化时字段名为"last_login_time"
	Uid           string `gorm:"column:uid" json:"uid"`                         // Uid字段是一个字符串，在数据库中的列名为"uid"，在JSON序列化时字段名为"uid"
}
