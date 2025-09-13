package api

import "github.com/Meduzz/quickapi/model"

type (
	Create struct {
		Entity any
	}

	Read struct {
		ID      string
		Preload map[string]string
	}

	Update struct {
		ID     string
		Entity any
		Hooks  []model.Hook
	}

	Delete struct {
		ID    string
		Hooks []model.Hook
	}

	Search struct {
		Skip    int
		Take    int
		Where   map[string]string
		Sort    map[string]string
		Preload map[string]string
		Hooks   []model.Hook
	}

	Patch struct {
		ID      string
		Data    map[string]any
		Preload map[string]string
		Hooks   []model.Hook
	}
)

func NewCreate(it any) *Create {
	return &Create{Entity: it}
}

func NewRead(id string, preload map[string]string) *Read {
	return &Read{ID: id, Preload: preload}
}

func NewUpate(id string, it any, hooks []model.Hook) *Update {
	return &Update{ID: id, Entity: it, Hooks: hooks}
}

func NewDelete(id string, hooks []model.Hook) *Delete {
	return &Delete{ID: id, Hooks: hooks}
}

func NewSearch(skip int, take int, where map[string]string, sort map[string]string, preload map[string]string, hooks []model.Hook) *Search {
	return &Search{Skip: skip, Take: take, Where: where, Sort: sort, Preload: preload, Hooks: hooks}
}

func NewPatch(id string, data map[string]any, preload map[string]string, hooks []model.Hook) *Patch {
	return &Patch{ID: id, Data: data, Preload: preload, Hooks: hooks}
}
