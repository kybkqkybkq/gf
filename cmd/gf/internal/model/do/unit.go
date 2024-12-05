// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Unit is the golang structure of table unit for DAO operations like Where/Data.
type Unit struct {
	g.Meta   `orm:"table:unit, do:true"`
	Id       interface{} // 编号
	Name     interface{} // 名称
	Address  interface{} // 地址
	BossId   interface{} // 老板
	CreateAt interface{} // 创建时间
	UpdateAt interface{} // 修改时间
	DeleteAt interface{} // 删除时间
}
