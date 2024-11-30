package storage

import (
	"github.com/Meduzz/quickapi/model"
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
