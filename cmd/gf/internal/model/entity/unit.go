// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Unit is the golang structure for table unit.
type Unit struct {
	Id       int64     `json:"id"       orm:"id"        description:"编号"`   // 编号
	Name     string    `json:"name"     orm:"name"      description:"名称"`   // 名称
	Address  string    `json:"address"  orm:"address"   description:"地址"`   // 地址
	BossId   int64     `json:"bossId"   orm:"boss_id"   description:"老板"`   // 老板
	CreateAt time.Time `json:"createAt" orm:"create_at" description:"创建时间"` // 创建时间
	UpdateAt time.Time `json:"updateAt" orm:"update_at" description:"修改时间"` // 修改时间
	DeleteAt time.Time `json:"deleteAt" orm:"delete_at" description:"删除时间"` // 删除时间
}
