// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UnitDao is the data access object for table unit.
type UnitDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns UnitColumns // columns contains all the column names of Table for convenient usage.
}

// UnitColumns defines and stores column names for table unit.
type UnitColumns struct {
	Id       string // 编号
	Name     string // 名称
	Address  string // 地址
	BossId   string // 老板
	CreateAt string // 创建时间
	UpdateAt string // 修改时间
	DeleteAt string // 删除时间
}

// unitColumns holds the columns for table unit.
var unitColumns = UnitColumns{
	Id:       "id",
	Name:     "name",
	Address:  "address",
	BossId:   "boss_id",
	CreateAt: "create_at",
	UpdateAt: "update_at",
	DeleteAt: "delete_at",
}

// NewUnitDao creates and returns a new DAO object for table data access.
func NewUnitDao() *UnitDao {
	return &UnitDao{
		group:   "default",
		table:   "unit",
		columns: unitColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UnitDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UnitDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UnitDao) Columns() UnitColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UnitDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UnitDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UnitDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
