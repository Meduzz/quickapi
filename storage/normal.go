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
	normalStorage struct {
		db     *gorm.DB
		entity model.Entity
	}
)

func NewStorer(db *gorm.DB, entity model.Entity) Storer {
	return &normalStorage{db, entity}
}

func (s *normalStorage) Create(entity any) (any, error) {
	err := s.db.
		Table(s.entity.Name()).
		Create(entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *normalStorage) Read(id string, preload map[string]string) (any, error) {
	entity := s.entity.Create()

	query := s.preloadQuery(s.db, preload)
	err := query.
		Table(s.entity.Name()).
		First(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	return entity, nil
}

func (s *normalStorage) Update(id string, entity any) (any, error) {
	query := s.db.Session(&gorm.Session{FullSaveAssociations: true})
	err := query.
		Table(s.entity.Name()).
		Save(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	return entity, nil
}

func (s *normalStorage) Delete(id string) error {
	entity := s.entity.Create()
	query := s.db.
		Table(s.entity.Name()).
		Select(clause.Associations)

	err := query.
		Table(s.entity.Name()).
		Delete(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return herror.ErrNotFound
		}

		return err
	}

	return nil
}

func (s *normalStorage) Search(skip, take int, where map[string]string, sort map[string]string, preload map[string]string, hooks ...model.Hook) (any, error) {
	data := s.entity.CreateArray()

	query := s.db.
		Table(s.entity.Name()).
		Offset(skip).
		Limit(take)

	query = s.preloadQuery(query, preload)

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

func (s *normalStorage) Patch(id string, data map[string]any, preload map[string]string) (any, error) {
	entity := s.entity.Create()
	err := s.db.
		Table(s.entity.Name()).
		Model(entity).
		Where("id = ?", id).
		Updates(data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	query := s.preloadQuery(s.db, preload)
	err = query.
		Table(s.entity.Name()).
		Find(entity, id).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *normalStorage) preloadQuery(query *gorm.DB, preload map[string]string) *gorm.DB {
	preloadSupport, ok := s.entity.(model.PreloadSupport)

	if ok {
		// each pair represents a named preload with an optional value into a condition
		for name, conditionValue := range preload {
			// fetch the named preload from the entity
			querySpecs := preloadSupport.Preload(name)

			// for each field in the spec apply Preload
			for field, config := range querySpecs {
				if config.Converter == nil {
					config.Converter = func(s string) any { return s }
				}

				if config.Condition == "" {
					config.Condition = "1 = 1"
				}

				query = query.Preload(field, config.Condition, config.Converter(conditionValue))
			}
		}
	}

	return query
}
