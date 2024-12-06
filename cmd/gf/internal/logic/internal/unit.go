package internal

import (
	"context"
	"gf/internal/dao"
	"gf/internal/model/entity"
	"log"
)

func init() {
	// service.RegisterUnit(NewsUnit())
}

func NewsUnit() *sUnit {
	return &sUnit{}
}

type sUnit struct {
}

func (*sUnit) CreateUnit(ctx context.Context, in entity.Unit) error {
	var unit entity.Unit
	// unit.Name = in.Name
	// unit.Height = in.Height
	// unit.Longitude = in.Longitude
	// unit.Latitude = in.Latitude
	// unit.Head = in.Head
	// unit.Pitch = in.Pitch
	// unit.Roll = in.Roll

	_, err := dao.Unit.Ctx(ctx).Insert(in)
	log.Println("unit====", unit)
	log.Println("err====", err)
	return err
}
