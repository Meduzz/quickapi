package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	JsonTable struct {
		ID      int64          `json:"id" gorm:"autoIncrement"`
		Created int64          `json:"created" gorm:"autoCreateTime:milli"`
		Updated int64          `json:"updated" gorm:"autoUpdateTime:milli"`
		Data    datatypes.JSON `json:"data" gorm:"serializer:json"`
	}

	jsonStorage struct {
		db     *gorm.DB
		entity model.Entity
	}
)

func NewJsonStore(db *gorm.DB, entity model.Entity) Storer {
	return &jsonStorage{
		db:     db,
		entity: entity,
	}
}

func (j *jsonStorage) Create(data any) (any, error) {
	bs, err := json.Marshal(data)

	if err != nil {
		return nil, herror.ErrBadRequest
	}

	entity := &JsonTable{}
	entity.Data = bs

	err = j.db.
		Table(j.entity.Name()).
		Create(entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (j *jsonStorage) Read(id string, preload map[string]string) (any, error) {
	entity := &JsonTable{}
	err := j.db.
		Table(j.entity.Name()).
		First(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	return entity, nil
}

func (j *jsonStorage) Update(id string, data any) (any, error) {
	entity := &JsonTable{}
	err := j.db.
		Table(j.entity.Name()).
		First(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	bs, err := json.Marshal(data)

	if err != nil {
		return nil, herror.ErrBadRequest
	}

	entity.Data = bs

	err = j.db.
		Table(j.entity.Name()).
		Save(entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (j *jsonStorage) Delete(id string) error {
	entity := &JsonTable{}
	err := j.db.
		Table(j.entity.Name()).
		Delete(entity, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return herror.ErrNotFound
		}

		return err
	}

	return nil
}

func (j *jsonStorage) Search(skip int, take int, where map[string]string, sort map[string]string, preload map[string]string, hooks ...model.Hook) (any, error) {
	data := make([]*JsonTable, 0)

	query := j.db.
		Table(j.entity.Name()).
		Offset(skip).
		Limit(take)

	if len(where) > 0 {
		jsonQuery := datatypes.JSONQuery("data")

		for field, value := range where {
			jsonQuery = jsonQuery.Equals(value, field)
		}

		query = query.Where(jsonQuery)
	}

	if len(sort) > 0 {
		sorting := make([]string, 0)

		for k, v := range sort {
			if slice.Contains([]string{"id", "created", "updated"}, k) {
				sorting = append(sorting, fmt.Sprintf("%s %s", k, v))
			}
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

func (j *jsonStorage) Patch(id string, data map[string]any, preload map[string]string) (any, error) {
	entity := &JsonTable{}
	query := j.db.
		Table(j.entity.Name()).
		Model(entity).
		Where("id = ?", id)

	jsonQuery := datatypes.JSONSet("data")

	for field, value := range data {
		jsonQuery = jsonQuery.Set(field, value)
	}

	err := query.UpdateColumn("data", jsonQuery).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, herror.ErrNotFound
		}

		return nil, err
	}

	err = j.db.
		Table(j.entity.Name()).
		Find(entity, id).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}
