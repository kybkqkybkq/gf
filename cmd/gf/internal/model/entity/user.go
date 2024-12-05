// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// User is the golang structure for table user.
type User struct {
	Id       int64     `json:"id"       orm:"id"        description:"编号"`   // 编号
	Name     string    `json:"name"     orm:"name"      description:"名称"`   // 名称
	Birthday time.Time `json:"birthday" orm:"birthday"  description:"生日"`   // 生日
	Password string    `json:"password" orm:"password"  description:"密码"`   // 密码
	UnitId   int64     `json:"unitId"   orm:"unit_id"   description:"单位"`   // 单位
	Phone    string    `json:"phone"    orm:"phone"     description:"电话"`   // 电话
	Email    string    `json:"email"    orm:"email"     description:"邮箱"`   // 邮箱
	CreateAt time.Time `json:"createAt" orm:"create_at" description:"创建时间"` // 创建时间
	UpdateAt time.Time `json:"updateAt" orm:"update_at" description:"修改时间"` // 修改时间
	DeleteAt time.Time `json:"deleteAt" orm:"delete_at" description:"删除时间"` // 删除时间
}
