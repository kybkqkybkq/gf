// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package logic

import (
	"github.com/gogf/gf/cmd/gf/v2/internal/logic/internal"
)

// internalUserLogic is internal type for wrapping internal DAO implements.
type internalUserLogic = *internal.UserLogic

// userLogic is the data access object for table user.
// You can define custom methods on it to extend its functionality as you wish.
type userLogic struct {
	internalUserLogic
}

var (
	// User is globally public accessible object for table user operations.
	User = userLogic{
		internal.NewUserLogic(),
	}
)

// Fill with you ideas below.
