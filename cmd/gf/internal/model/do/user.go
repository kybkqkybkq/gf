// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta   `orm:"table:user, do:true"`
	Id       interface{} // 编号
	Name     interface{} // 名称
	Birthday interface{} // 生日
	Password interface{} // 密码
	UnitId   interface{} // 单位
	Phone    interface{} // 电话
	Email    interface{} // 邮箱
	CreateAt interface{} // 创建时间
	UpdateAt interface{} // 修改时间
	DeleteAt interface{} // 删除时间
}
