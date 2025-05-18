package storage

import (
	"fmt"

	"github.com/Meduzz/quickapi/model"
	"gorm.io/gorm"
)

type (
	Storer interface {
		Create(any) (any, error)
		Read(string, map[string]string) (any, error)
		Update(string, any) (any, error)
		Delete(string) error
		Search(int, int, map[string]string, map[string]string, map[string]string, ...model.Hook) (any, error)
		Patch(string, map[string]any, map[string]string) (any, error)
	}
)

func CreateStorage(db *gorm.DB, entity model.Entity) (Storer, error) {
	if entity.Kind() == model.NormalKind {
		return NewStorer(db, entity), nil
	} else if entity.Kind() == model.JsonKind {
		return NewJsonStore(db, entity), nil
	} else {
		return nil, fmt.Errorf("unknown kind: %s", entity.Kind())
	}
}
