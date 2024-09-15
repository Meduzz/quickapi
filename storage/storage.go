package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	Storer interface {
		Create(any) (any, error)
		Read(string) (any, error)
		Update(any) (any, error)
		Delete(string) error
		Search(int, int, map[string]string, map[string]string, ...model.Hook) (any, error)
		Patch(string, map[string]any) (any, error)
	}

	storage struct {
		db     *gorm.DB
		entity model.Entity
	}
)

func NewStorer(db *gorm.DB, entity model.Entity) Storer {
	return &storage{db, entity}
}

func (s *storage) Create(entity any) (any, error) {
	err := s.db.Create(entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *storage) Read(id string) (any, error) {
	entity := s.entity.Create()

	err := s.db.First(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	return entity, nil
}

func (s *storage) Update(entity any) (any, error) {
	query := s.db.Session(&gorm.Session{FullSaveAssociations: true})

	err := query.Save(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	return entity, nil
}

func (s *storage) Delete(id string) error {
	entity := s.entity.Create()
	query := s.db.Select(clause.Associations)

	err := query.Delete(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return herror.ErrNotFound
		}

		return err
	}

	return nil
}

func (s *storage) Search(skip, take int, where map[string]string, sort map[string]string, hooks ...model.Hook) (any, error) {
	data := s.entity.CreateArray()

	query := s.db.
		Offset(skip).
		Limit(take)

	if len(where) > 0 {
		query = query.Where(where)
	}

	if len(sort) > 0 {
		sorting := make([]string, 0)

		for k, v := range sort {
			sorting = append(sorting, fmt.Sprintf("%s %s", k, v))
		}

		query.Order(strings.Join(sorting, ", "))
	}

	slice.ForEach(hooks, func(hook model.Hook) {
		query = query.Scopes(hook)
	})

	err := query.Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *storage) Patch(id string, data map[string]any) (any, error) {
	entity := s.entity.Create()

	err := s.db.Model(entity).
		Where("id = ?", id).
		Updates(data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	query := s.db

	err = query.Find(entity, id).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}
