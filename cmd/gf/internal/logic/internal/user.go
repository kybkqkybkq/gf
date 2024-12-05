// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserLogic is the data access object for table user.
type UserLogic struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns UserColumns // columns contains all the column names of Table for convenient usage.
}

// UserColumns defines and stores column names for table user.
type UserColumns struct {
	Id       string // 编号
	Name     string // 名称
	Birthday string // 生日
	Password string // 密码
	UnitId   string // 单位
	Phone    string // 电话
	Email    string // 邮箱
	CreateAt string // 创建时间
	UpdateAt string // 修改时间
	DeleteAt string // 删除时间
}

// userColumns holds the columns for table user.
var userColumns = UserColumns{
	Id:       "id",
	Name:     "name",
	Birthday: "birthday",
	Password: "password",
	UnitId:   "unit_id",
	Phone:    "phone",
	Email:    "email",
	CreateAt: "create_at",
	UpdateAt: "update_at",
	DeleteAt: "delete_at",
}

// NewUserLogic creates and returns a new DAO object for table data access.
func NewUserLogic() *UserLogic {
	return &UserLogic{
		group:   "default",
		table:   "user",
		columns: userColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (logic *UserLogic) DB() gdb.DB {
	return g.DB(logic.group)
}

// Table returns the table name of current logic.
func (logic *UserLogic) Table() string {
	return logic.table
}

// Columns returns all column names of current logic.
func (logic *UserLogic) Columns() UserColumns {
	return logic.columns
}

// Group returns the configuration group name of database of current logic.
func (logic *UserLogic) Group() string {
	return logic.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (logic *UserLogic) Ctx(ctx context.Context) *gdb.Model {
	return logic.DB().Model(logic.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (logic *UserLogic) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return logic.Ctx(ctx).Transaction(ctx, f)
}
