package internal

import (
	"context"
	"log"

	"github.com/gogf/gf/cmd/gf/v2/internal/dao"
	"github.com/gogf/gf/cmd/gf/v2/internal/model/entity"
)

func init() {
	// service.RegisterUser(NewsUser())
}

func NewsUser() *sUser {
	return &sUser{}
}

type sUser struct {
}

func (*sUser) CreateUser(ctx context.Context, in entity.User) error {
	var user entity.User
	// user.Name = in.Name
	// user.Height = in.Height
	// user.Longitude = in.Longitude
	// user.Latitude = in.Latitude
	// user.Head = in.Head
	// user.Pitch = in.Pitch
	// user.Roll = in.Roll

	_, err := dao.User.Ctx(ctx).Insert(user)
	log.Println("user====", user)
	log.Println("err====", err)
	return err
}
