package storage

import (
	"fmt"

	"github.com/Meduzz/quickapi/api"
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

	Storage interface {
		Create(*api.Create) (any, error)
		Read(*api.Read) (any, error)
		Update(*api.Update) (any, error)
		Delete(*api.Delete) error
		Search(*api.Search) (any, error)
		Patch(*api.Patch) (any, error)
	}

	genericStorage struct {
		storer Storer
	}
)

var (
	_ Storage = (*genericStorage)(nil)
)

func CreateStorage(db *gorm.DB, entity model.Entity) (Storage, error) {
	var storer Storer

	if entity.Kind() == model.NormalKind {
		storer = NewStorer(db, entity)
	} else if entity.Kind() == model.JsonKind {
		storer = NewJsonStore(db, entity)
	} else {
		return nil, fmt.Errorf("unknown kind: %s", entity.Kind())
	}

	return &genericStorage{storer}, nil
}

func (gs *genericStorage) Create(create *api.Create) (any, error) {
	return gs.storer.Create(create.Entity)
}

func (gs *genericStorage) Read(read *api.Read) (any, error) {
	return gs.storer.Read(read.ID, read.Preload)
}

func (gs *genericStorage) Update(update *api.Update) (any, error) {
	return gs.storer.Update(update.ID, update.Entity)
}

func (gs *genericStorage) Delete(delete *api.Delete) error {
	return gs.storer.Delete(delete.ID)
}

func (gs *genericStorage) Search(search *api.Search) (any, error) {
	return gs.storer.Search(search.Skip, search.Take, search.Where, search.Sort, search.Preload, search.Hooks...)
}

func (gs *genericStorage) Patch(patch *api.Patch) (any, error) {
	return gs.storer.Patch(patch.ID, patch.Data, patch.Preload)
}
